package routers

import (
	"io"
	"log"
	"path/filepath"

	"github.com/lzjluzijie/base36"

	"github.com/lzjluzijie/6tu/onedrive"

	"github.com/lzjluzijie/6tu/models"
	"gopkg.in/macaron.v1"
)

type UploadResponse struct {
	Size int64
	URL  string
}

func Upload(ctx *macaron.Context) {
	var r io.Reader
	fr, fh, err := ctx.GetFile("tu")
	if err != nil {
		ctx.Error(403, err.Error())
		return
	}
	r = fr

	//upload
	id, hash, err := onedrive.Upload(r)
	if err != nil {
		ctx.Error(403, err.Error())
		return
	}

	ext := filepath.Ext(fh.Filename)
	short := base36.Encode(hash)[:16] + ext

	//rename
	url, err := onedrive.Rename(id, short)
	if err != nil {
		ctx.Error(403, err.Error())
		return
	}

	image := &models.Image{
		Name:       fh.Filename,
		Size:       fh.Size,
		Hash:       hash,
		Short:      short,
		OneDriveID: id,
		URL:        url,
	}

	err = models.InsertImage(image)
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

	log.Println(image)

	resp := &UploadResponse{
		Size: fh.Size,
		URL:  "https://6tu.halu.lu/i/" + short,
	}
	ctx.JSON(200, resp)
}
