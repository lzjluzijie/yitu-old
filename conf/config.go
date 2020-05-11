package conf

import (
	"encoding/json"
	"io/ioutil"

	"github.com/lzjluzijie/yitu/db"

	"github.com/lzjluzijie/yitu/node/onedrive"
)

type Config struct {
	HttpAddr string
	Database *db.Config
	OneDrive *onedrive.Config
}

func LoadConfig() (config *Config) {
	data, err := ioutil.ReadFile("yitu.json")
	if err != nil {
		panic(err.Error())
	}

	config = &Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		panic(err.Error())
	}
	return
}

func (config *Config) Save() {
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		panic(err.Error())
	}

	err = ioutil.WriteFile("yitu.json", data, 0600)
	if err != nil {
		panic(err.Error())
	}
}
