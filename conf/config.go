package conf

import "C"
import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/lzjluzijie/yitu/onedrive"
)

type Config struct {
	HttpAddr  string
	HttpsAddr string
	Cert      string
	Key       string
	Database  DatabaseConfig
	OneDrive  onedrive.Config
}

type DatabaseConfig struct {
	Driver string
	Source string
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

	onedrive.SetConfig(config.OneDrive)

	go func() {
		for {
			oc, err := onedrive.Refresh()

			for err != nil {
				log.Println(err.Error())
				time.Sleep(5 * time.Second)
				oc, err = onedrive.Refresh()
			}

			config.OneDrive = oc
			config.Save()
			time.Sleep(59 * time.Minute)
		}
	}()
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
