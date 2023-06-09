package internal

const (
	EventTypePortScanRunnerBusy = 1
	EventTypePortScanRequest    = 2
	EventTypePortScanResult     = 3
)

type PortScanRunnerBusyEvent struct {
	Busy bool
}

func (event PortScanRunnerBusyEvent) Type() uint32 {
	return EventTypePortScanRunnerBusy
}

type PortScanRequestEvent struct {
	PortScanRequest
}

func (event PortScanRequestEvent) Type() uint32 {
	return EventTypePortScanRequest
}

type PortScanResultEvent struct {
	PortScanResult
}

func (event PortScanResultEvent) Type() uint32 {
	return EventTypePortScanResult
}
