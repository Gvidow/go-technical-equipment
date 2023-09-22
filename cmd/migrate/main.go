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

	err = db.AutoMigrate(&ds.Equipment{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Migrate success")
}
