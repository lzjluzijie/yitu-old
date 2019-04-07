package models

import "time"

type Tu struct {
	ID   int64 `xorm:"pk autoincr"`
	Name string
	Size int64
	Hash string

	OneDriveURL string

	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
}
