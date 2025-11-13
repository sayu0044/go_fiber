package mongo

import (
	model "go-fiber/app/model/mongo"
	service "go-fiber/app/service/mongo"
	middleware "go-fiber/middleware/mongo"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

// swagger:ignore
var (
	_ model.LoginRequest
	_ model.LoginResponse
	_ model.GetProfileResponse
	_ model.GetAllAlumniResponse
	_ model.GetAlumniByIDResponse
	_ model.CreateAlumniRequest
	_ model.CreateAlumniResponse
	_ model.UpdateAlumniRequest
	_ model.UpdateAlumniResponse
	_ model.DeleteAlumniResponse
	_ model.CheckAlumniResponse
	_ model.ListRolesResponse
	_ model.GetRoleByIDResponse
	_ model.CreateRoleRequest
	_ model.CreateRoleResponse
	_ model.UpdateRoleRequest
	_ model.UpdateRoleResponse
	_ model.DeleteRoleResponse
)

func AlumniRoutes(app *fiber.App, db *mongo.Database) {
	api := app.Group("/go-fiber-mongo")

	api.Post("/login", loginHandler(db))

	protected := api.Group("", middleware.AuthRequired())
	protected.Get("/profile", profileHandler(db))

	alumni := protected.Group("/alumni")
	alumni.Get("/", middleware.UserAndAdmin(), getAllAlumniHandler(db))
	alumni.Get("/check", middleware.UserAndAdmin(), checkAlumniHandler(db))
	alumni.Get("/:id", middleware.UserAndAdmin(), getAlumniByIDHandler(db))
	alumni.Post("/", middleware.AdminOnly(), createAlumniHandler(db))
	alumni.Put("/:id", middleware.AdminOnly(), updateAlumniHandler(db))
	alumni.Delete("/:id", middleware.AdminOnly(), deleteAlumniHandler(db))

	roles := protected.Group("/roles")
	roles.Get("/", middleware.UserAndAdmin(), listRolesHandler(db))
	roles.Get("/:id", middleware.UserAndAdmin(), getRoleByIDHandler(db))
	roles.Post("/", middleware.AdminOnly(), createRoleHandler(db))
	roles.Put("/:id", middleware.AdminOnly(), updateRoleHandler(db))
	roles.Delete("/:id", middleware.AdminOnly(), deleteRoleHandler(db))
}

// @Summary Login user (Mongo)
// @Description Mengembalikan token JWT untuk mengakses API Mongo
// @Tags Auth (Mongo)
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Login request"
// @Success 200 {object} model.LoginResponse
// @Failure 400 {object} fiber.Map
// @Failure 401 {object} fiber.Map
// @Router /login [post]
func loginHandler(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return service.LoginService(c, db)
	}
}

// @Summary Profil pengguna saat ini
// @Description Mengambil profil user berdasarkan token
// @Tags Profile (Mongo)
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.GetProfileResponse
// @Failure 401 {object} fiber.Map
// @Router /profile [get]
func profileHandler(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return service.GetProfileService(c, db)
	}
}

// @Summary Daftar alumni
// @Description Mengambil daftar alumni dengan pagination, sorting, dan pencarian
// @Tags Alumni (Mongo)
// @Produce json
// @Security BearerAuth
// @Param page query int false "Halaman"
// @Param limit query int false "Jumlah per halaman"
// @Param sortBy query string false "Kolom sortir"
// @Param order query string false "Urutan asc/desc"
// @Param search query string false "Kata kunci pencarian"
// @Success 200 {object} model.GetAllAlumniResponse
// @Failure 401 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /alumni [get]
func getAllAlumniHandler(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return service.GetAllAlumniService(c, db)
	}
}

// @Summary Detail alumni
// @Description Mengambil data alumni berdasarkan ID
// @Tags Alumni (Mongo)
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID Alumni"
// @Success 200 {object} model.GetAlumniByIDResponse
// @Failure 400 {object} model.GetAlumniByIDResponse
// @Failure 404 {object} model.GetAlumniByIDResponse
// @Router /alumni/{id} [get]
func getAlumniByIDHandler(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return service.GetAlumniByIDService(c, db)
	}
}

// @Summary Tambah alumni
// @Description Membuat data alumni baru
// @Tags Alumni (Mongo)
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.CreateAlumniRequest true "Permintaan pembuatan alumni"
// @Success 201 {object} model.CreateAlumniResponse
// @Failure 400 {object} model.CreateAlumniResponse
// @Failure 500 {object} model.CreateAlumniResponse
// @Router /alumni [post]
func createAlumniHandler(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return service.CreateAlumniService(c, db)
	}
}

// @Summary Perbarui alumni
// @Description Memperbarui data alumni berdasarkan ID
// @Tags Alumni (Mongo)
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID Alumni"
// @Param request body model.UpdateAlumniRequest true "Permintaan pembaruan alumni"
// @Success 200 {object} model.UpdateAlumniResponse
// @Failure 400 {object} model.UpdateAlumniResponse
// @Failure 404 {object} model.UpdateAlumniResponse
// @Router /alumni/{id} [put]
func updateAlumniHandler(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return service.UpdateAlumniService(c, db)
	}
}

// @Summary Hapus alumni
// @Description Menghapus data alumni berdasarkan ID
// @Tags Alumni (Mongo)
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID Alumni"
// @Success 200 {object} model.DeleteAlumniResponse
// @Failure 400 {object} model.DeleteAlumniResponse
// @Failure 404 {object} model.DeleteAlumniResponse
// @Router /alumni/{id} [delete]
func deleteAlumniHandler(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return service.DeleteAlumniService(c, db)
	}
}

// @Summary Cek alumni berdasarkan NIM
// @Description Mengecek status alumni menggunakan API key legacy
// @Tags Alumni (Mongo)
// @Produce json
// @Security BearerAuth
// @Param key query string true "API key"
// @Param nim query string true "NIM Mahasiswa"
// @Success 200 {object} model.CheckAlumniResponse
// @Failure 400 {object} fiber.Map
// @Failure 401 {object} fiber.Map
// @Router /alumni/check [get]
func checkAlumniHandler(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return service.CheckAlumniService(c, db)
	}
}

// @Summary Daftar role
// @Description Mengambil seluruh role yang tersedia
// @Tags Roles (Mongo)
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.ListRolesResponse
// @Failure 500 {object} fiber.Map
// @Router /roles [get]
func listRolesHandler(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return service.ListRolesService(c, db)
	}
}

// @Summary Detail role
// @Description Mengambil data role berdasarkan ID
// @Tags Roles (Mongo)
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID Role"
// @Success 200 {object} model.GetRoleByIDResponse
// @Failure 404 {object} model.GetRoleByIDResponse
// @Router /roles/{id} [get]
func getRoleByIDHandler(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return service.GetRoleByIDService(c, db)
	}
}

// @Summary Tambah role
// @Description Membuat role baru
// @Tags Roles (Mongo)
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.CreateRoleRequest true "Permintaan pembuatan role"
// @Success 201 {object} model.CreateRoleResponse
// @Failure 400 {object} model.CreateRoleResponse
// @Router /roles [post]
func createRoleHandler(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return service.CreateRoleService(c, db)
	}
}

// @Summary Perbarui role
// @Description Memperbarui data role berdasarkan ID
// @Tags Roles (Mongo)
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID Role"
// @Param request body model.UpdateRoleRequest true "Permintaan pembaruan role"
// @Success 200 {object} model.UpdateRoleResponse
// @Failure 400 {object} model.UpdateRoleResponse
// @Failure 404 {object} model.UpdateRoleResponse
// @Router /roles/{id} [put]
func updateRoleHandler(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return service.UpdateRoleService(c, db)
	}
}

// @Summary Hapus role
// @Description Menghapus role berdasarkan ID
// @Tags Roles (Mongo)
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID Role"
// @Success 200 {object} model.DeleteRoleResponse
// @Failure 404 {object} model.DeleteRoleResponse
// @Router /roles/{id} [delete]
func deleteRoleHandler(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return service.DeleteRoleService(c, db)
	}
}
