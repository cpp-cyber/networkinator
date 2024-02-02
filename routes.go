package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func addPublicRoutes(g *gin.RouterGroup) {
    g.GET("/", basicAuth)
    g.GET("/ws/agent", wsAgent)
    g.POST("/api/agents/add", AddAgent)
}

func addPrivateRoutes(g *gin.RouterGroup) {
    g.GET("/home", index)
    g.GET("/api/agents/get", GetAgents)
    g.GET("/ws/web", wsWeb)
}

func index(c *gin.Context) {
    c.HTML(200, "index.html", gin.H{})
}
