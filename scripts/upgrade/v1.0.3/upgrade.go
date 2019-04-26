package main

import (
	"log"
	"strings"

	"github.com/lzjluzijie/yitu/models"
	"github.com/lzjluzijie/yitu/onedrive"
)

/*
upgrade from v1.0.2 to v1.0.3
re-share images and update urls
*/

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	onedrive.LoadConfig()
	models.PrepareEngine()
	x := models.Engine()

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

		log.Println(tu.OneDriveID)

		if tu.OneDriveURL != "" {
			tu.OneDriveURL = CutURL(tu.OneDriveURL)
		}

		if tu.OneDriveWebPURL != "" {
			tu.OneDriveWebPURL = CutURL(tu.OneDriveWebPURL)
		}

		if tu.OneDriveFHDURL != "" {
			tu.OneDriveFHDURL = CutURL(tu.OneDriveFHDURL)
		}

		if tu.OneDriveFHDWebPURL != "" {
			tu.OneDriveFHDWebPURL = CutURL(tu.OneDriveFHDWebPURL)
		}

		_, err = x.Table(new(models.Tu)).Where("sha256=?", tu.SHA256).Update(map[string]interface{}{
			"one_drive_url":          tu.OneDriveURL,
			"one_drive_webp_url":     tu.OneDriveWebPURL,
			"one_drive_fhd_url":      tu.OneDriveFHDURL,
			"one_drive_fhd_webp_url": tu.OneDriveFHDWebPURL,
		})
		if err != nil {
			panic(err)
		}

		log.Printf("%d: %s", tu.ID, tu.OneDriveURL)
		id = tu.ID
	}
	return
}

func CutURL(url string) string {
	p := strings.LastIndexByte(url, '&')
	if p <= 0 {
		return url
	}

	return url[:p]
}
