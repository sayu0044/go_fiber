package repository

import (
	"database/sql"
	"fmt"
	"go-fiber/app/model"
	"log"
	"strings"
	"time"
)

// Pekerjaan Alumni Repository Functions

func GetAllPekerjaan(db *sql.DB) ([]model.PekerjaanAlumni, error) {
	query := `SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at, is_delete FROM pekerjaan_alumni WHERE is_delete IS NULL ORDER BY created_at DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pekerjaan []model.PekerjaanAlumni
	for rows.Next() {
		var p model.PekerjaanAlumni
		err := rows.Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt, &p.IsDeleted)
		if err != nil {
			return nil, err
		}
		pekerjaan = append(pekerjaan, p)
	}
	return pekerjaan, nil
}

func GetPekerjaanByID(db *sql.DB, id int) (*model.PekerjaanAlumni, error) {
	pekerjaan := new(model.PekerjaanAlumni)
	query := `SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at, is_delete FROM pekerjaan_alumni WHERE id = $1 AND is_delete IS NULL`
	err := db.QueryRow(query, id).Scan(&pekerjaan.ID, &pekerjaan.AlumniID, &pekerjaan.NamaPerusahaan, &pekerjaan.PosisiJabatan, &pekerjaan.BidangIndustri, &pekerjaan.LokasiKerja, &pekerjaan.GajiRange, &pekerjaan.TanggalMulaiKerja, &pekerjaan.TanggalSelesaiKerja, &pekerjaan.StatusPekerjaan, &pekerjaan.DeskripsiPekerjaan, &pekerjaan.CreatedAt, &pekerjaan.UpdatedAt, &pekerjaan.IsDeleted)
	if err != nil {
		return nil, err
	}
	return pekerjaan, nil
}

func GetPekerjaanByAlumniID(db *sql.DB, alumniID int) ([]model.PekerjaanAlumni, error) {
	query := `SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at, is_delete FROM pekerjaan_alumni WHERE alumni_id = $1 AND is_delete IS NULL ORDER BY tanggal_mulai_kerja DESC`
	rows, err := db.Query(query, alumniID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pekerjaan []model.PekerjaanAlumni
	for rows.Next() {
		var p model.PekerjaanAlumni
		err := rows.Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt, &p.IsDeleted)
		if err != nil {
			return nil, err
		}
		pekerjaan = append(pekerjaan, p)
	}
	return pekerjaan, nil
}

func CreatePekerjaan(db *sql.DB, req *model.CreatePekerjaanAlumniRepositoryRequest) (*model.PekerjaanAlumni, error) {
	query := `INSERT INTO pekerjaan_alumni (alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id, created_at, updated_at`
	
	now := time.Now()
	var id int
	var createdAt, updatedAt time.Time
	
	err := db.QueryRow(query, req.AlumniID, req.NamaPerusahaan, req.PosisiJabatan, req.BidangIndustri, req.LokasiKerja, req.GajiRange, req.TanggalMulaiKerja, req.TanggalSelesaiKerja, req.StatusPekerjaan, req.DeskripsiPekerjaan, now, now).
		Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	pekerjaan := &model.PekerjaanAlumni{
		ID:                  id,
		AlumniID:            req.AlumniID,
		NamaPerusahaan:      req.NamaPerusahaan,
		PosisiJabatan:       req.PosisiJabatan,
		BidangIndustri:      req.BidangIndustri,
		LokasiKerja:         req.LokasiKerja,
		GajiRange:           req.GajiRange,
		TanggalMulaiKerja:   req.TanggalMulaiKerja,
		TanggalSelesaiKerja: req.TanggalSelesaiKerja,
		StatusPekerjaan:     req.StatusPekerjaan,
		DeskripsiPekerjaan:  req.DeskripsiPekerjaan,
		CreatedAt:           createdAt,
		UpdatedAt:           updatedAt,
	}
	return pekerjaan, nil
}

func UpdatePekerjaan(db *sql.DB, id int, req *model.UpdatePekerjaanAlumniRepositoryRequest) (*model.PekerjaanAlumni, error) {
	// Build dynamic query based on provided fields
	setParts := []string{
		"nama_perusahaan = $1",
		"posisi_jabatan = $2",
		"bidang_industri = $3",
		"lokasi_kerja = $4",
		"gaji_range = $5",
		"tanggal_mulai_kerja = $6",
		"tanggal_selesai_kerja = $7",
		"status_pekerjaan = $8",
		"deskripsi_pekerjaan = $9",
		"updated_at = $10",
	}
	
	args := []interface{}{
		req.NamaPerusahaan,
		req.PosisiJabatan,
		req.BidangIndustri,
		req.LokasiKerja,
		req.GajiRange,
		req.TanggalMulaiKerja,
		req.TanggalSelesaiKerja,
		req.StatusPekerjaan,
		req.DeskripsiPekerjaan,
		time.Now(),
		id,
	}

	query := "UPDATE pekerjaan_alumni SET " + strings.Join(setParts, ", ") + " WHERE id = $11 RETURNING id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at"

	pekerjaan := new(model.PekerjaanAlumni)
	err := db.QueryRow(query, args...).Scan(&pekerjaan.ID, &pekerjaan.AlumniID, &pekerjaan.NamaPerusahaan, &pekerjaan.PosisiJabatan, &pekerjaan.BidangIndustri, &pekerjaan.LokasiKerja, &pekerjaan.GajiRange, &pekerjaan.TanggalMulaiKerja, &pekerjaan.TanggalSelesaiKerja, &pekerjaan.StatusPekerjaan, &pekerjaan.DeskripsiPekerjaan, &pekerjaan.CreatedAt, &pekerjaan.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return pekerjaan, nil
}

func DeletePekerjaan(db *sql.DB, id int) error {
	query := `DELETE FROM pekerjaan_alumni WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

func SoftDeletePekerjaan(db *sql.DB, id int) error {
    query := `UPDATE pekerjaan_alumni SET is_delete = $1 WHERE id = $2`
    _, err := db.Exec(query, time.Now(), id)
    return err
}

func RestorePekerjaan(db *sql.DB, id int) error {
    query := `UPDATE pekerjaan_alumni SET is_delete = NULL WHERE id = $1`
    _, err := db.Exec(query, id)
    return err
}

func HardDeletePekerjaan(db *sql.DB, id int) error {
    query := `DELETE FROM pekerjaan_alumni WHERE id = $1`
    _, err := db.Exec(query, id)
    return err
}



func GetPekerjaanWithDeletedByID(db *sql.DB, id int) (*model.PekerjaanAlumni, error) {
    p := new(model.PekerjaanAlumni)
    query := `SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
                     lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
                     status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at, is_delete
              FROM pekerjaan_alumni WHERE id = $1`
    err := db.QueryRow(query, id).Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan,
        &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange,
        &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan,
        &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt, &p.IsDeleted)
    if err != nil {
        return nil, err
    }
    return p, nil
}





// GetPekerjaanRepo -> ambil data pekerjaan alumni dari DB dengan pagination, sorting, dan search
func GetPekerjaanRepo(db *sql.DB, search, sortBy, order string, limit, offset int) ([]model.PekerjaanAlumni, error) {
	query := fmt.Sprintf(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at, is_delete
		FROM pekerjaan_alumni
		WHERE (nama_perusahaan ILIKE $1 OR posisi_jabatan ILIKE $1 OR bidang_industri ILIKE $1 OR lokasi_kerja ILIKE $1) AND is_delete IS NULL
		ORDER BY %s %s
		LIMIT $2 OFFSET $3
	`, sortBy, order)

	rows, err := db.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		log.Println("Query error:", err)
		return nil, err
	}
	defer rows.Close()

	var pekerjaan []model.PekerjaanAlumni
	for rows.Next() {
		var p model.PekerjaanAlumni
		err := rows.Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt, &p.IsDeleted)
		if err != nil {
			return nil, err
		}
		pekerjaan = append(pekerjaan, p)
	}
	return pekerjaan, nil
}

// GetDeletedPekerjaanRepo -> ambil data pekerjaan alumni yang sudah dihapus dengan pagination
func GetDeletedPekerjaanRepo(db *sql.DB, offset, limit int) ([]model.PekerjaanAlumni, int, error) {
	// Query untuk mengambil data yang sudah dihapus
	query := `SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at, is_delete
		FROM pekerjaan_alumni
		WHERE is_delete IS NOT NULL
		ORDER BY is_delete DESC
		LIMIT $1 OFFSET $2`

	rows, err := db.Query(query, limit, offset)
	if err != nil {
		log.Println("Query error:", err)
		return nil, 0, err
	}
	defer rows.Close()

	var pekerjaan []model.PekerjaanAlumni
	for rows.Next() {
		var p model.PekerjaanAlumni
		err := rows.Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt, &p.IsDeleted)
		if err != nil {
			return nil, 0, err
		}
		pekerjaan = append(pekerjaan, p)
	}

	// Query untuk menghitung total data yang sudah dihapus
	countQuery := `SELECT COUNT(*) FROM pekerjaan_alumni WHERE is_delete IS NOT NULL`
	var total int
	err = db.QueryRow(countQuery).Scan(&total)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}

	return pekerjaan, total, nil
}

// CountPekerjaanRepo -> hitung total data untuk pagination
func CountPekerjaanRepo(db *sql.DB, search string) (int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM pekerjaan_alumni WHERE (nama_perusahaan ILIKE $1 OR posisi_jabatan ILIKE $1 OR bidang_industri ILIKE $1 OR lokasi_kerja ILIKE $1) AND is_delete IS NULL`
	err := db.QueryRow(countQuery, "%"+search+"%").Scan(&total)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return total, nil
}

