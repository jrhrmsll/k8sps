package actor

import (
	"context"
	"log"
	"sync"

	"github.com/jrhrmsll/k8sps/pkg/internal"
	"github.com/jrhrmsll/k8sps/pkg/util"
	"github.com/kelindar/event"
)

type PortScanRunner struct {
	ctx     context.Context
	logger  *log.Logger
	scanner internal.PortScanner

	requests chan internal.PortScanRequest
	abort    chan struct{}

	busy bool
	mu   sync.RWMutex
}

func NewPortScanRunner(
	ctx context.Context,
	logger *log.Logger,
	scanner internal.PortScanner,
) *PortScanRunner {
	requests := make(chan internal.PortScanRequest, 1)

	runner := &PortScanRunner{
		ctx:      ctx,
		logger:   logger,
		scanner:  scanner,
		requests: requests,
		abort:    make(chan struct{}),
	}

	event.Subscribe(internal.PubSub, func(e internal.PortScanRequestEvent) {
		// if busy discard jobs
		if runner.isBusy() {
			return
		}

		requests <- e.PortScanRequest
	})

	return runner
}

func (runner *PortScanRunner) isBusy() bool {
	runner.mu.RLock()
	defer runner.mu.RUnlock()

	return runner.busy
}

func (runner *PortScanRunner) lock() {
	runner.mu.Lock()
	defer runner.mu.Unlock()

	runner.busy = true
	event.Publish(internal.PubSub, internal.PortScanRunnerBusyEvent{Busy: true})
}

func (runner *PortScanRunner) unlock() {
	runner.mu.Lock()
	defer runner.mu.Unlock()

	runner.busy = false
	event.Publish(internal.PubSub, internal.PortScanRunnerBusyEvent{Busy: false})
}

func (runner *PortScanRunner) Start() error {
	for {
		select {
		case request := <-runner.requests:
			runner.lock()

			entries := make(map[string][]string)
			for node, ips := range request.NodesIPAddresses {
				runner.logger.Printf("scanning ports from '%s' node\n", node)

				ports, err := runner.scanner.Scan(ips)
				if err != nil {
					runner.logger.Println(err)
				}

				entries[node] = util.MapToSlice(ports)
			}

			event.Publish(internal.PubSub, internal.PortScanResultEvent{
				PortScanResult: internal.PortScanResult{
					NodesPorts: entries,
				},
			})

			runner.unlock()
		case <-runner.abort:
			return nil
		}
	}
}

func (runner *PortScanRunner) Stop(err error) {
	close(runner.abort)
}
