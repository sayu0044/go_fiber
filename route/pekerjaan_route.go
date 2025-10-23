package route

import (
	"database/sql"
	"go-fiber/app/service"
	"go-fiber/middleware"

	"github.com/gofiber/fiber/v2"
)

func PekerjaanRoutes(app *fiber.App, db *sql.DB) {
	api := app.Group("/go-fiber")
	protected := api.Group("", middleware.AuthRequired())

	pekerjaan := protected.Group("/pekerjaan")
	pekerjaan.Get("/", middleware.UserAndAdmin(), func(c *fiber.Ctx) error {
		return service.GetAllPekerjaanService(c, db)
	})
	pekerjaan.Get("/trash", middleware.UserAndAdmin(), func(c *fiber.Ctx) error {
		return service.ListDeletedPekerjaanService(c, db)
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
	pekerjaan.Put("/soft-delete/:id", middleware.UserAndAdmin(), func(c *fiber.Ctx) error {
		return service.SoftDeletePekerjaanService(c, db)
	})
	pekerjaan.Put("/restore/:id", middleware.UserAndAdmin(), func(c *fiber.Ctx) error {
		return service.RestorePekerjaanService(c, db)
	})
	pekerjaan.Delete("/hard-delete/:id", middleware.UserAndAdmin(), func(c *fiber.Ctx) error {
		return service.HardDeletePekerjaanService(c, db)
	})
}
