package onedrive

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type ShareResponse struct {
	ID   string
	Link SharedLink
}

type SharedLink struct {
	Scope  string
	Type   string
	WebURL string
}

func Share(id string) (url string, err error) {
	req, err := NewRequest("POST", fmt.Sprintf("https://graph.microsoft.com/v1.0/me/drive/items/%s/createLink", id), strings.NewReader(`{"type":"view","scope":"anonymous"}`))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	shareResponse := &ShareResponse{}
	err = json.Unmarshal(data, shareResponse)
	if err != nil {
		return
	}

	url = shareResponse.Link.WebURL
	return
}
