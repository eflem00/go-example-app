package http

import (
	"net/http"

	"github.com/eflem00/go-example-app/common"
	"github.com/go-chi/chi/v5"
)

type HttpController struct {
	settings      *common.Settings
	logger        *common.Logger
	resultHandler *ResultHandler
	healthHandler *HealthHandler
}

func NewHttpController(settings *common.Settings, logger *common.Logger, resultHandler *ResultHandler, healthHandler *HealthHandler) *HttpController {
	return &HttpController{
		settings,
		logger,
		resultHandler,
		healthHandler,
	}
}

func (controller *HttpController) Start() error {
	controller.logger.Info("Starting http controller")

	r := chi.NewRouter()
	r.Get("/", controller.healthHandler.Health)
	r.Get("/health", controller.healthHandler.Health)
	r.Get("/results/{id}", controller.resultHandler.GetResultById)

	port := controller.settings.Port

	controller.logger.Infof("listening on %v", port)

	err := http.ListenAndServe(port, r)

	controller.logger.Err(err, "error in http controller")

	return err
}

func (controller *HttpController) Exit() {
	controller.logger.Error("detected exit in http controller")
}
