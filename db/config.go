package db

import (
	"github.com/jinzhu/gorm"
	"github.com/lzjluzijie/yitu/node/onedrive"
)

type Config struct {
	Driver string
	Source string
}

func LoadDB(c *Config) {
	database, err := gorm.Open(c.Driver, c.Source)
	if err != nil {
		panic(err)
	}

	db = database
	db = db.Set("gorm:auto_preload", true)
	db.AutoMigrate(&Tu{}, &onedrive.File{})
}
