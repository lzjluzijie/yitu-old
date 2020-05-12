package main

import (
	"log"
	"net/http"

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
	n := onedrive.NewNode(config.OneDrive)
	route.SetN(n)

	r := route.NewEngine()
	err := http.ListenAndServe(config.HttpAddr, r)
	if err != nil {
		panic(err)
	}
}
