package mongo

import (
	"go-fiber/app/model/mongo"
	repository "go-fiber/app/repository/mongo"

	"github.com/gofiber/fiber/v2"
	mongoDB "go.mongodb.org/mongo-driver/mongo"
)

func CreateRoleService(c *fiber.Ctx, db *mongoDB.Database) error {
	var req mongo.CreateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(mongo.CreateRoleResponse{
			Success: false,
			Message: "Request tidak valid",
			Data:    mongo.Role{},
		})
	}
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(mongo.CreateRoleResponse{
			Success: false,
			Message: "Nama role wajib diisi",
			Data:    mongo.Role{},
		})
	}
	role, err := repository.CreateRole(db, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(mongo.CreateRoleResponse{
			Success: false,
			Message: "Gagal membuat role",
			Data:    mongo.Role{},
		})
	}
	return c.Status(fiber.StatusCreated).JSON(mongo.CreateRoleResponse{
		Success: true,
		Message: "Berhasil membuat role",
		Data:    *role,
	})
}

func GetRoleByIDService(c *fiber.Ctx, db *mongoDB.Database) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(mongo.GetRoleByIDResponse{
			Success: false,
			Message: "ID tidak valid",
			Data:    mongo.Role{},
		})
	}
	role, err := repository.GetRoleByID(db, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(mongo.GetRoleByIDResponse{
			Success: false,
			Message: "Gagal mengambil role",
			Data:    mongo.Role{},
		})
	}
	if role == nil {
		return c.Status(fiber.StatusNotFound).JSON(mongo.GetRoleByIDResponse{
			Success: false,
			Message: "Role tidak ditemukan",
			Data:    mongo.Role{},
		})
	}
	return c.JSON(mongo.GetRoleByIDResponse{
		Success: true,
		Message: "Berhasil mengambil role",
		Data:    *role,
	})
}

func ListRolesService(c *fiber.Ctx, db *mongoDB.Database) error {
	roles, err := repository.ListRoles(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(mongo.ListRolesResponse{
			Success: false,
			Message: "Gagal mengambil roles",
			Data:    []mongo.Role{},
		})
	}
	return c.JSON(mongo.ListRolesResponse{
		Success: true,
		Message: "Berhasil mengambil roles",
		Data:    roles,
	})
}

func UpdateRoleService(c *fiber.Ctx, db *mongoDB.Database) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(mongo.UpdateRoleResponse{
			Success: false,
			Message: "ID tidak valid",
			Data:    mongo.Role{},
		})
	}
	var req mongo.UpdateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(mongo.UpdateRoleResponse{
			Success: false,
			Message: "Request tidak valid",
			Data:    mongo.Role{},
		})
	}
	role, err := repository.UpdateRole(db, id, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(mongo.UpdateRoleResponse{
			Success: false,
			Message: "Gagal mengupdate role",
			Data:    mongo.Role{},
		})
	}
	if role == nil {
		return c.Status(fiber.StatusNotFound).JSON(mongo.UpdateRoleResponse{
			Success: false,
			Message: "Role tidak ditemukan",
			Data:    mongo.Role{},
		})
	}
	return c.JSON(mongo.UpdateRoleResponse{
		Success: true,
		Message: "Berhasil mengupdate role",
		Data:    *role,
	})
}

func DeleteRoleService(c *fiber.Ctx, db *mongoDB.Database) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(mongo.DeleteRoleResponse{
			Success: false,
			Message: "ID tidak valid",
		})
	}
	if err := repository.DeleteRole(db, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(mongo.DeleteRoleResponse{
			Success: false,
			Message: "Gagal menghapus role",
		})
	}
	return c.JSON(mongo.DeleteRoleResponse{
		Success: true,
		Message: "Berhasil menghapus role",
	})
}
