package main

import (
	"go-fiber/config"
	"go-fiber/database"
	"go-fiber/route"
	"os"
)

func main() {
	config.LoadEnv()
	db := database.ConnectDB()
	app := config.NewApp(db)
	route.AlumniRoutes(app, db)
	route.PekerjaanRoutes(app, db)
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	app.Listen(":" + port)
}
