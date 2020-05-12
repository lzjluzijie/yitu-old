package db

import "github.com/lzjluzijie/yitu/node/onedrive"

type Tu struct {
	Model
	Name   string
	Size   int64
	MD5    string
	SHA256 string

	IP string

	Requests uint64

	Width  int
	Height int

	OneDrive []onedrive.File `gorm:"foreignkey:TuID"`
}
