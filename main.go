package main

import (
	"log"

	"github.com/astaxie/beego"
	_ "github.com/b0ralgin/test-beego/routers"
	"github.com/b0ralgin/test-beego/utilities"
)

func main() {
	err := utilities.MongoStartup()
	if err != nil {
		log.Fatal(err.Error())
	}
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
