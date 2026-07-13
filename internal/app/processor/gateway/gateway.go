package pgateway

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/andreyloginov-afk/catalog-service/internal/app/config/section"
	"github.com/andreyloginov-afk/catalog-service/internal/app/processor"
	catalogv1 "github.com/andreyloginov-afk/catalog-service/internal/pkg/grpc/gen/catalog/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	gatewayReadHeaderTimeout = 5 * time.Second
	gatewayReadTimeout       = 30 * time.Second
	gatewayWriteTimeout      = 30 * time.Second
	gatewayIdleTimeout       = 120 * time.Second
	gatewayShutdownTimeout   = 5 * time.Second
)

type gatewayProc struct {
	server   http.Server
	addr     string
	grpcAddr string
}

func NewGateway(
	cfgGateway section.ProcessorGateway,
	cfgGrpc section.ProcessorGrpc,
) processor.Processor {
	addr := fmt.Sprintf(":%d", cfgGateway.ListenPort)
	grpcAddr := net.JoinHostPort("localhost", strconv.Itoa(int(cfgGrpc.ListenPort)))
	return &gatewayProc{
		addr:     addr,
		grpcAddr: grpcAddr,
		server: http.Server{
			Addr:              addr,
			ReadHeaderTimeout: gatewayReadHeaderTimeout,
			ReadTimeout:       gatewayReadTimeout,
			WriteTimeout:      gatewayWriteTimeout,
			IdleTimeout:       gatewayIdleTimeout,
		},
	}
}

func (p *gatewayProc) StartAsync(ctx context.Context, wg *sync.WaitGroup) {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := catalogv1.RegisterCatalogServiceHandlerFromEndpoint(ctx, mux, p.grpcAddr, opts); err != nil {
		log.Fatal().Err(err).Msg("failed to register gateway handler")
		return
	}
	p.server.Handler = mux

	l, err := (&net.ListenConfig{}).Listen(ctx, "tcp", p.addr)
	if err != nil {
		log.Fatal().Err(err).Msg("gateway cannot start without listener")
		return
	}
	log.Info().Str("addr", p.addr).Msg("gateway server listening on")

	go func() {
		if err := p.server.Serve(l); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("gateway server stopped")
		}
	}()

	processor.WatchForShutdown(ctx, wg, l)
	processor.WatchForShutdown(ctx, wg, processor.NewCloserContextFunc(p.server.Shutdown, ctx, gatewayShutdownTimeout))
}
