package main

import (
	"go-fiber/config"
	appConfig "go-fiber/config/postgre"
	"go-fiber/database"
	route "go-fiber/route/postgre"
	"os"
)

func main() {
	config.LoadEnv()
	db := database.ConnectDB()
	app := appConfig.NewApp(db)
	route.AlumniRoutes(app, db)
	route.PekerjaanRoutes(app, db)
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	app.Listen(":" + port)
}
