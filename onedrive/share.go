package onedrive

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

func Share(id string) (url string) {
	req, err := NewRequest("POST", fmt.Sprintf("https://graph.microsoft.com/v1.0/me/drive/items/%s/createLink", id), strings.NewReader(`{"type":"view","scope":"anonymous"}`))
	if err != nil {
		return ""
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err.Error())
		return ""
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return ""
	}

	log.Println(string(data))

	shareResponse := &ShareResponse{}
	err = json.Unmarshal(data, shareResponse)
	if err != nil {
		log.Println(err.Error())
		return ""
	}

	url = shareResponse.Link.WebURL
	return
}
