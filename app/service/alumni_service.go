package service

import (
	"database/sql"
	"go-fiber/app/model"
	"go-fiber/app/repository"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"strconv"
	"strings"
    "go-fiber/utils"
)

// Alumni Services

func GetAllAlumniService(c *fiber.Ctx, db *sql.DB) error {
	// Debug semua query parameters
	log.Printf("=== All Query Parameters ===")
	for key, value := range c.Queries() {
		log.Printf("Query param '%s': '%s' (len: %d)", key, value, len(value))
	}
	
	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "id")
	order := c.Query("order", "asc")
	search := c.Query("search", "")

	// Debug logging
	log.Printf("=== Service Debug ===")
	log.Printf("Service received parameters - page: %d, limit: %d, sortBy: '%s', order: '%s', search: '%s' (len: %d)", page, limit, sortBy, order, search, len(search))
	log.Printf("Search is empty: %t", search == "")

	// Hitung offset untuk pagination
	offset := (page - 1) * limit

	// Validasi input
	sortByWhitelist := map[string]bool{"id": true, "nama": true, "email": true, "jurusan": true, "angkatan": true, "tahun_lulus": true, "created_at": true}
	if !sortByWhitelist[sortBy] {
		sortBy = "id"
	}
	if strings.ToLower(order) != "desc" {
		order = "asc"
	}

	// Ambil data dari repository
	alumni, err := repository.GetAlumniRepo(db, search, sortBy, order, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch alumni",
		})
	}

	total, err := repository.CountAlumniRepo(db, search)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to count alumni",
		})
	}

	// Buat response pakai model
	response := model.GetAllAlumniResponse{
		Success: true,
		Message: "Berhasil mengambil data alumni",
		Data: model.AlumniData{
			Items: alumni,
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

func GetAlumniByIDService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.GetAlumniByIDResponse{
			Success: false,
			Message: "ID tidak valid",
			Data:    model.Alumni{},
		})
	}

	alumni, err := repository.GetAlumniByID(db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(model.GetAlumniByIDResponse{
				Success: false,
				Message: "Alumni tidak ditemukan",
				Data:    model.Alumni{},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.GetAlumniByIDResponse{
			Success: false,
			Message: "Gagal mengambil data alumni: " + err.Error(),
			Data:    model.Alumni{},
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.GetAlumniByIDResponse{
		Success: true,
		Message: "Berhasil mengambil data alumni",
		Data:    *alumni,
	})
}

func CreateAlumniService(c *fiber.Ctx, db *sql.DB) error {
	var req model.CreateAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.CreateAlumniResponse{
			Success: false,
			Message: "Format data tidak valid: " + err.Error(),
			Data:    model.Alumni{},
		})
	}

	if req.NIM == "" || req.Nama == "" || req.Jurusan == "" || req.Email == "" || req.Password == "" || req.RoleID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(model.CreateAlumniResponse{
			Success: false,
			Message: "NIM, nama, jurusan, email, password, dan role_id wajib diisi",
			Data:    model.Alumni{},
		})
	}

	if req.Angkatan <= 0 || req.TahunLulus <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(model.CreateAlumniResponse{
			Success: false,
			Message: "Angkatan dan tahun lulus harus lebih dari 0",
			Data:    model.Alumni{},
		})
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.CreateAlumniResponse{
			Success: false,
			Message: "Gagal memproses password",
			Data:    model.Alumni{},
		})
	}

	// Konversi ke repository request
	repoReq := &model.CreateAlumniRepositoryRequest{
		NIM:        req.NIM,
		Nama:       req.Nama,
		Jurusan:    req.Jurusan,
		Angkatan:   req.Angkatan,
		TahunLulus: req.TahunLulus,
		Email:      req.Email,
		Password:   hashed,
		RoleID:     req.RoleID,
		NoTelepon:  req.NoTelepon,
		Alamat:     req.Alamat,
	}

	alumni, err := repository.CreateAlumni(db, repoReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.CreateAlumniResponse{
			Success: false,
			Message: "Gagal membuat alumni: " + err.Error(),
			Data:    model.Alumni{},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.CreateAlumniResponse{
		Success: true,
		Message: "Berhasil membuat alumni",
		Data:    *alumni,
	})
}

func UpdateAlumniService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.UpdateAlumniResponse{
			Success: false,
			Message: "ID tidak valid",
			Data:    model.Alumni{},
		})
	}

	var req model.UpdateAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.UpdateAlumniResponse{
			Success: false,
			Message: "Format data tidak valid: " + err.Error(),
			Data:    model.Alumni{},
		})
	}

	// Check if alumni exists
	_, err = repository.GetAlumniByID(db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(model.UpdateAlumniResponse{
				Success: false,
				Message: "Alumni tidak ditemukan",
				Data:    model.Alumni{},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.UpdateAlumniResponse{
			Success: false,
			Message: "Gagal mengambil data alumni: " + err.Error(),
			Data:    model.Alumni{},
		})
	}

	// Hash password if provided
	repoReq := &model.UpdateAlumniRepositoryRequest{
		NIM:        req.NIM,
		Nama:       req.Nama,
		Jurusan:    req.Jurusan,
		Angkatan:   req.Angkatan,
		TahunLulus: req.TahunLulus,
		Email:      req.Email,
		RoleID:     req.RoleID,
		NoTelepon:  req.NoTelepon,
		Alamat:     req.Alamat,
	}

	if req.Password != nil && *req.Password != "" {
		hashed, err := utils.HashPassword(*req.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(model.UpdateAlumniResponse{
				Success: false,
				Message: "Gagal memproses password",
				Data:    model.Alumni{},
			})
		}
		repoReq.Password = &hashed
	}

	alumni, err := repository.UpdateAlumni(db, id, repoReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.UpdateAlumniResponse{
			Success: false,
			Message: "Gagal mengupdate alumni: " + err.Error(),
			Data:    model.Alumni{},
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.UpdateAlumniResponse{
		Success: true,
		Message: "Berhasil mengupdate alumni",
		Data:    *alumni,
	})
}

func DeleteAlumniService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.DeleteAlumniResponse{
			Success: false,
			Message: "ID tidak valid",
		})
	}

	// Check if alumni exists
	_, err = repository.GetAlumniByID(db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(model.DeleteAlumniResponse{
				Success: false,
				Message: "Alumni tidak ditemukan",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.DeleteAlumniResponse{
			Success: false,
			Message: "Gagal mengambil data alumni: " + err.Error(),
		})
	}

	err = repository.DeleteAlumni(db, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.DeleteAlumniResponse{
			Success: false,
			Message: "Gagal menghapus alumni: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.DeleteAlumniResponse{
		Success: true,
		Message: "Berhasil menghapus alumni",
	})
}


// Legacy function for backward compatibility
func CheckAlumniService(c *fiber.Ctx, db *sql.DB) error {
	key := c.Params("key")
	if key != os.Getenv("API_KEY") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Key tidak valid",
			"success": false,
		})
	}
	nim := c.FormValue("nim")
	if nim == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "NIM wajib diisi",
			"success": false,
		})
	}
	alumni, err := repository.CheckAlumniByNim(db, nim)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message": "Mahasiswa bukan alumni",
				"success": true,
				"isAlumni": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal cek alumni karena " + err.Error(),
			"success": false,
		})
	}
	return c.Status(fiber.StatusOK).JSON(model.CheckAlumniResponse{
		Success:  true,
		Message:  "Berhasil mendapatkan data alumni",
		IsAlumni: true,
		Alumni:   alumni,
	})
}

// Get Alumni Employment Status Service
func GetAlumniEmploymentStatusService(c *fiber.Ctx, db *sql.DB) error {
	// Parse query parameters
	req := &model.AlumniEmploymentStatusRequest{
		Page:  1,
		Limit: 20,
	}

	// Parse query parameters
	if idStr := c.Query("id"); idStr != "" {
		if id, err := strconv.Atoi(idStr); err == nil {
			req.ID = &id
		}
	}
	if nama := c.Query("nama"); nama != "" {
		req.Nama = &nama
	}
	if jurusan := c.Query("jurusan"); jurusan != "" {
		req.Jurusan = &jurusan
	}
	if angkatanStr := c.Query("angkatan"); angkatanStr != "" {
		if angkatan, err := strconv.Atoi(angkatanStr); err == nil {
			req.Angkatan = &angkatan
		}
	}
	if bidangIndustri := c.Query("bidang_industri"); bidangIndustri != "" {
		req.BidangIndustri = &bidangIndustri
	}
	if namaPerusahaan := c.Query("nama_perusahaan"); namaPerusahaan != "" {
		req.NamaPerusahaan = &namaPerusahaan
	}
	if posisiJabatan := c.Query("posisi_jabatan"); posisiJabatan != "" {
		req.PosisiJabatan = &posisiJabatan
	}
	if lebihDari1TahunStr := c.Query("lebih_dari_1_tahun"); lebihDari1TahunStr != "" {
		if lebihDari1Tahun, err := strconv.Atoi(lebihDari1TahunStr); err == nil {
			req.LebihDari1Tahun = &lebihDari1Tahun
		}
	}
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			req.Page = page
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			req.Limit = limit
		}
	}

	// Get data from repository
	results, err := repository.GetAlumniEmploymentStatus(db, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data status pekerjaan alumni: " + err.Error(),
			"success": false,
		})
	}

	// Calculate pagination info
	totalRecords := len(results)
	totalPages := (totalRecords + req.Limit - 1) / req.Limit

	return c.Status(fiber.StatusOK).JSON(model.GetAlumniEmploymentStatusResponse{
		Success: true,
		Message: "Berhasil mengambil data status pekerjaan alumni",
		Data:    results,
		Pagination: model.PaginationInfo{
			CurrentPage:  req.Page,
			PerPage:      req.Limit,
			TotalRecords: totalRecords,
			TotalPages:   totalPages,
		},
	})
}
