package main

import (
	"e-klinik/cmd/rest"
	"e-klinik/config"
	"e-klinik/pkg"
	"e-klinik/pkg/logging"
	"encoding/gob"
	"log"
)

func main() {

	cfg := config.NewConfig()
	logger := logging.NewLogger(cfg)
	pg := NewPostgre(cfg)
	gob.Register([]interface{}{})          // If any slice of interface is used
	gob.Register(map[string]interface{}{}) // If any m

	rmq, err := pkg.NewRabbit(cfg)

	failOnError(err, "rabbit failed")
	// defer rmq.Conn.Close()

	// err = rmq.ExchangeDeclare()
	// failOnError(err, "rabbit failed")
	// err  = rmq.QueueDeclare()
	// failOnError(err, "rabbit failed")
	go rest.HttpServer(cfg, rmq, pg)
	go rest.RabbitConsumer(rmq, cfg, pg, logger)
	select {}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
func NewPostgre(cfg *config.Config) *pkg.Postgres {
	pg, err := pkg.NewPgx(cfg)
	if err != nil {
		log.Fatal(err)
	}
	return pg
}
