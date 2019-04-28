package routers

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func RegisterRouters(router *gin.Engine) {
	router.Use(static.Serve("/", static.LocalFile("frontend/dist", false)))

	router.GET("/t/:id", GetTu)
	router.GET("/t/:id/*type", GetTu)

	api := router.Group("/api")
	api.POST("/upload", Upload)
	api.GET("/delete/:dc", Delete)
	api.GET("/history", GetUploadHistory)

	return
}
