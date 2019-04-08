package routers

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

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
	if tu != nil {
		resp := &UploadResponse{
			Name: tu.Name,
			Size: tu.Size,
			Hash: tu.Hash,
			URL:  fmt.Sprintf("https://t.halu.lu/t/%d", tu.ID),
		}
		c.JSON(200, resp)
		return
	}

	//upload
	id, parent, err := onedrive.Upload(fmt.Sprintf(`/yitu/%s/%s/%s`, time.Now().Format("20160102"), hash, name), data)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

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
	tu = &models.Tu{
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
		Name: name,
		Size: size,
		Hash: hash,
		URL:  fmt.Sprintf("https://t.halu.lu/t/%d", tu.ID),
	}
	c.JSON(200, resp)
	return
}
