package build

import (
	"fmt"
	"net/http"

	mdlwr "github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	promGo "github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	exporters "go.opentelemetry.io/otel/exporters/prometheus"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"go.opentelemetry.io/otel/trace"
)

var DurationBucketsInMilliseconds = []float64{1, 5, 10, 25, 50, 100, 250, 500, 1000, 3000, 5000, 10000, 20000}

type restServer interface {
	Context() *mdlwr.Context
}

func NewRouter(
	appName string,
	reg promGo.Registerer,
	router *mux.Router,
	api restServer,
) (func(next http.Handler) http.Handler, error) {
	provider, err := prometheusProvider(appName, reg)
	if err != nil {
		return nil, fmt.Errorf("creating prometheus provider: %w", err)
	}

	router.Use(otelhttp.NewMiddleware(
		appName,
		otelhttp.WithMeterProvider(provider),
	))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attr := semconv.HTTPRouteKey.String(pathTemplate(r, api))

			span := trace.SpanFromContext(r.Context())
			span.SetAttributes(attr)

			if labeler, ok := otelhttp.LabelerFromContext(r.Context()); ok {
				labeler.Add(attr)
			}

			next.ServeHTTP(w, r)
		})
	}, nil
}

func prometheusProvider(appName string, reg promGo.Registerer) (*sdkmetric.MeterProvider, error) {
	resources, err := resource.Merge(resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(appName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("creating new resource: %w", err)
	}

	exporter, err := exporters.New(
		exporters.WithNamespace(appName),
		exporters.WithRegisterer(reg),
	)
	if err != nil {
		return nil, fmt.Errorf("creating prometheus exporter: %w", err)
	}

	//nolint:exhaustruct
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resources),
		sdkmetric.WithReader(exporter),
		sdkmetric.WithView(sdkmetric.NewView(
			sdkmetric.Instrument{
				Name:        appName,
				Description: "processing duration",
				Kind:        sdkmetric.InstrumentKindHistogram,
			},
			sdkmetric.Stream{
				Aggregation: sdkmetric.AggregationExplicitBucketHistogram{
					Boundaries: DurationBucketsInMilliseconds,
					NoMinMax:   false,
				},
			},
		)),
	)

	return provider, nil
}

func pathTemplate(r *http.Request, api restServer) string {
	uri := r.URL.EscapedPath()

	if route, _, ok := api.Context().RouteInfo(r); ok {
		uri = route.PathPattern
	}

	return uri
}
