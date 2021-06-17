package main

import (
	"forklift-bat-backend/model"
	"forklift-bat-backend/router"
)

func main() {
	router := router.InitRouter()
	model.SyncDatabase()
	router.Run(":8088")
}
