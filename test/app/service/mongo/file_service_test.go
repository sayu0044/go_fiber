package mongo_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	service "go-fiber/app/service/mongo"

	"github.com/gofiber/fiber/v2"
)

func TestUploadPhotoService_MissingUserID(t *testing.T) {
	// Test with optional parameter that can be empty
	app := fiber.New()
	app.Post("/upload/:id?", func(c *fiber.Ctx) error {
		// This will get empty string if id is not provided
		return service.UploadPhotoService(c, nil)
	})

	req := httptest.NewRequest(http.MethodPost, "/upload/", nil)
	req.Header.Set("Content-Type", "multipart/form-data")
	resp, _ := app.Test(req)

	// Should fail because user ID is required (empty string)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestUploadPhotoService_NoFile(t *testing.T) {
	app := fiber.New()
	app.Post("/upload/:id", func(c *fiber.Ctx) error {
		return service.UploadPhotoService(c, nil)
	})

	req := httptest.NewRequest(http.MethodPost, "/upload/507f1f77bcf86cd799439011", nil)
	req.Header.Set("Content-Type", "multipart/form-data")
	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

// TestUploadPhotoService_InvalidUserID is skipped because it requires database access
// before ObjectID validation. The invalid user ID validation happens after file
// is saved and repository is created, which requires a valid database connection.
// This scenario is better tested with integration tests that have a real database.

func TestUploadCertificateService_MissingUserID(t *testing.T) {
	// Test with empty string as ID using custom handler
	app := fiber.New()
	app.Post("/upload/*", func(c *fiber.Ctx) error {
		c.Params("id", "")
		return service.UploadCertificateService(c, nil)
	})

	req := httptest.NewRequest(http.MethodPost, "/upload/", nil)
	req.Header.Set("Content-Type", "multipart/form-data")
	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestUploadCertificateService_NoFile(t *testing.T) {
	app := fiber.New()
	app.Post("/upload/:id", func(c *fiber.Ctx) error {
		return service.UploadCertificateService(c, nil)
	})

	req := httptest.NewRequest(http.MethodPost, "/upload/507f1f77bcf86cd799439011", nil)
	req.Header.Set("Content-Type", "multipart/form-data")
	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

// TestUploadCertificateService_InvalidUserID is skipped because it requires database access
// before ObjectID validation. The invalid user ID validation happens after file
// is saved and repository is created, which requires a valid database connection.
// This scenario is better tested with integration tests that have a real database.

// Test helper functions by accessing them through reflection or by testing indirectly
// Since isAllowed, extensionForType, and sniffContentType are private,
// we test them indirectly through the public functions

func TestUploadPhotoService_FileTooLarge(t *testing.T) {
	app := fiber.New()
	app.Post("/upload/:id", func(c *fiber.Ctx) error {
		// Use nil database - will fail at file size validation before DB access
		return service.UploadPhotoService(c, nil)
	})

	// Create a file larger than 1MB (max size for photo)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	part, _ := writer.CreateFormFile("file", "large.jpg")
	// Write 2MB of data (larger than 1MB limit)
	largeData := make([]byte, 2*1024*1024)
	part.Write(largeData)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload/507f1f77bcf86cd799439011", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, _ := app.Test(req)

	// Should fail at file size validation (2MB > 1MB limit)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestUploadCertificateService_FileTooLarge(t *testing.T) {
	app := fiber.New()
	app.Post("/upload/:id", func(c *fiber.Ctx) error {
		// Use nil database - will fail at file size validation before DB access
		return service.UploadCertificateService(c, nil)
	})

	// Create a file larger than 2MB (max size for certificate)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	part, _ := writer.CreateFormFile("file", "large.pdf")
	// Write 3MB of data (larger than 2MB limit)
	largeData := make([]byte, 3*1024*1024)
	part.Write(largeData)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload/507f1f77bcf86cd799439011", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, _ := app.Test(req)

	// Should fail at file size validation (3MB > 2MB limit)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestUploadPhotoService_InvalidFileType(t *testing.T) {
	app := fiber.New()
	app.Post("/upload/:id", func(c *fiber.Ctx) error {
		// Use nil database - will fail at file type validation before DB access
		return service.UploadPhotoService(c, nil)
	})

	// Create a PDF file (not allowed for photo upload)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	part, _ := writer.CreateFormFile("file", "test.pdf")
	part.Write([]byte("%PDF-1.4"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload/507f1f77bcf86cd799439011", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, _ := app.Test(req)

	// Should fail at file type validation (PDF not allowed for photo)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestUploadCertificateService_InvalidFileType(t *testing.T) {
	app := fiber.New()
	app.Post("/upload/:id", func(c *fiber.Ctx) error {
		// Use nil database - will fail at file type validation before DB access
		return service.UploadCertificateService(c, nil)
	})

	// Create a JPEG file (not allowed for certificate upload)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	part, _ := writer.CreateFormFile("file", "test.jpg")
	part.Write([]byte{0xFF, 0xD8, 0xFF, 0xE0}) // JPEG magic bytes
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload/507f1f77bcf86cd799439011", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, _ := app.Test(req)

	// Should fail at file type validation (JPEG not allowed for certificate)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

// TestUploadService_CreatesUploadDirectory is skipped because it requires database access.
// Directory creation happens after file validation and before database operations.
// This scenario is better tested with integration tests that have a real database.

