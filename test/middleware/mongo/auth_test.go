package mongo_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	model "go-fiber/app/model/mongo"
	mw "go-fiber/middleware/mongo"
	utils "go-fiber/utils/mongo"

	"github.com/gofiber/fiber/v2"
)

func setupAuthApp() *fiber.App {
	app := fiber.New()
	app.Get("/", mw.AuthRequired(), func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"ok": true})
	})
	return app
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	app := setupAuthApp()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	app := setupAuthApp()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "NotBearer abc")
	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	app := setupAuthApp()

	token, err := utils.GenerateToken(model.User{Username: "john", Role: "admin"})
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}


