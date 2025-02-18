package webserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/victormacedo996/poc-onnx/internal/config"
	"github.com/victormacedo996/poc-onnx/internal/webserver/rest/controllers"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
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
	r.router.Use(otelhttp.NewMiddleware("my-server"))
	r.router.Use(middleware.Logger)
	r.router.Use(middleware.Recoverer)
	r.router.Get("/health", health_controller.StatusOk)
	r.router.Post("/predict", model_controller.Predict)
}

func (r *Router) Start(health_controller controllers.HealthController, model_controller controllers.ModelController) {
	meter_provider, trace_provider := r.configureOtel()
	ctx := context.Background()
	defer func() {
		if err := trace_provider.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()
	defer func() {
		if err := meter_provider.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()
	r.configureRoutes(health_controller, model_controller)
	port := fmt.Sprintf(":%v", r.config.WEBSERVER_PORT)
	err := http.ListenAndServe(port, r.router)
	if err != nil {
		panic(err)
	}

}

func (r *Router) configureOtel() (*trace.TracerProvider, *metric.MeterProvider) {
	ctx := context.Background()

	// Initialize OpenTelemetry Resource
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("my-chi-app"),
		),
	)
	if err != nil {
		panic(err)
	}

	// Initialize Trace Exporter
	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("localhost:5000"),
	)
	if err != nil {
		panic(err)
	}

	// Create Trace Provider
	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter),
		trace.WithResource(res),
	)

	// Set global Trace Provider
	otel.SetTracerProvider(traceProvider)

	// Initialize Metrics Exporter
	metricExporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint("localhost:5000"),
	)
	if err != nil {
		panic(err)
	}

	// Create Meter Provider
	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(metricExporter)),
	)

	// Set global Meter Provider
	otel.SetMeterProvider(meterProvider)

	// Set global propagator to tracecontext (to propagate trace headers)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return traceProvider, meterProvider
}
