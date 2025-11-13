package mongo

import (
	model "go-fiber/app/model/mongo"
	service "go-fiber/app/service/mongo"
	middleware "go-fiber/middleware/mongo"

	"github.com/gofiber/fiber/v2"
	goMongo "go.mongodb.org/mongo-driver/mongo"
)

// swagger:ignore
var _ model.FileUploadResponse

func FileRoutes(app *fiber.App, db *goMongo.Database) {
	api := app.Group("/go-fiber-mongo")

	files := api.Group("/users/:id/upload", middleware.AuthRequired(), middleware.UserSelfOrAdmin())
	files.Post("/photo", uploadPhotoHandler(db))
	files.Post("/certificate", uploadCertificateHandler(db))
}

// @Summary Upload foto profil
// @Description Upload file foto jpeg/jpg/png maksimal 1MB untuk user tertentu
// @Tags Files (Mongo)
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID User"
// @Param file formData file true "File foto (jpeg/jpg/png, max 1MB)"
// @Success 201 {object} model.FileUploadResponse
// @Failure 400 {object} fiber.Map
// @Failure 403 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /users/{id}/upload/photo [post]
func uploadPhotoHandler(db *goMongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return service.UploadPhotoService(c, db)
	}
}

// @Summary Upload sertifikat/ijazah
// @Description Upload file pdf maksimal 2MB untuk user tertentu
// @Tags Files (Mongo)
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID User"
// @Param file formData file true "File sertifikat (pdf, max 2MB)"
// @Success 201 {object} model.FileUploadResponse
// @Failure 400 {object} fiber.Map
// @Failure 403 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /users/{id}/upload/certificate [post]
func uploadCertificateHandler(db *goMongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return service.UploadCertificateService(c, db)
	}
}
