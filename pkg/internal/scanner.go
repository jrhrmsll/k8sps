package internal

type PortScanner interface {
	Scan(ips []string) (map[string]struct{}, error)
}
