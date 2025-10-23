package mongo

import (
	"context"
	"time"

	"go-fiber/app/model/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDB "go.mongodb.org/mongo-driver/mongo"
)

func CreateRole(db *mongoDB.Database, req *mongo.CreateRoleRequest) (*mongo.Role, error) {
	collection := db.Collection("roles")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	role := &mongo.Role{
		Name: req.Name,
	}

	result, err := collection.InsertOne(ctx, role)
	if err != nil {
		return nil, err
	}

	role.ID = result.InsertedID.(primitive.ObjectID)
	return role, nil
}

func GetRoleByID(db *mongoDB.Database, id string) (*mongo.Role, error) {
	collection := db.Collection("roles")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var role mongo.Role
	filter := bson.M{"_id": objID}
	err = collection.FindOne(ctx, filter).Decode(&role)
	if err != nil {
		if err == mongoDB.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func GetRoleByName(db *mongoDB.Database, name string) (*mongo.Role, error) {
	collection := db.Collection("roles")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var role mongo.Role
	filter := bson.M{"name": name}
	err := collection.FindOne(ctx, filter).Decode(&role)
	if err != nil {
		if err == mongoDB.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func ListRoles(db *mongoDB.Database) ([]mongo.Role, error) {
	collection := db.Collection("roles")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var roles []mongo.Role
	if err = cursor.All(ctx, &roles); err != nil {
		return nil, err
	}
	return roles, nil
}

func UpdateRole(db *mongoDB.Database, id string, req *mongo.UpdateRoleRequest) (*mongo.Role, error) {
	collection := db.Collection("roles")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Build update document
	update := bson.M{}
	if req.Name != nil {
		update["name"] = *req.Name
	}

	if len(update) == 0 {
		return GetRoleByID(db, id)
	}

	filter := bson.M{"_id": objID}
	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return nil, err
	}

	// Return updated document
	var role mongo.Role
	err = collection.FindOne(ctx, filter).Decode(&role)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func DeleteRole(db *mongoDB.Database, id string) error {
	collection := db.Collection("roles")
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
