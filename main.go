package main

import (
	"bookmanager/models"
	"bookmanager/routers"
)

func main(){
	models.InitDb()
	models.InitClient()
	routers.InitRouter()
}