package models

import (
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/go-xorm/core.v0"
	"gopkg.in/go-xorm/xorm.v0"
)

var x *xorm.Engine

func PrepareEngine() {
	engine, err := xorm.NewEngine("sqlite3", "./yitu.db?parseTime=true&loc=Local")
	if err != nil {
		panic(err)
	}

	engine.SetLogger(xorm.NewSimpleLogger(os.Stdout))
	engine.SetMapper(core.GonicMapper{})

	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	engine.TZLocation = location

	err = engine.Sync2(new(Tu))
	if err != nil {
		panic(err)
	}

	x = engine
}

func Engine() *xorm.Engine {
	return x
}
