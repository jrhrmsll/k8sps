package actor

import (
	"context"
	"log"
	"net/http"

	"github.com/jrhrmsll/k8sps/pkg/controller"
	"github.com/jrhrmsll/k8sps/pkg/k8s"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type HttpApplication struct {
	ctx    context.Context
	srv    *http.Server
	logger *log.Logger
	agent  bool
}

func NewHttpApplication(
	ctx context.Context,
	logger *log.Logger,
) *HttpApplication {

	app := echo.New()
	{
		app.Use(middleware.Recover())
	}

	internal := app.Group("internal")
	{
		internal.GET("/health", func(c echo.Context) error {
			return c.JSON(http.StatusOK, "ok")
		})
	}

	api := app.Group("api")
	{
		api.Use(middleware.Logger())

		kubernetesClient, err := k8s.NewKubernetesClient()
		if err != nil {
			log.Fatal(err)
		}

		portScanController := controller.NewPortScanController(logger, kubernetesClient)

		api.POST("/scan", portScanController.Scan)
		api.GET("/report", portScanController.Report)
	}

	return &HttpApplication{
		ctx: ctx,
		srv: &http.Server{
			Addr:    ":8080",
			Handler: app,
		},
		logger: logger,
		agent:  true,
	}
}

func (app *HttpApplication) Start() error {
	err := app.srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (app *HttpApplication) Stop(error) {
	err := app.srv.Shutdown(app.ctx)
	if err != nil {
		log.Println(err)
	}
}
