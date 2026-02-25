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

	_ = r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()

		if path != "" && len(methods) > 0 {
			log.Info().
				Str("path", path).
				Strs("methods", methods).
				Msg("registered route")
		}
		return nil
	})

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
