package router

import (
	"e-klinik/api/handler"

	"github.com/gin-gonic/gin"
)

func Kehadiran(group *gin.RouterGroup, h *handler.KehadiranHandlerImpl) {

	//Kehadiran
	group.POST("", h.CreateKehadiran)
	group.GET("", h.ListKehadiran)
	group.PUT("/:id", h.UpdateKehadiran)
	group.DELETE("/:id", h.DeleteKehadiran)
	group.GET("/user/status", h.CheckKehadiran)
	group.GET("pembimbing-klinik", h.GetKehadiranByPembimbingStatus)
	group.GET("/mahasiswa", h.GetKehadiranByMahasiswaStatus)
	group.GET("/users", h.ListDistinctUserKehadiran)
}
