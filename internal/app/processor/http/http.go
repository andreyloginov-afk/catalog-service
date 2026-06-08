package rprocessor

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/andreyloginov-afk/catalog-service/internal/app/config/section"
	rhandler "github.com/andreyloginov-afk/catalog-service/internal/app/handler/http"
	"github.com/andreyloginov-afk/catalog-service/internal/app/processor"
	"github.com/andreyloginov-afk/catalog-service/internal/app/util"
	"github.com/andreyloginov-afk/catalog-service/internal/pkg/http/httph"
	"github.com/andreyloginov-afk/catalog-service/internal/pkg/http/mzerolog"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type httpProc struct {
	server http.Server
	addr   string
}

type gracefulServer struct {
	srv *http.Server
}

func (gs gracefulServer) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return gs.srv.Shutdown(ctx)
}

func NewHttp(hHealth rhandler.Health, hCategory rhandler.Category, hProduct rhandler.Product, cfg section.ProcessorWebServer) processor.Processor {

	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(handlerNotFound)

	debugLogger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()
	debugLogger.Debug().Msg("Custom logger for mzerolog") // это сообщение вы должны увидеть

	r.Use(
		httph.NewErrorMiddlewear(),
		mzerolog.NewMiddleware(
			mzerolog.WithLogger(debugLogger),
			mzerolog.WithSkipper(util.IsFilteredHttpRoute),
		),
	)

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
func (p *httpProc) StartAsync(ctx context.Context, wg *sync.WaitGroup) {

	l, err := (&net.ListenConfig{}).Listen(ctx, "tcp", p.addr)
	if err != nil {
		log.Fatal().Err(err).Msg("server cannot start without listener")
	}
	log.Info().Str("addr", p.addr).Msg("HTTP server listening on")

	go p.serve(l)

	processor.WatchForShutdown(ctx, wg, l)

	processor.WatchForShutdown(ctx, wg, gracefulServer{srv: &p.server})
}

func (p *httpProc) serve(l net.Listener) {
	_ = p.server.Serve(l)
}
