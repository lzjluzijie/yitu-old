package main

import (
	"github.com/lzjluzijie/6tu/routers"
	"gopkg.in/macaron.v1"
	"github.com/go-macaron/pongo2"
	"net/http"
	"fmt"
)

func main() {
	m := macaron.Classic()
	m.Use(pongo2.Pongoer())

	routers.RegisterRouters(m)

	//m.Run()
	go func() {
		err := http.ListenAndServe(":80", m)
		if err != nil {
			fmt.Println(err.Error())
		}
	}()

	err := http.ListenAndServeTLS(":443","cert","key",m)
	if err != nil {
		fmt.Println(err.Error())
	}
}
