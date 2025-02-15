package dto

type PredictRequest struct {
	Sepal_length float32 `json:"sepal_length"`
	Sepal_width  float32 `json:"sepal_width"`
	Petal_length float32 `json:"petal_length"`
	Petal_width  float32 `json:"petal_width"`
}
