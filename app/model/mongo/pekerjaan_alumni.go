package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PekerjaanAlumni struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	AlumniID            primitive.ObjectID `bson:"alumni_id" json:"alumni_id"`
	NamaPerusahaan      string             `bson:"nama_perusahaan" json:"nama_perusahaan"`
	PosisiJabatan       string             `bson:"posisi_jabatan" json:"posisi_jabatan"`
	BidangIndustri      string             `bson:"bidang_industri" json:"bidang_industri"`
	LokasiKerja         string             `bson:"lokasi_kerja" json:"lokasi_kerja"`
	GajiRange           *string            `bson:"gaji_range,omitempty" json:"gaji_range,omitempty"`
	TanggalMulaiKerja   time.Time          `bson:"tanggal_mulai_kerja" json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja *time.Time         `bson:"tanggal_selesai_kerja,omitempty" json:"tanggal_selesai_kerja,omitempty"`
	StatusPekerjaan     string             `bson:"status_pekerjaan" json:"status_pekerjaan"`
	DeskripsiPekerjaan  *string            `bson:"deskripsi_pekerjaan,omitempty" json:"deskripsi_pekerjaan,omitempty"`
	CreatedAt           time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt           time.Time          `bson:"updated_at" json:"updated_at"`
}

// Service Layer Request (tanggal sebagai string)
type CreatePekerjaanAlumniRequest struct {
	AlumniID            string  `json:"alumni_id" validate:"required"`
	NamaPerusahaan      string  `json:"nama_perusahaan" validate:"required"`
	PosisiJabatan       string  `json:"posisi_jabatan" validate:"required"`
	BidangIndustri      string  `json:"bidang_industri" validate:"required"`
	LokasiKerja         string  `json:"lokasi_kerja" validate:"required"`
	GajiRange           *string `json:"gaji_range,omitempty"`
	TanggalMulaiKerja   string  `json:"tanggal_mulai_kerja" validate:"required"`
	TanggalSelesaiKerja *string `json:"tanggal_selesai_kerja,omitempty"`
	StatusPekerjaan     string  `json:"status_pekerjaan" validate:"required,oneof=aktif selesai resigned"`
	DeskripsiPekerjaan  *string `json:"deskripsi_pekerjaan,omitempty"`
}

// Repository Layer Request (tanggal sebagai time.Time)
type CreatePekerjaanAlumniRepositoryRequest struct {
	AlumniID            primitive.ObjectID `bson:"alumni_id"`
	NamaPerusahaan      string             `bson:"nama_perusahaan"`
	PosisiJabatan       string             `bson:"posisi_jabatan"`
	BidangIndustri      string             `bson:"bidang_industri"`
	LokasiKerja         string             `bson:"lokasi_kerja"`
	GajiRange           *string            `bson:"gaji_range,omitempty"`
	TanggalMulaiKerja   time.Time          `bson:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja *time.Time         `bson:"tanggal_selesai_kerja,omitempty"`
	StatusPekerjaan     string             `bson:"status_pekerjaan"`
	DeskripsiPekerjaan  *string            `bson:"deskripsi_pekerjaan,omitempty"`
}

// Service Layer Request (tanggal sebagai string)
type UpdatePekerjaanAlumniRequest struct {
	NamaPerusahaan      string  `json:"nama_perusahaan" validate:"required"`
	PosisiJabatan       string  `json:"posisi_jabatan" validate:"required"`
	BidangIndustri      string  `json:"bidang_industri" validate:"required"`
	LokasiKerja         string  `json:"lokasi_kerja" validate:"required"`
	GajiRange           *string `json:"gaji_range,omitempty"`
	TanggalMulaiKerja   string  `json:"tanggal_mulai_kerja" validate:"required"`
	TanggalSelesaiKerja *string `json:"tanggal_selesai_kerja,omitempty"`
	StatusPekerjaan     string  `json:"status_pekerjaan" validate:"required,oneof=aktif selesai resigned"`
	DeskripsiPekerjaan  *string `json:"deskripsi_pekerjaan,omitempty"`
}

// Repository Layer Request (tanggal sebagai time.Time)
type UpdatePekerjaanAlumniRepositoryRequest struct {
	NamaPerusahaan      string     `bson:"nama_perusahaan"`
	PosisiJabatan       string     `bson:"posisi_jabatan"`
	BidangIndustri      string     `bson:"bidang_industri"`
	LokasiKerja         string     `bson:"lokasi_kerja"`
	GajiRange           *string    `bson:"gaji_range,omitempty"`
	TanggalMulaiKerja   time.Time  `bson:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja *time.Time `bson:"tanggal_selesai_kerja,omitempty"`
	StatusPekerjaan     string     `bson:"status_pekerjaan"`
	DeskripsiPekerjaan  *string    `bson:"deskripsi_pekerjaan,omitempty"`
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
