package upload

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"io"

	"github.com/lzjluzijie/multipartreader"
)

type SMMSUploader struct {
}

type SMMSResponse struct {
	Code string
	Data SMMSResponseData
}

type SMMSResponseData struct {
	URL string
}

func (uploader *SMMSUploader) Upload(r io.Reader, name string, size int64) (url string, err error) {
	mr := multipartreader.NewMultipartReader()
	mr.AddFormReader(r, "smfile", name, size)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://sm.ms/api/upload", mr)
	if err != nil {
		return
	}

	mr.SetupHTTPRequest(req)
	req.Header.Add("User-Agent", "6tu")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	//log.Println(string(data))

	uResp := &SMMSResponse{}
	err = json.Unmarshal(data, uResp)
	if err != nil {
		return
	}

	if uResp.Code != "success" {
		err = errors.New(fmt.Sprintf("unknown code %s", uResp.Code))
		return
	}

	url = uResp.Data.URL
	return
}
