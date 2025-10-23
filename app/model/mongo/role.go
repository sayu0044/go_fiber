package mongo

import "go.mongodb.org/mongo-driver/bson/primitive"

type Role struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name string             `bson:"name" json:"name"`
}

type CreateRoleRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateRoleRequest struct {
	Name *string `json:"name,omitempty"`
}

// Response Structs
type GetRoleByIDResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    Role   `json:"data"`
}

type ListRolesResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    []Role `json:"data"`
}

type CreateRoleResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    Role   `json:"data"`
}

type UpdateRoleResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    Role   `json:"data"`
}

type DeleteRoleResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
