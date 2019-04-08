package models

import (
	"errors"
	"time"
)

type Tu struct {
	ID   int64 `xorm:"pk autoincr"`
	Name string
	Size int64
	Hash string

	OneDriveFolderID string
	OneDriveID       string
	OneDriveURL      string

	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
}

func InsertTu(tu *Tu) (err error) {
	_, err = x.Insert(tu)
	return
}

func GetTuByID(id uint64) (tu *Tu, err error) {
	tu = new(Tu)
	has, err := x.ID(id).Get(tu)
	if err != nil {
		return
	}

	if !has {
		err = errors.New("not found")
		return
	}
	return
}

func GetTuByHash(hash string) (tu *Tu, err error) {
	tu = &Tu{Hash: hash}
	has, err := x.Get(tu)
	if err != nil {
		return
	}

	if !has {
		err = errors.New("not found")
		return
	}
	return
}
