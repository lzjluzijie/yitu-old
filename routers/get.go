package routers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lzjluzijie/yitu/models"
)

func GetTu(c *gin.Context) {
	id := c.Param("id")
	t := c.Param("type")

	var tu *models.Tu
	var err error
	var has bool

	if len(id) == 32 {
		has, tu, err = models.GetTu(&models.Tu{MD5: id})
	} else if len(id) == 64 {
		has, tu, err = models.GetTu(&models.Tu{SHA256: id})
	} else {
		tid, er := strconv.ParseInt(id, 10, 64)
		if er != nil {
			c.String(http.StatusBadRequest, er.Error())
			return
		}
		has, tu, err = models.GetTu(&models.Tu{ID: tid})
	}

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if !has {
		c.String(http.StatusNotFound, "not found")
		return
	}

	if tu.OneDriveURL == "" {
		c.String(http.StatusNotFound, "not found")
		return
	}

	c.Header("Cache-Control", "public, max-age=3110400")

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
