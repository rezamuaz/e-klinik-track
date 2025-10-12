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
	AuthHandler *handler.AuthHandlerImpl
	MainHandler *handler.MainHandlerImpl
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
	{

		auth := web.Group("/auth")
		router.Auth(auth, h.AuthHandler)
		main := web.Group("/main")
		main.Use(middleware.JwtAuth(cfg.JWT.AccessTokenSecret))
		// main.Use(middleware.RbacAuthzMiddleware(cb, rdb))
		router.Main(main, h.MainHandler)
		rbac := web.Group("/rbac")
		rbac.Use(middleware.JwtAuth(cfg.JWT.AccessTokenSecret))
		router.Role(rbac, h.AuthHandler)

	}

	return &pkg.Server{Router: r}

}
