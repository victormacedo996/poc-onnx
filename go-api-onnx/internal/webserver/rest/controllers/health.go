package controllers

import (
	"net/http"

	"github.com/victormacedo996/poc-onnx/internal/webserver/rest/dto"
	"github.com/victormacedo996/poc-onnx/internal/webserver/rest/response"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (h *HealthController) StatusOk(w http.ResponseWriter, r *http.Request) {
	resp := dto.HealthResponse{Message: "ok"}
	response.StatusOk(w, r, resp)

}
func (h *HealthController) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	response.StatusMethodNotAllowed(w, r)
}
