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

	Width  int
	Height int

	OneDriveFolderID string
	OneDriveID       string
	OneDriveURL      string
	OneDriveWebpID   string
	OneDriveWebpURL  string

	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
}

func InsertTu(tu *Tu) (err error) {
	_, err = x.Insert(tu)
	return
}

func UpdateTu(tu *Tu) (err error) {
	_, err = x.ID(tu.ID).Update(tu)
	return
}

func DeleteTu(tu *Tu) (err error) {
	_, err = x.ID(tu.ID).Delete(tu)
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
	_, err = x.Get(tu)
	return
}
