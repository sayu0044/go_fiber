package postgre

import (
	"database/sql"
	"fmt"
	model "go-fiber/app/model/postgre"
	"log"
	"strings"
	"time"
)

// Alumni Repository Functions

func GetAllAlumni(db *sql.DB) ([]model.Alumni, error) {
	query := `SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at FROM alumni ORDER BY nama`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alumni []model.Alumni
	for rows.Next() {
		var a model.Alumni
		err := rows.Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		alumni = append(alumni, a)
	}
	return alumni, nil
}

// GetAlumniRepo -> ambil data alumni dari DB dengan pagination, sorting, dan search
func GetAlumniRepo(db *sql.DB, search, sortBy, order string, limit, offset int) ([]model.Alumni, error) {
	// Debug logging
	log.Printf("=== GetAlumniRepo Debug ===")
	log.Printf("Search parameters - search: '%s' (len: %d), sortBy: '%s', order: '%s', limit: %d, offset: %d", search, len(search), sortBy, order, limit, offset)
	log.Printf("Search is empty: %t", search == "")

	// Buat query yang lebih sederhana untuk testing
	var query string
	var args []interface{}

	if search == "" {
		// Jika search kosong, tampilkan semua data
		query = fmt.Sprintf(`
            SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
			FROM alumni
			ORDER BY %s %s
			LIMIT $1 OFFSET $2
		`, sortBy, order)
		args = []interface{}{limit, offset}
		log.Printf("Using query without search filter")
	} else {
		// Jika ada search, gunakan filter
		searchPattern := "%" + search + "%"
		query = fmt.Sprintf(`
            SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
			FROM alumni
			WHERE nama ILIKE $1 OR email ILIKE $1 OR jurusan ILIKE $1 OR nim ILIKE $1
			ORDER BY %s %s
			LIMIT $2 OFFSET $3
		`, sortBy, order)
		args = []interface{}{searchPattern, limit, offset}
		log.Printf("Using query with search filter - pattern: '%s'", searchPattern)
	}

	log.Printf("Final query: %s", query)
	log.Printf("Query parameters: %v", args)

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Printf("Query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var alumni []model.Alumni
	for rows.Next() {
		var a model.Alumni
		err := rows.Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		alumni = append(alumni, a)
	}
	return alumni, nil
}

// CountAlumniRepo -> hitung total data untuk pagination
func CountAlumniRepo(db *sql.DB, search string) (int, error) {
	var total int
	var countQuery string
	var args []interface{}

	if search == "" {
		countQuery = `SELECT COUNT(*) FROM alumni`
		args = []interface{}{}
	} else {
		countQuery = `SELECT COUNT(*) FROM alumni WHERE nama ILIKE $1 OR email ILIKE $1 OR jurusan ILIKE $1 OR nim ILIKE $1`
		args = []interface{}{"%" + search + "%"}
	}

	err := db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return total, nil
}

func GetAlumniByID(db *sql.DB, id int) (*model.Alumni, error) {
	alumni := new(model.Alumni)
	query := `SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, role_id, no_telepon, alamat, password, created_at, updated_at FROM alumni WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&alumni.ID, &alumni.NIM, &alumni.Nama, &alumni.Jurusan, &alumni.Angkatan, &alumni.TahunLulus, &alumni.Email, &alumni.RoleID, &alumni.NoTelepon, &alumni.Alamat, &alumni.Password, &alumni.CreatedAt, &alumni.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return alumni, nil
}

func CreateAlumni(db *sql.DB, req *model.CreateAlumniRepositoryRequest) (*model.Alumni, error) {
	query := `INSERT INTO alumni (nim, nama, jurusan, angkatan, tahun_lulus, email, role_id, no_telepon, alamat, password, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id, created_at, updated_at`

	now := time.Now()
	var id int
	var createdAt, updatedAt time.Time

	err := db.QueryRow(query, req.NIM, req.Nama, req.Jurusan, req.Angkatan, req.TahunLulus, req.Email, req.RoleID, req.NoTelepon, req.Alamat, req.Password, now, now).
		Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	alumni := &model.Alumni{
		ID:         id,
		NIM:        req.NIM,
		Nama:       req.Nama,
		Jurusan:    req.Jurusan,
		Angkatan:   req.Angkatan,
		TahunLulus: req.TahunLulus,
		Email:      req.Email,
		RoleID:     req.RoleID,
		NoTelepon:  req.NoTelepon,
		Alamat:     req.Alamat,
		Password:   req.Password,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
	return alumni, nil
}

func UpdateAlumni(db *sql.DB, id int, req *model.UpdateAlumniRepositoryRequest) (*model.Alumni, error) {
	// Build dynamic query based on provided fields
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.NIM != nil {
		setParts = append(setParts, "nim = $"+fmt.Sprintf("%d", argIndex))
		args = append(args, *req.NIM)
		argIndex++
	}
	if req.Nama != nil {
		setParts = append(setParts, "nama = $"+fmt.Sprintf("%d", argIndex))
		args = append(args, *req.Nama)
		argIndex++
	}
	if req.Jurusan != nil {
		setParts = append(setParts, "jurusan = $"+fmt.Sprintf("%d", argIndex))
		args = append(args, *req.Jurusan)
		argIndex++
	}
	if req.Angkatan != nil {
		setParts = append(setParts, "angkatan = $"+fmt.Sprintf("%d", argIndex))
		args = append(args, *req.Angkatan)
		argIndex++
	}
	if req.TahunLulus != nil {
		setParts = append(setParts, "tahun_lulus = $"+fmt.Sprintf("%d", argIndex))
		args = append(args, *req.TahunLulus)
		argIndex++
	}
	if req.Email != nil {
		setParts = append(setParts, "email = $"+fmt.Sprintf("%d", argIndex))
		args = append(args, *req.Email)
		argIndex++
	}
	if req.RoleID != nil {
		setParts = append(setParts, "role_id = $"+fmt.Sprintf("%d", argIndex))
		args = append(args, *req.RoleID)
		argIndex++
	}
	if req.Password != nil {
		setParts = append(setParts, "password = $"+fmt.Sprintf("%d", argIndex))
		args = append(args, *req.Password)
		argIndex++
	}
	if req.NoTelepon != nil {
		setParts = append(setParts, "no_telepon = $"+fmt.Sprintf("%d", argIndex))
		args = append(args, *req.NoTelepon)
		argIndex++
	}
	if req.Alamat != nil {
		setParts = append(setParts, "alamat = $"+fmt.Sprintf("%d", argIndex))
		args = append(args, *req.Alamat)
		argIndex++
	}

	if len(setParts) == 0 {
		return GetAlumniByID(db, id)
	}

	setParts = append(setParts, "updated_at = $"+fmt.Sprintf("%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	args = append(args, id)

	query := "UPDATE alumni SET " + strings.Join(setParts, ", ") + " WHERE id = $" + fmt.Sprintf("%d", argIndex) + " RETURNING id, nim, nama, jurusan, angkatan, tahun_lulus, email, role_id, no_telepon, alamat, password, created_at, updated_at"

	alumni := new(model.Alumni)
	err := db.QueryRow(query, args...).Scan(&alumni.ID, &alumni.NIM, &alumni.Nama, &alumni.Jurusan, &alumni.Angkatan, &alumni.TahunLulus, &alumni.Email, &alumni.RoleID, &alumni.NoTelepon, &alumni.Alamat, &alumni.Password, &alumni.CreatedAt, &alumni.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return alumni, nil
}

func DeleteAlumni(db *sql.DB, id int) error {
	query := `DELETE FROM alumni WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

// Legacy function for backward compatibility
func CheckAlumniByNim(db *sql.DB, nim string) (*model.Alumni, error) {
	alumni := new(model.Alumni)
	query := `SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, role_id, no_telepon, alamat, password, created_at, updated_at FROM alumni WHERE nim = $1`
	err := db.QueryRow(query, nim).Scan(&alumni.ID, &alumni.NIM, &alumni.Nama, &alumni.Jurusan, &alumni.Angkatan, &alumni.TahunLulus, &alumni.Email, &alumni.RoleID, &alumni.NoTelepon, &alumni.Alamat, &alumni.Password, &alumni.CreatedAt, &alumni.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return alumni, nil
}

// Get Alumni Employment Status with filtering and pagination
func GetAlumniEmploymentStatus(db *sql.DB, req *model.AlumniEmploymentStatusRequest) ([]model.AlumniEmploymentStatus, error) {
	// Set default pagination
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}

	offset := (req.Page - 1) * req.Limit

	// Build WHERE clause based on filters
	whereConditions := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.ID != nil {
		whereConditions = append(whereConditions, "a.id = $"+fmt.Sprintf("%d", argIndex))
		args = append(args, *req.ID)
		argIndex++
	}
	if req.Nama != nil {
		whereConditions = append(whereConditions, "LOWER(a.nama) LIKE LOWER($"+fmt.Sprintf("%d", argIndex)+")")
		args = append(args, "%"+*req.Nama+"%")
		argIndex++
	}
	if req.Jurusan != nil {
		whereConditions = append(whereConditions, "LOWER(a.jurusan) LIKE LOWER($"+fmt.Sprintf("%d", argIndex)+")")
		args = append(args, "%"+*req.Jurusan+"%")
		argIndex++
	}
	if req.Angkatan != nil {
		whereConditions = append(whereConditions, "a.angkatan = $"+fmt.Sprintf("%d", argIndex))
		args = append(args, *req.Angkatan)
		argIndex++
	}
	if req.BidangIndustri != nil {
		whereConditions = append(whereConditions, "LOWER(le.bidang_industri) LIKE LOWER($"+fmt.Sprintf("%d", argIndex)+")")
		args = append(args, "%"+*req.BidangIndustri+"%")
		argIndex++
	}
	if req.NamaPerusahaan != nil {
		whereConditions = append(whereConditions, "LOWER(le.nama_perusahaan) LIKE LOWER($"+fmt.Sprintf("%d", argIndex)+")")
		args = append(args, "%"+*req.NamaPerusahaan+"%")
		argIndex++
	}
	if req.PosisiJabatan != nil {
		whereConditions = append(whereConditions, "LOWER(le.posisi_jabatan) LIKE LOWER($"+fmt.Sprintf("%d", argIndex)+")")
		args = append(args, "%"+*req.PosisiJabatan+"%")
		argIndex++
	}
	if req.LebihDari1Tahun != nil {
		whereConditions = append(whereConditions, "CASE WHEN le.tanggal_mulai_kerja <= (CURRENT_DATE - INTERVAL '1 year') THEN 1 ELSE 0 END = $"+fmt.Sprintf("%d", argIndex))
		args = append(args, *req.LebihDari1Tahun)
		argIndex++
	}

	// Build WHERE clause
	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Add pagination parameters
	args = append(args, req.Limit, offset)

	query := `
		-- Subquery: pekerjaan terbaru per alumni (berdasarkan tanggal_mulai_kerja terbesar)
		WITH latest_employment AS (
			SELECT p.*
			FROM pekerjaan_alumni p
			JOIN (
				SELECT alumni_id, MAX(tanggal_mulai_kerja) AS latest_start
				FROM pekerjaan_alumni
				GROUP BY alumni_id
			) t ON p.alumni_id = t.alumni_id AND p.tanggal_mulai_kerja = t.latest_start
		),
		employment_counts AS (
			SELECT alumni_id, COUNT(*) AS employment_count
			FROM pekerjaan_alumni
			GROUP BY alumni_id
		)
		SELECT
			a.id,
			a.nama,
			a.jurusan,
			a.angkatan,
			le.bidang_industri,
			le.nama_perusahaan,
			le.posisi_jabatan,
			le.tanggal_mulai_kerja,
			le.gaji_range,
			CASE
				WHEN le.tanggal_mulai_kerja <= (CURRENT_DATE - INTERVAL '1 year') THEN 1
				ELSE 0
			END AS lebih_dari_1_tahun,
			COALESCE(ec.employment_count, 0) AS employment_count
		FROM alumni a
		LEFT JOIN latest_employment le ON a.id = le.alumni_id
		LEFT JOIN employment_counts ec ON a.id = ec.alumni_id
		` + whereClause + `
		ORDER BY a.nama
		LIMIT $` + fmt.Sprintf("%d", argIndex) + ` OFFSET $` + fmt.Sprintf("%d", argIndex+1)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.AlumniEmploymentStatus
	for rows.Next() {
		var result model.AlumniEmploymentStatus
		err := rows.Scan(
			&result.ID,
			&result.Nama,
			&result.Jurusan,
			&result.Angkatan,
			&result.BidangIndustri,
			&result.NamaPerusahaan,
			&result.PosisiJabatan,
			&result.TanggalMulaiKerja,
			&result.GajiRange,
			&result.LebihDari1Tahun,
			&result.EmploymentCount,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}
