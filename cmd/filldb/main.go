package main

import (
	"log"

	"github.com/gvidow/go-technical-equipment/internal/app/ds"
	"github.com/gvidow/go-technical-equipment/internal/app/dsn"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()))
	if err != nil {
		log.Fatal(err)
	}
	equipments := []ds.Equipment{
		{Title: "Лазерная ручка", Picture: "/upload/equipment/lizer.png", Description: "ручка", Status: "active", Count: 0, AvailableNow: 0},
		{Title: "Проектор", Picture: "/upload/equipment/projector.png", Description: "проектор", Status: "active", Count: 0, AvailableNow: 0},
		{Title: "Экран", Picture: "/upload/equipment/display.png", Description: "экран", Status: "active", Count: 0, AvailableNow: 0},
		{Title: "Лазерная ручка", Picture: "/upload/equipment/lizer.png", Description: "ручка", Status: "active", Count: 0, AvailableNow: 0},
		{Title: "Проектор", Picture: "/upload/equipment/projector.png", Description: "проектор", Status: "active", Count: 0, AvailableNow: 0},
		{Title: "Экран", Picture: "/upload/equipment/display.png", Description: "экран", Status: "active", Count: 0, AvailableNow: 0},
	}
	db.Create(equipments)
	log.Println("OK")
}
