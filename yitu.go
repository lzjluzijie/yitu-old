package main

import (
	"fmt"
	"net/http"

	"github.com/lzjluzijie/yitu/onedrive"

	"github.com/lzjluzijie/yitu/routers"
	"gopkg.in/macaron.v1"
)

func main() {
	m := macaron.Classic()
	m.Use(macaron.Renderer())

	routers.RegisterRouters(m)
	onedrive.LoadConfig()
	onedrive.Refresh()

	//m.Run()
	go func() {
		err := http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://6tu.halu.lu"+r.URL.String(), http.StatusMovedPermanently)
		}))
		if err != nil {
			fmt.Println(err.Error())
		}
	}()

	err := http.ListenAndServeTLS(":443", "cert", "key", m)
	if err != nil {
		fmt.Println(err.Error())
	}
}
