package mongo

import (
	service "go-fiber/app/service/mongo"
	middleware "go-fiber/middleware/mongo"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func PekerjaanRoutes(app *fiber.App, db *mongo.Database) {
	api := app.Group("/go-fiber-mongo")
	protected := api.Group("", middleware.AuthRequired())

	pekerjaan := protected.Group("/pekerjaan")
	pekerjaan.Get("/", middleware.UserAndAdmin(), func(c *fiber.Ctx) error {
		return service.GetAllPekerjaanService(c, db)
	})
	pekerjaan.Get("/:id", middleware.UserAndAdmin(), func(c *fiber.Ctx) error {
		return service.GetPekerjaanByIDService(c, db)
	})
	pekerjaan.Get("/alumni/:alumni_id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.GetPekerjaanByAlumniIDService(c, db)
	})
	pekerjaan.Post("/", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.CreatePekerjaanService(c, db)
	})
	pekerjaan.Put("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.UpdatePekerjaanService(c, db)
	})
}
