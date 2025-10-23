package postgre

import (
	"database/sql"
	"fmt"
	model "go-fiber/app/model/postgre"
	repository "go-fiber/app/repository/postgre"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Pekerjaan Alumni Services

func GetAllPekerjaanService(c *fiber.Ctx, db *sql.DB) error {
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
		return c.Status(500).JSON(model.GetAllPekerjaanResponse{
			Success: false,
			Message: "Failed to fetch pekerjaan",
			Data: model.PekerjaanData{
				Items: []model.PekerjaanAlumni{},
				Meta:  model.MetaInfo{},
			},
		})
	}

	total, err := repository.CountPekerjaanRepo(db, search)
	if err != nil {
		return c.Status(500).JSON(model.GetAllPekerjaanResponse{
			Success: false,
			Message: "Failed to count pekerjaan",
			Data: model.PekerjaanData{
				Items: []model.PekerjaanAlumni{},
				Meta:  model.MetaInfo{},
			},
		})
	}

	// Buat response pakai model
	response := model.GetAllPekerjaanResponse{
		Success: true,
		Message: "Berhasil mengambil data pekerjaan",
		Data: model.PekerjaanData{
			Items: pekerjaan,
			Meta: model.MetaInfo{
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

func GetPekerjaanByIDService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.GetPekerjaanAlumniByIDResponse{
			Success: false,
			Message: "ID tidak valid",
			Data:    model.PekerjaanAlumni{},
		})
	}

	pekerjaan, err := repository.GetPekerjaanByID(db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(model.GetPekerjaanAlumniByIDResponse{
				Success: false,
				Message: "Pekerjaan tidak ditemukan",
				Data:    model.PekerjaanAlumni{},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.GetPekerjaanAlumniByIDResponse{
			Success: false,
			Message: "Gagal mengambil data pekerjaan: " + err.Error(),
			Data:    model.PekerjaanAlumni{},
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.GetPekerjaanAlumniByIDResponse{
		Success: true,
		Message: "Berhasil mengambil data pekerjaan",
		Data:    *pekerjaan,
	})
}

func GetPekerjaanByAlumniIDService(c *fiber.Ctx, db *sql.DB) error {
	alumniIDStr := c.Params("alumni_id")
	alumniID, err := strconv.Atoi(alumniIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.GetPekerjaanAlumniByAlumniIDResponse{
			Success: false,
			Message: "ID alumni tidak valid",
			Data:    []model.PekerjaanAlumni{},
		})
	}

	// Check if alumni exists
	_, err = repository.GetAlumniByID(db, alumniID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(model.GetPekerjaanAlumniByAlumniIDResponse{
				Success: false,
				Message: "Alumni tidak ditemukan",
				Data:    []model.PekerjaanAlumni{},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.GetPekerjaanAlumniByAlumniIDResponse{
			Success: false,
			Message: "Gagal mengambil data alumni: " + err.Error(),
			Data:    []model.PekerjaanAlumni{},
		})
	}

	pekerjaan, err := repository.GetPekerjaanByAlumniID(db, alumniID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.GetPekerjaanAlumniByAlumniIDResponse{
			Success: false,
			Message: "Gagal mengambil data pekerjaan: " + err.Error(),
			Data:    []model.PekerjaanAlumni{},
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.GetPekerjaanAlumniByAlumniIDResponse{
		Success: true,
		Message: "Berhasil mengambil data pekerjaan alumni",
		Data:    pekerjaan,
	})
}

func CreatePekerjaanService(c *fiber.Ctx, db *sql.DB) error {
	var req model.CreatePekerjaanAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.CreatePekerjaanAlumniResponse{
			Success: false,
			Message: "Format data tidak valid: " + err.Error(),
			Data:    model.PekerjaanAlumni{},
		})
	}

	// Basic validation
	if req.NamaPerusahaan == "" || req.PosisiJabatan == "" || req.BidangIndustri == "" || req.LokasiKerja == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.CreatePekerjaanAlumniResponse{
			Success: false,
			Message: "Nama perusahaan, posisi jabatan, bidang industri, dan lokasi kerja wajib diisi",
			Data:    model.PekerjaanAlumni{},
		})
	}

	if req.StatusPekerjaan != "aktif" && req.StatusPekerjaan != "selesai" && req.StatusPekerjaan != "resigned" {
		return c.Status(fiber.StatusBadRequest).JSON(model.CreatePekerjaanAlumniResponse{
			Success: false,
			Message: "Status pekerjaan harus aktif, selesai, atau resigned",
			Data:    model.PekerjaanAlumni{},
		})
	}

	// Check if alumni exists
	_, err := repository.GetAlumniByID(db, req.AlumniID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(model.CreatePekerjaanAlumniResponse{
				Success: false,
				Message: "Alumni tidak ditemukan",
				Data:    model.PekerjaanAlumni{},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.CreatePekerjaanAlumniResponse{
			Success: false,
			Message: "Gagal mengambil data alumni: " + err.Error(),
			Data:    model.PekerjaanAlumni{},
		})
	}

	// Parse tanggal mulai kerja
	tanggalMulai, err := time.Parse("2006-01-02", req.TanggalMulaiKerja)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.CreatePekerjaanAlumniResponse{
			Success: false,
			Message: "Format tanggal mulai kerja tidak valid. Gunakan format YYYY-MM-DD",
			Data:    model.PekerjaanAlumni{},
		})
	}

	// Parse tanggal selesai kerja jika ada
	var tanggalSelesai *time.Time
	if req.TanggalSelesaiKerja != nil && *req.TanggalSelesaiKerja != "" {
		parsed, err := time.Parse("2006-01-02", *req.TanggalSelesaiKerja)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.CreatePekerjaanAlumniResponse{
				Success: false,
				Message: "Format tanggal selesai kerja tidak valid. Gunakan format YYYY-MM-DD",
				Data:    model.PekerjaanAlumni{},
			})
		}
		tanggalSelesai = &parsed
	}

	// Konversi ke repository request
	repoReq := &model.CreatePekerjaanAlumniRepositoryRequest{
		AlumniID:            req.AlumniID,
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
		return c.Status(fiber.StatusInternalServerError).JSON(model.CreatePekerjaanAlumniResponse{
			Success: false,
			Message: "Gagal membuat pekerjaan: " + err.Error(),
			Data:    model.PekerjaanAlumni{},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.CreatePekerjaanAlumniResponse{
		Success: true,
		Message: "Berhasil membuat pekerjaan",
		Data:    *pekerjaan,
	})
}

func UpdatePekerjaanService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.UpdatePekerjaanAlumniResponse{
			Success: false,
			Message: "ID tidak valid",
			Data:    model.PekerjaanAlumni{},
		})
	}

	var req model.UpdatePekerjaanAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.UpdatePekerjaanAlumniResponse{
			Success: false,
			Message: "Format data tidak valid: " + err.Error(),
			Data:    model.PekerjaanAlumni{},
		})
	}

	// Check if pekerjaan exists
	_, err = repository.GetPekerjaanByID(db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(model.UpdatePekerjaanAlumniResponse{
				Success: false,
				Message: "Pekerjaan tidak ditemukan",
				Data:    model.PekerjaanAlumni{},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.UpdatePekerjaanAlumniResponse{
			Success: false,
			Message: "Gagal mengambil data pekerjaan: " + err.Error(),
			Data:    model.PekerjaanAlumni{},
		})
	}

	// Validate status
	if req.StatusPekerjaan != "aktif" && req.StatusPekerjaan != "selesai" && req.StatusPekerjaan != "resigned" {
		return c.Status(fiber.StatusBadRequest).JSON(model.UpdatePekerjaanAlumniResponse{
			Success: false,
			Message: "Status pekerjaan harus aktif, selesai, atau resigned",
			Data:    model.PekerjaanAlumni{},
		})
	}

	// Parse tanggal mulai kerja
	tanggalMulai, err := time.Parse("2006-01-02", req.TanggalMulaiKerja)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.UpdatePekerjaanAlumniResponse{
			Success: false,
			Message: "Format tanggal mulai kerja tidak valid. Gunakan format YYYY-MM-DD",
			Data:    model.PekerjaanAlumni{},
		})
	}

	// Parse tanggal selesai kerja jika ada
	var tanggalSelesai *time.Time
	if req.TanggalSelesaiKerja != nil && *req.TanggalSelesaiKerja != "" {
		parsed, err := time.Parse("2006-01-02", *req.TanggalSelesaiKerja)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(model.UpdatePekerjaanAlumniResponse{
				Success: false,
				Message: "Format tanggal selesai kerja tidak valid. Gunakan format YYYY-MM-DD",
				Data:    model.PekerjaanAlumni{},
			})
		}
		tanggalSelesai = &parsed
	}

	// Konversi ke repository request
	repoReq := &model.UpdatePekerjaanAlumniRepositoryRequest{
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

	pekerjaan, err := repository.UpdatePekerjaan(db, id, repoReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.UpdatePekerjaanAlumniResponse{
			Success: false,
			Message: "Gagal mengupdate pekerjaan: " + err.Error(),
			Data:    model.PekerjaanAlumni{},
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.UpdatePekerjaanAlumniResponse{
		Success: true,
		Message: "Berhasil mengupdate pekerjaan",
		Data:    *pekerjaan,
	})
}

func DeletePekerjaanService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.HardDeletePekerjaanAlumniResponse{
			Success: false,
			Message: "ID tidak valid",
		})
	}

	// Check if pekerjaan exists
	_, err = repository.GetPekerjaanByID(db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(model.HardDeletePekerjaanAlumniResponse{
				Success: false,
				Message: "Pekerjaan tidak ditemukan",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.HardDeletePekerjaanAlumniResponse{
			Success: false,
			Message: "Gagal mengambil data pekerjaan: " + err.Error(),
		})
	}

	err = repository.DeletePekerjaan(db, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HardDeletePekerjaanAlumniResponse{
			Success: false,
			Message: "Gagal menghapus pekerjaan: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.HardDeletePekerjaanAlumniResponse{
		Success: true,
		Message: "Berhasil menghapus pekerjaan",
	})
}

func ListDeletedPekerjaanService(c *fiber.Ctx, db *sql.DB) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Get user info for access control
	userIDInterface := c.Locals("user_id")
	userRoleInterface := c.Locals("role")

	if userIDInterface == nil || userRoleInterface == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "User tidak terautentikasi",
			"success": false,
		})
	}

	userID, ok := userIDInterface.(int)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User ID tidak valid",
			"success": false,
		})
	}

	userRole, ok := userRoleInterface.(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Role tidak valid",
			"success": false,
		})
	}

	// Debug logging
	fmt.Printf("DEBUG: userID=%d, userRole=%s, offset=%d, limit=%d\n", userID, userRole, offset, limit)

	pekerjaan, total, err := repository.GetDeletedPekerjaanRepo(db, offset, limit)
	if err != nil {
		fmt.Printf("DEBUG: Error from repository: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil daftar pekerjaan yang dihapus: " + err.Error(),
			"success": false,
		})
	}

	fmt.Printf("DEBUG: Retrieved %d pekerjaan, total=%d\n", len(pekerjaan), total)

	// Filter data based on user role
	var filteredPekerjaan []model.PekerjaanAlumni
	if userRole == "admin" {
		// Admin can see all deleted pekerjaan
		filteredPekerjaan = pekerjaan
		fmt.Printf("DEBUG: Admin access - showing all %d pekerjaan\n", len(filteredPekerjaan))
	} else {
		// User can only see their own deleted pekerjaan
		for _, p := range pekerjaan {
			if p.AlumniID == userID {
				filteredPekerjaan = append(filteredPekerjaan, p)
			}
		}
		// Recalculate total for user
		total = len(filteredPekerjaan)
		fmt.Printf("DEBUG: User access - showing %d pekerjaan for user %d\n", len(filteredPekerjaan), userID)
	}

	response := model.GetSoftDeletedPekerjaanAlumniResponse{
		Success: true,
		Message: "Berhasil mengambil daftar pekerjaan yang dihapus",
		Data:    filteredPekerjaan,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func RestorePekerjaanService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}
	userIDInterface := c.Locals("user_id")
	userRoleInterface := c.Locals("role")

	if userIDInterface == nil || userRoleInterface == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "User tidak terautentikasi",
			"success": false,
		})
	}

	userID, ok := userIDInterface.(int)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User ID tidak valid",
			"success": false,
		})
	}

	userRole, ok := userRoleInterface.(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Role tidak valid",
			"success": false,
		})
	}
	pekerjaan, err := repository.GetPekerjaanWithDeletedByID(db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(model.RestorePekerjaanAlumniResponse{
				Success: false,
				Message: "Pekerjaan tidak ditemukan",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.RestorePekerjaanAlumniResponse{
			Success: false,
			Message: "Gagal mengambil data pekerjaan: " + err.Error(),
		})
	}
	if pekerjaan.IsDeleted == nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.RestorePekerjaanAlumniResponse{
			Success: false,
			Message: "Pekerjaan tidak dalam status terhapus",
		})
	}
	if userRole == "user" {
		if pekerjaan.AlumniID != userID {
			return c.Status(fiber.StatusForbidden).JSON(model.RestorePekerjaanAlumniResponse{
				Success: false,
				Message: "Anda hanya dapat merestore pekerjaan milik sendiri",
			})
		}
	}
	// Admin can restore any pekerjaan (no additional check needed)
	if err := repository.RestorePekerjaan(db, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.RestorePekerjaanAlumniResponse{
			Success: false,
			Message: "Gagal merestore pekerjaan: " + err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(model.RestorePekerjaanAlumniResponse{
		Success: true,
		Message: "Berhasil merestore pekerjaan",
	})
}

func HardDeletePekerjaanService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}
	userIDInterface := c.Locals("user_id")
	userRoleInterface := c.Locals("role")

	if userIDInterface == nil || userRoleInterface == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "User tidak terautentikasi",
			"success": false,
		})
	}

	userID, ok := userIDInterface.(int)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User ID tidak valid",
			"success": false,
		})
	}

	userRole, ok := userRoleInterface.(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Role tidak valid",
			"success": false,
		})
	}
	pekerjaan, err := repository.GetPekerjaanWithDeletedByID(db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(model.HardDeletePekerjaanAlumniResponse{
				Success: false,
				Message: "Pekerjaan tidak ditemukan",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.HardDeletePekerjaanAlumniResponse{
			Success: false,
			Message: "Gagal mengambil data pekerjaan: " + err.Error(),
		})
	}
	if pekerjaan.IsDeleted == nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.HardDeletePekerjaanAlumniResponse{
			Success: false,
			Message: "Pekerjaan harus di-soft delete terlebih dahulu sebelum hard delete",
		})
	}
	if userRole == "user" {
		if pekerjaan.AlumniID != userID {
			return c.Status(fiber.StatusForbidden).JSON(model.HardDeletePekerjaanAlumniResponse{
				Success: false,
				Message: "Anda hanya dapat menghapus permanen pekerjaan milik sendiri",
			})
		}
	}
	// Admin can hard delete any pekerjaan (no additional check needed)
	if err := repository.HardDeletePekerjaan(db, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HardDeletePekerjaanAlumniResponse{
			Success: false,
			Message: "Gagal menghapus permanen pekerjaan: " + err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(model.HardDeletePekerjaanAlumniResponse{
		Success: true,
		Message: "Berhasil menghapus permanen pekerjaan",
	})
}

func SoftDeletePekerjaanService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.SoftDeletePekerjaanAlumniResponse{
			Success: false,
			Message: "ID tidak valid",
		})
	}

	// Get user info from context
	userIDInterface := c.Locals("user_id")
	userRoleInterface := c.Locals("role")

	if userIDInterface == nil || userRoleInterface == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(model.SoftDeletePekerjaanAlumniResponse{
			Success: false,
			Message: "User tidak terautentikasi",
		})
	}

	userID, ok := userIDInterface.(int)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(model.SoftDeletePekerjaanAlumniResponse{
			Success: false,
			Message: "User ID tidak valid",
		})
	}

	userRole, ok := userRoleInterface.(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(model.SoftDeletePekerjaanAlumniResponse{
			Success: false,
			Message: "Role tidak valid",
		})
	}

	// Get pekerjaan data to check ownership
	pekerjaan, err := repository.GetPekerjaanWithDeletedByID(db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(model.SoftDeletePekerjaanAlumniResponse{
				Success: false,
				Message: "Pekerjaan tidak ditemukan",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.SoftDeletePekerjaanAlumniResponse{
			Success: false,
			Message: "Gagal mengambil data pekerjaan: " + err.Error(),
		})
	}

	// Check if already soft deleted
	if pekerjaan.IsDeleted != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.SoftDeletePekerjaanAlumniResponse{
			Success: false,
			Message: "Pekerjaan sudah dihapus sebelumnya",
		})
	}

	// Role-based access control
	if userRole == "user" {
		// User can only delete their own pekerjaan
		if pekerjaan.AlumniID != userID {
			return c.Status(fiber.StatusForbidden).JSON(model.SoftDeletePekerjaanAlumniResponse{
				Success: false,
				Message: "Anda hanya dapat menghapus pekerjaan yang dimiliki sendiri",
			})
		}
	}
	// Admin can delete any pekerjaan (no additional check needed)

	// Perform soft delete
	err = repository.SoftDeletePekerjaan(db, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.SoftDeletePekerjaanAlumniResponse{
			Success: false,
			Message: "Gagal menghapus pekerjaan: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.SoftDeletePekerjaanAlumniResponse{
		Success: true,
		Message: "Berhasil menghapus pekerjaan",
	})
}
