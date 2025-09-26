package build

import (
	"context"
	"fmt"
	"net/http"

	"otusdelivery/internal/restapi"
	"otusdelivery/internal/restapi/operations"
	deliveryOps "otusdelivery/internal/restapi/operations/delivery"
	"otusdelivery/internal/restapi/operations/other"
	"otusdelivery/internal/service/api/delivery"

	httpMW "otusdelivery/pkg/http"

	"github.com/go-openapi/loads"
	mdlwr "github.com/go-openapi/runtime/middleware"
	"github.com/pkg/errors"
)

func (b *Builder) buildAPI() (*operations.RestServerAPI, *loads.Document, error) {
	swaggerSpec, err := loads.Spec("api/swagger/file.yaml")
	if err != nil {
		return nil, nil, fmt.Errorf("load swagger specs: %w", err)
	}

	api := operations.NewRestServerAPI(swaggerSpec)

	psql, err := b.PostgresClient()
	if err != nil {
		return nil, nil, fmt.Errorf("creating postgres client: %w", err)
	}

	repo := b.NewRepo(psql.DB)

	deliverySrv := delivery.NewService(psql.DB, repo)

	handler := restapi.NewHandler(deliverySrv)

	api.OtherGetHealthHandler = other.GetHealthHandlerFunc(
		handler.GetHealth,
	)

	api.DeliveryGetCheckDeliveryStatusOrderIDHandler = deliveryOps.GetCheckDeliveryStatusOrderIDHandlerFunc(
		handler.CheckOrderDelivery,
	)

	api.DeliveryPostReserveSlotHandler = deliveryOps.PostReserveSlotHandlerFunc(
		handler.ReserveSlotForOrder,
	)

	api.DeliveryDeleteUnreserveSlotOrderIDHandler = deliveryOps.DeleteUnreserveSlotOrderIDHandlerFunc(
		handler.UnreserveSlotForOrder,
	)

	return api, swaggerSpec, nil
}

//nolint:funlen
func (b *Builder) RestAPIServer(ctx context.Context) (*http.Server, error) {
	server, err := b.HTTPServer(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating http server: %w", err)
	}

	router := b.httpRouter()

	api, swaggerSpec, err := b.buildAPI()
	if err != nil {
		return nil, errors.Wrap(err, "building API")
	}

	apiEndpoint := swaggerSpec.BasePath()
	apiRouter := router.Name("api").Subrouter()

	metricsMW, err := NewRouter(b.config.App.Name, b.prometheusRegistry, router, api)
	if err != nil {
		return nil, fmt.Errorf("creating metrics middleware: %w", err)
	}

	apiRouter.Use(metricsMW)
	apiRouter.Use(httpMW.HTTPRequestBodyLoggerWithContext(ctx))

	swaggerUIOpts := mdlwr.SwaggerUIOpts{ //nolint:exhaustruct
		BasePath: apiEndpoint,
		SpecURL:  fmt.Sprintf("%s/swagger.json", apiEndpoint),
	}

	apiRouter.PathPrefix(apiEndpoint).Handler(
		func() http.Handler {
			api.Init()

			return mdlwr.Spec(
				apiEndpoint,
				swaggerSpec.Raw(),
				mdlwr.SwaggerUI(
					swaggerUIOpts,
					api.Context().RoutesHandler(metricsMW),
				),
			)
		}(),
	)

	return server, nil
}
