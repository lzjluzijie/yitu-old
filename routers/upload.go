package routers

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/h2non/bimg.v1"

	"github.com/gin-gonic/gin"

	"github.com/lzjluzijie/yitu/models"

	"github.com/lzjluzijie/yitu/onedrive"
)

const fhdWidth = 1920

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

	//bimg check image
	image := bimg.NewImage(data)
	is, err := image.Size()
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	width := is.Width
	height := is.Height

	//insert to database
	tu = &models.Tu{
		Name:   name,
		Size:   size,
		Hash:   hash,
		Width:  width,
		Height: height,
	}
	err = models.InsertTu(tu)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	resp := &UploadResponse{
		Name: name,
		Size: size,
		Hash: hash,
		URL:  fmt.Sprintf("https://t.halu.lu/t/%d", tu.ID),
	}
	c.JSON(200, resp)

	//async upload
	path := fmt.Sprintf(`/yitu/%s/%s/`, time.Now().Format("20060102"), hash)
	ext := filepath.Ext(name)
	noext := strings.TrimSuffix(name, ext)
	go func() {
		id, parent, url, err := onedrive.UploadAndShare(path+name, data)
		if err != nil {
			log.Println(err.Error())
			err = models.DeleteTu(tu)
			return
		}

		tu.OneDriveFolderID = parent
		tu.OneDriveID = id
		tu.OneDriveURL = url

		//update
		err = models.UpdateTu(tu)
		if err != nil {
			log.Println(err.Error())
			return
		}

	}()

	//WebP
	go func() {
		webp, err := image.Process(bimg.Options{
			Type:    bimg.WEBP,
			Quality: 95,
		})
		if err != nil {
			log.Println(err.Error())
			return
		}

		webpID, _, webpURL, err := onedrive.UploadAndShare(path+noext+".webp", webp)
		if err != nil {
			log.Println(err.Error())
			return
		}
		tu.OneDriveWebPID = webpID
		tu.OneDriveWebPURL = webpURL

		//update webp
		err = models.UpdateTu(tu)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}()

	if width > fhdWidth {
		fhdHeight := int(height * fhdWidth / width)

		//FHD
		go func() {
			fhd, err := image.Process(bimg.Options{
				Width:   fhdWidth,
				Height:  fhdHeight,
				Quality: 95,
			})
			if err != nil {
				log.Println(err.Error())
				return
			}

			fhdID, _, fhdURL, err := onedrive.UploadAndShare(path+noext+".fhd"+ext, fhd)
			if err != nil {
				log.Println(err.Error())
				return
			}
			tu.OneDriveFHDID = fhdID
			tu.OneDriveFHDURL = fhdURL

			err = models.UpdateTu(tu)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}()

		//FHD WebP
		go func() {
			fhd, err := image.Process(bimg.Options{
				Width:   fhdWidth,
				Height:  fhdHeight,
				Type:    bimg.WEBP,
				Quality: 95,
			})
			if err != nil {
				log.Println(err.Error())
				return
			}

			fhdWebpID, _, fhdWebpURL, err := onedrive.UploadAndShare(path+noext+".fhd.webp", fhd)
			if err != nil {
				log.Println(err.Error())
				return
			}
			tu.OneDriveFHDWebPID = fhdWebpID
			tu.OneDriveFHDWebPURL = fhdWebpURL

			err = models.UpdateTu(tu)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}()
	}

	return
}
