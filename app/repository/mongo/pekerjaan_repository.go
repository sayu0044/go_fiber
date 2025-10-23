package mongo

import (
	"context"
	"log"
	"strings"
	"time"

	"go-fiber/app/model/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDB "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Pekerjaan Alumni Repository Functions

func GetAllPekerjaan(db *mongoDB.Database) ([]mongo.PekerjaanAlumni, error) {
	collection := db.Collection("pekerjaan_alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}
	opts := options.Find().SetSort(bson.M{"created_at": -1})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pekerjaan []mongo.PekerjaanAlumni
	if err = cursor.All(ctx, &pekerjaan); err != nil {
		return nil, err
	}
	return pekerjaan, nil
}

func GetPekerjaanByID(db *mongoDB.Database, id string) (*mongo.PekerjaanAlumni, error) {
	collection := db.Collection("pekerjaan_alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var pekerjaan mongo.PekerjaanAlumni
	filter := bson.M{"_id": objID}
	err = collection.FindOne(ctx, filter).Decode(&pekerjaan)
	if err != nil {
		if err == mongoDB.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &pekerjaan, nil
}

func GetPekerjaanByAlumniID(db *mongoDB.Database, alumniID string) ([]mongo.PekerjaanAlumni, error) {
	collection := db.Collection("pekerjaan_alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(alumniID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"alumni_id": objID}
	opts := options.Find().SetSort(bson.M{"tanggal_mulai_kerja": -1})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pekerjaan []mongo.PekerjaanAlumni
	if err = cursor.All(ctx, &pekerjaan); err != nil {
		return nil, err
	}
	return pekerjaan, nil
}

func CreatePekerjaan(db *mongoDB.Database, req *mongo.CreatePekerjaanAlumniRepositoryRequest) (*mongo.PekerjaanAlumni, error) {
	collection := db.Collection("pekerjaan_alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()
	pekerjaan := &mongo.PekerjaanAlumni{
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
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	result, err := collection.InsertOne(ctx, pekerjaan)
	if err != nil {
		return nil, err
	}

	pekerjaan.ID = result.InsertedID.(primitive.ObjectID)
	return pekerjaan, nil
}

func UpdatePekerjaan(db *mongoDB.Database, id string, req *mongo.UpdatePekerjaanAlumniRepositoryRequest) (*mongo.PekerjaanAlumni, error) {
	collection := db.Collection("pekerjaan_alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			"nama_perusahaan":       req.NamaPerusahaan,
			"posisi_jabatan":        req.PosisiJabatan,
			"bidang_industri":       req.BidangIndustri,
			"lokasi_kerja":          req.LokasiKerja,
			"gaji_range":            req.GajiRange,
			"tanggal_mulai_kerja":   req.TanggalMulaiKerja,
			"tanggal_selesai_kerja": req.TanggalSelesaiKerja,
			"status_pekerjaan":      req.StatusPekerjaan,
			"deskripsi_pekerjaan":   req.DeskripsiPekerjaan,
			"updated_at":            time.Now(),
		},
	}

	filter := bson.M{"_id": objID}
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	// Return updated document
	var pekerjaan mongo.PekerjaanAlumni
	err = collection.FindOne(ctx, filter).Decode(&pekerjaan)
	if err != nil {
		return nil, err
	}
	return &pekerjaan, nil
}

func DeletePekerjaan(db *mongoDB.Database, id string) error {
	collection := db.Collection("pekerjaan_alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	_, err = collection.DeleteOne(ctx, filter)
	return err
}

// GetPekerjaanRepo -> ambil data pekerjaan alumni dari DB dengan pagination, sorting, dan search
func GetPekerjaanRepo(db *mongoDB.Database, search, sortBy, order string, limit, offset int) ([]mongo.PekerjaanAlumni, error) {
	collection := db.Collection("pekerjaan_alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build filter
	filter := bson.M{}
	if search != "" {
		filter["$or"] = []bson.M{
			{"nama_perusahaan": bson.M{"$regex": search, "$options": "i"}},
			{"posisi_jabatan": bson.M{"$regex": search, "$options": "i"}},
			{"bidang_industri": bson.M{"$regex": search, "$options": "i"}},
			{"lokasi_kerja": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	// Build sort
	sortOrder := 1
	if strings.ToLower(order) == "desc" {
		sortOrder = -1
	}
	sort := bson.M{sortBy: sortOrder}

	// Set options
	opts := options.Find().
		SetSort(sort).
		SetSkip(int64(offset)).
		SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		log.Println("Query error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var pekerjaan []mongo.PekerjaanAlumni
	if err = cursor.All(ctx, &pekerjaan); err != nil {
		return nil, err
	}
	return pekerjaan, nil
}

// CountPekerjaanRepo -> hitung total data untuk pagination
func CountPekerjaanRepo(db *mongoDB.Database, search string) (int, error) {
	collection := db.Collection("pekerjaan_alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}
	if search != "" {
		filter["$or"] = []bson.M{
			{"nama_perusahaan": bson.M{"$regex": search, "$options": "i"}},
			{"posisi_jabatan": bson.M{"$regex": search, "$options": "i"}},
			{"bidang_industri": bson.M{"$regex": search, "$options": "i"}},
			{"lokasi_kerja": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
