package main

import (
	"context"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/jrhrmsll/k8sps/pkg/actor"
	"github.com/jrhrmsll/k8sps/pkg/scanner"
	"github.com/oklog/run"
)

func main() {
	ctx := context.Background()

	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	logger.Println("Starting Port Scanner")

	actors := run.Group{}
	actors.Add(run.SignalHandler(ctx, os.Interrupt, os.Kill, syscall.SIGTERM))

	httpApplication := actor.NewHttpApplication(ctx, logger)
	actors.Add(httpApplication.Start, httpApplication.Stop)

	portWhitelist := []string{}
	v := os.Getenv("PS_PORT_WHITELIST")
	if v != "" {
		portWhitelist = strings.Split(v, ",")
	}

	portScanner := scanner.NewNmapPortScanner(portWhitelist)
	portScanRunner := actor.NewPortScanRunner(ctx, logger, portScanner)
	actors.Add(portScanRunner.Start, portScanRunner.Stop)

	portScanReporter := actor.NewPortScanReporter(ctx, logger)
	actors.Add(portScanReporter.Start, portScanReporter.Stop)

	err := actors.Run()
	if err != nil {
		logger.Println(err)
	}
}
