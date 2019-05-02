package onedrive

import (
	"fmt"
	"io"
	"net/http"
)

var config Config

type Config struct {
	ClientID     string
	ClientSecret string
	AccessToken  string
	RefreshToken string
}

func SetConfig(c Config) {
	config = c
}

func NewRequest(method, url string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, url, body)
	if err != nil {
		return
	}

	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", config.AccessToken))
	req.Header.Add("User-Agent", "6tu")
	return
}
