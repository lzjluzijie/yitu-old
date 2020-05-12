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

func (n *Node) Refresh() (err error) {
	v := url.Values{
		"client_id":     {n.ClientID},
		"client_secret": {n.ClientSecret},
		"redirect_uri":  {n.RedirectURI},
		"grant_type":    {"refresh_token"},
		"refresh_token": {n.RefreshToken},
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

	n.AccessToken = refreshResponse.AccessToken
	n.RefreshToken = refreshResponse.RefreshToken
	return
}
