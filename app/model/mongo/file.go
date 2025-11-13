package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// File represents stored file metadata in MongoDB
type File struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	AlumniID     primitive.ObjectID `bson:"alumni_id" json:"alumni_id"`
	Category     string             `bson:"category" json:"category"` // photo | certificate
	FileName     string             `bson:"file_name" json:"file_name"`
	OriginalName string             `bson:"original_name" json:"original_name"`
	FilePath     string             `bson:"file_path" json:"file_path"`
	FileType     string             `bson:"file_type" json:"file_type"`
	FileSize     int64              `bson:"file_size" json:"file_size"`
	UploadedAt   time.Time          `bson:"uploaded_at" json:"uploaded_at"`
}

// FileResponse is a trimmed response for clients
type FileResponse struct {
	ID           string `json:"id"`
	Category     string `json:"category"`
	FileName     string `json:"file_name"`
	OriginalName string `json:"original_name"`
	FilePath     string `json:"file_path"`
	FileType     string `json:"file_type"`
	FileSize     int64  `json:"file_size"`
}

// FileUploadResponse merepresentasikan response standar untuk upload file
type FileUploadResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    FileResponse `json:"data"`
}
