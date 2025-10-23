package postgre

import (
	"database/sql"
	service "go-fiber/app/service/postgre"
	middleware "go-fiber/middleware/postgre"

	"github.com/gofiber/fiber/v2"
)

func AlumniRoutes(app *fiber.App, db *sql.DB) {
	api := app.Group("/go-fiber-postgre")

	api.Post("/login", func(c *fiber.Ctx) error {
		return service.LoginService(c, db)
	})

	protected := api.Group("", middleware.AuthRequired())

	protected.Get("/profile", func(c *fiber.Ctx) error {
		return service.GetProfileService(c, db)
	})

	alumni := protected.Group("/alumni")
	alumni.Get("/", middleware.UserAndAdmin(), func(c *fiber.Ctx) error {
		return service.GetAllAlumniService(c, db)
	})
	alumni.Get("/:id", middleware.UserAndAdmin(), func(c *fiber.Ctx) error {
		return service.GetAlumniByIDService(c, db)
	})
	alumni.Post("/", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.CreateAlumniService(c, db)
	})
	alumni.Put("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.UpdateAlumniService(c, db)
	})
	alumni.Delete("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.DeleteAlumniService(c, db)
	})
	alumni.Post("/check/:key", func(c *fiber.Ctx) error {
		return service.CheckAlumniService(c, db)
	})

	roles := protected.Group("/roles")
	roles.Get("/", middleware.UserAndAdmin(), func(c *fiber.Ctx) error {
		return service.ListRolesService(c, db)
	})
	roles.Get("/:id", middleware.UserAndAdmin(), func(c *fiber.Ctx) error {
		return service.GetRoleByIDService(c, db)
	})
	roles.Post("/", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.CreateRoleService(c, db)
	})
	roles.Put("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.UpdateRoleService(c, db)
	})
	roles.Delete("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.DeleteRoleService(c, db)
	})
}
