package mongo_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	service "go-fiber/app/service/mongo"

	"github.com/gofiber/fiber/v2"
)

func TestLoginService_InvalidJSON(t *testing.T) {
	app := fiber.New()
	app.Post("/login", func(c *fiber.Ctx) error { return service.LoginService(c, nil) })

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("{invalid-json"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestLoginService_MissingFields(t *testing.T) {
	app := fiber.New()
	app.Post("/login", func(c *fiber.Ctx) error { return service.LoginService(c, nil) })

	body := []byte(`{"email":"","password":""}`)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}


