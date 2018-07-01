package models

import (
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/lzjluzijie/6tu/upload"
	"github.com/lzjluzijie/base36"
	"golang.org/x/crypto/sha3"
)

type Image struct {
	ID    int64 `xorm:"pk autoincr"`
	Name  string
	Size  int64
	Hash  []byte `json:"-"`
	Short string
	URL   string

	//todo
	U string

	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`

	//r io.Reader `json:"-"`
}

func NewImageFromStream(r io.Reader, name string, size int64) (image *Image, err error) {
	//todo check image

	h := sha3.New512()
	tr := io.TeeReader(r, h)

	//upload to smms
	uploader := &upload.SMMSUploader{}
	u, err := uploader.Upload(tr, name, size)
	if err != nil {
		return
	}

	//upload ok, get hash and short
	hash := h.Sum(nil)
	short := base36.Encode(hash)[:10]

	image = &Image{
		Name:  name,
		Size:  size,
		Hash:  hash,
		Short: short,
		URL:   fmt.Sprintf("https://6tu.halu.lu/i/%s", short),

		U: u,
	}

	//insert data
	_, err = x.Insert(image)
	if err != nil {
		return
	}
	log.Println(fmt.Sprintf("inserted %v", image))
	return
}

func GetImageFromShort(short string) (image *Image, err error) {
	if len(short) != 10 {
		err = errors.New(fmt.Sprintf("short length != 10: %s", short))
		return
	}

	image = &Image{
		Short: short,
	}

	has, err := x.Get(image)
	if err != nil {
		return
	}
	if !has {
		err = errors.New(fmt.Sprintf("short not found: %s", short))
		return
	}
	return
}

func (image *Image) GetRedirectURL() (url string) {
	if image.U == "" {
		return ""
	}

	return image.U
}
