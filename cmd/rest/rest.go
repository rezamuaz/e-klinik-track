package rest

import (
	"context"
	"e-klinik/config"
	"e-klinik/infra/pg"
	"e-klinik/infra/types"
	"e-klinik/infra/worker"
	"e-klinik/internal/di"
	"e-klinik/pkg"
	"e-klinik/pkg/constant"
	"e-klinik/pkg/logging"
	"e-klinik/utils"
	"fmt"
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
	r = sub, obj, act, a, s, t

	[policy_definition]
	p = sub, obj, act, a, s, t

	[role_definition]
	g = _, _  
	g2 = _, _ 

	[policy_effect]
	e = some(where (p.eft == allow))

	[matchers]
	m = g(r.sub, "1") || (g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act)
	`)

	casbin, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		log.Fatalf("Failed to create adapter: %v", err)
	}
	err = casbin.LoadPolicy()
	if err != nil {
		log.Fatalf("Failed to load policy: %v", err)
	}
	casbin.EnableAutoSave(true)
	casbin.EnableLog(true)

	//Initialize redis
	rdb := pkg.NewRedisCache(cfg)

	err = LoadResourceMappingsIntoRedis(rdb, pg)
	if err != nil {
		log.Print("Failed to load data:", err)
	}
	//Dependency Injection
	init := di.Injector(cfg, pubCh, pg, rdb, casbin)
	server := &http.Server{
		Addr:         _defaultAddr,
		Handler:      init.Router,
		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
	}
	server.Addr = net.JoinHostPort("", cfg.Server.ExternalPort)

	// ✅ MENAMPILKAN PORT SAAT INI DI LOG
	log.Printf("🚀 Starting HTTP server on %s...", server.Addr)

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

func LoadResourceMappingsIntoRedis(rdb *pkg.RedisCache, postgre *pkg.Postgres) error {
	log.Println("[Cache] Memuat pemetaan resource dari database (SQLC) ke Redis...")

	ctx := context.Background()
	queries := pg.New(postgre.Pool)

	// 1. Ambil semua data pemetaan dari DB
	mappings, err := queries.ListResourceMappings(ctx)
	if err != nil {
		return fmt.Errorf("gagal mengambil pemetaan dari DB (SQLC): %w", err)
	}

	// 2. Masukkan ke Redis dalam satu transaksi (Pipeline)
	pipe := rdb.Client.Pipeline()
	var validMappingsCount int // Counter untuk baris yang berhasil diproses

	for _, m := range mappings {
		// Cek Keamanan: Pastikan Path dan Method TIDAK nil.
		// Hanya resource API yang memiliki Path dan Method yang harus dicache.
		if m.Path == nil || m.Method == nil {
			// Log baris data yang diabaikan (mungkin itu adalah resource VIEW/Menu)
			log.Printf("[Cache] Mengabaikan resource key: %s (Path/Method nil)", m.ResourceKey)
			continue
		}

		// DEREFERENCE POINTER dengan aman (gunakan * di depan)
		// Kita yakin pointer tidak nil karena telah diperiksa di atas.
		path := *m.Path
		method := *m.Method

		// Format Key Redis: /api/articles:POST
		redisKey := fmt.Sprintf("%s:%s", path, method)

		// Format Value Redis: resourceKey:action (e.g., data:article:create)
		redisValue := fmt.Sprintf("%s:%s", m.ResourceKey, m.Action)

		pipe.Set(ctx, redisKey, redisValue, 0)
		validMappingsCount++
	}

	// 3. Eksekusi Pipeline
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("gagal mengeksekusi Redis pipeline: %w", err)
	}

	// Gunakan validMappingsCount, bukan len(mappings)
	log.Printf("[Cache] ✅ Berhasil memuat %d pemetaan resource API ke Redis. (%d total baris dari DB)", validMappingsCount, len(mappings))
	return nil
}
