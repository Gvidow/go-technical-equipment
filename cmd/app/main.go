package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gvidow/go-technical-equipment/internal/app"
	conf "github.com/gvidow/go-technical-equipment/internal/app/config"
	"github.com/gvidow/go-technical-equipment/logger"
	"github.com/joho/godotenv"
	"go.uber.org/config"
)

//	@title			Swagger Equipment API
//	@version		1.0
//	@description	This is a backend server for project with bmstu equipments.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.basic	BasicAuth

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	app, err := app.New(ctx, log, appCfg)
	if err != nil {
		log.Fatal(err)
	}

	if err = app.Run(); err != nil {
		log.Fatal(err)
	}
}
