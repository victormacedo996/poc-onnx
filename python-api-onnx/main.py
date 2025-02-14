from fastapi import FastAPI
import uvicorn
import os
from dto.response import HealthResponse, PredictResponse
from dto.request import PredictRequest
import onnxruntime as rt

app = FastAPI()

WEBSERVER_PORT = int(os.getenv("WEBSERVER_PORT", 5000))

PREDICTION_MAPPING = {
    0: "iris-setosa",
    1: "iris-versicolor",
    2: "iris-virginica"
}

sess = rt.InferenceSession("../models/logistic_regression_iris.onnx")
input_name = sess.get_inputs()[0].name

def onnx_predict(sepal_length: float, sepal_width: float, petal_length: float, petal_width: float) -> str:
    prediction = sess.run(None, {input_name: [[sepal_length, sepal_width, petal_length, petal_width]]})
    return PREDICTION_MAPPING[prediction[0].tolist()[0]]


@app.get("/health")
async def health():
    return HealthResponse(message="ok")


@app.post('/predict')
async def predict(model_input: PredictRequest) -> PredictResponse:
    prediction = onnx_predict(model_input.sepal_length, model_input.sepal_width, model_input.petal_length, model_input.petal_width)
    return PredictResponse(prediction=prediction)


if __name__ == '__main__':
    uvicorn.run(app, port=WEBSERVER_PORT, host="0.0.0.0")
