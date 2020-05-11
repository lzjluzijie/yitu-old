package db

type Tu struct {
	Model
	Name   string
	Size   int64
	MD5    string
	SHA256 string

	IP         string
	DeleteCode string

	Requests uint64

	Width  int
	Height int

	OneDriveFolderID string
	OneDriveID       string
	OneDriveURL      string
}
