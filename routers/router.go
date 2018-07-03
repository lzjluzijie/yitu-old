package routers

import (
	"net/http"

	"fmt"
	"log"

	"github.com/lzjluzijie/6tu/models"
	"gopkg.in/macaron.v1"
)

func RegisterRouters(m *macaron.Macaron) {
	m.Get("/", func(ctx *macaron.Context) {
		ctx.HTML(200, "home")
	})

	m.Get("/i/:short", GetImage)

	m.Group("/api", func() {
		m.Post("/upload", Upload)
	})

	log.Println("routers ok")
}

func GetImage(ctx *macaron.Context) {
	short := ctx.Params(":short")
	image, err := models.GetImageFromShort(short)
	if err != nil {
		ctx.Error(403, err.Error())
		return
	}

	u := image.GetRedirectURL()
	if u == "" {
		ctx.Error(404, fmt.Sprintf("not found: %s", short))
		return
	}
	ctx.Redirect(u, http.StatusTemporaryRedirect)
}

func Upload(ctx *macaron.Context) {
	image := &models.Image{}
	//upload from url
	u := ctx.Query("url")
	if u != "" {
		resp, err := http.Get(u)
		if err != nil {
			ctx.Error(403, err.Error())
			return
		}

		image, err = models.NewImageFromStream(resp.Body, resp.Request.URL.Path, resp.ContentLength)
		if err != nil {
			ctx.Error(403, err.Error())
			return
		}
	} else {
		r, fh, err := ctx.GetFile("image")
		if err != nil {
			ctx.Error(403, err.Error())
			return
		}

		image, err = models.NewImageFromStream(r, fh.Filename, fh.Size)
		if err != nil {
			ctx.Error(403, err.Error())
			return
		}
	}

	log.Println(image)
	ctx.JSON(200, image)
}
