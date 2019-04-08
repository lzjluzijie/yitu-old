package routers

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/h2non/bimg.v1"

	"github.com/gin-gonic/gin"

	"github.com/lzjluzijie/yitu/models"

	"github.com/lzjluzijie/yitu/onedrive"
)

type UploadResponse struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Hash string `json:"hash"`
	URL  string `json:"url"`
}

func Upload(c *gin.Context) {
	//original file
	f, err := c.FormFile("tu")
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	name := f.Filename

	if f.Size >= 50*1024*1024 {
		c.String(http.StatusBadRequest, "file too big")
		return
	}

	file, err := f.Open()
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	//hash
	h := sha256.New()
	r := io.TeeReader(file, h)
	data, err := ioutil.ReadAll(r)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	hash := fmt.Sprintf("%x", h.Sum(nil))

	//check
	size := int64(len(data))
	if f.Size != size {
		c.String(http.StatusBadRequest, fmt.Sprintf("file size does not match: %d, %d", f.Size, size))
		return
	}

	tu, err := models.GetTuByHash(hash)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if tu.ID != 0 {
		resp := &UploadResponse{
			Name: tu.Name,
			Size: tu.Size,
			Hash: tu.Hash,
			URL:  fmt.Sprintf("https://t.halu.lu/t/%d", tu.ID),
		}
		c.JSON(200, resp)
		return
	}

	//bimg
	image := bimg.NewImage(data)
	is, err := image.Size()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	//upload
	path := fmt.Sprintf(`/yitu/%s/%s/`, time.Now().Format("20060102"), hash)
	id, parent, err := onedrive.Upload(path+name, data)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	//webp
	webp, err := image.Process(bimg.Options{
		Type:    bimg.WEBP,
		Quality: 95,
	})
	webpID, _, err := onedrive.Upload(path+strings.TrimSuffix(name, filepath.Ext(name))+".webp", webp)
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

	webpURL, err := onedrive.Share(webpID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	webpURL += "?download=1"

	//insert to database
	tu = &models.Tu{
		Name:             name,
		Size:             size,
		Hash:             hash,
		Width:            is.Width,
		Height:           is.Height,
		OneDriveFolderID: parent,
		OneDriveID:       id,
		OneDriveURL:      url,
		OneDriveWebpID:   webpID,
		OneDriveWebpURL:  webpURL,
	}
	err = models.InsertTu(tu)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	//finish
	resp := &UploadResponse{
		Name: name,
		Size: size,
		Hash: hash,
		URL:  fmt.Sprintf("https://t.halu.lu/t/%d", tu.ID),
	}
	c.JSON(200, resp)
	return
}
