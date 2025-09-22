//go:build wireinject
// +build wireinject

package di

import (
	"e-klinik/api/handler"
	"e-klinik/config"
	"e-klinik/infra/worker"

	"e-klinik/internal/api"
	"e-klinik/internal/usecase"
	"e-klinik/pkg"

	"github.com/casbin/casbin/v2"
	"github.com/google/wire"
	"github.com/streadway/amqp"
)

// var repositorySet = wire.NewSet(
//
//	pgr.NewArtistRepository,
//	wire.Bind(new(pgr.ArtistRepository), new(*pgr.ArtistRepositoryImpl)),
//	pgr.NewGenreRepository,
//	wire.Bind(new(pgr.GenreRepository), new(*pgr.GenreRepositoryImpl)),
//	pgr.NewPostRepository,
//	wire.Bind(new(pgr.PostRepository), new(*pgr.PostRepositoryImpl)),
//	pgr.NewUowRepository,
//	wire.Bind(new(pgr.UowRepository), new(*pgr.UowRepositoryImpl)),
//
// )

var usecaseSet = wire.NewSet(
	usecase.NewUserUsecase,
	wire.Bind(new(usecase.UserUsecase), new(*usecase.UserUsecaseImpl)),
	usecase.NewMainUsecase,
	wire.Bind(new(usecase.MainUsecase), new(*usecase.MainUsecaseImpl)),
)

var handlerSet = wire.NewSet(
	handler.NewAuthHandler,
	wire.Bind(new(handler.AuthHandler), new(*handler.AuthHandlerImpl)),
	handler.NewMainHandler,
	wire.Bind(new(handler.MainHandler), new(*handler.MainHandlerImpl)),
)

// InitServer is the injector entry po int.
func Injector(cfg *config.Config, ch *amqp.Channel, pg *pkg.Postgres, cache *pkg.RistrettoCache, casbin *casbin.Enforcer) *pkg.Server {
	wire.Build(
		// repositorySet,
		usecaseSet,
		handlerSet,
		worker.NewQueueService,
		api.NewApiRouter,
		wire.Struct(new(api.Initialized), "*"),
	)
	return nil

}
