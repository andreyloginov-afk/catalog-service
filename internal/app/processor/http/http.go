package rprocessor

import (
	"fmt"
	"net/http"
	"time"

	"github.com/andreyloginov-afk/catalog-service/internal/app/config/section"
	rhandler "github.com/andreyloginov-afk/catalog-service/internal/app/handler/http"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type httpProc struct {
	server http.Server
	addr   string
}

func NewHttp(hHealth rhandler.Health, hCategory rhandler.Category, hProduct rhandler.Product, cfg section.ProcessorWebServer) *httpProc {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(handlerNotFound)
	vGenericRegHealthCheck(r, hHealth)
	// API v1
	rV1 := r.PathPrefix("/v1").Subrouter()

	if hCategory != nil {
		v1RegCategoryHandler(rV1, hCategory)
	}
	if hProduct != nil {
		v1RegProductHandler(rV1, hProduct)
	}

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
