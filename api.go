package main

import (
	"log"
	"net/http"
	"networkinator/models"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func GetAgents(c *gin.Context) {
    agents, err := GetAllAgents()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    agentArr := make([][]string, len(agents))
    for i := 0; i < len(agents); i++ {
        agentArr[i] = []string{agents[i].Hostname, agents[i].HostOS, agents[i].IP, agents[i].ID}
    }

    c.JSON(http.StatusOK, agentArr)
}

func AddAgent(c *gin.Context) {
    jsonData := make(map[string]interface{})
    err := c.ShouldBindJSON(&jsonData)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    hostname := jsonData["Hostname"].(string)
    hostOS := jsonData["HostOS"].(string)
    id := jsonData["ID"].(string)
    ip := strings.Split(c.ClientIP(), ":")[0]

    agent := models.Agent{}
    tx := db.First(&agent, "Hostname = ?", hostname)
    if tx.Error == nil {
        c.JSON(http.StatusOK, gin.H{"message": "Agent already exists"})
        return
    }

    err = AddAgentToDB(id, hostname, hostOS, ip)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Agent added"})
}

func AgentStatus(jsonData []byte) {
    for client := range webClients {
        err := client.WriteMessage(websocket.TextMessage, jsonData)
        if err != nil {
            log.Println(err)
            client.Close()
            delete(webClients, client)
        }
    }
}
