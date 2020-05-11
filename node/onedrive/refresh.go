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

func (node *Node) Refresh() (err error) {
	v := url.Values{
		"client_id":     {node.ClientID},
		"client_secret": {node.ClientSecret},
		"redirect_uri":  {node.RedirectURI},
		"grant_type":    {"refresh_token"},
		"refresh_token": {node.RefreshToken},
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

	node.AccessToken = refreshResponse.AccessToken
	node.RefreshToken = refreshResponse.RefreshToken
	return
}
