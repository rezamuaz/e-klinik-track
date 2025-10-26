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
	usecase.NewFasilitasUseCase,
	wire.Bind(new(usecase.FasilitasUsecase), new(*usecase.FasilitasUsecaseImpl)),
	usecase.NewKontrakUsecase,
	wire.Bind(new(usecase.KontrakUsecase), new(*usecase.KontrakUsecaseImpl)),
	usecase.NewRuanganUsecase,
	wire.Bind(new(usecase.RuanganUsecase), new(*usecase.RuanganUsecaseImpl)),
	usecase.NewMataKuliahUsecase,
	wire.Bind(new(usecase.MataKuliahUsecase), new(*usecase.MataKuliahUsecaseImpl)),
	usecase.NewKehadiranUsecase,
	wire.Bind(new(usecase.KehadiranUsecase), new(*usecase.KehadiranUsecaseImpl)),
	usecase.NewSkpKehadiranUsecase,
	wire.Bind(new(usecase.SkpKehadiranUsecase), new(*usecase.SkpKehadiranUsecaseImpl)),
	usecase.NewSkpUsecase,
	wire.Bind(new(usecase.SkpUsecase), new(*usecase.SkpUsecaseImpl)),
	usecase.NewActorUsecase,
	wire.Bind(new(usecase.ActorUsecase), new(*usecase.ActorUsecaseImpl)),
	usecase.NewSummaryUsecase,
	wire.Bind(new(usecase.SummaryUsecase), new(*usecase.SummaryUsecaseImpl)),
)

var handlerSet = wire.NewSet(
	handler.NewAuthHandler,
	wire.Bind(new(handler.AuthHandler), new(*handler.AuthHandlerImpl)),
	handler.NewUserHandler,
	wire.Bind(new(handler.UserHandler), new(*handler.UserHandlerImpl)),
	handler.NewFasilitasHandler,
	wire.Bind(new(handler.FasilitasHandler), new(*handler.FasilitasHandlerImpl)),
	handler.NewKontrakHandler,
	wire.Bind(new(handler.KontrakHandler), new(*handler.KontrakHandlerImpl)),
	handler.NewRuanganHandler,
	wire.Bind(new(handler.RuanganHandler), new(*handler.RuanganHandlerImpl)),
	handler.NewMataKuliahHandler,
	wire.Bind(new(handler.MataKuliahHandler), new(*handler.MataKuliahHandlerImpl)),
	handler.NewKehadiranHandler,
	wire.Bind(new(handler.KehadiranHandler), new(*handler.KehadiranHandlerImpl)),
	handler.NewSkpKehadiranHandler,
	wire.Bind(new(handler.SkpKehadiranHandler), new(*handler.SkpKehadiranHandlerImpl)),
	handler.NewSkpHandler,
	wire.Bind(new(handler.SkpHandler), new(*handler.SkpHandlerImpl)),
	handler.NewActorHandler,
	wire.Bind(new(handler.ActorHandler), new(*handler.ActorHandlerImpl)),
	handler.NewSummaryHandler,
	wire.Bind(new(handler.SummaryHandler), new(*handler.SummaryHandlerImpl)),
	handler.NewPermissionHandler,
	wire.Bind(new(handler.PermissionHandler), new(*handler.PermissionHandlerImpl)),
)

// InitServer is the injector entry po int.
func Injector(cfg *config.Config, ch *amqp.Channel, pg *pkg.Postgres, cache *pkg.RedisCache, casbin *casbin.Enforcer) *pkg.Server {
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
