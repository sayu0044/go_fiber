package mongo

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	model "go-fiber/app/model/mongo"
	repository "go-fiber/app/repository/mongo"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	goMongo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	categoryPhoto       = "photo"
	categoryCertificate = "certificate"
)

func UploadPhotoService(c *fiber.Ctx, db *goMongo.Database) error {
	return handleUpload(c, db, categoryPhoto, 1*1024*1024, []string{"image/jpeg", "image/png"})
}

func UploadCertificateService(c *fiber.Ctx, db *goMongo.Database) error {
	return handleUpload(c, db, categoryCertificate, 2*1024*1024, []string{"application/pdf"})
}

func handleUpload(c *fiber.Ctx, db *goMongo.Database, category string, maxSize int64, allowedTypes []string) error {
	userIDParam := c.Params("id")
	if userIDParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "User id is required"})
	}

    fileHeader, err := c.FormFile("file")
    if err != nil || fileHeader == nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "No file uploaded. Use form-data with key 'file' (type File)",
        })
    }

	if fileHeader.Size > maxSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "File too large"})
	}

	contentType, err := sniffContentType(fileHeader)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	if !isAllowed(contentType, allowedTypes) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "File type not allowed"})
	}

	// Build destination path
	userDir := filepath.Join("uploads", userIDParam, category)
	if err := os.MkdirAll(userDir, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to create upload directory"})
	}

	// Create destination filename
	ext := extensionForType(contentType, fileHeader.Filename)
	newName := uuid.New().String() + ext
	destPath := filepath.Join(userDir, newName)

	// Save file
	if err := saveUploadedFile(fileHeader, destPath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	// Persist metadata (upsert behavior for photo: keep one latest)
	repo := repository.NewFileRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	alumniOID, err := primitive.ObjectIDFromHex(userIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid user id"})
	}

	// If category is photo, remove previous file if any (single latest policy)
	if category == categoryPhoto {
		if existing, err := repo.FindByAlumniAndCategory(ctx, alumniOID, category); err == nil && existing != nil {
			_ = os.Remove(existing.FilePath)
			_ = repo.DeleteByID(ctx, existing.ID)
		}
	}

	record := &model.File{
		AlumniID:     alumniOID,
		Category:     category,
		FileName:     newName,
		OriginalName: fileHeader.Filename,
		FilePath:     destPath,
		FileType:     contentType,
		FileSize:     fileHeader.Size,
	}
	if err := repo.Create(ctx, record); err != nil {
		_ = os.Remove(destPath)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to save metadata"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "File uploaded successfully",
		"data": model.FileResponse{
			ID:           record.ID.Hex(),
			Category:     record.Category,
			FileName:     record.FileName,
			OriginalName: record.OriginalName,
			FilePath:     "/" + filepath.ToSlash(record.FilePath),
			FileType:     record.FileType,
			FileSize:     record.FileSize,
		},
	})
}

func isAllowed(ct string, allowed []string) bool {
	for _, a := range allowed {
		if strings.EqualFold(ct, a) {
			return true
		}
	}
	return false
}

func sniffContentType(hdr *multipart.FileHeader) (string, error) {
	f, err := hdr.Open()
	if err != nil { return "", err }
	defer f.Close()

	buf := make([]byte, 512)
	n, err := f.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	return http.DetectContentType(buf[:n]), nil
}

func extensionForType(ct string, fallbackFromName string) string {
	switch ct {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "application/pdf":
		return ".pdf"
	default:
		return filepath.Ext(fallbackFromName)
	}
}

func saveUploadedFile(hdr *multipart.FileHeader, destPath string) error {
	src, err := hdr.Open()
	if err != nil { return err }
	defer src.Close()

	out, err := os.Create(destPath)
	if err != nil { return err }
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}


