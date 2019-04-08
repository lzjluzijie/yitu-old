package routers

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lzjluzijie/yitu/models"

	"github.com/lzjluzijie/yitu/onedrive"
)

type UploadResponse struct {
	Size int64  `json:"size"`
	URL  string `json:"url"`
}

func Upload(c *gin.Context) {
	f, err := c.FormFile("tu")

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	name := f.Filename
	size := f.Size

	if size >= 50*1024*1024 {
		c.String(http.StatusBadRequest, "file too big")
		return
	}

	file, err := f.Open()

	//hash
	h := sha256.New()
	r := io.TeeReader(file, h)

	//upload
	id, parent, err := onedrive.Upload(name, size, r)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	hash := fmt.Sprintf("%x", h.Sum(nil))

	//rename parent folder
	err = onedrive.Rename(parent, hash)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	//share
	url, err := onedrive.Share(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	url += "?download=1"

	//insert to database
	tu := &models.Tu{
		Name:             name,
		Size:             size,
		Hash:             hash,
		OneDriveFolderID: parent,
		OneDriveID:       id,
		OneDriveURL:      url,
	}
	err = models.InsertTu(tu)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	//finish
	resp := &UploadResponse{
		Size: size,
		URL:  fmt.Sprintf("https://t.halu.lu/t/%d/%s", tu.ID, name),
	}
	c.JSON(200, resp)
}
