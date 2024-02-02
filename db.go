package main

import (
	"networkinator/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectToSQLite() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("network.db"), &gorm.Config{})
	if err != nil {
        panic(err)
	}
    return db
}

func AddAgentToDB(ID string, hostname, hostOS, ip string) error {
    agent := models.Agent{ID: ID, Hostname: hostname, HostOS: hostOS, IP: ip, Status: "Alive"}
    result := db.Create(&agent)
    if result.Error != nil {
        return result.Error
    }
    return nil
}

func GetAllAgents() ([]models.Agent, error) {
    var agents []models.Agent
    result := db.Find(&agents)
    if result.Error != nil {
        return nil, result.Error
    }
    return agents, nil
}

func GetAgentByIP(ip string) (models.Agent, error) {
    var agent models.Agent
    result := db.Where("IP = ?", ip).First(&agent)
    if result.Error != nil {
        return agent, result.Error
    }
    return agent, nil
}
