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
	_ model.GetAllPekerjaanResponse
	_ model.GetPekerjaanAlumniByIDResponse
	_ model.GetPekerjaanAlumniByAlumniIDResponse
	_ model.CreatePekerjaanAlumniRequest
	_ model.CreatePekerjaanAlumniResponse
	_ model.UpdatePekerjaanAlumniRequest
	_ model.UpdatePekerjaanAlumniResponse
)

func PekerjaanRoutes(app *fiber.App, db *mongo.Database) {
	api := app.Group("/go-fiber-mongo")
	protected := api.Group("", middleware.AuthRequired())

	pekerjaan := protected.Group("/pekerjaan")
	pekerjaan.Get("/", middleware.UserAndAdmin(), getAllPekerjaanHandler(db))
	pekerjaan.Get("/:id", middleware.UserAndAdmin(), getPekerjaanByIDHandler(db))
	pekerjaan.Get("/alumni/:alumni_id", middleware.AdminOnly(), getPekerjaanByAlumniIDHandler(db))
	pekerjaan.Post("/", middleware.AdminOnly(), createPekerjaanHandler(db))
	pekerjaan.Put("/:id", middleware.AdminOnly(), updatePekerjaanHandler(db))
	pekerjaan.Delete("/:id", middleware.AdminOnly(), deletePekerjaanHandler(db))
}

// @Summary Daftar pekerjaan alumni
// @Description Mengambil daftar pekerjaan alumni dengan pagination, sorting, dan pencarian
// @Tags Pekerjaan (Mongo)
// @Produce json
// @Security BearerAuth
// @Param page query int false "Halaman"
// @Param limit query int false "Jumlah per halaman"
// @Param sortBy query string false "Kolom sortir (id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, tanggal_mulai_kerja, status_pekerjaan, created_at)"
// @Param order query string false "Urutan asc/desc"
// @Param search query string false "Kata kunci pencarian"
// @Success 200 {object} model.GetAllPekerjaanResponse
// @Failure 500 {object} fiber.Map
// @Router /pekerjaan [get]
func getAllPekerjaanHandler(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return service.GetAllPekerjaanService(c, db)
	}
}

// @Summary Detail pekerjaan alumni
// @Description Mengambil detail pekerjaan alumni berdasarkan ID pekerjaan
// @Tags Pekerjaan (Mongo)
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID Pekerjaan"
// @Success 200 {object} model.GetPekerjaanAlumniByIDResponse
// @Failure 400 {object} model.GetPekerjaanAlumniByIDResponse
// @Failure 404 {object} model.GetPekerjaanAlumniByIDResponse
// @Router /pekerjaan/{id} [get]
func getPekerjaanByIDHandler(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return service.GetPekerjaanByIDService(c, db)
	}
}

// @Summary Daftar pekerjaan per alumni
// @Description Mengambil seluruh pekerjaan milik alumni tertentu
// @Tags Pekerjaan (Mongo)
// @Produce json
// @Security BearerAuth
// @Param alumni_id path string true "ID Alumni"
// @Success 200 {object} model.GetPekerjaanAlumniByAlumniIDResponse
// @Failure 400 {object} model.GetPekerjaanAlumniByAlumniIDResponse
// @Failure 500 {object} model.GetPekerjaanAlumniByAlumniIDResponse
// @Router /pekerjaan/alumni/{alumni_id} [get]
func getPekerjaanByAlumniIDHandler(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return service.GetPekerjaanByAlumniIDService(c, db)
	}
}

// @Summary Tambah pekerjaan alumni
// @Description Membuat data pekerjaan baru untuk alumni
// @Tags Pekerjaan (Mongo)
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.CreatePekerjaanAlumniRequest true "Data pekerjaan alumni"
// @Success 201 {object} model.CreatePekerjaanAlumniResponse
// @Failure 400 {object} model.CreatePekerjaanAlumniResponse
// @Failure 500 {object} model.CreatePekerjaanAlumniResponse
// @Router /pekerjaan [post]
func createPekerjaanHandler(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return service.CreatePekerjaanService(c, db)
	}
}

// @Summary Perbarui pekerjaan alumni
// @Description Memperbarui data pekerjaan alumni berdasarkan ID pekerjaan
// @Tags Pekerjaan (Mongo)
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID Pekerjaan"
// @Param request body model.UpdatePekerjaanAlumniRequest true "Data pekerjaan alumni yang diperbarui"
// @Success 200 {object} model.UpdatePekerjaanAlumniResponse
// @Failure 400 {object} model.UpdatePekerjaanAlumniResponse
// @Failure 500 {object} model.UpdatePekerjaanAlumniResponse
// @Router /pekerjaan/{id} [put]
func updatePekerjaanHandler(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return service.UpdatePekerjaanService(c, db)
	}
}

// @Summary Hapus pekerjaan alumni
// @Description Menghapus data pekerjaan alumni berdasarkan ID pekerjaan
// @Tags Pekerjaan (Mongo)
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID Pekerjaan"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /pekerjaan/{id} [delete]
func deletePekerjaanHandler(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return service.DeletePekerjaanService(c, db)
	}
}
