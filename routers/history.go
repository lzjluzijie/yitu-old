package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lzjluzijie/yitu/models"
)

func GetUploadHistory(c *gin.Context) {
	ip := c.ClientIP()
	t, err := models.GetUploadHistory(ip)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	data := make([]UploadResponse, len(t))
	for i, tu := range t {
		data[i] = GetUploadResponse(tu)
	}

	c.JSON(http.StatusOK, data)
	return
}
