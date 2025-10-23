package mongo

import (
	service "go-fiber/app/service/mongo"
	"go-fiber/middleware"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewApp(db *mongo.Database) *fiber.App {
	app := fiber.New()
	app.Use(middleware.LoggerMiddleware)
	app.Post("/check/:key", func(c *fiber.Ctx) error {
		return service.CheckAlumniService(c, db)
	})
	return app
}
