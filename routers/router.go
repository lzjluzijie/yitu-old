package routers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/lzjluzijie/yitu/models"
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
}

func GetTu(c *gin.Context) {
	id := c.Param("id")
	t := c.Param("type")

	tid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	tu, err := models.GetTuByID(tid)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if (t == "/webp" || t == "webp") && tu.OneDriveWebPURL != "" {
		c.Redirect(http.StatusMovedPermanently, tu.OneDriveWebPURL)
		return
	}

	if (t == "/fhd" || t == "fhd") && tu.OneDriveFHDURL != "" {
		c.Redirect(http.StatusMovedPermanently, tu.OneDriveFHDURL)
		return
	}
	if (t == "/fhdwebp" || t == "fhdwebp") && tu.OneDriveFHDWebPURL != "" {
		c.Redirect(http.StatusMovedPermanently, tu.OneDriveFHDWebPURL)
		return
	}

	c.Redirect(http.StatusMovedPermanently, tu.OneDriveURL)
	return
}
