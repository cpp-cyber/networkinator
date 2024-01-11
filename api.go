package main

import (
	"fmt"
	"net/http"
	"networkinator/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetConnections(c *gin.Context) {
    db, err := ConnectToSQLite()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    connections, err := GetAllConnections(db)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    connectionMap := make(map[string][]string)
    for _, connection := range connections {
        connectionMap[connection.ID] = []string{connection.Src, connection.Dst, strconv.Itoa(connection.Port), strconv.Itoa(connection.Count)}
    }

    c.JSON(http.StatusOK, connectionMap)
}

func GetFilteredConnections(c *gin.Context) {
    filter := c.Param("filter")
    fmt.Println(filter)
    filterList := strings.Split(filter, ",")

	db, err := ConnectToSQLite()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

    connections := []models.Connection{}
    for _, ip := range filterList {
        connList, err := GetConnectionsByIP(db, ip)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        connections = append(connections, connList...)
    }

	connectionMap := make(map[string][]string)
	for _, connection := range connections {
		connectionMap[connection.ID] = []string{connection.Src, connection.Dst, strconv.Itoa(connection.Port)}
	}

	c.JSON(http.StatusOK, connectionMap)
}

func AddConnection(c *gin.Context) {
	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	src := jsonData["Src"].(string)
	dst := jsonData["Dst"].(string)
	port := jsonData["Port"].(string)

	portInt, err := strconv.Atoi(port)
	if err != nil || portInt < 0 || portInt > 65535 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not convert port to int"})
		return
	}

	db, err := ConnectToSQLite()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

    connection := models.Connection{}
    tx := db.First(&connection, "Src = ? AND Dst = ? AND Port = ?", src, dst, portInt)
	if tx.Error == nil {
        IncrementConnectionCount(db, connection.ID)
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
		return
	}

	err = AddConnectionToDB(db, src, dst, portInt, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

    for client := range clients {
        err := client.WriteJSON(jsonData)
        if err != nil {
            fmt.Println(err)
            client.Close()
            delete(clients, client)
        }
    }

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
