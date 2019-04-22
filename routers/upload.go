package routers

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
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
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	MD5       string `json:"md5"`
	URL       string `json:"url"`
	DeleteURL string `json:"delete_url"`
}

func GetUploadResponse(tu *models.Tu) (resp UploadResponse) {
	return UploadResponse{
		Name:      tu.Name,
		Size:      tu.Size,
		MD5:       tu.MD5,
		URL:       fmt.Sprintf("https://t.halu.lu/t/%d", tu.ID),
		DeleteURL: fmt.Sprintf("https://t.halu.lu/api/delete/%s", tu.DeleteCode),
	}
}

func RandomDeleteCode() (dc string) {
	r := make([]byte, 8)
	_, err := rand.Read(r)
	if err != nil {
		panic(err)
	}

	dc = fmt.Sprintf("%x", r)
	return
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
	m := md5.New()
	s := sha256.New()
	r := io.TeeReader(file, io.MultiWriter(m, s))
	data, err := ioutil.ReadAll(r)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	MD5 := fmt.Sprintf("%x", m.Sum(nil))
	SHA256 := fmt.Sprintf("%x", s.Sum(nil))

	//check file size
	size := int64(len(data))
	if f.Size != size {
		c.String(http.StatusBadRequest, fmt.Sprintf("file size does not match: %d, %d", f.Size, size))
		return
	}

	//check sha256, insert new if already exist
	tu, err := models.GetTuBySHA256(SHA256)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if tu.ID != 0 {
		tu.ID = 0
		tu.DeleteCode = RandomDeleteCode()

		err = models.InsertTu(tu)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(200, GetUploadResponse(tu))
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

	dc := RandomDeleteCode()

	//upload original image
	path := fmt.Sprintf(`/yitu/%s/%s/`, time.Now().Format("20060102"), SHA256)
	ext := filepath.Ext(name)
	noext := strings.TrimSuffix(name, ext)
	id, parent, url, err := onedrive.UploadAndShare(path+name, data)
	if err != nil {
		log.Println(err.Error())
		err = models.DeleteTu(tu)
		return
	}

	g, err := onedrive.GetGuestURL(url)
	if err != nil {
		log.Println(err.Error())
	} else {
		url = g
	}

	//insert to database
	tu = &models.Tu{
		Name:       name,
		Size:       size,
		MD5:        MD5,
		SHA256:     SHA256,
		IP:         c.ClientIP(),
		DeleteCode: dc,
		Width:      width,
		Height:     height,

		OneDriveFolderID: parent,
		OneDriveID:       id,
		OneDriveURL:      url,
	}

	err = models.InsertTu(tu)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(200, GetUploadResponse(tu))

	//async upload
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

		id, _, url, err := onedrive.UploadAndShare(path+noext+".webp", webp)
		if err != nil {
			log.Println(err.Error())
			return
		}

		g, err := onedrive.GetGuestURL(url)
		if err != nil {
			log.Println(err.Error())
		} else {
			url = g
		}

		tu.OneDriveWebPID = id
		tu.OneDriveWebPURL = url

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

			id, _, url, err := onedrive.UploadAndShare(path+noext+".fhd"+ext, fhd)
			if err != nil {
				log.Println(err.Error())
				return
			}

			g, err := onedrive.GetGuestURL(url)
			if err != nil {
				log.Println(err.Error())
			} else {
				url = g
			}

			tu.OneDriveFHDID = id
			tu.OneDriveFHDURL = url

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

			id, _, url, err := onedrive.UploadAndShare(path+noext+".fhd.webp", fhd)
			if err != nil {
				log.Println(err.Error())
				return
			}

			g, err := onedrive.GetGuestURL(url)
			if err != nil {
				log.Println(err.Error())
			} else {
				url = g
			}

			tu.OneDriveFHDWebPID = id
			tu.OneDriveFHDWebPURL = url

			err = models.UpdateTu(tu)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}()
	}

	return
}
