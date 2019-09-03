package main

import (
	"fmt"
	"github.com/lzjluzijie/yitu/conf"
	"log"
	"net/http"

	"github.com/lzjluzijie/yitu/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/lzjluzijie/yitu/routers"
)

const VERSION = `v1.1.0`

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("yitu %s by halulu", VERSION)

	config := conf.LoadConfig()
	models.PrepareEngine(config.Database.Driver, config.Database.Source)

	engine := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowHeaders = []string{"*"}
	engine.Use(cors.New(corsConfig))

	routers.RegisterRouters(engine)

	//go func() {
	//	err := http.ListenAndServeTLS(config.HttpsAddr, "cert", "key", engine)
	//	if err != nil {
	//		fmt.Println(err.Error())
	//	}
	//}()

	//err := http.ListenAndServe(config.HttpAddr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	http.Redirect(w, r, "https://t.halu.lu"+r.URL.String(), http.StatusMovedPermanently)
	//}))
	//if err != nil {
	//	fmt.Println(err.Error())
	//}

	err := http.ListenAndServe(config.HttpAddr, engine)
	if err != nil {
		fmt.Println(err.Error())
	}
}
