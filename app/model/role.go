package model

type Role struct {
    ID   int    `json:"id" db:"id"`
    Name string `json:"name" db:"name"`
}

type CreateRoleRequest struct {
    Name string `json:"name" validate:"required"`
}

type UpdateRoleRequest struct {
    Name *string `json:"name"`
}


