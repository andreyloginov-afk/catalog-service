package pgrpc

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/andreyloginov-afk/catalog-service/internal/app/config/section"
	"github.com/andreyloginov-afk/catalog-service/internal/app/processor"
	catalogv1 "github.com/andreyloginov-afk/catalog-service/internal/pkg/grpc/gen/catalog/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcProc struct {
	server *grpc.Server
	addr   string
}

func NewGRPC(
	catalogV1 catalogv1.CatalogServiceServer,
	cfg section.ProcessorGrpc,
) processor.Processor {
	srv := grpc.NewServer()
	catalogv1.RegisterCatalogServiceServer(srv, catalogV1)
	reflection.Register(srv)
	return &grpcProc{server: srv, addr: fmt.Sprintf(":%d", cfg.ListenPort)}
}

func (p *grpcProc) StartAsync(ctx context.Context, wg *sync.WaitGroup) {
	l, err := (&net.ListenConfig{}).Listen(ctx, "tcp", p.addr)
	if err != nil {
		log.Fatal().Err(err).Msg("gRPC server cannot start without listener")
	}
	log.Info().Str("addr", p.addr).Msg("gRPC server listening on")

	go func() {
		if err := p.server.Serve(l); err != nil {
			log.Error().Err(err).Msg("gRPC server stopped")
		}
	}()

	processor.WatchForShutdown(ctx, wg, processor.CloserFunc(func() error {
		p.server.GracefulStop()
		return nil
	}))
}
