package mongo_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	service "go-fiber/app/service/mongo"

	"github.com/gofiber/fiber/v2"
)

func TestGetPekerjaanByID_EmptyID(t *testing.T) {
	app := fiber.New()
	app.Get("/pekerjaan/:id?", func(c *fiber.Ctx) error { return service.GetPekerjaanByIDService(c, nil) })
	req := httptest.NewRequest(http.MethodGet, "/pekerjaan/", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestGetPekerjaanByAlumniID_EmptyID(t *testing.T) {
	app := fiber.New()
	app.Get("/pekerjaan/alumni/:alumni_id?", func(c *fiber.Ctx) error { return service.GetPekerjaanByAlumniIDService(c, nil) })
	req := httptest.NewRequest(http.MethodGet, "/pekerjaan/alumni/", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestCreatePekerjaan_MissingFields(t *testing.T) {
	app := fiber.New()
	app.Post("/pekerjaan", func(c *fiber.Ctx) error { return service.CreatePekerjaanService(c, nil) })
	// minimal invalid body (missing required fields)
	body := []byte(`{"alumni_id":"","nama_perusahaan":""}`)
	req := httptest.NewRequest(http.MethodPost, "/pekerjaan", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestCreatePekerjaan_InvalidStatus(t *testing.T) {
	app := fiber.New()
	app.Post("/pekerjaan", func(c *fiber.Ctx) error { return service.CreatePekerjaanService(c, nil) })
	body := []byte(`{
		"alumni_id":"64b1f0c2c2c2c2c2c2c2c2c2",
		"nama_perusahaan":"Acme",
		"posisi_jabatan":"Dev",
		"bidang_industri":"IT",
		"lokasi_kerja":"Bandung",
		"tanggal_mulai_kerja":"2024-01-01",
		"status_pekerjaan":"unknown"
	}`)
	req := httptest.NewRequest(http.MethodPost, "/pekerjaan", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}


