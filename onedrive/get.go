package onedrive

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type GetURLResponse struct {
	URL string `json:"@microsoft.graph.downloadUrl"`
}

func GetURL(id string) (url string) {
	req, err := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/me/drive/items/"+id, nil)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
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

	jResp := &GetURLResponse{}
	err = json.Unmarshal(data, jResp)
	if err != nil {
		log.Println(err.Error())
		return ""
	}

	url = jResp.URL
	return
}
