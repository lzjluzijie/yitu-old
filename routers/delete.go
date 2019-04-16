package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lzjluzijie/yitu/models"
)

func Delete(c *gin.Context) {
	dc := c.Param("dc")
	if dc == "" {
		c.String(http.StatusBadRequest, "delete code is empty")
		return
	}

	err := models.DeleteByCode(dc)
	if err != nil {
		if err.Error() == "not found" {
			c.String(http.StatusNotFound, err.Error())
			return
		}

		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, "ok")
	return
}
