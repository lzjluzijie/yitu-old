package onedrive

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type RenameResponse struct {
	URL string `json:"@microsoft.graph.downloadUrl"`
}

func Rename(id, name string) (url string, err error) {
	req, err := http.NewRequest("PATCH", "https://graph.microsoft.com/v1.0/me/drive/items/"+id, strings.NewReader("{\"name\": \""+name+"\"}"))
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")
	setHeader(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	renameResponse := &RenameResponse{}
	err = json.Unmarshal(data, renameResponse)
	if err != nil {
		return
	}

	url = renameResponse.URL
	return
}
