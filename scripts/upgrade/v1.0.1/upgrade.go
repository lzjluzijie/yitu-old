package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"

	"github.com/lzjluzijie/yitu/onedrive"

	"github.com/lzjluzijie/yitu/models"
)

/*
upgrade from v1.0.0 to v1.0.1
need to rename column hash => sha256 first
then this script will insert md5 for every tu
*/

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	onedrive.LoadConfig()
	models.PrepareEngine()
	x := models.Engine()

	client := http.DefaultClient
	cookieJar, _ := cookiejar.New(nil)
	client.Jar = cookieJar

	id := int64(0)
	for {
		tu := &models.Tu{}
		has, err := x.Where("id>?", id).Get(tu)
		if err != nil {
			panic(err)
		}

		if !has {
			log.Println("done")
			break
		}

		if tu.MD5 != "" {
			id = tu.ID
			continue
		}

		if tu.OneDriveURL == "" {
			err = models.DeleteTu(tu)
			if err != nil {
				panic(err)
			}
		}

		log.Println(tu.OneDriveURL)
		m := md5.New()
		resp, err := client.Get(tu.OneDriveURL)
		if err != nil {
			panic(err)
		}

		_, err = io.Copy(m, resp.Body)
		if err != nil {
			panic(err)
		}

		MD5 := fmt.Sprintf("%x", m.Sum(nil))

		_, err = x.Table(new(models.Tu)).Where("sha256=?", tu.SHA256).Update(map[string]interface{}{"md5": MD5})
		if err != nil {
			panic(err)
		}

		log.Printf("%s: %s", tu.SHA256, MD5)
		id = tu.ID
	}
	return
}
