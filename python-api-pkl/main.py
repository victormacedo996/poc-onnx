from fastapi import FastAPI
import uvicorn
import os
from dto.response import HealthResponse, PredictResponse
from dto.request import PredictRequest
import pickle
from sklearn.linear_model import LogisticRegression
from opentelemetry.instrumentation.fastapi import FastAPIInstrumentor
from opentelemetry.sdk.resources import SERVICE_NAME, Resource
from opentelemetry import metrics, trace
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.metrics import MeterProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.sdk.metrics.export import PeriodicExportingMetricReader
from opentelemetry.exporter.otlp.proto.grpc.metric_exporter import OTLPMetricExporter
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter

app = FastAPI()

WEBSERVER_PORT = int(os.getenv("WEBSERVER_PORT", 5000))
OTEL_SERVICE_NAME = os.getenv("OTEL_SERVICE_NAME", "poc-python-api-pkl")
OTEL_COLLECTOR_OTLP_ENDPOINT = os.getenv("OTEL_COLLECTOR_OTLP_ENDPOINT", "otelcol:4317")

otel_resource = Resource(attributes={
    SERVICE_NAME: OTEL_SERVICE_NAME,
})

otlp_exporter = OTLPMetricExporter(endpoint=f"http://{OTEL_COLLECTOR_OTLP_ENDPOINT}")
otlp_metric_reader = PeriodicExportingMetricReader(otlp_exporter, export_interval_millis=3000)
meter_provider = MeterProvider(resource=otel_resource, metric_readers=[otlp_metric_reader])
metrics.set_meter_provider(meter_provider)
meter = metrics.get_meter(__name__)
meter_provider = meter_provider

otlp_span_exporter = OTLPSpanExporter(endpoint=f"http://{OTEL_COLLECTOR_OTLP_ENDPOINT}",)
span_processor = BatchSpanProcessor(otlp_span_exporter)
trace_provider = TracerProvider(resource=otel_resource)
trace_provider.add_span_processor(span_processor)
trace.set_tracer_provider(tracer_provider=trace_provider)
tracer = trace.get_tracer(__name__)
tracer_provider = trace_provider

PREDICTION_MAPPING = {
    0: "iris-setosa",
    1: "iris-versicolor",
    2: "iris-virginica"
}

with open("../models/logistic_regression_iris.pkl", 'rb') as file:
    model: LogisticRegression = pickle.load(file)

@app.get("/health")
async def health():
    return HealthResponse(message="ok")


@app.post('/predict')
async def predict(model_input: PredictRequest) -> PredictResponse:
    prediction = model.predict([[model_input.sepal_length, model_input.sepal_width, model_input.petal_length, model_input.petal_width]])
    return PredictResponse(prediction=PREDICTION_MAPPING[prediction[0]])

FastAPIInstrumentor.instrument_app(app, tracer_provider=trace_provider, meter_provider=meter_provider)

if __name__ == '__main__':
    uvicorn.run(app, port=WEBSERVER_PORT, host="0.0.0.0")
