package onedrive

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

var config *Config

type Config struct {
	ClientID     string
	ClientSecret string
	AccessToken  string
	RefreshToken string
}

func setHeader(req *http.Request) *http.Request {
	req.Header.Add("Authorization", "bearer "+config.AccessToken)
	req.Header.Add("User-Agent", "6tu")
	return req
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
