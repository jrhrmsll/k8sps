package controller

import (
	"log"
	"net/http"
	"os"

	"github.com/jrhrmsll/k8sps/pkg/internal"
	"github.com/jrhrmsll/k8sps/pkg/k8s"
	"github.com/kelindar/event"
	"github.com/labstack/echo"
)

type PortScanController struct {
	logger             *log.Logger
	kubernetesClient   *k8s.KubernetesClient
	portScanRunnerBusy bool
	reportsLocation    string
}

func NewPortScanController(
	logger *log.Logger,
	kubernetesClient *k8s.KubernetesClient,
) *PortScanController {
	reportsLocation := os.Getenv("PS_REPORTS_LOCATION")
	if reportsLocation == "" {
		reportsLocation = internal.PortScanDefaultReportLocation
	}

	controller := &PortScanController{
		logger:           logger,
		kubernetesClient: kubernetesClient,
		reportsLocation:  reportsLocation,
	}

	event.Subscribe(internal.PubSub, func(e internal.PortScanRunnerBusyEvent) {
		controller.portScanRunnerBusy = e.Busy
	})

	return controller
}

func (controller *PortScanController) Scan(c echo.Context) error {
	if controller.portScanRunnerBusy {
		controller.logger.Println("server busy")
		return echo.NewHTTPError(http.StatusTooManyRequests)
	}

	nodesIPAddresses, err := controller.kubernetesClient.NodesIPAddresses()
	if err != nil {
		controller.logger.Printf("kubernetes client error: %s\n", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	controller.logger.Println("sending portscan request")

	event.Publish(internal.PubSub, internal.PortScanRequestEvent{
		PortScanRequest: internal.PortScanRequest{
			NodesIPAddresses: nodesIPAddresses,
		},
	})

	return c.NoContent(http.StatusNoContent)
}

func (controller *PortScanController) Report(c echo.Context) error {
	return c.File(controller.reportsLocation)
}
