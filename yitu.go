package main

import (
	"log"
	"net/http"
	"time"

	"github.com/lzjluzijie/yitu/node/onedrive"

	"github.com/lzjluzijie/yitu/db"

	"github.com/lzjluzijie/yitu/conf"
	"github.com/lzjluzijie/yitu/route"
)

const VERSION = `v0.2.0`

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("yitu %s by halulu", VERSION)

	config := conf.LoadConfig()

	db.LoadDB(config.Database)
	n := &onedrive.Node{Config: config.OneDrive}
	route.SetN(n)
	go func() {
		for {
			err := n.Refresh()

			if err != nil {
				log.Println(err.Error())
				time.Sleep(5 * time.Second)
			} else {
				config.Save()
				time.Sleep(50 * time.Minute)
			}
		}
	}()

	r := route.NewEngine()
	err := http.ListenAndServe(config.HttpAddr, r)
	if err != nil {
		panic(err)
	}
}
