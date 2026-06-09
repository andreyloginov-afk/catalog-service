package builder

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"

	"github.com/andreyloginov-afk/catalog-service/internal/app/config"
	rhandler "github.com/andreyloginov-afk/catalog-service/internal/app/handler/http"
	hcategory "github.com/andreyloginov-afk/catalog-service/internal/app/handler/http/category"
	rhealth "github.com/andreyloginov-afk/catalog-service/internal/app/handler/http/health"
	hproduct "github.com/andreyloginov-afk/catalog-service/internal/app/handler/http/product"
	"github.com/andreyloginov-afk/catalog-service/internal/app/processor"
	rprocessor "github.com/andreyloginov-afk/catalog-service/internal/app/processor/http"
	pprocessor "github.com/andreyloginov-afk/catalog-service/internal/app/processor/other"
	"github.com/andreyloginov-afk/catalog-service/internal/app/repository"
	pcategory "github.com/andreyloginov-afk/catalog-service/internal/app/repository/category"
	rcpostgres "github.com/andreyloginov-afk/catalog-service/internal/app/repository/conn/postgres"
	pproduct "github.com/andreyloginov-afk/catalog-service/internal/app/repository/product"
	"github.com/andreyloginov-afk/catalog-service/internal/app/service"
	scategory "github.com/andreyloginov-afk/catalog-service/internal/app/service/category"
	sproduct "github.com/andreyloginov-afk/catalog-service/internal/app/service/product"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

type Builder struct {
	cCtx *cli.Context
	ctx  context.Context
	wg   sync.WaitGroup
	err  error
	cfg  config.Config

	chErrors chan error

	connPostgres *rcpostgres.Client

	categoryRepo repository.Category
	productRepo  repository.Product

	categoryService service.Category
	productService  service.Product

	healthHandler   rhandler.Health
	categoryHandler rhandler.Category
	productHandler  rhandler.Product
	processors      []processor.Processor
}

func NewBuilder(cCtx *cli.Context) *Builder {
	b := Builder{

		cCtx:     cCtx,
		chErrors: make(chan error, 4096),
	}
	// отменяемый контекст
	ctx, cancelFunc := context.WithCancel(context.Background())
	b.ctx = ctx
	//канал сигналов
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	//Запускаем горутину
	go b.waitForSignal(sig, cancelFunc)

	//горутина вывода ошибок
	go b.printErrors()

	b.healthHandler = rhealth.NewHandler()
	return &b
}

func (b *Builder) waitForSignal(sig chan os.Signal, cancelFunc func()) {
	s := <-sig
	log.Info().Str("signal", s.String()).Msg("Shutdown is requested")
	cancelFunc()
}

func (b *Builder) printErrors() {
	for err := range b.chErrors {
		log.Error().Err(err).Msg("Got new error")
	}
}

func (b *Builder) BuildConfig() {
	b.exec(true, func(b *Builder) {
		b.buildConfig()
	})
}
func (b *Builder) Run() {
	if b.ctx.Err() != nil {
		log.Info().Msg("Shutdown during initialization")
		return
	}

	if b.err != nil {
		log.Fatal().Err(b.err).Msg("Failed to initialize application")
	}
	log.Info().Msg("Application initialized")
	defer log.Info().Msg("Application completed")

	for _, proc := range b.processors {
		proc.StartAsync(b.ctx, &b.wg)
	}
	b.wg.Wait()
}

////////////////////////////////////////////////////////////////////////////////
///// REPOSITORY CONNECTIONS ///////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (b *Builder) BuildRepoConnPostgres() {
	b.exec(true, func(b *Builder) {
		client, err := rcpostgres.NewConn(b.ctx, b.cfg.Repository.Postgres)
		if err != nil {
			b.err = err
			return
		}
		b.connPostgres = client
	})

}

func (b *Builder) BuildRepoConnMigrator() {
	b.exec(b.connPostgres != nil, func(b *Builder) {
		migrator := pprocessor.NewMigrator(b.connPostgres)
		b.processors = append(b.processors, migrator)
	}, b.connPostgres)
}

////////////////////////////////////////////////////////////////////////////////
///// REPOSITORIES /////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (b *Builder) BuildRepoCategory() {
	b.exec(true, func(b *Builder) {
		b.categoryRepo = pcategory.NewRepoFromPostgres(b.connPostgres)
	}, b.connPostgres)

}

func (b *Builder) BuildRepoProduct() {
	b.exec(true, func(b *Builder) {
		b.productRepo = pproduct.NewRepoFromPostgres(b.connPostgres)
	}, b.connPostgres)
}

////////////////////////////////////////////////////////////////////////////////
///// SERVICES /////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (b *Builder) BuildServiceCategory() {
	b.exec(true, func(b *Builder) {
		b.categoryService = scategory.NewService(b.categoryRepo, b.productRepo)

	}, b.categoryRepo, b.productRepo)
}

func (b *Builder) BuildServiceProduct() {
	b.exec(true, func(b *Builder) {
		b.productService = sproduct.NewService(b.productRepo, b.categoryRepo)
	}, b.productRepo, b.categoryRepo)
}

////////////////////////////////////////////////////////////////////////////////
///// HANDLERS /////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (b *Builder) BuildHandlerHttpCategory() {
	b.exec(true, func(b *Builder) {
		b.categoryHandler = hcategory.NewHandler(b.categoryService)

	}, b.categoryService)
}

func (b *Builder) BuildHandlerHttpProduct() {
	b.exec(true, func(b *Builder) {
		b.productHandler = hproduct.NewHandler(b.productService)
	}, b.productService)
}

////////////////////////////////////////////////////////////////////////////////
///// PROCESSORS ///////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (b *Builder) BuildProcHttp() {
	b.exec(true, func(b *Builder) {
		proc := rprocessor.NewHttp(b.healthHandler,
			b.categoryHandler,
			b.productHandler,
			b.cfg.Processor.WebServer,
		)
		b.processors = append(b.processors, proc)
	}, b.healthHandler)
}

////////////////////////////////////////////////////////////////////////////////
///// PRIVATE //////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (b *Builder) exec(preCond bool, cb func(b *Builder), requiredArgs ...any) {
	if !preCond || b.err != nil || b.ctx.Err() != nil {
		return
	}

	for i, requiredArg := range requiredArgs {
		rv := reflect.ValueOf(requiredArg)
		if !rv.IsValid() {
			b.err = fmt.Errorf("BUG: required argument #%d is nil (check dependencies)", i)
			return
		}
		if rv.Type().Kind() == reflect.Struct || !rv.IsZero() {
			continue
		}
		b.err = fmt.Errorf("BUG: required %s, but empty", rv.Type().String())
		return
	}

	cb(b)
}

func (b *Builder) buildConfig() {
	args := config.LoadArgs{
		Output:          b.cCtx.App.Writer,
		EnableSimpleLog: b.cCtx.Bool("no-json"),
	}
	config.Load(args)
	b.cfg = config.Root

}
