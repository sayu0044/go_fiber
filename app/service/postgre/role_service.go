package postgre

import (
	"database/sql"
	model "go-fiber/app/model/postgre"
	repository "go-fiber/app/repository/postgre"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CreateRoleService(c *fiber.Ctx, db *sql.DB) error {
	var req model.CreateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.CreateRoleResponse{
			Success: false,
			Message: "Request tidak valid",
			Data:    model.Role{},
		})
	}
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.CreateRoleResponse{
			Success: false,
			Message: "Nama role wajib diisi",
			Data:    model.Role{},
		})
	}
	role, err := repository.CreateRole(db, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.CreateRoleResponse{
			Success: false,
			Message: "Gagal membuat role",
			Data:    model.Role{},
		})
	}
	return c.Status(fiber.StatusCreated).JSON(model.CreateRoleResponse{
		Success: true,
		Message: "Berhasil membuat role",
		Data:    *role,
	})
}

func GetRoleByIDService(c *fiber.Ctx, db *sql.DB) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.GetRoleByIDResponse{
			Success: false,
			Message: "ID tidak valid",
			Data:    model.Role{},
		})
	}
	role, err := repository.GetRoleByID(db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(model.GetRoleByIDResponse{
				Success: false,
				Message: "Role tidak ditemukan",
				Data:    model.Role{},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.GetRoleByIDResponse{
			Success: false,
			Message: "Gagal mengambil role",
			Data:    model.Role{},
		})
	}
	return c.JSON(model.GetRoleByIDResponse{
		Success: true,
		Message: "Berhasil mengambil role",
		Data:    *role,
	})
}

func ListRolesService(c *fiber.Ctx, db *sql.DB) error {
	roles, err := repository.ListRoles(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ListRolesResponse{
			Success: false,
			Message: "Gagal mengambil roles",
			Data:    []model.Role{},
		})
	}
	return c.JSON(model.ListRolesResponse{
		Success: true,
		Message: "Berhasil mengambil roles",
		Data:    roles,
	})
}

func UpdateRoleService(c *fiber.Ctx, db *sql.DB) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.UpdateRoleResponse{
			Success: false,
			Message: "ID tidak valid",
			Data:    model.Role{},
		})
	}
	var req model.UpdateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.UpdateRoleResponse{
			Success: false,
			Message: "Request tidak valid",
			Data:    model.Role{},
		})
	}
	role, err := repository.UpdateRole(db, id, &req)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(model.UpdateRoleResponse{
				Success: false,
				Message: "Role tidak ditemukan",
				Data:    model.Role{},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.UpdateRoleResponse{
			Success: false,
			Message: "Gagal mengupdate role",
			Data:    model.Role{},
		})
	}
	return c.JSON(model.UpdateRoleResponse{
		Success: true,
		Message: "Berhasil mengupdate role",
		Data:    *role,
	})
}

func DeleteRoleService(c *fiber.Ctx, db *sql.DB) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.DeleteRoleResponse{
			Success: false,
			Message: "ID tidak valid",
		})
	}
	if err := repository.DeleteRole(db, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.DeleteRoleResponse{
			Success: false,
			Message: "Gagal menghapus role",
		})
	}
	return c.JSON(model.DeleteRoleResponse{
		Success: true,
		Message: "Berhasil menghapus role",
	})
}
