package build

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

func (b *Builder) HTTPServer(ctx context.Context) (*http.Server, error) {
	const timeout = time.Millisecond * 25

	router := b.httpRouter()
	router.Handle(metricsEndpoint, promhttp.HandlerFor(b.prometheus(), promhttp.HandlerOpts{})) //nolint:exhaustruct

	//nolint:exhaustruct
	server := http.Server{
		Addr:              b.config.HTTPAddr(),
		ReadHeaderTimeout: timeout,
		Handler:           router,
		ErrorLog:          log.New(zerolog.Nop(), "", 0),
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	return &server, nil
}

func (b *Builder) httpRouter() *mux.Router {
	if b.http.router != nil {
		return b.http.router
	}

	b.http.router = mux.NewRouter()

	return b.http.router
}
