package routers

import (
	"fmt"

	"github.com/lzjluzijie/6tu/core"
	"github.com/lzjluzijie/6tu/upload"
	"gopkg.in/macaron.v1"
	"log"
)

func RegisterRouters(m *macaron.Macaron) {
	m.Get("/", func(ctx *macaron.Context) {
		ctx.HTML(200, "home")
	})

	m.Post("/api/upload", Upload)
}

func Upload(ctx *macaron.Context) {
	f, fh, err := ctx.GetFile("image")
	if err != nil {
		ctx.Error(403, err.Error())
		return
	}

	image := &core.Image{
		Name: fh.Filename,
		Size: fh.Size,

		Reader: f,
	}

	u := &upload.SMMSUploader{}

	err = u.Upload(image)
	if err != nil {
		ctx.Error(403, err.Error())
		return
	}

	ctx.WriteHeader(200)
	log.Println(image.URL)
	_, err = ctx.Write([]byte(image.URL))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
