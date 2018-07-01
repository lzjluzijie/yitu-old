package models

import (
	"log"

	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

var x *xorm.Engine

func PrepareEngine() {
	engine, err := xorm.NewEngine("sqlite3", "./6tu.db")
	if err != nil {
		panic(err)
	}

	engine.SetLogger(nil)
	engine.SetMapper(core.GonicMapper{})

	err = engine.Sync2(new(Image))
	if err != nil {
		panic(err)
	}

	x = engine
	log.Println("engine ok")
}
