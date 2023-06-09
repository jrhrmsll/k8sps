package actor

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jrhrmsll/k8sps/pkg/internal"
	"github.com/kelindar/event"
)

type PortScanReporter struct {
	ctx             context.Context
	logger          *log.Logger
	reportsLocation string

	results chan internal.PortScanResult
	abort   chan struct{}
}

func NewPortScanReporter(
	ctx context.Context,
	logger *log.Logger,
) *PortScanReporter {
	reportsLocation := os.Getenv("PS_REPORTS_LOCATION")
	if reportsLocation == "" {
		reportsLocation = internal.PortScanDefaultReportLocation
	}

	reporter := &PortScanReporter{
		ctx:             ctx,
		logger:          logger,
		reportsLocation: reportsLocation,
		results:         make(chan internal.PortScanResult),
		abort:           make(chan struct{}),
	}

	event.Subscribe(internal.PubSub, func(e internal.PortScanResultEvent) {
		reporter.results <- e.PortScanResult
	})

	return reporter
}

func (reporter *PortScanReporter) store(result internal.PortScanResult) error {
	file, err := os.Create(reporter.reportsLocation)
	if err != nil {
		return err
	}

	defer file.Close()

	for node, ports := range result.NodesPorts {
		_, err := file.WriteString(fmt.Sprintf("%s: [%s]\n", node, strings.Join(ports, ",")))
		if err != nil {
			return err
		}
	}

	return nil
}

func (reporter *PortScanReporter) Start() error {
	for {
		select {
		case result := <-reporter.results:
			reporter.logger.Println("saving portscan report")

			err := reporter.store(result)
			if err != nil {
				reporter.logger.Println(err)
			}

		case <-reporter.abort:
			return nil
		}
	}
}

func (reporter *PortScanReporter) Stop(err error) {
	close(reporter.abort)
}
