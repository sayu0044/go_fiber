package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Alumni struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	NIM        string             `bson:"nim" json:"nim"`
	Nama       string             `bson:"nama" json:"nama"`
	Jurusan    string             `bson:"jurusan" json:"jurusan"`
	Angkatan   int                `bson:"angkatan" json:"angkatan"`
	TahunLulus int                `bson:"tahun_lulus" json:"tahun_lulus"`
	Email      string             `bson:"email" json:"email"`
	RoleID     primitive.ObjectID `bson:"role_id" json:"role_id"`
	NoTelepon  *string            `bson:"no_telepon,omitempty" json:"no_telepon,omitempty"`
	Alamat     *string            `bson:"alamat,omitempty" json:"alamat,omitempty"`
	Password   string             `bson:"password" json:"-"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

// Service Layer Request DTOs
type CreateAlumniRequest struct {
	NIM        string  `json:"nim" validate:"required"`
	Nama       string  `json:"nama" validate:"required"`
	Jurusan    string  `json:"jurusan" validate:"required"`
	Angkatan   int     `json:"angkatan" validate:"required"`
	TahunLulus int     `json:"tahun_lulus" validate:"required"`
	Email      string  `json:"email" validate:"required,email"`
	Password   string  `json:"password" validate:"required"`
	RoleID     string  `json:"role_id" validate:"required"`
	NoTelepon  *string `json:"no_telepon,omitempty"`
	Alamat     *string `json:"alamat,omitempty"`
}

// Repository Layer Request
type CreateAlumniRepositoryRequest struct {
	NIM        string             `bson:"nim"`
	Nama       string             `bson:"nama"`
	Jurusan    string             `bson:"jurusan"`
	Angkatan   int                `bson:"angkatan"`
	TahunLulus int                `bson:"tahun_lulus"`
	Email      string             `bson:"email"`
	Password   string             `bson:"password"`
	RoleID     primitive.ObjectID `bson:"role_id"`
	NoTelepon  *string            `bson:"no_telepon,omitempty"`
	Alamat     *string            `bson:"alamat,omitempty"`
}

type UpdateAlumniRequest struct {
	NIM        *string `json:"nim,omitempty"`
	Nama       *string `json:"nama,omitempty"`
	Jurusan    *string `json:"jurusan,omitempty"`
	Angkatan   *int    `json:"angkatan,omitempty"`
	TahunLulus *int    `json:"tahun_lulus,omitempty"`
	Email      *string `json:"email,omitempty"`
	Password   *string `json:"password,omitempty"`
	RoleID     *string `json:"role_id,omitempty"`
	NoTelepon  *string `json:"no_telepon,omitempty"`
	Alamat     *string `json:"alamat,omitempty"`
}

// Repository Layer Request
type UpdateAlumniRepositoryRequest struct {
	NIM        *string             `bson:"nim,omitempty"`
	Nama       *string             `bson:"nama,omitempty"`
	Jurusan    *string             `bson:"jurusan,omitempty"`
	Angkatan   *int                `bson:"angkatan,omitempty"`
	TahunLulus *int                `bson:"tahun_lulus,omitempty"`
	Email      *string             `bson:"email,omitempty"`
	Password   *string             `bson:"password,omitempty"`
	RoleID     *primitive.ObjectID `bson:"role_id,omitempty"`
	NoTelepon  *string             `bson:"no_telepon,omitempty"`
	Alamat     *string             `bson:"alamat,omitempty"`
}

// Alumni Employment Status Response
type AlumniEmploymentStatus struct {
	ID                primitive.ObjectID `bson:"_id" json:"id"`
	Nama              string             `bson:"nama" json:"nama"`
	Jurusan           string             `bson:"jurusan" json:"jurusan"`
	Angkatan          int                `bson:"angkatan" json:"angkatan"`
	BidangIndustri    *string            `bson:"bidang_industri,omitempty" json:"bidang_industri,omitempty"`
	NamaPerusahaan    *string            `bson:"nama_perusahaan,omitempty" json:"nama_perusahaan,omitempty"`
	PosisiJabatan     *string            `bson:"posisi_jabatan,omitempty" json:"posisi_jabatan,omitempty"`
	TanggalMulaiKerja *time.Time         `bson:"tanggal_mulai_kerja,omitempty" json:"tanggal_mulai_kerja,omitempty"`
	GajiRange         *string            `bson:"gaji_range,omitempty" json:"gaji_range,omitempty"`
	LebihDari1Tahun   int                `bson:"lebih_dari_1_tahun" json:"lebih_dari_1_tahun"`
	EmploymentCount   int                `bson:"employment_count" json:"employment_count"`
}

// Request for filtering alumni employment status
type AlumniEmploymentStatusRequest struct {
	ID              *string `json:"id,omitempty" query:"id"`
	Nama            *string `json:"nama,omitempty" query:"nama"`
	Jurusan         *string `json:"jurusan,omitempty" query:"jurusan"`
	Angkatan        *int    `json:"angkatan,omitempty" query:"angkatan"`
	BidangIndustri  *string `json:"bidang_industri,omitempty" query:"bidang_industri"`
	NamaPerusahaan  *string `json:"nama_perusahaan,omitempty" query:"nama_perusahaan"`
	PosisiJabatan   *string `json:"posisi_jabatan,omitempty" query:"posisi_jabatan"`
	LebihDari1Tahun *int    `json:"lebih_dari_1_tahun,omitempty" query:"lebih_dari_1_tahun"`
	Page            int     `json:"page" query:"page"`
	Limit           int     `json:"limit" query:"limit"`
}

// Response Structs
type GetAlumniByIDResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    Alumni `json:"data"`
}

type CreateAlumniResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    Alumni `json:"data"`
}

type UpdateAlumniResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    Alumni `json:"data"`
}

type DeleteAlumniResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type CheckAlumniResponse struct {
	Success  bool    `json:"success"`
	Message  string  `json:"message"`
	IsAlumni bool    `json:"isAlumni"`
	Alumni   *Alumni `json:"alumni,omitempty"`
}

type GetAlumniEmploymentStatusResponse struct {
	Success    bool                     `json:"success"`
	Message    string                   `json:"message"`
	Data       []AlumniEmploymentStatus `json:"data"`
	Pagination PaginationInfo           `json:"pagination"`
}

type PaginationInfo struct {
	CurrentPage  int `json:"current_page"`
	PerPage      int `json:"per_page"`
	TotalRecords int `json:"total_records"`
	TotalPages   int `json:"total_pages"`
}
