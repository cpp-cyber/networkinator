package main

import (
	"net/http"
    "fmt"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}


func addPublicRoutes(g *gin.RouterGroup) {
	g.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})
	g.GET("/status", status)
    g.GET("/api/hosts", GetHosts)
    g.GET("/api/hosts/:filter", GetHosts)
	g.GET("/api/connections", GetConnections)
    g.GET("/ws", ws)

	g.POST("/api/hosts", AddHost)
	g.POST("/api/connections", AddConnection)
}

func status(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}

func ws(c *gin.Context) {
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    clients[conn] = true

    go handleWebSocketConnection(conn)
}


func handleWebSocketConnection(conn *websocket.Conn) {
  for {
    _, _, err := conn.ReadMessage()
    if err != nil {
      fmt.Println(err)
      conn.Close()
      delete(clients, conn)
      break
    }
  }
}
