package postgre

import "time"

type PekerjaanAlumni struct {
	ID                  int        `json:"id"`
	AlumniID            int        `json:"alumni_id"`
	NamaPerusahaan      string     `json:"nama_perusahaan"`
	PosisiJabatan       string     `json:"posisi_jabatan"`
	BidangIndustri      string     `json:"bidang_industri"`
	LokasiKerja         string     `json:"lokasi_kerja"`
	GajiRange           *string    `json:"gaji_range"`
	TanggalMulaiKerja   time.Time  `json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja *time.Time `json:"tanggal_selesai_kerja"`
	StatusPekerjaan     string     `json:"status_pekerjaan"`
	DeskripsiPekerjaan  *string    `json:"deskripsi_pekerjaan"`
	IsDeleted           *time.Time `json:"is_delete"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

// Service Layer Request (tanggal sebagai string)
type CreatePekerjaanAlumniRequest struct {
	AlumniID            int     `json:"alumni_id" validate:"required"`
	NamaPerusahaan      string  `json:"nama_perusahaan" validate:"required"`
	PosisiJabatan       string  `json:"posisi_jabatan" validate:"required"`
	BidangIndustri      string  `json:"bidang_industri" validate:"required"`
	LokasiKerja         string  `json:"lokasi_kerja" validate:"required"`
	GajiRange           *string `json:"gaji_range"`
	TanggalMulaiKerja   string  `json:"tanggal_mulai_kerja" validate:"required"`
	TanggalSelesaiKerja *string `json:"tanggal_selesai_kerja"`
	StatusPekerjaan     string  `json:"status_pekerjaan" validate:"required,oneof=aktif selesai resigned"`
	DeskripsiPekerjaan  *string `json:"deskripsi_pekerjaan"`
}

// Repository Layer Request (tanggal sebagai time.Time)
type CreatePekerjaanAlumniRepositoryRequest struct {
	AlumniID            int        `json:"alumni_id"`
	NamaPerusahaan      string     `json:"nama_perusahaan"`
	PosisiJabatan       string     `json:"posisi_jabatan"`
	BidangIndustri      string     `json:"bidang_industri"`
	LokasiKerja         string     `json:"lokasi_kerja"`
	GajiRange           *string    `json:"gaji_range"`
	TanggalMulaiKerja   time.Time  `json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja *time.Time `json:"tanggal_selesai_kerja"`
	StatusPekerjaan     string     `json:"status_pekerjaan"`
	DeskripsiPekerjaan  *string    `json:"deskripsi_pekerjaan"`
}

// Service Layer Request (tanggal sebagai string)
type UpdatePekerjaanAlumniRequest struct {
	NamaPerusahaan      string  `json:"nama_perusahaan" validate:"required"`
	PosisiJabatan       string  `json:"posisi_jabatan" validate:"required"`
	BidangIndustri      string  `json:"bidang_industri" validate:"required"`
	LokasiKerja         string  `json:"lokasi_kerja" validate:"required"`
	GajiRange           *string `json:"gaji_range"`
	TanggalMulaiKerja   string  `json:"tanggal_mulai_kerja" validate:"required"`
	TanggalSelesaiKerja *string `json:"tanggal_selesai_kerja"`
	StatusPekerjaan     string  `json:"status_pekerjaan" validate:"required,oneof=aktif selesai resigned"`
	DeskripsiPekerjaan  *string `json:"deskripsi_pekerjaan"`
}

// Repository Layer Request (tanggal sebagai time.Time)
type UpdatePekerjaanAlumniRepositoryRequest struct {
	NamaPerusahaan      string     `json:"nama_perusahaan"`
	PosisiJabatan       string     `json:"posisi_jabatan"`
	BidangIndustri      string     `json:"bidang_industri"`
	LokasiKerja         string     `json:"lokasi_kerja"`
	GajiRange           *string    `json:"gaji_range"`
	TanggalMulaiKerja   time.Time  `json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja *time.Time `json:"tanggal_selesai_kerja"`
	StatusPekerjaan     string     `json:"status_pekerjaan"`
	DeskripsiPekerjaan  *string    `json:"deskripsi_pekerjaan"`
}

// Response Structs
type GetPekerjaanAlumniByIDResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    PekerjaanAlumni `json:"data"`
}

type GetPekerjaanAlumniByAlumniIDResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    []PekerjaanAlumni `json:"data"`
}

type CreatePekerjaanAlumniResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    PekerjaanAlumni `json:"data"`
}

type UpdatePekerjaanAlumniResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    PekerjaanAlumni `json:"data"`
}

type SoftDeletePekerjaanAlumniResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type HardDeletePekerjaanAlumniResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type GetSoftDeletedPekerjaanAlumniResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    []PekerjaanAlumni `json:"data"`
}

type RestorePekerjaanAlumniResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
