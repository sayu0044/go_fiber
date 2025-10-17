package route

import (
    "github.com/gofiber/fiber/v2"
    "database/sql"
    "go-fiber/app/service"
    "go-fiber/middleware"
)

func RegisterRoutes(app *fiber.App, db *sql.DB) {
	// API Group
	api := app.Group("/go-fiber")

	// Public Routes (tidak perlu login)
	api.Post("/login", func(c *fiber.Ctx) error {
		return service.LoginService(c, db)
	})

	// Protected Routes (perlu login)
	protected := api.Group("", middleware.AuthRequired())
	
	// Profile route
	protected.Get("/profile", func(c *fiber.Ctx) error {
		return service.GetProfileService(c, db)
	})

    // Alumni Routes with Access Control
    alumni := protected.Group("/alumni", middleware.UserAndAdmin())
	
	// GET routes (Admin and User access)
	alumni.Get("/", func(c *fiber.Ctx) error {
		return service.GetAllAlumniService(c, db)
	}) 
	
	// Test endpoint untuk debugging
	alumni.Get("/test", func(c *fiber.Ctx) error {
		search := c.Query("search", "")
		return c.JSON(fiber.Map{
			"search": search,
			"search_length": len(search),
			"is_empty": search == "",
		})
	})
	alumni.Get("/employment-status", func(c *fiber.Ctx) error {
		return service.GetAlumniEmploymentStatusService(c, db)
	}) 
	alumni.Get("/:id", func(c *fiber.Ctx) error {
		return service.GetAlumniByIDService(c, db)
	}) 

	// POST, PUT, DELETE routes (Admin-only access)
	alumni.Post("/", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.CreateAlumniService(c, db)
	}) 
	alumni.Put("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.UpdateAlumniService(c, db)
	}) 
	alumni.Delete("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.DeleteAlumniService(c, db)
	}) 

    // Pekerjaan Alumni Routes with Access Control
    pekerjaan := protected.Group("/pekerjaan", middleware.UserAndAdmin())
	
	// GET routes (Admin and User access)
	pekerjaan.Get("/", func(c *fiber.Ctx) error {
		return service.GetAllPekerjaanService(c, db)
	}) 
	
	// Trash management routes - HARUS SEBELUM /:id
	pekerjaan.Get("/trash", func(c *fiber.Ctx) error {
		return service.ListDeletedPekerjaanService(c, db)
	}) 
	
	pekerjaan.Get("/:id", func(c *fiber.Ctx) error {
		return service.GetPekerjaanByIDService(c, db)
	}) 

	// GET route (Admin-only access)
	pekerjaan.Get("/alumni/:alumni_id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.GetPekerjaanByAlumniIDService(c, db)
	}) 

	// POST, PUT, DELETE routes (Admin-only access)
	pekerjaan.Post("/", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.CreatePekerjaanService(c, db)
	}) 
	pekerjaan.Put("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.UpdatePekerjaanService(c, db)
	}) 
	pekerjaan.Delete("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.DeletePekerjaanService(c, db)
	}) 

	// Soft delete route (User and Admin access with role-based restrictions)
	pekerjaan.Put("/soft-delete/:id", func(c *fiber.Ctx) error {
		return service.SoftDeletePekerjaanService(c, db)
	}) 

	// Restore and Hard Delete routes
	pekerjaan.Put("/restore/:id", func(c *fiber.Ctx) error {
		return service.RestorePekerjaanService(c, db)
	}) 
	
	pekerjaan.Delete("/hard-delete/:id", func(c *fiber.Ctx) error {
		return service.HardDeletePekerjaanService(c, db)
	}) 

	roles := protected.Group("/roles", middleware.AdminOnly())
	roles.Get("/", func(c *fiber.Ctx) error { return service.ListRolesService(c, db) })
	roles.Get("/:id", func(c *fiber.Ctx) error { return service.GetRoleByIDService(c, db) })
	roles.Post("/", func(c *fiber.Ctx) error { return service.CreateRoleService(c, db) })
	roles.Put("/:id", func(c *fiber.Ctx) error { return service.UpdateRoleService(c, db) })
	roles.Delete("/:id", func(c *fiber.Ctx) error { return service.DeleteRoleService(c, db) })

	


}
