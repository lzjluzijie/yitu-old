package node

import (
	"time"
)

type Node interface {
	Handle(*Object) interface{}
}

type Model struct {
	ID        uint64 `gorm:"primary_key"`
	TuID      uint64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Object struct {
	Name   string
	Size   int64
	MD5    string
	SHA256 string
	Data   []byte
}
