package route

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"path"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lzjluzijie/yitu/db"
)

const MAXSIZE = 50 * 1024 * 1024

type UploadResponse struct {
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	MD5       string `json:"md5"`
	SHA256    string `json:"sha256"`
	URL       string `json:"url"`
	DeleteURL string `json:"delete_url"`
}

func NewUploadResponse(tu *db.Tu) (resp UploadResponse) {
	return UploadResponse{
		Name:      tu.Name,
		Size:      tu.Size,
		MD5:       tu.MD5,
		SHA256:    tu.SHA256,
		URL:       fmt.Sprintf("https://t.halu.lu/t/%d", tu.ID),
		DeleteURL: fmt.Sprintf("https://t.halu.lu/api/delete/%s", tu.DeleteCode),
	}
}

func GetTu(c *gin.Context) {
	id := c.Param("id")
	//t := c.Param("type")

	tu := &db.Tu{}

	if len(id) == 32 {
		tu.MD5 = id
	} else if len(id) == 64 {
		tu.SHA256 = id
	} else {
		tid, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		tu.ID = tid
	}

	db.GetDB().Where(tu).First(tu)

	if tu.ID == 0 {
		c.String(http.StatusNotFound, "not found")
		return
	}

	//c.Header("Cache-Control", "public, max-age=3110400")
	c.String(200, tu.SHA256)
	//c.Redirect(http.StatusMovedPermanently, tu.OneDriveURL)
	return
}

func UploadTu(c *gin.Context) {
	var name string
	var size int64
	var file io.Reader
	tu := &db.Tu{}

	url := c.PostForm("url")
	if url != "" {
		client := http.DefaultClient
		jar, err := cookiejar.New(nil)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		client.Jar = jar

		resp, err := client.Get(url)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		name = path.Base(resp.Request.URL.String())
		size = resp.ContentLength
		file = resp.Body
	} else {
		f, err := c.FormFile("tu")
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		name = f.Filename
		size = f.Size
		file, err = f.Open()
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
	}

	if size >= MAXSIZE {
		c.String(http.StatusBadRequest, "file too big")
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
	if size != int64(len(data)) {
		c.String(http.StatusBadRequest, fmt.Sprintf("file size does not match: %d, %d", size, int64(len(data))))
		return
	}

	//check sha256, insert new if already exist
	tu.SHA256 = SHA256
	db.GetDB().Where(tu).First(tu)

	if tu.ID != 0 {
		tu.ID = 0
		tu.Name = name
		c.JSON(200, NewUploadResponse(tu))
		return
	}

	//bimg check image
	//image := bimg.NewImage(data)
	//is, err := image.Size()
	//if err != nil {
	//	c.String(http.StatusBadRequest, err.Error())
	//	return
	//}

	///check orientation
	//metadata, err := image.Metadata()
	//if err != nil {
	//	c.String(http.StatusBadRequest, err.Error())
	//	return
	//}

	//width := is.Width
	//height := is.Height
	//if metadata.Orientation == 6 {
	//	width = is.Height
	//	height = is.Width
	//}

	//upload original image
	path := fmt.Sprintf(`/yitu/%s/%s/`, time.Now().Format("20060102"), SHA256)
	//ext := filepath.Ext(name)
	//noext := strings.TrimSuffix(name, ext)
	id, parent, url, err := n.UploadAndShare(path+name, data)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	//insert to database
	tu = &db.Tu{
		Name:   name,
		Size:   size,
		MD5:    MD5,
		SHA256: SHA256,
		IP:     c.ClientIP(),
		//DeleteCode: DeleteCode(SHA256),
		//Width:      width,
		//Height:     height,

		OneDriveFolderID: parent,
		OneDriveID:       id,
		OneDriveURL:      url,
	}

	db.GetDB().Create(tu)
	c.JSON(200, NewUploadResponse(tu))
}
