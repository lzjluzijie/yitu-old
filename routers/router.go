package routers

import (
	"log"
	"strconv"

	"github.com/lzjluzijie/yitu/models"

	"gopkg.in/macaron.v1"
)

func RegisterRouters(m *macaron.Macaron) {
	m.Get("/", func(ctx *macaron.Context) {
		ctx.HTML(200, "home")
	})

	m.Get("/t/:id/:name", GetTu)

	m.Group("/api", func() {
		m.Post("/upload", Upload)
	})

	log.Println("routers ok")
}

func GetTu(ctx *macaron.Context) {
	t := ctx.Params(":id")

	id, err := strconv.ParseUint(t, 10, 64)
	if err != nil {
		ctx.Error(500, err.Error())
		return
	}

	tu, err := models.GetTuByID(id)
	if err != nil {
		ctx.Error(500, err.Error())
		return
	}

	ctx.Redirect(tu.OneDriveURL, 301)
}
