package db

import "github.com/jinzhu/gorm"

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
	db.AutoMigrate(&Tu{})
}
