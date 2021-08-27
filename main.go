package main

import (
	"bookmanager/models"
	"bookmanager/routers"
)

func main(){
	models.Ms.Sum=0
	models.InitDb()
	models.InitClient()
	routers.InitRouter()

}