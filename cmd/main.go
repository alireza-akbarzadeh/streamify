package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/techies/streamify/internal/app"
	"github.com/techies/streamify/internal/handler"
	"github.com/techies/streamify/internal/routes"
	"github.com/techies/streamify/internal/server"
)

func bootstrap() error {
	appCfg, err := app.New()
	if err != nil {
		return err
	}
	defer appCfg.Close()

	h := handler.NewHandler(appCfg)
	router := routes.SetupRoutes(h, appCfg)
	appCfg.Server.Handler = router

	return server.Run(appCfg.Server)
}

// @title           Streamify API
// @version         1.0
// @description     Core authentication and streaming service for Streamify.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.streamify.com/support
// @contact.email  support@streamify.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Type "Bearer" followed by a space and then your token.

func main() {
	_ = godotenv.Load()

	if err := bootstrap(); err != nil {
		log.Fatal(err)
	}
}
