package scanner

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Ullaakut/nmap/v3"
)

type NmapPortScanner struct {
	portWhitelist []string
}

func NewNmapPortScanner(portWhitelist []string) *NmapPortScanner {
	return &NmapPortScanner{
		portWhitelist: portWhitelist,
	}
}

func (ps *NmapPortScanner) Scan(ips []string) (map[string]struct{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	scanner, err := nmap.NewScanner(
		ctx,
		nmap.WithTargets(ips...),
		nmap.WithPorts("0-65535"),
	)

	if len(ps.portWhitelist) > 0 {
		scanner.AddOptions(nmap.WithPortExclusions(ps.portWhitelist...))
	}

	if err != nil {
		log.Fatalf("unable to create nmap scanner: %v", err)
	}

	result, warnings, err := scanner.Run()
	if len(*warnings) > 0 {
		log.Printf("run finished with warnings: %s\n", *warnings) // Warnings are non-critical errors from nmap.
	}
	if err != nil {
		log.Fatalf("unable to run nmap scan: %v", err.Error())
	}

	entries := make(map[string]struct{})
	for _, host := range result.Hosts {
		if len(host.Ports) == 0 || len(host.Addresses) == 0 {
			continue
		}

		for _, port := range host.Ports {
			entries[fmt.Sprintf("%d", port.ID)] = struct{}{}
		}
	}

	return entries, nil
}
