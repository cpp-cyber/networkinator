package main

import (
	"networkinator/models"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectToSQLite() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("network.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func AddConnectionToDB(db *gorm.DB, src, dst string, port, count int) error {
	id := uuid.New().String()
	connection := models.Connection{ID: id, Src: src, Dst: dst, Port: port, Count: count}
	result := db.Create(&connection)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetAllConnections(db *gorm.DB) ([]models.Connection, error) {
	var connections []models.Connection
	result := db.Find(&connections)
	if result.Error != nil {
		return nil, result.Error
	}
	return connections, nil
}

func GetConnectionsByIP(db *gorm.DB, host string) ([]models.Connection, error) {
    var connections []models.Connection
    result := db.Where("src = ? OR dst = ?", host, host).Find(&connections)
    if result.Error != nil {
        return nil, result.Error
    }
    return connections, nil
}

func IncrementConnectionCount(db *gorm.DB, id string) error {
    var connection models.Connection
    result := db.Where("id = ?", id).First(&connection)
    if result.Error != nil {
        return result.Error
    }
    connection.Count++
    result = db.Save(&connection)
    if result.Error != nil {
        return result.Error
    }
    return nil
}
