package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/victormacedo996/poc-onnx/internal/onnx"
	"github.com/victormacedo996/poc-onnx/internal/webserver/rest/dto"
	"github.com/victormacedo996/poc-onnx/internal/webserver/rest/response"
	"github.com/yalue/onnxruntime_go"
)

type ModelController struct {
	sess *onnxruntime_go.DynamicAdvancedSession
}

func NewModelController(sess *onnxruntime_go.DynamicAdvancedSession) *ModelController {
	return &ModelController{
		sess: sess,
	}
}

func (m *ModelController) Predict(w http.ResponseWriter, r *http.Request) {
	var incoming_prediction dto.PredictRequest

	err := json.NewDecoder(r.Body).Decode(&incoming_prediction)
	if err != nil {
		response.StatusUnprocessableEntity(w, r, fmt.Errorf("error parsing JSON"))
		return
	}
	defer r.Body.Close()

	prediction, err := onnx.Predict([]float32{incoming_prediction.Sepal_length, incoming_prediction.Sepal_width, incoming_prediction.Petal_length, incoming_prediction.Petal_width}, m.sess, 1, 4)
	if err != nil {
		response.StatusInternalServerError(w, r, err)
		return
	}

	response.StatusOk(w, r, dto.ModelPredictResponse{Prediction: prediction})
}
func (m *ModelController) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	response.StatusMethodNotAllowed(w, r)
}
