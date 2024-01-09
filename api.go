package main

import (
	"net/http"
	"networkinator/models"
	"strconv"
    "fmt"

	"github.com/gin-gonic/gin"
)

func GetHosts(c *gin.Context) {

	db, err := connectToSQLite()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	hosts, err := getHostsEntries(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	hostMap := make(map[int]string)
	for _, host := range hosts {
		hostMap[host.ID] = host.IP
	}

	c.JSON(http.StatusOK, hostMap)
}

func GetConnections(c *gin.Context) {
	db, err := connectToSQLite()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	connections, err := getConnections(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	connectionMap := make(map[string][]string)
	for _, connection := range connections {
		connectionMap[connection.ID] = []string{connection.Src, connection.Dst, strconv.Itoa(connection.Port)}
	}

	c.JSON(http.StatusOK, connectionMap)
}

func AddHost(c *gin.Context) {
	var jsonData models.Host
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ip := jsonData.IP
	db, err := connectToSQLite()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = createHost(db, ip)
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

	db, err := connectToSQLite()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	srcHost := models.Host{}
	dstHost := models.Host{}

	tx := db.First(&srcHost, "IP = ?", src)
	if tx.Error != nil {
		createHost(db, src)
		db.First(&srcHost, "IP = ?", src)
	}

	tx = db.First(&dstHost, "IP = ?", dst)
	if tx.Error != nil {
		createHost(db, dst)
		db.First(&dstHost, "IP = ?", dst)
	}

	tx = db.First(&models.Connection{}, "Src = ? AND Dst = ? AND Port = ?", srcHost.IP, dstHost.IP, portInt)
	if tx.Error == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Connection already exists"})
		return
	}

	err = createConnection(db, srcHost.IP, dstHost.IP, portInt)
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
