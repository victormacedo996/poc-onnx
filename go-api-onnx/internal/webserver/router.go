package webserver

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/victormacedo996/poc-onnx/internal/config"
	"github.com/victormacedo996/poc-onnx/internal/webserver/rest/controllers"
)

type Router struct {
	router *chi.Mux
	config config.Config
}

func NewRouter(config config.Config) *Router {

	return &Router{
		router: chi.NewRouter(),
		config: config,
	}
}

func (r *Router) configureRoutes(health_controller controllers.HealthController, model_controller controllers.ModelController) {
	r.router.Use(middleware.Logger)
	r.router.Use(middleware.Recoverer)
	r.router.Get("/health", health_controller.StatusOk)
	r.router.Post("/predict", model_controller.Predict)
}

func (r *Router) Start(health_controller controllers.HealthController, model_controller controllers.ModelController) {
	r.configureRoutes(health_controller, model_controller)
	port := fmt.Sprintf(":%v", r.config.WEBSERVER_PORT)
	err := http.ListenAndServe(port, r.router)
	if err != nil {
		panic(err)
	}

}
