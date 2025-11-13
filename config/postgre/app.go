package postgre

import (
	"database/sql"
	service "go-fiber/app/service/postgre"
	"go-fiber/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func NewApp(db *sql.DB) *fiber.App {
	app := fiber.New(fiber.Config{
		BodyLimit: 2 * 1024 * 1024, // 2MB to support PDF; photo is checked in handler
	})
	app.Use(middleware.LoggerMiddleware)
	app.Use(cors.New())
	app.Post("/check/:key", func(c *fiber.Ctx) error {
		return service.CheckAlumniService(c, db)
	})
	return app
}
