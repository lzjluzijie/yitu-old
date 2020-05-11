package route

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lzjluzijie/yitu/node/onedrive"
)

var n *onedrive.Node

func SetN(node *onedrive.Node) {
	n = node
}

func NewEngine() (r *gin.Engine) {
	r = gin.Default()

	// CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowHeaders = []string{"*"}
	r.Use(cors.New(corsConfig))

	r.GET("/t/:id", GetTu)
	r.GET("/t/:id/*type", GetTu)

	api := r.Group("/api")
	api.POST("/upload", UploadTu)
	//api.GET("/history", GetUploadHistory)

	return
}
