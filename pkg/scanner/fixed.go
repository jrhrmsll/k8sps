package scanner

type FixedPortScanner struct {
	portWhitelist []string
}

func NewFixedPortScanner(portWhitelist []string) *FixedPortScanner {
	return &FixedPortScanner{
		portWhitelist: portWhitelist,
	}
}

func (ps *FixedPortScanner) Scan(ips []string) (map[string]struct{}, error) {
	ports := map[string]struct{}{
		"8080": {},
		"53":   {},
		"9100": {},
		"5432": {},
		"22":   {},
	}

	// remove whitelisted ports
	for _, ignore := range ps.portWhitelist {
		delete(ports, ignore)
	}

	return ports, nil
}
