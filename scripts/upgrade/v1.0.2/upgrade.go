package main

import (
	"log"

	"github.com/lzjluzijie/yitu/models"
	"github.com/lzjluzijie/yitu/onedrive"
)

/*
upgrade from v1.0.1 to v1.0.2
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

		tu.OneDriveURL, err = ReShare(tu.OneDriveID)
		if err != nil {
			panic(err)
		}

		if tu.OneDriveWebPID != "" {
			tu.OneDriveWebPURL, err = ReShare(tu.OneDriveWebPID)
			if err != nil {
				panic(err)
			}
		}

		if tu.OneDriveFHDID != "" {
			tu.OneDriveFHDURL, err = ReShare(tu.OneDriveFHDID)
			if err != nil {
				panic(err)
			}
		}

		if tu.OneDriveFHDWebPID != "" {
			tu.OneDriveFHDWebPURL, err = ReShare(tu.OneDriveFHDWebPID)
			if err != nil {
				panic(err)
			}
		}

		_, err = x.Table(new(models.Tu)).Where("one_drive_id=?", tu.OneDriveID).Update(map[string]interface{}{
			"one_drive_url": tu.OneDriveURL,
			"one_drive_webp_url": tu.OneDriveWebPURL,
			"one_drive_fhd_url": tu.OneDriveFHDURL,
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

func ReShare(id string) (url string, err error) {
	url, err = onedrive.Share(id)
	if err != nil {
		return
	}

	url, err = onedrive.GetGuestURL(url)
	if err != nil {
		return
	}
	return
}
