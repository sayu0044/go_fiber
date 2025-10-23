package postgre

import (
	"database/sql"
	"fmt"
	model "go-fiber/app/model/postgre"
	"strings"
)

func CreateRole(db *sql.DB, req *model.CreateRoleRequest) (*model.Role, error) {
	var role model.Role
	err := db.QueryRow(`INSERT INTO roles (name) VALUES ($1) RETURNING id, name`, req.Name).Scan(&role.ID, &role.Name)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func GetRoleByID(db *sql.DB, id int) (*model.Role, error) {
	var role model.Role
	err := db.QueryRow(`SELECT id, name FROM roles WHERE id = $1`, id).Scan(&role.ID, &role.Name)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func ListRoles(db *sql.DB) ([]model.Role, error) {
	rows, err := db.Query(`SELECT id, name FROM roles ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var roles []model.Role
	for rows.Next() {
		var r model.Role
		if err := rows.Scan(&r.ID, &r.Name); err != nil {
			return nil, err
		}
		roles = append(roles, r)
	}
	return roles, nil
}

func UpdateRole(db *sql.DB, id int, req *model.UpdateRoleRequest) (*model.Role, error) {
	setParts := []string{}
	args := []interface{}{}
	idx := 1
	if req.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", idx))
		args = append(args, *req.Name)
		idx++
	}
	if len(setParts) == 0 {
		return GetRoleByID(db, id)
	}
	args = append(args, id)
	query := "UPDATE roles SET " + strings.Join(setParts, ", ") + fmt.Sprintf(" WHERE id = $%d RETURNING id, name", idx)
	var role model.Role
	if err := db.QueryRow(query, args...).Scan(&role.ID, &role.Name); err != nil {
		return nil, err
	}
	return &role, nil
}

func DeleteRole(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM roles WHERE id = $1`, id)
	return err
}
