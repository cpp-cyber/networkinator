package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func addPublicRoutes(g *gin.RouterGroup) {
    g.GET("/", index)
    g.GET("/connections", connections)
    g.GET("/filter", filter)
    g.GET("/agents", agents)

	g.GET("/api/connections/get", GetConnections)
    g.GET("/api/agents/get", GetAgents)
    g.POST("/api/agents/add", AddAgent)
    g.GET("/ws", ws)
    g.GET("/ws/agent/status", wsAgentStatus)
}

func index(c *gin.Context) {
    c.HTML(200, "index.html", gin.H{})
}

func connections(c *gin.Context) {
    c.HTML(200, "connections.html", gin.H{})
}

func filter(c *gin.Context) {
    c.HTML(200, "filter.html", gin.H{})
}

func agents(c *gin.Context) {
    c.HTML(200, "agents.html", gin.H{})
}

func ws(c *gin.Context) {
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    clients[conn] = true

    go handleAgentConnectionSocket(conn)
}

func wsAgentStatus(c *gin.Context) {
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    agentStatusClients[conn] = true

    go handleAgentStatusSocket(conn)
}

func handleAgentConnectionSocket(conn *websocket.Conn) {
  for {
    _, msg, err := conn.ReadMessage()
    if err != nil {
      fmt.Println(err)
      conn.Close()
      delete(clients, conn)
      break
    }
    AddConnection(msg)
  }
}

func handleAgentStatusSocket(conn *websocket.Conn) {
    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            fmt.Println(err)

            deadClient := strings.Split(conn.NetConn().RemoteAddr().String(), ":")[0]

            deadAgent, err := GetAgentByIP(deadClient)
            if err != nil {
                fmt.Println(err)
                return
            }

            jsonData := []byte(fmt.Sprintf(`{"ID": "%s", "Status": "Dead"}`, deadAgent.ID))
            for client := range agentStatusClients {
                client.WriteMessage(websocket.TextMessage, jsonData)
                UpdateAgentStatus(deadClient, "Dead")
            }

            conn.Close()
            delete(agentStatusClients, conn)
            break
        }
        AgentStatus(msg)
    }
}
