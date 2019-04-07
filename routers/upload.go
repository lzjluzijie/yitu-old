package routers

import (
	"crypto/sha256"
	"fmt"
	"io"

	"github.com/lzjluzijie/yitu/onedrive"

	"gopkg.in/macaron.v1"
)

type UploadResponse struct {
	Size int64  `json:"size"`
	URL  string `json:"url"`
}

func Upload(ctx *macaron.Context) {
	file, fh, err := ctx.GetFile("tu")

	if err != nil {
		ctx.Error(400, err.Error())
		return
	}

	name := fh.Filename
	size := fh.Size

	if size >= 50*1024*1024 {
		ctx.Error(400, "file too big")
		return
	}

	//hash
	h := sha256.New()
	r := io.TeeReader(file, h)

	//upload
	id, parent, err := onedrive.Upload(name, size, r)
	if err != nil {
		ctx.Error(500, err.Error())
		return
	}

	hash := fmt.Sprintf("%x", h.Sum(nil))

	//rename parent folder
	err = onedrive.Rename(parent, hash)
	if err != nil {
		ctx.Error(500, err.Error())
		return
	}

	//share
	url := onedrive.Share(id)

	resp := &UploadResponse{
		Size: size,
		URL:  url + "?download=1",
	}
	ctx.JSON(200, resp)
}
