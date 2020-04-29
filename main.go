package main

import (
	"github.com/akhettar/app-features-manager/api"
	_ "github.com/akhettar/app-features-manager/docs"
	"github.com/akhettar/app-features-manager/features"
	"github.com/akhettar/app-features-manager/repository"
	"github.com/labstack/gommon/log"
)

// @BasePath /
// @title App Status API
// @version 1.0

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /
func main() {

	log.Info("Starting up the server..")
	router := api.NewAppStatusHandler(repository.NewRepository(),
		features.NewUnleashClient()).CreateRouter()
	// Start server
	router.Logger.Fatal(router.Start(":1323"))
	log.Info("Shutting down the server..")
}
