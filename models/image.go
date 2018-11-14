package models

import (
	"errors"
	"fmt"
	"log"
	"time"
)

type Image struct {
	ID         int64 `xorm:"pk autoincr"`
	Name       string
	Size       int64
	Hash       []byte `json:"-"`
	Short      string
	OneDriveID string
	URL        string

	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`

	//r io.Reader `json:"-"`
}

func InsertImage(image *Image) (err error) {
	_, err = x.Insert(image)
	if err != nil {
		return
	}
	log.Println(fmt.Sprintf("inserted %v", image))
	return
}

func UpdateImage(image *Image) (err error) {
	_, err = x.ID(image.ID).Update(image)
	if err != nil {
		return
	}
	log.Println(fmt.Sprintf("updated %v", image))
	return
}

func GetImageFromShort(short string) (image *Image, err error) {
	//if len(short) != 16 {
	//	err = errors.New(fmt.Sprintf("short length != 16: %s", short))
	//	return
	//}

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
