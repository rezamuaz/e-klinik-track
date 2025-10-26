package router

import (
	"e-klinik/api/handler"

	"github.com/gin-gonic/gin"
)

func KehadiranSkp(group *gin.RouterGroup, h *handler.SkpKehadiranHandlerImpl) {

	//SKP Kehadiran
	group.POST("", h.CreateSyncKehadiranSkp)
	group.GET("", h.ListKehadiranSkp)
	group.GET("/active", h.SkpByKehadiranId)
	group.GET("/active-nama", h.IntervensiKehadiranId)
	group.PUT("", h.UpdateKehadiranSkp)
	group.DELETE("", h.DeleteKehadiranSkp)
	group.POST("/approve", h.ApproveKehadiranSkp)

}
