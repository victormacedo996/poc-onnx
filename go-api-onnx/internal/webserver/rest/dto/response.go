package dto

type HealthResponse struct {
	Message string `json:"message"`
}

type ModelPredictResponse struct {
	Prediction string `json:"prediction"`
}
