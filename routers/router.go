package routers

import (
	"fmt"

	"encoding/base64"
	"io"
	"log"
	"net/http"

	"github.com/lzjluzijie/6tu/core"
	"github.com/lzjluzijie/6tu/upload"
	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/crypto/sha3"
	"gopkg.in/macaron.v1"
)

var db *leveldb.DB

func RegisterRouters(m *macaron.Macaron) {
	d, err := leveldb.OpenFile("./db", nil)
	if err != nil {
		panic(err)
	}
	db = d

	m.Get("/", func(ctx *macaron.Context) {
		ctx.HTML(200, "home")
	})

	m.Get("/i/:hash", GetImage)

	m.Post("/api/upload", Upload)
}

func GetImage(ctx *macaron.Context) {
	hashString := ctx.Params(":hash")
	log.Println(hashString)
	hash, err := base64.URLEncoding.DecodeString(hashString)
	if err != nil {
		ctx.Error(403, err.Error())
		return
	}

	url, err := db.Get(hash, nil)
	if err != nil {
		ctx.Error(403, err.Error())
		return
	}

	ctx.Redirect(string(url), http.StatusTemporaryRedirect)
}

func Upload(ctx *macaron.Context) {
	fr, fh, err := ctx.GetFile("image")
	if err != nil {
		ctx.Error(403, err.Error())
		return
	}
	h := sha3.New512()

	tr := io.TeeReader(fr, h)

	image := &core.Image{
		Name: fh.Filename,
		Size: fh.Size,

		Reader: tr,
	}

	u := &upload.SMMSUploader{}

	err = u.Upload(image)
	if err != nil {
		ctx.Error(403, err.Error())
		return
	}

	//upload ok
	hash := h.Sum(nil)
	err = db.Put(hash, []byte(image.URL), nil)
	if err != nil {
		ctx.Error(403, err.Error())
		return
	}

	image.URL = fmt.Sprintf("https://6tu.halu.lu/i/%s", base64.URLEncoding.EncodeToString(hash))
	image.Hash = hash
	fmt.Println(image)

	ctx.JSON(200, image)
}
