package core

import "io"

type Image struct {
	Name string
	Size int64

	URL string

	Reader io.ReadSeeker
}
