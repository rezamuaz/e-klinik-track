package pkg

import (
	"context"
	"e-klinik/config"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
)

const (
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

// Postgres -.
type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration
	Pool         *pgxpool.Pool
}

func NewPgx(cfg *config.Config) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:  cfg.Postgre.PG_PoolMax,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&TimeZone=Asia/Jakarta",
		cfg.Postgre.PG_Username,
		cfg.Postgre.PG_Password,
		cfg.Postgre.PG_Host,
		cfg.Postgre.PG_Port,
		cfg.Postgre.PG_Name,
	)

	poolConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
	}

	// Use the AfterConnect hook to register the type map.
	poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())
		return nil
	}

	poolConfig.MaxConns = int32(pg.maxPoolSize)

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
	}

	return pg, nil
}

// Close -.
func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
