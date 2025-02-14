from fastapi import FastAPI
import uvicorn
import os
from dto.response import HealthResponse, PredictResponse
from dto.request import PredictRequest
import pickle

app = FastAPI()

WEBSERVER_PORT = int(os.getenv("WEBSERVER_PORT", 5000))

PREDICTION_MAPPING = {
    0: "iris-setosa",
    1: "iris-versicolor",
    2: "iris-virginica"
}

with open("../models/logistic_regression_iris.pkl", 'rb') as file:
    model = pickle.load(file)


@app.get("/health")
async def health():
    return HealthResponse(message="ok")


@app.post('/predict')
async def predict(model_input: PredictRequest) -> PredictResponse:
    prediction = model.predict([[model_input.sepal_length, model_input.sepal_width, model_input.petal_length, model_input.petal_width]])
    return PredictResponse(prediction=PREDICTION_MAPPING[prediction[0]])


if __name__ == '__main__':
    uvicorn.run(app, port=WEBSERVER_PORT, host="0.0.0.0")
