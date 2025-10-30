package mongo

import (
	service "go-fiber/app/service/mongo"
	middleware "go-fiber/middleware/mongo"

	"github.com/gofiber/fiber/v2"
	goMongo "go.mongodb.org/mongo-driver/mongo"
)

func FileRoutes(app *fiber.App, db *goMongo.Database) {
	api := app.Group("/go-fiber-mongo")

	files := api.Group("/users/:id/upload", middleware.AuthRequired(), middleware.UserSelfOrAdmin())
	files.Post("/photo", func(c *fiber.Ctx) error { return service.UploadPhotoService(c, db) })
	files.Post("/certificate", func(c *fiber.Ctx) error { return service.UploadCertificateService(c, db) })
}


