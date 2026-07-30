package rprocessor

import (
	"net/http"

	rhandler "github.com/andreyloginov-afk/catalog-service/internal/app/handler/http"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func vGenericRegHealthCheck(r *mux.Router, h rhandler.Health) {
	reg(r, http.MethodGet, "/health", http.HandlerFunc(h.LastCheck))
	reg(r, http.MethodGet, "/metrics", promhttp.Handler())
}

func handlerNotFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
