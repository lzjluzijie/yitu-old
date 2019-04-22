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

	if len(id) == 64 {
		tu, err = models.GetTuByHash(id)
	} else {
		tid, er := strconv.ParseUint(id, 10, 64)
		if er != nil {
			c.String(http.StatusBadRequest, er.Error())
			return
		}
		tu, err = models.GetTuByID(tid)
	}

	if err != nil {
		if err.Error() == "not found" {
			c.String(http.StatusNotFound, err.Error())
			return
		}

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

	if tu.OneDriveURL == "" {
		c.String(http.StatusNotFound, "not found")
		return
	}

	c.Redirect(http.StatusMovedPermanently, tu.OneDriveURL)
	return
}
