package postgre

import (
	"database/sql"
	service "go-fiber/app/service/postgre"
	"go-fiber/middleware"

	"github.com/gofiber/fiber/v2"
)

func NewApp(db *sql.DB) *fiber.App {
	app := fiber.New()
	app.Use(middleware.LoggerMiddleware)
	app.Post("/check/:key", func(c *fiber.Ctx) error {
		return service.CheckAlumniService(c, db)
	})
	return app
}
