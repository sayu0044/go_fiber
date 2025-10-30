package mongo

import (
	"context"
	model "go-fiber/app/model/mongo"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FileRepository interface {
	Create(ctx context.Context, file *model.File) error
	FindByAlumniAndCategory(ctx context.Context, alumniID primitive.ObjectID, category string) (*model.File, error)
	ListByAlumni(ctx context.Context, alumniID primitive.ObjectID) ([]model.File, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
}

type fileRepository struct {
	collection *mongo.Collection
}

func NewFileRepository(db *mongo.Database) FileRepository {
	return &fileRepository{collection: db.Collection("files")}
}

func (r *fileRepository) Create(ctx context.Context, file *model.File) error {
	file.UploadedAt = time.Now()
	res, err := r.collection.InsertOne(ctx, file)
	if err != nil {
		return err
	}
	file.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *fileRepository) FindByAlumniAndCategory(ctx context.Context, alumniID primitive.ObjectID, category string) (*model.File, error) {
	var f model.File
	err := r.collection.FindOne(ctx, bson.M{"alumni_id": alumniID, "category": category}).Decode(&f)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *fileRepository) ListByAlumni(ctx context.Context, alumniID primitive.ObjectID) ([]model.File, error) {
	cur, err := r.collection.Find(ctx, bson.M{"alumni_id": alumniID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var files []model.File
	if err := cur.All(ctx, &files); err != nil {
		return nil, err
	}
	return files, nil
}

func (r *fileRepository) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
