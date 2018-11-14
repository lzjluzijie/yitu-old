package routers

import (
	"fmt"
	"time"

	"github.com/lzjluzijie/6tu/onedrive"

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

	if image.UpdatedAt.Add(time.Minute * 59).Before(time.Now()) {
		url := onedrive.GetURL(image.OneDriveID)
		if url == "" {
			ctx.Error(404, fmt.Sprintf("not found: %s", short))
			return
		}

		image.URL = url
		models.UpdateImage(image)
	}

	ctx.Resp.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Resp.Header().Add("Pragma", "no-cache")
	ctx.Resp.Header().Add("Expire", "0")
	ctx.Redirect(image.URL, 302)
}
