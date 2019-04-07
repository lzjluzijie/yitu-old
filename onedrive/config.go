package onedrive

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var config *Config
var date = time.Now().Format("20060102")

type Config struct {
	ClientID     string
	ClientSecret string
	AccessToken  string
	RefreshToken string
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

func LoadConfig() {
	data, err := ioutil.ReadFile("onedrive.json")
	if err != nil {
		panic(err.Error())
	}

	config = &Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		panic(err.Error())
	}
}

func SaveConfig() {
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		panic(err.Error())
	}

	err = ioutil.WriteFile("onedrive.json", data, 0600)
	if err != nil {
		panic(err.Error())
	}

	log.Println("saved")
}
