package onedrive

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
)

type UploadResponse struct {
	ID              string
	parentReference ParentReference
}

type ParentReference struct {
	ID   string
	Path string
}

type CreateSessionResponse struct {
	UploadUrl string
}

func Upload(name string, size int64, r io.Reader) (id, parent string, err error) {
	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/me/drive/root:/yitu/%s/%d/%s:/createUploadSession", date, rand.Uint64(), name)

	req, err := NewRequest("POST", url, nil)
	if err != nil {
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	createSessionResponse := &CreateSessionResponse{}
	err = json.Unmarshal(data, createSessionResponse)
	if err != nil {
		return
	}

	uploadURL := createSessionResponse.UploadUrl

	req, err = NewRequest("PUT", uploadURL, r)
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

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	uploadResponse := &UploadResponse{}
	err = json.Unmarshal(data, uploadResponse)
	if err != nil {
		return
	}

	id = uploadResponse.ID
	parent = uploadResponse.parentReference.ID
	return
}
