package main

import (
	"go-fiber/config"
	appConfig "go-fiber/config/postgre"
	"go-fiber/database"
	mongoRoute "go-fiber/route/mongo"
	route "go-fiber/route/postgre"
	"log"
	"os"
)

func main() {
	config.LoadEnv()

	// PostgreSQL setup
	db := database.ConnectDB()
	app := appConfig.NewApp(db)
	route.AlumniRoutes(app, db)
	route.PekerjaanRoutes(app, db)

	// MongoDB setup
	mongoDB := database.ConnectMongoDB()

	// Run MongoDB migrations
	if err := database.RunMigrations(mongoDB); err != nil {
		log.Fatalf("MongoDB migration failed: %v", err)
	}

	// Register MongoDB routes ke app yang sama
	mongoRoute.AlumniRoutes(app, mongoDB)
	mongoRoute.PekerjaanRoutes(app, mongoDB)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	// Start server di satu port saja
	app.Listen(":" + port)
}
