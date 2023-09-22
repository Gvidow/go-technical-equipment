package main

import (
	"fmt"
	"os"

	"github.com/gvidow/go-technical-equipment/internal/app"
	conf "github.com/gvidow/go-technical-equipment/internal/app/config"
	"github.com/gvidow/go-technical-equipment/logger"
	"github.com/joho/godotenv"
	"go.uber.org/config"
)

func main() {
	log, err := logger.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.NewYAML(config.File("configs/config.yml"))
	if err != nil {
		log.Fatal(err)
	}

	appCfg := conf.New(cfg)
	app, err := app.New(log, appCfg)
	if err != nil {
		log.Fatal(err)
	}

	if err = app.Run(); err != nil {
		log.Fatal(err)
	}
}
