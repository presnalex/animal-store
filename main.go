package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	pb "github.com/presnalex/animal-store-proto/go/animal_store_proto"
	pbswagger "github.com/presnalex/animal-store-proto/swagger"
	hdl "github.com/presnalex/animal-store/handler"
	"github.com/presnalex/animal-store/handler/encoders"
	strg "github.com/presnalex/animal-store/storage"
	"github.com/presnalex/go-micro/v3/database/postgres"
	dbwrapper "github.com/presnalex/go-micro/v3/database/wrapper"
	"github.com/presnalex/go-micro/v3/rest"
	"github.com/presnalex/statscheck"
	swagger "github.com/presnalex/swagger-ui-route"
	jsoncodec "go.unistack.org/micro-codec-json/v3"
	consulconfig "go.unistack.org/micro-config-consul/v3"
	envconfig "go.unistack.org/micro-config-env/v3"
	httpServer "go.unistack.org/micro-server-http/v3"
	micro "go.unistack.org/micro/v3"
	"go.unistack.org/micro/v3/config"
	"go.unistack.org/micro/v3/logger"
	"go.unistack.org/micro/v3/server"
)

const appName = "animal-store"

var (
	BuildDate  string
	AppVersion string
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-ch
		logger.Infof(ctx, "handle signal %v, exiting", sig)
		cancel()
	}()

	cfg := newConfig(appName, AppVersion)
	consulcfg := consulconfig.NewConfig(
		config.Struct(cfg),
		config.Codec(jsoncodec.NewCodec()),
		config.BeforeLoad(func(ctx context.Context, c config.Config) error {
			if len(cfg.Consul.NamespacePath) == 0 {
				cfg.Consul.NamespacePath = "go-micro-layouts"
			}
			if len(cfg.Consul.AppPath) == 0 {
				cfg.Consul.AppPath = "animal-store"
			}
			logger.Infof(ctx, "Consul Address: %s", cfg.Consul.Addr)
			logger.Infof(ctx, "Consul Path: %s", filepath.Join(cfg.Consul.NamespacePath, cfg.Consul.AppPath))
			return c.Init(
				consulconfig.Address(cfg.Consul.Addr),
				consulconfig.Token(cfg.Consul.Token),
				consulconfig.Path(filepath.Join(cfg.Consul.NamespacePath, cfg.Consul.AppPath)),
			)
		}),
	)
	if err := config.DefaultBeforeLoad(ctx, consulcfg); err != nil {
		logger.Fatalf(ctx, "failed to load config: %v", err)
	}

	err := config.Load(ctx,
		[]config.Config{
			config.NewConfig(
				config.Struct(cfg),
			),
			envconfig.NewConfig(
				config.Struct(cfg),
			),
			consulcfg},
	)
	if err != nil {
		logger.Fatalf(ctx, "failed to load config: %v", err)
	}

	srv := httpServer.NewServer(
		server.Address(cfg.Server.Addr),
		server.Name(server.DefaultServer.Options().Name),
		server.Version(server.DefaultServer.Options().Version),
		server.ID(server.DefaultServer.Options().ID),
		server.Codec("application/json", jsoncodec.NewCodec()),
	)

	svc := micro.NewService(
		micro.Context(ctx),
		micro.Server(srv),
	)

	if err = svc.Init(); err != nil {
		logger.Fatal(ctx, err)
	}

	dbPrimary, err := postgres.Connect(cfg.PostgresPrimary)
	if err != nil {
		logger.Fatal(ctx, "cannot connect to postgres primary: ", err)
		return
	}

	storage := strg.NewStorage(dbwrapper.NewWrapper(
		dbPrimary,
		dbwrapper.DBHost(cfg.PostgresPrimary.Addr),
		dbwrapper.DBName(cfg.PostgresPrimary.DBName),
		dbwrapper.ServiceName(svc.Server().Options().Name),
		dbwrapper.ServiceVersion(svc.Server().Options().Version),
		dbwrapper.ServiceID(svc.Server().Options().ID)))

	writer, err := hdl.NewWriter(encoders.NewJSONProto())
	if err != nil {
		logger.Fatalf(ctx, "writer init failed: %v", err)
		return
	}

	handler, err := hdl.NewHandler(writer, storage)
	if err != nil {
		logger.Fatalf(ctx, "handler init failed: %v", err)
		return
	}

	statsOpts := append([]statscheck.Option{},
		statscheck.WithDefaultHealth(),
		statscheck.WithMetrics(),
		statscheck.WithVersionDate(AppVersion, BuildDate),
	)

	if cfg.Core.Profile {
		statsOpts = append(statsOpts, statscheck.WithProfile())
	}

	healthServer := statscheck.NewServer(statsOpts...)
	mux := healthServer.Mux()

	if err = rest.Register(mux, handler, pb.NewAnimalStoreServiceEndpoints()); err != nil {
		logger.Fatal(ctx, "rest.Register failed: %v", err)
	}

	swagger.Register(mux, appName, pbswagger.Asset, pbswagger.AssetDir, pbswagger.AssetInfo)
	if err = srv.Handle(srv.NewHandler(mux)); err != nil {
		logger.Fatal(ctx, err)
	}

	if err = svc.Run(); err != nil {
		logger.Fatal(ctx, "svc run: ", err)
	}
}
