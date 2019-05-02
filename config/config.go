package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/lzjluzijie/yitu/onedrive"
)

var C *Config

type Config struct {
	Database DatabaseConfig
	OneDrive onedrive.Config
}

type DatabaseConfig struct {
	Driver string
	Source string
}

func LoadConfig() {
	data, err := ioutil.ReadFile("yitu.json")
	if err != nil {
		panic(err.Error())
	}

	C = &Config{}
	err = json.Unmarshal(data, C)
	if err != nil {
		panic(err.Error())
	}

	onedrive.SetConfig(C.OneDrive)

	go func() {
		for {
			oc, err := onedrive.Refresh()

			for err != nil {
				log.Println(err.Error())
				time.Sleep(5 * time.Second)
				oc, err = onedrive.Refresh()
			}

			C.OneDrive = oc
			SaveConfig()
			time.Sleep(59 * time.Minute)
		}
	}()
}

func SaveConfig() {
	data, err := json.MarshalIndent(C, "", "    ")
	if err != nil {
		panic(err.Error())
	}

	err = ioutil.WriteFile("yitu.json", data, 0600)
	if err != nil {
		panic(err.Error())
	}
}
