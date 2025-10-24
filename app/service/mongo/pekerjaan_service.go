package mongo

import (
	"fmt"
	"go-fiber/app/model/mongo"
	repository "go-fiber/app/repository/mongo"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDB "go.mongodb.org/mongo-driver/mongo"
)

// Pekerjaan Alumni Services

func GetAllPekerjaanService(c *fiber.Ctx, db *mongoDB.Database) error {
	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "id")
	order := c.Query("order", "asc")
	search := c.Query("search", "")

	// Hitung offset untuk pagination
	offset := (page - 1) * limit

	// Validasi input
	sortByWhitelist := map[string]bool{"id": true, "nama_perusahaan": true, "posisi_jabatan": true, "bidang_industri": true, "lokasi_kerja": true, "tanggal_mulai_kerja": true, "status_pekerjaan": true, "created_at": true}
	if !sortByWhitelist[sortBy] {
		sortBy = "id"
	}
	if strings.ToLower(order) != "desc" {
		order = "asc"
	}

	// Ambil data dari repository
	pekerjaan, err := repository.GetPekerjaanRepo(db, search, sortBy, order, limit, offset)
	if err != nil {
		return c.Status(500).JSON(mongo.GetAllPekerjaanResponse{
			Success: false,
			Message: "Failed to fetch pekerjaan",
			Data: mongo.PekerjaanData{
				Items: []mongo.PekerjaanAlumni{},
				Meta:  mongo.MetaInfo{},
			},
		})
	}

	total, err := repository.CountPekerjaanRepo(db, search)
	if err != nil {
		return c.Status(500).JSON(mongo.GetAllPekerjaanResponse{
			Success: false,
			Message: "Failed to count pekerjaan",
			Data: mongo.PekerjaanData{
				Items: []mongo.PekerjaanAlumni{},
				Meta:  mongo.MetaInfo{},
			},
		})
	}

	// Buat response pakai model
	response := mongo.GetAllPekerjaanResponse{
		Success: true,
		Message: "Berhasil mengambil data pekerjaan",
		Data: mongo.PekerjaanData{
			Items: pekerjaan,
			Meta: mongo.MetaInfo{
				Page:   page,
				Limit:  limit,
				Total:  total,
				Pages:  (total + limit - 1) / limit,
				SortBy: sortBy,
				Order:  order,
				Search: search,
			},
		},
	}

	return c.JSON(response)
}

func GetPekerjaanByIDService(c *fiber.Ctx, db *mongoDB.Database) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(mongo.GetPekerjaanAlumniByIDResponse{
			Success: false,
			Message: "ID tidak valid",
			Data:    mongo.PekerjaanAlumni{},
		})
	}

	pekerjaan, err := repository.GetPekerjaanByID(db, idStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(mongo.GetPekerjaanAlumniByIDResponse{
			Success: false,
			Message: "Gagal mengambil data pekerjaan: " + err.Error(),
			Data:    mongo.PekerjaanAlumni{},
		})
	}

	if pekerjaan == nil {
		return c.Status(fiber.StatusNotFound).JSON(mongo.GetPekerjaanAlumniByIDResponse{
			Success: false,
			Message: "Pekerjaan tidak ditemukan",
			Data:    mongo.PekerjaanAlumni{},
		})
	}

	return c.Status(fiber.StatusOK).JSON(mongo.GetPekerjaanAlumniByIDResponse{
		Success: true,
		Message: "Berhasil mengambil data pekerjaan",
		Data:    *pekerjaan,
	})
}

func GetPekerjaanByAlumniIDService(c *fiber.Ctx, db *mongoDB.Database) error {
	alumniIDStr := c.Params("alumni_id")
	if alumniIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(mongo.GetPekerjaanAlumniByAlumniIDResponse{
			Success: false,
			Message: "ID alumni tidak valid",
			Data:    []mongo.PekerjaanAlumni{},
		})
	}

	// Check if alumni exists
	_, err := repository.GetAlumniByID(db, alumniIDStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(mongo.GetPekerjaanAlumniByAlumniIDResponse{
			Success: false,
			Message: "Gagal mengambil data alumni: " + err.Error(),
			Data:    []mongo.PekerjaanAlumni{},
		})
	}

	pekerjaan, err := repository.GetPekerjaanByAlumniID(db, alumniIDStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(mongo.GetPekerjaanAlumniByAlumniIDResponse{
			Success: false,
			Message: "Gagal mengambil data pekerjaan: " + err.Error(),
			Data:    []mongo.PekerjaanAlumni{},
		})
	}

	return c.Status(fiber.StatusOK).JSON(mongo.GetPekerjaanAlumniByAlumniIDResponse{
		Success: true,
		Message: "Berhasil mengambil data pekerjaan alumni",
		Data:    pekerjaan,
	})
}

func CreatePekerjaanService(c *fiber.Ctx, db *mongoDB.Database) error {
	var req mongo.CreatePekerjaanAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(mongo.CreatePekerjaanAlumniResponse{
			Success: false,
			Message: "Format data tidak valid: " + err.Error(),
			Data:    mongo.PekerjaanAlumni{},
		})
	}

	// Debug logging untuk melihat data yang diterima
	fmt.Printf("Received request data: %+v\n", req)

	// Basic validation dengan pesan yang lebih detail
	if req.AlumniID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(mongo.CreatePekerjaanAlumniResponse{
			Success: false,
			Message: "Alumni ID wajib diisi",
			Data:    mongo.PekerjaanAlumni{},
		})
	}
	if req.NamaPerusahaan == "" {
		return c.Status(fiber.StatusBadRequest).JSON(mongo.CreatePekerjaanAlumniResponse{
			Success: false,
			Message: "Nama perusahaan wajib diisi",
			Data:    mongo.PekerjaanAlumni{},
		})
	}
	if req.PosisiJabatan == "" {
		return c.Status(fiber.StatusBadRequest).JSON(mongo.CreatePekerjaanAlumniResponse{
			Success: false,
			Message: "Posisi jabatan wajib diisi",
			Data:    mongo.PekerjaanAlumni{},
		})
	}
	if req.BidangIndustri == "" {
		return c.Status(fiber.StatusBadRequest).JSON(mongo.CreatePekerjaanAlumniResponse{
			Success: false,
			Message: "Bidang industri wajib diisi",
			Data:    mongo.PekerjaanAlumni{},
		})
	}
	if req.LokasiKerja == "" {
		return c.Status(fiber.StatusBadRequest).JSON(mongo.CreatePekerjaanAlumniResponse{
			Success: false,
			Message: "Lokasi kerja wajib diisi",
			Data:    mongo.PekerjaanAlumni{},
		})
	}
	if req.TanggalMulaiKerja == "" {
		return c.Status(fiber.StatusBadRequest).JSON(mongo.CreatePekerjaanAlumniResponse{
			Success: false,
			Message: "Tanggal mulai kerja wajib diisi",
			Data:    mongo.PekerjaanAlumni{},
		})
	}
	if req.StatusPekerjaan == "" {
		return c.Status(fiber.StatusBadRequest).JSON(mongo.CreatePekerjaanAlumniResponse{
			Success: false,
			Message: "Status pekerjaan wajib diisi",
			Data:    mongo.PekerjaanAlumni{},
		})
	}
	if req.StatusPekerjaan != "aktif" && req.StatusPekerjaan != "selesai" && req.StatusPekerjaan != "resigned" {
		return c.Status(fiber.StatusBadRequest).JSON(mongo.CreatePekerjaanAlumniResponse{
			Success: false,
			Message: "Status pekerjaan harus aktif, selesai, atau resigned",
			Data:    mongo.PekerjaanAlumni{},
		})
	}

	// Check if alumni exists
	_, err := repository.GetAlumniByID(db, req.AlumniID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(mongo.CreatePekerjaanAlumniResponse{
			Success: false,
			Message: "Gagal mengambil data alumni: " + err.Error(),
			Data:    mongo.PekerjaanAlumni{},
		})
	}

	// Parse tanggal mulai kerja
	tanggalMulai, err := time.Parse("2006-01-02", req.TanggalMulaiKerja)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(mongo.CreatePekerjaanAlumniResponse{
			Success: false,
			Message: "Format tanggal mulai kerja tidak valid. Gunakan format YYYY-MM-DD",
			Data:    mongo.PekerjaanAlumni{},
		})
	}

	// Parse tanggal selesai kerja jika ada
	var tanggalSelesai *time.Time
	if req.TanggalSelesaiKerja != nil && *req.TanggalSelesaiKerja != "" {
		parsed, err := time.Parse("2006-01-02", *req.TanggalSelesaiKerja)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(mongo.CreatePekerjaanAlumniResponse{
				Success: false,
				Message: "Format tanggal selesai kerja tidak valid. Gunakan format YYYY-MM-DD",
				Data:    mongo.PekerjaanAlumni{},
			})
		}
		tanggalSelesai = &parsed
	}

	// Convert AlumniID string to ObjectID
	alumniID, err := primitive.ObjectIDFromHex(req.AlumniID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(mongo.CreatePekerjaanAlumniResponse{
			Success: false,
			Message: "Alumni ID tidak valid",
			Data:    mongo.PekerjaanAlumni{},
		})
	}

	// Konversi ke repository request
	repoReq := &mongo.CreatePekerjaanAlumniRepositoryRequest{
		AlumniID:            alumniID,
		NamaPerusahaan:      req.NamaPerusahaan,
		PosisiJabatan:       req.PosisiJabatan,
		BidangIndustri:      req.BidangIndustri,
		LokasiKerja:         req.LokasiKerja,
		GajiRange:           req.GajiRange,
		TanggalMulaiKerja:   tanggalMulai,
		TanggalSelesaiKerja: tanggalSelesai,
		StatusPekerjaan:     req.StatusPekerjaan,
		DeskripsiPekerjaan:  req.DeskripsiPekerjaan,
	}

	pekerjaan, err := repository.CreatePekerjaan(db, repoReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(mongo.CreatePekerjaanAlumniResponse{
			Success: false,
			Message: "Gagal membuat pekerjaan: " + err.Error(),
			Data:    mongo.PekerjaanAlumni{},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(mongo.CreatePekerjaanAlumniResponse{
		Success: true,
		Message: "Berhasil membuat pekerjaan",
		Data:    *pekerjaan,
	})
}

func UpdatePekerjaanService(c *fiber.Ctx, db *mongoDB.Database) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(mongo.UpdatePekerjaanAlumniResponse{
			Success: false,
			Message: "ID tidak valid",
			Data:    mongo.PekerjaanAlumni{},
		})
	}

	var req mongo.UpdatePekerjaanAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(mongo.UpdatePekerjaanAlumniResponse{
			Success: false,
			Message: "Format data tidak valid: " + err.Error(),
			Data:    mongo.PekerjaanAlumni{},
		})
	}

	// Check if pekerjaan exists
	_, err := repository.GetPekerjaanByID(db, idStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(mongo.UpdatePekerjaanAlumniResponse{
			Success: false,
			Message: "Gagal mengambil data pekerjaan: " + err.Error(),
			Data:    mongo.PekerjaanAlumni{},
		})
	}

	// Validate status
	if req.StatusPekerjaan != "aktif" && req.StatusPekerjaan != "selesai" && req.StatusPekerjaan != "resigned" {
		return c.Status(fiber.StatusBadRequest).JSON(mongo.UpdatePekerjaanAlumniResponse{
			Success: false,
			Message: "Status pekerjaan harus aktif, selesai, atau resigned",
			Data:    mongo.PekerjaanAlumni{},
		})
	}

	// Parse tanggal mulai kerja
	tanggalMulai, err := time.Parse("2006-01-02", req.TanggalMulaiKerja)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(mongo.UpdatePekerjaanAlumniResponse{
			Success: false,
			Message: "Format tanggal mulai kerja tidak valid. Gunakan format YYYY-MM-DD",
			Data:    mongo.PekerjaanAlumni{},
		})
	}

	// Parse tanggal selesai kerja jika ada
	var tanggalSelesai *time.Time
	if req.TanggalSelesaiKerja != nil && *req.TanggalSelesaiKerja != "" {
		parsed, err := time.Parse("2006-01-02", *req.TanggalSelesaiKerja)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(mongo.UpdatePekerjaanAlumniResponse{
				Success: false,
				Message: "Format tanggal selesai kerja tidak valid. Gunakan format YYYY-MM-DD",
				Data:    mongo.PekerjaanAlumni{},
			})
		}
		tanggalSelesai = &parsed
	}

	// Konversi ke repository request
	repoReq := &mongo.UpdatePekerjaanAlumniRepositoryRequest{
		NamaPerusahaan:      req.NamaPerusahaan,
		PosisiJabatan:       req.PosisiJabatan,
		BidangIndustri:      req.BidangIndustri,
		LokasiKerja:         req.LokasiKerja,
		GajiRange:           req.GajiRange,
		TanggalMulaiKerja:   tanggalMulai,
		TanggalSelesaiKerja: tanggalSelesai,
		StatusPekerjaan:     req.StatusPekerjaan,
		DeskripsiPekerjaan:  req.DeskripsiPekerjaan,
	}

	pekerjaan, err := repository.UpdatePekerjaan(db, idStr, repoReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(mongo.UpdatePekerjaanAlumniResponse{
			Success: false,
			Message: "Gagal mengupdate pekerjaan: " + err.Error(),
			Data:    mongo.PekerjaanAlumni{},
		})
	}

	return c.Status(fiber.StatusOK).JSON(mongo.UpdatePekerjaanAlumniResponse{
		Success: true,
		Message: "Berhasil mengupdate pekerjaan",
		Data:    *pekerjaan,
	})
}

func DeletePekerjaanService(c *fiber.Ctx, db *mongoDB.Database) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid",
		})
	}

	// Check if pekerjaan exists
	_, err := repository.GetPekerjaanByID(db, idStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data pekerjaan: " + err.Error(),
		})
	}

	err = repository.DeletePekerjaan(db, idStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menghapus pekerjaan: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Berhasil menghapus pekerjaan",
	})
}
