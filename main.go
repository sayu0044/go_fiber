package main

import (
	"go-fiber/config"
	appConfig "go-fiber/config/postgre"
	"go-fiber/database"
	_ "go-fiber/docs"
	mongoRoute "go-fiber/route/mongo"
	route "go-fiber/route/postgre"
	"log"
	"os"

	fiberSwagger "github.com/gofiber/swagger"
)

// @title Go Fiber Mongo API
// @version 1.0
// @description Dokumentasi API untuk layanan MongoDB pada aplikasi Go Fiber.
// @host localhost:3000
// @BasePath /go-fiber-mongo
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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
	mongoRoute.FileRoutes(app, mongoDB)

	// Swagger documentation
	app.Get("/swagger/*", fiberSwagger.HandlerDefault)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	// Serve uploaded files statically
	app.Static("/uploads", "./uploads")

	// Start server di satu port saja
	app.Listen(":" + port)
}
