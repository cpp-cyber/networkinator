package main

import (
	"networkinator/models"
	"log"

	"github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

var HostCount int

var clients = make(map[*websocket.Conn]bool)

func main() {

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.Static("/assets", "./assets/")

	public := router.Group("/")
	addPublicRoutes(public)

	db, err := ConnectToSQLite()
	if err != nil {
		log.Fatalln(err)
	}

	err = db.AutoMigrate(&models.Connection{})
	if err != nil {
		log.Fatalln(err)
	}

    log.Fatalln(router.Run(":80"))
}
