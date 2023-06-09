package internal

const PortScanDefaultReportLocation = "/mnt/portscan/report"

type PortScanRequest struct {
	NodesIPAddresses map[string][]string
}

type PortScanResult struct {
	NodesPorts map[string][]string
}
