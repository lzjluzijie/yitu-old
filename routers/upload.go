package routers

import (
	"github.com/lzjluzijie/6tu/onedrive"

	"gopkg.in/macaron.v1"
)

type UploadResponse struct {
	Size int64  `json:"size"`
	URL  string `json:"url"`
}

func Upload(ctx *macaron.Context) {
	r, fh, err := ctx.GetFile("tu")

	if err != nil {
		ctx.Error(400, err.Error())
		return
	}

	name := fh.Filename
	size := fh.Size
	id := ""

	if size >= 50*1024*1024 {
		ctx.Error(400, "file too big")
		return
	} else if size >= 4*1024*1024 {
		id, err = onedrive.UploadLarge(name, size, r)
		if err != nil {
			ctx.Error(500, err.Error())
			return
		}
	} else {
		id, err = onedrive.Upload(name, r)
		if err != nil {
			ctx.Error(500, err.Error())
			return
		}
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

	url := onedrive.Share(id)

	resp := &UploadResponse{
		Size: size,
		URL:  url + "?download=1",
	}
	ctx.JSON(200, resp)
}
