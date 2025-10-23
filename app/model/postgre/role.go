package postgre

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
