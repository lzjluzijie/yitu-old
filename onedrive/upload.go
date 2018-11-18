package onedrive

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type UploadResponse struct {
	ID string
}

type CreateSessionResponse struct {
	UploadUrl string
}

var st int64
var n int

func init() {
	st = time.Now().Unix()
	n = 1
}

func Upload(name string, r io.Reader) (id string, err error) {
	u := fmt.Sprintf("https://graph.microsoft.com/v1.0/me/drive/root:/6tu/%d/%d/%s:/content", st, n, name)
	n++

	req, err := http.NewRequest("PUT", u, r)
	if err != nil {
		return
	}

	setHeader(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	//log.Println(string(data))

	uResp := &UploadResponse{}
	err = json.Unmarshal(data, uResp)
	if err != nil {
		return
	}

	id = uResp.ID
	return
}

func UploadLarge(name string, size int64, r io.Reader) (id string, err error) {
	u := fmt.Sprintf("https://graph.microsoft.com/v1.0/me/drive/root:/6tu/%d/%d/%s:/createUploadSession", st, n, name)

	fmt.Println(u)

	req, err := http.NewRequest("POST", u, nil)
	if err != nil {
		return
	}

	setHeader(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	fmt.Println(string(data))

	cResp := &CreateSessionResponse{}
	err = json.Unmarshal(data, cResp)
	if err != nil {
		return
	}

	uploadURL := cResp.UploadUrl

	req, err = http.NewRequest("PUT", uploadURL, r)
	if err != nil {
		return
	}

	//不知道为什么设置header没有用
	req.Header.Add("Content-Length", fmt.Sprintf("%d", size))
	req.Header.Add("Content-Range", fmt.Sprintf("bytes 0-%d/%d", size-1, size))
	req.ContentLength = size
	setHeader(req)
	//fmt.Println(req.Header)

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	//fmt.Println(string(data))

	uResp := &UploadResponse{}
	err = json.Unmarshal(data, uResp)
	if err != nil {
		return
	}

	id = uResp.ID
	return
}
