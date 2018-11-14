package onedrive

import (
	"crypto/rand"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/lzjluzijie/base36"
	"golang.org/x/crypto/sha3"
)

type UploadResponse struct {
	ID string
}

func Upload(r io.Reader) (id string, hash []byte, err error) {
	tmp := Random()
	u := "https://graph.microsoft.com/v1.0/me/drive/root:/6tu/" + tmp + ".png:/content"

	h := sha3.New256()
	tr := io.TeeReader(r, h)

	req, err := http.NewRequest("PUT", u, tr)
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

	hash = h.Sum(nil)
	id = jResp.ID
	return
}

func Random() (random string) {
	r := make([]byte, 256)
	_, err := rand.Read(r)

	if err != nil {
		return ""
	}

	s := base36.Encode(r)
	return s[:32]
}
