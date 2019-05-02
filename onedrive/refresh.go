package onedrive

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func Refresh() (c Config, err error) {
	v := url.Values{
		"client_id":     {config.ClientID},
		"client_secret": {config.ClientSecret},
		"redirect_uri":  {config.RedirectURI},
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

	refreshResponse := &RefreshResponse{}
	err = json.Unmarshal(data, refreshResponse)
	if err != nil {
		return
	}

	config.AccessToken = refreshResponse.AccessToken
	config.RefreshToken = refreshResponse.RefreshToken

	c = config
	return
}
