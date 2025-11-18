package mongo_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"os"

	service "go-fiber/app/service/mongo"

	"github.com/gofiber/fiber/v2"
)

func TestGetAlumniByIDService_EmptyID(t *testing.T) {
	app := fiber.New()
	app.Get("/alumni/:id?", func(c *fiber.Ctx) error { return service.GetAlumniByIDService(c, nil) })

	req := httptest.NewRequest(http.MethodGet, "/alumni/", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestGetAlumniByIDService_BadFormat(t *testing.T) {
	app := fiber.New()
	app.Get("/alumni/:id", func(c *fiber.Ctx) error { return service.GetAlumniByIDService(c, nil) })

	req := httptest.NewRequest(http.MethodGet, "/alumni/not-a-hex-id", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestCreateAlumniService_MissingRequired(t *testing.T) {
	app := fiber.New()
	app.Post("/alumni", func(c *fiber.Ctx) error { return service.CreateAlumniService(c, nil) })

	body := []byte(`{"nama": "John"}`)
	req := httptest.NewRequest(http.MethodPost, "/alumni", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestUpdateAlumniService_InvalidID(t *testing.T) {
	app := fiber.New()
	app.Put("/alumni/:id", func(c *fiber.Ctx) error { return service.UpdateAlumniService(c, nil) })

	req := httptest.NewRequest(http.MethodPut, "/alumni/bad-id", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestDeleteAlumniService_InvalidID(t *testing.T) {
	app := fiber.New()
	app.Delete("/alumni/:id", func(c *fiber.Ctx) error { return service.DeleteAlumniService(c, nil) })

	req := httptest.NewRequest(http.MethodDelete, "/alumni/bad-id", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestCheckAlumniService_InvalidKeyOrNIM(t *testing.T) {
	// ensure deterministic key check
	_ = os.Setenv("API_KEY", "test-key")

	app := fiber.New()
	app.Get("/alumni/check", func(c *fiber.Ctx) error { return service.CheckAlumniService(c, nil) })

	// missing key
	req1 := httptest.NewRequest(http.MethodGet, "/alumni/check", nil)
	resp1, _ := app.Test(req1)
	if resp1.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp1.StatusCode)
	}

	// with valid key but missing nim -> bad request
	req2 := httptest.NewRequest(http.MethodGet, "/alumni/check?key=test-key", nil)
	resp2, _ := app.Test(req2)
	if resp2.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp2.StatusCode)
	}
}


