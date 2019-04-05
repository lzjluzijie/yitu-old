package onedrive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type ShareResponse struct {
 	ID string
 	Link SharedLink
}

type SharedLink struct {
	Scope string
	Type string
	WebURL string
}

func Share(id string) (url string) {
	req, err := http.NewRequest("POST", fmt.Sprintf("https://graph.microsoft.com/v1.0/me/drive/items/%s/createLink", id), bytes.NewBuffer([]byte(`{"type":"view","scope":"anonymous"}`)))
	if err != nil {
		log.Println(err.Error())
		return ""
	}

	req.Header.Set("Content-Type", "application/json")
	setHeader(req)

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

	jResp := &ShareResponse{}
	err = json.Unmarshal(data, jResp)
	if err != nil {
		log.Println(err.Error())
		return ""
	}

	url = jResp.Link.WebURL
	return
}
