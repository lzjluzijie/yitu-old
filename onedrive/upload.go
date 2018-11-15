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

	jResp := &UploadResponse{}
	json.Unmarshal(data, jResp)

	id = jResp.ID
	return
}
