package middleware

import (
	"e-klinik/pkg"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RbacAuthzMiddleware(e *casbin.Enforcer, rdb *pkg.RedisCache) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		path := c.Request.URL.Path
		method := c.Request.Method

		// 1. Buat Key Redis dari Path dan Method
		redisKey := fmt.Sprintf("%s:%s", path, method)

		// 2. Akses Redis untuk mendapatkan Pemetaan Resource Key dan Action
		val, err := rdb.GetRaw(ctx, redisKey)

		if err == redis.Nil {
			// Resource tidak terdaftar di cache. Ini mungkin rute publik atau error konfigurasi.
			log.Printf("Resource not mapped in Redis: %s", redisKey)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Authorization resource not configured."})
			return
		} else if err != nil {
			// Error Redis (Timeout, Koneksi)
			log.Printf("Redis error: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Authorization service unavailable."})
			return
		}

		// Asumsi Value Redis: "resourceKey:action" (e.g., "data:article:create")
		parts := strings.Split(val, ":")
		if len(parts) != 2 {
			log.Printf("Invalid Redis value format: %s", val)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid resource configuration."})
			return
		}

		resourceKey := parts[0]  // e.g., "data:article"
		action := parts[1]       // e.g., "create"
		userID, _ := c.Get("id") // Dapatkan Subjek

		// 3. Eksekusi Casbin Enforce
		allowed, casbinErr := e.Enforce(userID, resourceKey, action)
		if casbinErr != nil {
			log.Printf("Casbin error: %v", casbinErr)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Casbin check failed."})
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Akses ditolak."})
			return
		}

		// Lanjut ke handler
		c.Next()
	}
}
