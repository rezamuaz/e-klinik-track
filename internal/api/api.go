package api

import (
	"e-klinik/api/handler"
	"e-klinik/api/middleware"
	"e-klinik/api/router"
	"e-klinik/config"
	"e-klinik/pkg"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

type Initialized struct {
	ActorHandler        *handler.ActorHandlerImpl
	AuthHandler         *handler.AuthHandlerImpl
	FasilitasHandler    *handler.FasilitasHandlerImpl
	KehadiranHandler    *handler.KehadiranHandlerImpl
	KontrakHandler      *handler.KontrakHandlerImpl
	MataKuliahHandler   *handler.MataKuliahHandlerImpl
	RuanganHandler      *handler.RuanganHandlerImpl
	SkpHandler          *handler.SkpHandlerImpl
	SkpKehadiranHandler *handler.SkpKehadiranHandlerImpl
	SummaryHandler      *handler.SummaryHandlerImpl
	UserHandler         *handler.UserHandlerImpl
	PermissionHandler   *handler.PermissionHandlerImpl
}

func NewApiRouter(cfg *config.Config, h *Initialized, cb *casbin.Enforcer, rdb *pkg.RedisCache) *pkg.Server {

	// arangoC := pkg.NewArangoDatabase(cfg)
	gin.SetMode(cfg.Server.RunMode)
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!!!")
	})

	// Gin Route Initialized
	api := r.Group("/api")
	v1 := api.Group("/v1")
	web := v1.Group("/web")
	web.Use(middleware.ErrorHandler())
	{

		auth := web.Group("/auth")
		router.Auth(auth, h.AuthHandler)
		main := web.Group("/main")
		main.Use(middleware.JwtAuth(cfg.JWT.AccessTokenSecret))
		fasilitas := main.Group("/fasilitas")
		router.Fasilitas(fasilitas, h.FasilitasHandler)
		kontrak := main.Group("/kontrak")
		router.Kontrak(kontrak, h.KontrakHandler)
		ruangan := main.Group("/ruangan")
		router.Ruangan(ruangan, h.RuanganHandler)
		mataKuliah := main.Group("/mata-kuliah")
		router.MataKuliah(mataKuliah, h.MataKuliahHandler)
		kehadiran := main.Group("/kehadiran")
		router.Kehadiran(kehadiran, h.KehadiranHandler)
		kehadiranSkp := main.Group("/kehadiran-skp")
		router.KehadiranSkp(kehadiranSkp, h.SkpKehadiranHandler)
		skp := main.Group("/skp")
		router.Skp(skp, h.SkpHandler)
		actor := main.Group("/actor")
		router.Actor(actor, h.ActorHandler)
		user := main.Group("/users")
		router.User(user, h.UserHandler)
		summary := main.Group("/summary")
		router.Summary(summary, h.SummaryHandler)
		permission := main.Group("/permissions")
		router.Permission(permission, h.PermissionHandler)

	}

	return &pkg.Server{Router: r}

}
