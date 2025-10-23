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

// Alumni Repository Functions

func GetAllAlumni(db *mongoDB.Database) ([]mongo.Alumni, error) {
	collection := db.Collection("alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"nama": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var alumni []mongo.Alumni
	if err = cursor.All(ctx, &alumni); err != nil {
		return nil, err
	}
	return alumni, nil
}

// GetAlumniRepo -> ambil data alumni dari DB dengan pagination, sorting, dan search
func GetAlumniRepo(db *mongoDB.Database, search, sortBy, order string, limit, offset int) ([]mongo.Alumni, error) {
	// Debug logging
	log.Printf("=== GetAlumniRepo Debug ===")
	log.Printf("Search parameters - search: '%s' (len: %d), sortBy: '%s', order: '%s', limit: %d, offset: %d", search, len(search), sortBy, order, limit, offset)
	log.Printf("Search is empty: %t", search == "")

	collection := db.Collection("alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build filter
	filter := bson.M{}
	if search != "" {
		searchPattern := bson.M{
			"$or": []bson.M{
				{"nama": bson.M{"$regex": search, "$options": "i"}},
				{"email": bson.M{"$regex": search, "$options": "i"}},
				{"jurusan": bson.M{"$regex": search, "$options": "i"}},
				{"nim": bson.M{"$regex": search, "$options": "i"}},
			},
		}
		filter = searchPattern
		log.Printf("Using query with search filter - pattern: '%s'", search)
	} else {
		log.Printf("Using query without search filter")
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

	log.Printf("Final filter: %v", filter)
	log.Printf("Sort: %v", sort)

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		log.Printf("Query error: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var alumni []mongo.Alumni
	if err = cursor.All(ctx, &alumni); err != nil {
		return nil, err
	}
	return alumni, nil
}

// CountAlumniRepo -> hitung total data untuk pagination
func CountAlumniRepo(db *mongoDB.Database, search string) (int, error) {
	collection := db.Collection("alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}
	if search != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"nama": bson.M{"$regex": search, "$options": "i"}},
				{"email": bson.M{"$regex": search, "$options": "i"}},
				{"jurusan": bson.M{"$regex": search, "$options": "i"}},
				{"nim": bson.M{"$regex": search, "$options": "i"}},
			},
		}
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func GetAlumniByID(db *mongoDB.Database, id string) (*mongo.Alumni, error) {
	collection := db.Collection("alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var alumni mongo.Alumni
	filter := bson.M{"_id": objID}
	err = collection.FindOne(ctx, filter).Decode(&alumni)
	if err != nil {
		if err == mongoDB.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &alumni, nil
}

func GetAlumniByEmail(db *mongoDB.Database, email string) (*mongo.Alumni, error) {
	collection := db.Collection("alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var alumni mongo.Alumni
	filter := bson.M{"email": email}
	err := collection.FindOne(ctx, filter).Decode(&alumni)
	if err != nil {
		if err == mongoDB.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &alumni, nil
}

func CreateAlumni(db *mongoDB.Database, req *mongo.CreateAlumniRepositoryRequest) (*mongo.Alumni, error) {
	collection := db.Collection("alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()
	alumni := &mongo.Alumni{
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
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	result, err := collection.InsertOne(ctx, alumni)
	if err != nil {
		return nil, err
	}

	alumni.ID = result.InsertedID.(primitive.ObjectID)
	return alumni, nil
}

func UpdateAlumni(db *mongoDB.Database, id string, req *mongo.UpdateAlumniRepositoryRequest) (*mongo.Alumni, error) {
	collection := db.Collection("alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Build update document
	update := bson.M{"$set": bson.M{"updated_at": time.Now()}}

	if req.NIM != nil {
		update["$set"].(bson.M)["nim"] = *req.NIM
	}
	if req.Nama != nil {
		update["$set"].(bson.M)["nama"] = *req.Nama
	}
	if req.Jurusan != nil {
		update["$set"].(bson.M)["jurusan"] = *req.Jurusan
	}
	if req.Angkatan != nil {
		update["$set"].(bson.M)["angkatan"] = *req.Angkatan
	}
	if req.TahunLulus != nil {
		update["$set"].(bson.M)["tahun_lulus"] = *req.TahunLulus
	}
	if req.Email != nil {
		update["$set"].(bson.M)["email"] = *req.Email
	}
	if req.RoleID != nil {
		update["$set"].(bson.M)["role_id"] = *req.RoleID
	}
	if req.Password != nil {
		update["$set"].(bson.M)["password"] = *req.Password
	}
	if req.NoTelepon != nil {
		update["$set"].(bson.M)["no_telepon"] = *req.NoTelepon
	}
	if req.Alamat != nil {
		update["$set"].(bson.M)["alamat"] = *req.Alamat
	}

	filter := bson.M{"_id": objID}
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	// Return updated document
	var alumni mongo.Alumni
	err = collection.FindOne(ctx, filter).Decode(&alumni)
	if err != nil {
		return nil, err
	}
	return &alumni, nil
}

func DeleteAlumni(db *mongoDB.Database, id string) error {
	collection := db.Collection("alumni")
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

// Legacy function for backward compatibility
func CheckAlumniByNim(db *mongoDB.Database, nim string) (*mongo.Alumni, error) {
	collection := db.Collection("alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var alumni mongo.Alumni
	filter := bson.M{"nim": nim}
	err := collection.FindOne(ctx, filter).Decode(&alumni)
	if err != nil {
		if err == mongoDB.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &alumni, nil
}

// Get Alumni Employment Status with filtering and pagination
func GetAlumniEmploymentStatus(db *mongoDB.Database, req *mongo.AlumniEmploymentStatusRequest) ([]mongo.AlumniEmploymentStatus, error) {
	// Set default pagination
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}

	offset := (req.Page - 1) * req.Limit

	// Build match stage for filtering
	matchStage := bson.M{}
	if req.ID != nil {
		objID, err := primitive.ObjectIDFromHex(*req.ID)
		if err == nil {
			matchStage["_id"] = objID
		}
	}
	if req.Nama != nil {
		matchStage["nama"] = bson.M{"$regex": *req.Nama, "$options": "i"}
	}
	if req.Jurusan != nil {
		matchStage["jurusan"] = bson.M{"$regex": *req.Jurusan, "$options": "i"}
	}
	if req.Angkatan != nil {
		matchStage["angkatan"] = *req.Angkatan
	}

	// Build aggregation pipeline
	pipeline := []bson.M{
		{"$match": matchStage},
		{"$lookup": bson.M{
			"from":         "pekerjaan_alumni",
			"localField":   "_id",
			"foreignField": "alumni_id",
			"as":           "pekerjaan",
		}},
		{"$unwind": bson.M{"path": "$pekerjaan", "preserveNullAndEmptyArrays": true}},
		{"$match": bson.M{"pekerjaan.is_delete": nil}},
		{"$group": bson.M{
			"_id":              "$_id",
			"nama":             bson.M{"$first": "$nama"},
			"jurusan":          bson.M{"$first": "$jurusan"},
			"angkatan":         bson.M{"$first": "$angkatan"},
			"latest_pekerjaan": bson.M{"$last": "$pekerjaan"},
			"employment_count": bson.M{"$sum": 1},
		}},
		{"$project": bson.M{
			"_id":                 1,
			"nama":                1,
			"jurusan":             1,
			"angkatan":            1,
			"bidang_industri":     "$latest_pekerjaan.bidang_industri",
			"nama_perusahaan":     "$latest_pekerjaan.nama_perusahaan",
			"posisi_jabatan":      "$latest_pekerjaan.posisi_jabatan",
			"tanggal_mulai_kerja": "$latest_pekerjaan.tanggal_mulai_kerja",
			"gaji_range":          "$latest_pekerjaan.gaji_range",
			"lebih_dari_1_tahun": bson.M{
				"$cond": bson.M{
					"if": bson.M{
						"$lte": []interface{}{"$latest_pekerjaan.tanggal_mulai_kerja", bson.M{"$subtract": []interface{}{"$$NOW", 365 * 24 * 60 * 60 * 1000}}},
					},
					"then": 1,
					"else": 0,
				},
			},
			"employment_count": 1,
		}},
		{"$sort": bson.M{"nama": 1}},
		{"$skip": offset},
		{"$limit": req.Limit},
	}

	collection := db.Collection("alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []mongo.AlumniEmploymentStatus
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
