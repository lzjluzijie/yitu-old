package routers

import (
	"time"

	"log"

	"github.com/lzjluzijie/6tu/onedrive"

	"gopkg.in/macaron.v1"
)

var cu map[string]string
var ct map[string]time.Time

func init() {
	cu = make(map[string]string)
	ct = make(map[string]time.Time)
}

func RegisterRouters(m *macaron.Macaron) {
	m.Get("/", func(ctx *macaron.Context) {
		ctx.HTML(200, "home")
	})

	m.Get("/i/:id/:name", GetImage)

	m.Group("/api", func() {
		m.Post("/upload", Upload)
	})

	log.Println("routers ok")
}

func GetImage(ctx *macaron.Context) {
	id := ctx.Params(":id")

	url := "https://6tu.halu.lu/"

	t, ok := ct[id]
	if !ok || t.Add(59*time.Minute).Before(time.Now()) {
		url = onedrive.GetURL(id)

		ct[id] = time.Now()
		cu[id] = url
	} else {
		url, _ = cu[id]
	}

	ctx.Resp.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Resp.Header().Add("Pragma", "no-cache")
	ctx.Resp.Header().Add("Expire", "0")
	ctx.Redirect(url, 302)
}
