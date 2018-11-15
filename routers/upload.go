package routers

import (
	"github.com/lzjluzijie/6tu/onedrive"

	"gopkg.in/macaron.v1"
)

type UploadResponse struct {
	Size int64
	URL  string
}

func Upload(ctx *macaron.Context) {
	r, fh, err := ctx.GetFile("tu")

	if err != nil {
		ctx.Error(403, err.Error())
		return
	}

	name := fh.Filename
	size := fh.Size

	//upload
	id, err := onedrive.Upload(name, r)
	if err != nil {
		ctx.Error(403, err.Error())
		return
	}

	////upload from url
	//u := ctx.Query("url")
	//if u != "" {
	//	resp, err := http.Get(u)
	//	if err != nil {
	//		ctx.Error(403, err.Error())
	//		return
	//	}
	//
	//	image, err = models.NewImageFromStream(resp.Body, resp.Request.URL.Path, resp.ContentLength)
	//	if err != nil {
	//		ctx.Error(403, err.Error())
	//		return
	//	}
	//} else {
	//	r, fh, err := ctx.GetFile("tu")
	//	if err != nil {
	//		ctx.Error(403, err.Error())
	//		return
	//	}
	//
	//	image, err = models.NewImageFromStream(r, fh.Filename, fh.Size)
	//	if err != nil {
	//		ctx.Error(403, err.Error())
	//		return
	//	}
	//}

	resp := &UploadResponse{
		Size: size,
		URL:  "https://6tu.halu.lu/i/" + id + "/" + name,
	}
	ctx.JSON(200, resp)
}
