package build

import (
	"context"
	"net/http"

	"otusdelivery/config"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

type Builder struct {
	config config.Config

	shutdown shutdown

	prometheusRegistry *prometheus.Registry

	http struct {
		router *mux.Router
		server *http.Server
	}
}

func New(ctx context.Context, conf config.Config) *Builder {
	b := Builder{config: conf} //nolint:exhaustruct

	return &b
}
