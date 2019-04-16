package routers

import (
	"github.com/gin-gonic/gin"
)

func RegisterRouters(router *gin.Engine) {
	router.StaticFile("/", "./frontend/dist/index.html")
	router.StaticFile("/index.html", "./frontend/dist/index.html")
	router.Static("/js", "./frontend/dist/js")
	router.Static("/css", "./frontend/dist/css")

	router.GET("/t/:id/*type", GetTu)
	//TODO
	router.GET("/t/:id", GetTu)

	api := router.Group("/api")
	api.POST("/upload", Upload)
	api.GET("/delete/:dc", Delete)

	return
}
