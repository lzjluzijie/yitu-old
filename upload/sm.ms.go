package upload

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/lzjluzijie/6tu/core"
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

func (uploader *SMMSUploader) Upload(image *core.Image) (err error) {
	mr := multipartreader.NewMultipartReader()
	mr.AddFormReader("smfile", image.Name, image.Reader, image.Size)
	//form := fmt.Sprintf("--%s\r\nContent-Disposition: form-data; name=smfile; filename=\"%s\"\r\n\r\n", mr.Boundary, image.Name)
	//mr.AddReader(strings.NewReader(form), int64(len(form)))
	//mr.AddReader(image.Reader, image.Size)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://sm.ms/api/upload", mr)
	if err != nil {
		return
	}

	mr.SetupHTTPRequest(req)
	req.Header.Add("User-Agent", "secimage")

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
		return errors.New(fmt.Sprintf("unknown code %s", uResp.Code))
	}

	//log.Println(uResp)
	image.URL = uResp.Data.URL
	return
}
