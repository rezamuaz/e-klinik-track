package rest

import (
	"context"
	"e-klinik/config"
	"e-klinik/infra/types"
	"e-klinik/infra/worker"
	"e-klinik/internal/di"
	"e-klinik/pkg"
	"e-klinik/pkg/constant"
	"e-klinik/pkg/logging"
	"e-klinik/utils"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"syscall"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	pgxadapter "github.com/gtoxlili/pgx-adapter"
	"github.com/typesense/typesense-go/v3/typesense/api"
)

const (
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultAddr            = ":80"
	_defaultShutdownTimeout = 3 * time.Second
)

func HttpServer(cfg *config.Config, rmq *pkg.RabbitMQ, pg *pkg.Postgres) {
	// Publisher channel
	//  _ = rmq.SetupExchange(pubCh)

	pubCh, err := rmq.NewChannel()
	if err != nil {
		log.Fatalf("Failed to create channel: %v", err)
		defer pubCh.Close()
	}
	if err := rmq.SetupExchange(pubCh); err != nil {
		log.Fatalf("Failed to setup exhange: %v", err)
	}

	adapter, err := pgxadapter.NewAdapter(context.Background(), pg.Pool)
	if err != nil {
		log.Fatalf("Failed to create adapter: %v", err)
	}

	m, err := model.NewModelFromString(`
	[request_definition]
	r = sub, obj, act

	[policy_definition]
	p = sub, obj, act

	[policy_effect]
	e = some(where (p.eft == allow))

	[matchers]
	m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
	`)

	casbin, err := casbin.NewEnforcer(m, adapter)

	cacheService := pkg.NewRistrettoCache()
	defer cacheService.Close()
	//Dependency Injection
	init := di.Injector(cfg, pubCh, pg, cacheService, casbin)
	server := &http.Server{
		Addr:         _defaultAddr,
		Handler:      init.Router,
		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
	}
	server.Addr = net.JoinHostPort("", cfg.Server.ExternalPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		log.Println("receive interrupt signal")
		if err := server.Close(); err != nil {
			log.Fatal("Http Server Close:", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("Http Server closed under request")
		} else {
			log.Fatal("Http Server closed unexpect")
		}
	}

	log.Println("Http Server exiting")
}

func RabbitConsumer(rmq *pkg.RabbitMQ, cfg *config.Config, pg *pkg.Postgres, logger logging.Logger) {
	ctx := context.Background()
	ts := pkg.NewTypeSense(cfg)

	schema := &api.CollectionSchema{
		Name: "pvsave",
		Fields: []api.Field{
			{Name: "id", Type: "string"},
			{Name: "title", Type: "string"},
			{Name: "alternative_title", Type: "string[]", Facet: utils.BoolPtr(true)},
			{Name: "artists", Type: "string[]", Facet: utils.BoolPtr(true)},      // artists' names as a string array
			{Name: "artist_ids", Type: "string[]", Facet: utils.BoolPtr(true)},   // artist IDs as a string array
			{Name: "artist_kanji", Type: "string[]", Facet: utils.BoolPtr(true)}, // artist Kanji as a string array
			{Name: "genres", Type: "string[]", Facet: utils.BoolPtr(true)},
			{Name: "genre_ids", Type: "string[]", Facet: utils.BoolPtr(true)}, // genre IDs as a string array
			{Name: "thumbnails", Type: "string"},
			{Name: "category", Type: "string[]", Facet: utils.BoolPtr(true)},
			{Name: "tags", Type: "string[]", Facet: utils.BoolPtr(true)},
			{Name: "status", Type: "string"},
			{Name: "views", Type: "int64", Facet: utils.BoolPtr(true), Sort: utils.BoolPtr(true)},
			{Name: "published_at", Type: "int64", Sort: utils.BoolPtr(true)},
			{Name: "published_by", Type: "string"},
			{Name: "created_at", Type: "int64"},
			{Name: "created_by", Type: "string"},
			{Name: "updated_at", Type: "int64"},
			{Name: "updated_by", Type: "string"},
		},
		DefaultSortingField: utils.StringPtr("published_at"),
	}
	err := ts.EnsureCollectionExists(ctx, schema)
	tsRepo := types.NewIndexRepository(ts)

	// postRepo := pgr.NewPostRepository(pg)
	// Consumer channel
	conCh, err := rmq.NewChannel()
	if err != nil {
		log.Fatalf("Failed to create channel: %v", err)
	}
	defer conCh.Close()
	if err := rmq.SetupQueue(conCh); err != nil {
		log.Fatalf("Failed to setup queue: %v", err)
	}
	// Setup context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	server := &worker.ConsumerService{
		Logger: logger, // assume your custom logger
		RMQ:    rmq,
		// PostRepo: postRepo,
		TsRepo: tsRepo, // your repo impl
		Ch:     conCh,
		Done:   make(chan struct{}),
	}

	if err := server.StartRabbitConsumer(ctx); err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigs:
		log.Println("Shutdown signal received.")
		cancel()
	case <-server.Done:
		log.Println("Consumer finished processing.")
	}

	// Wait for cleanup
	select {
	case <-server.Done:
	case <-time.After(5 * time.Second):
		log.Println("Timeout waiting for consumer to finish.")
	}
	_ = server.Ch.Cancel(constant.RMQConsumerName, false)
	_ = server.Ch.Close()
	_ = server.RMQ.Conn.Close()

	log.Println("Server shutdown gracefully.")
}
