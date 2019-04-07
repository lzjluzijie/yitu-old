package onedrive

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func Refresh() (err error) {
	v := url.Values{
		"client_id":     {config.ClientID},
		"client_secret": {config.ClientSecret},
		"redirect_uri":  {"http://127.0.0.1:23333"},
		"grant_type":    {"refresh_token"},
		"refresh_token": {config.RefreshToken},
	}

	resp, err := http.PostForm("https://login.microsoftonline.com/common/oauth2/v2.0/token", v)
	if err != nil {
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	log.Println(string(data))

	refreshResponse := &RefreshResponse{}
	err = json.Unmarshal(data, refreshResponse)
	if err != nil {
		return
	}

	config.AccessToken = refreshResponse.AccessToken
	config.RefreshToken = refreshResponse.RefreshToken

	SaveConfig()

	go func() {
		time.Sleep(59 * time.Minute)
		err = Refresh()

		for err != nil {
			time.Sleep(5 * time.Second)
			err = Refresh()
		}

	}()

	return
}
