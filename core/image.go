package core

import "io"

type Image struct {
	Name string
	Size int64
	Hash []byte `json:"-"`

	URL string

	Reader io.Reader `json:"-"`
}
