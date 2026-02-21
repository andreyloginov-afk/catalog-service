package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/andreyloginov-afk/catalog-service/internal/app/config/section"
	rhandler "github.com/andreyloginov-afk/catalog-service/internal/app/handler"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type httpProc struct {
	server http.Server
	addr   string
}

func NewHttp(hHealth rhandler.Health, cfg section.ProcessorWebServer) *httpProc {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(handlerNotFound)

	vGenericRegHealthCheck(r, hHealth)

	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			//nolint:nilerr // пропускаем маршруты без шаблона пути
			return nil
		}
		methods, err := route.GetMethods()
		if err != nil {
			methods = []string{"ANY"}
		}

		log.Info().
			Str("path", path).
			Strs("methods", methods).
			Msg("registered route")

		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to walk routes")
	}

	addr := fmt.Sprintf(":%d", cfg.ListenPort)

	p := httpProc{
		addr: addr,
		server: http.Server{
			Addr:              addr,
			Handler:           r,
			ReadHeaderTimeout: 5 * time.Second,
		},
	}

	return &p
}

func (p *httpProc) Serve() error {
	log.Info().
		Str("addr", p.addr).
		Msg("starting HTTP server")

	return p.server.ListenAndServe()
}
