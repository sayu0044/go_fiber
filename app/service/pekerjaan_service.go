package service

import (
	"database/sql"
	"fmt"
	"go-fiber/app/model"
	"go-fiber/app/repository"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
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
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch pekerjaan",
		})
	}

	total, err := repository.CountPekerjaanRepo(db, search)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to count pekerjaan",
		})
	}

	// Buat response pakai model
	response := model.PekerjaanResponse{
		Data: pekerjaan,
		Meta: model.MetaInfo{
			Page:   page,
			Limit:  limit,
			Total:  total,
			Pages:  (total + limit - 1) / limit,
			SortBy: sortBy,
			Order:  order,
			Search: search,
		},
	}

	return c.JSON(response)
}

func GetPekerjaanByIDService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}

	pekerjaan, err := repository.GetPekerjaanByID(db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Pekerjaan tidak ditemukan",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil data pekerjaan",
		"success": true,
		"data":    pekerjaan,
	})
}

func GetPekerjaanByAlumniIDService(c *fiber.Ctx, db *sql.DB) error {
	alumniIDStr := c.Params("alumni_id")
	alumniID, err := strconv.Atoi(alumniIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID alumni tidak valid",
			"success": false,
		})
	}

	// Check if alumni exists
	_, err = repository.GetAlumniByID(db, alumniID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Alumni tidak ditemukan",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data alumni: " + err.Error(),
			"success": false,
		})
	}

	pekerjaan, err := repository.GetPekerjaanByAlumniID(db, alumniID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil data pekerjaan alumni",
		"success": true,
		"data":    pekerjaan,
	})
}

func CreatePekerjaanService(c *fiber.Ctx, db *sql.DB) error {
	var req model.CreatePekerjaanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Format data tidak valid: " + err.Error(),
			"success": false,
		})
	}

	// Basic validation
	if req.NamaPerusahaan == "" || req.PosisiJabatan == "" || req.BidangIndustri == "" || req.LokasiKerja == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Nama perusahaan, posisi jabatan, bidang industri, dan lokasi kerja wajib diisi",
			"success": false,
		})
	}

	if req.StatusPekerjaan != "aktif" && req.StatusPekerjaan != "selesai" && req.StatusPekerjaan != "resigned" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Status pekerjaan harus aktif, selesai, atau resigned",
			"success": false,
		})
	}

	// Check if alumni exists
	_, err := repository.GetAlumniByID(db, req.AlumniID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Alumni tidak ditemukan",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data alumni: " + err.Error(),
			"success": false,
		})
	}

	pekerjaan, err := repository.CreatePekerjaan(db, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal membuat pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Berhasil membuat pekerjaan",
		"success": true,
		"data":    pekerjaan,
	})
}

func UpdatePekerjaanService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}

	var req model.UpdatePekerjaanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Format data tidak valid: " + err.Error(),
			"success": false,
		})
	}

	// Check if pekerjaan exists
	_, err = repository.GetPekerjaanByID(db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Pekerjaan tidak ditemukan",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	// If updating alumni_id, check if alumni exists
	if req.AlumniID != nil {
		_, err = repository.GetAlumniByID(db, *req.AlumniID)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"message": "Alumni tidak ditemukan",
					"success": false,
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Gagal mengambil data alumni: " + err.Error(),
				"success": false,
			})
		}
	}

	// Validate status if provided
	if req.StatusPekerjaan != nil {
		if *req.StatusPekerjaan != "aktif" && *req.StatusPekerjaan != "selesai" && *req.StatusPekerjaan != "resigned" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Status pekerjaan harus aktif, selesai, atau resigned",
				"success": false,
			})
		}
	}

	pekerjaan, err := repository.UpdatePekerjaan(db, id, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengupdate pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengupdate pekerjaan",
		"success": true,
		"data":    pekerjaan,
	})
}

func DeletePekerjaanService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}

	// Check if pekerjaan exists
	_, err = repository.GetPekerjaanByID(db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Pekerjaan tidak ditemukan",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	err = repository.DeletePekerjaan(db, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menghapus pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil menghapus pekerjaan",
		"success": true,
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
    
    response := model.PekerjaanResponse{
        Data: filteredPekerjaan,
        Meta: model.MetaInfo{
            Page:  page,
            Limit: limit,
            Total: total,
            Pages: (total + limit - 1) / limit,
        },
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
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "message": "Pekerjaan tidak ditemukan",
                "success": false,
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Gagal mengambil data pekerjaan: " + err.Error(),
            "success": false,
        })
    }
    if pekerjaan.IsDeleted == nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Pekerjaan tidak dalam status terhapus",
            "success": false,
        })
    }
    if userRole == "user" {
        if pekerjaan.AlumniID != userID {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "message": "Anda hanya dapat merestore pekerjaan milik sendiri",
                "success": false,
            })
        }
    }
    // Admin can restore any pekerjaan (no additional check needed)
    if err := repository.RestorePekerjaan(db, id); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Gagal merestore pekerjaan: " + err.Error(),
            "success": false,
        })
    }
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Berhasil merestore pekerjaan",
        "success": true,
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
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "message": "Pekerjaan tidak ditemukan",
                "success": false,
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Gagal mengambil data pekerjaan: " + err.Error(),
            "success": false,
        })
    }
    if pekerjaan.IsDeleted == nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Pekerjaan harus di-soft delete terlebih dahulu sebelum hard delete",
            "success": false,
        })
    }
    if userRole == "user" {
        if pekerjaan.AlumniID != userID {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "message": "Anda hanya dapat menghapus permanen pekerjaan milik sendiri",
                "success": false,
            })
        }
    }
    // Admin can hard delete any pekerjaan (no additional check needed)
    if err := repository.HardDeletePekerjaan(db, id); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Gagal menghapus permanen pekerjaan: " + err.Error(),
            "success": false,
        })
    }
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Berhasil menghapus permanen pekerjaan",
        "success": true,
    })
}

func SoftDeletePekerjaanService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}

	

	// Get user info from context
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

	// Get pekerjaan data to check ownership
	pekerjaan, err := repository.GetPekerjaanWithDeletedByID(db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Pekerjaan tidak ditemukan",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	// Check if already soft deleted
	if pekerjaan.IsDeleted != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Pekerjaan sudah dihapus sebelumnya",
			"success": false,
		})
	}

	// Role-based access control
	if userRole == "user" {
		// User can only delete their own pekerjaan
		if pekerjaan.AlumniID != userID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Anda hanya dapat menghapus pekerjaan yang dimiliki sendiri",
				"success": false,
			})
		}
	}
	// Admin can delete any pekerjaan (no additional check needed)

	// Perform soft delete
	err = repository.SoftDeletePekerjaan(db, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menghapus pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil menghapus pekerjaan",
		"success": true,
	})
}



