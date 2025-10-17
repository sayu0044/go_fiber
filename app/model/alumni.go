package model

import "time"

type Alumni struct {
	ID          int       `json:"id" db:"id"`
	NIM         string    `json:"nim" db:"nim"`
	Nama        string    `json:"nama" db:"nama"`
	Jurusan     string    `json:"jurusan" db:"jurusan"`
	Angkatan    int       `json:"angkatan" db:"angkatan"`
	TahunLulus  int       `json:"tahun_lulus" db:"tahun_lulus"`
	Email       string    `json:"email" db:"email"`
    RoleID      int       `json:"role_id" db:"role_id"`
	NoTelepon   *string   `json:"no_telepon" db:"no_telepon"`
	Alamat      *string   `json:"alamat" db:"alamat"`
    Password    string    `json:"-" db:"password"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	
	
}

// Request/Response DTOs
type CreateAlumniRequest struct {
	NIM       string  `json:"nim" validate:"required"`
	Nama      string  `json:"nama" validate:"required"`
	Jurusan   string  `json:"jurusan" validate:"required"`
	Angkatan  int     `json:"angkatan" validate:"required"`
	TahunLulus int    `json:"tahun_lulus" validate:"required"`
	Email     string  `json:"email" validate:"required,email"`
    Password  string  `json:"password" validate:"required"`
    RoleID    int     `json:"role_id" validate:"required"`
	NoTelepon *string `json:"no_telepon"`
	Alamat    *string `json:"alamat"`
}

type UpdateAlumniRequest struct {
	NIM       *string `json:"nim"`
	Nama      *string `json:"nama"`
	Jurusan   *string `json:"jurusan"`
	Angkatan  *int    `json:"angkatan"`
	TahunLulus *int   `json:"tahun_lulus"`
	Email     *string `json:"email"`
    Password  *string `json:"password"`
    RoleID    *int    `json:"role_id"`
	NoTelepon *string `json:"no_telepon"`
	Alamat    *string `json:"alamat"`
}

// Alumni Employment Status Response
type AlumniEmploymentStatus struct {
	ID                int       `json:"id" db:"id"`
	Nama              string    `json:"nama" db:"nama"`
	Jurusan           string    `json:"jurusan" db:"jurusan"`
	Angkatan          int       `json:"angkatan" db:"angkatan"`
	BidangIndustri    *string   `json:"bidang_industri" db:"bidang_industri"`
	NamaPerusahaan    *string   `json:"nama_perusahaan" db:"nama_perusahaan"`
	PosisiJabatan     *string   `json:"posisi_jabatan" db:"posisi_jabatan"`
	TanggalMulaiKerja *time.Time `json:"tanggal_mulai_kerja" db:"tanggal_mulai_kerja"`
	GajiRange         *string   `json:"gaji_range" db:"gaji_range"`
	LebihDari1Tahun   int       `json:"lebih_dari_1_tahun" db:"lebih_dari_1_tahun"`
	EmploymentCount   int       `json:"employment_count" db:"employment_count"`
}

// Request for filtering alumni employment status
type AlumniEmploymentStatusRequest struct {
	ID                *int    `json:"id" query:"id"`
	Nama              *string `json:"nama" query:"nama"`
	Jurusan           *string `json:"jurusan" query:"jurusan"`
	Angkatan          *int    `json:"angkatan" query:"angkatan"`
	BidangIndustri    *string `json:"bidang_industri" query:"bidang_industri"`
	NamaPerusahaan    *string `json:"nama_perusahaan" query:"nama_perusahaan"`
	PosisiJabatan     *string `json:"posisi_jabatan" query:"posisi_jabatan"`
	LebihDari1Tahun   *int    `json:"lebih_dari_1_tahun" query:"lebih_dari_1_tahun"`
	Page              int     `json:"page" query:"page"`
	Limit             int     `json:"limit" query:"limit"`
}
