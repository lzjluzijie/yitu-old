package onedrive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type UploadResponse struct {
	ID              string
	ParentReference ParentReference
}

type ParentReference struct {
	ID   string
	Path string
}

type CreateSessionResponse struct {
	UploadUrl string
}

func Upload(path string, data []byte) (id, parent string, err error) {
	size := int64(len(data))
	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/me/drive/root:%s:/createUploadSession", path)

	req, err := NewRequest("POST", url, nil)
	if err != nil {
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	createSessionResponse := &CreateSessionResponse{}
	err = json.NewDecoder(resp.Body).Decode(createSessionResponse)
	if err != nil {
		return
	}

	uploadURL := createSessionResponse.UploadUrl

	req, err = NewRequest("PUT", uploadURL, bytes.NewBuffer(data))
	if err != nil {
		return
	}

	req.Header.Add("Content-Length", fmt.Sprintf("%d", size))
	req.Header.Add("Content-Range", fmt.Sprintf("bytes 0-%d/%d", size-1, size))
	req.ContentLength = size

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	uploadResponse := &UploadResponse{}
	err = json.NewDecoder(resp.Body).Decode(uploadResponse)
	if err != nil {
		return
	}

	id = uploadResponse.ID
	parent = uploadResponse.ParentReference.ID
	return
}

func UploadAndShare(path string, data []byte) (id, parent, url string, err error) {
	id, parent, err = Upload(path, data)
	if err != nil {
		return
	}

	url, err = Share(id)
	if err != nil {
		return
	}
	url += "?download=1"
	return
}
