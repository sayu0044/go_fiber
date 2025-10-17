package service

import (
    "database/sql"
    "go-fiber/app/model"
    "go-fiber/app/repository"
    "github.com/gofiber/fiber/v2"
    "strconv"
)

func CreateRoleService(c *fiber.Ctx, db *sql.DB) error {
    var req model.CreateRoleRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Request tidak valid", "success": false})
    }
    if req.Name == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Nama role wajib diisi", "success": false})
    }
    role, err := repository.CreateRole(db, &req)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal membuat role", "success": false})
    }
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Berhasil membuat role", "success": true, "data": role})
}

func GetRoleByIDService(c *fiber.Ctx, db *sql.DB) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ID tidak valid", "success": false})
    }
    role, err := repository.GetRoleByID(db, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Role tidak ditemukan", "success": false})
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal mengambil role", "success": false})
    }
    return c.JSON(fiber.Map{"success": true, "data": role})
}

func ListRolesService(c *fiber.Ctx, db *sql.DB) error {
    roles, err := repository.ListRoles(db)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal mengambil roles", "success": false})
    }
    return c.JSON(fiber.Map{"success": true, "data": roles})
}

func UpdateRoleService(c *fiber.Ctx, db *sql.DB) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ID tidak valid", "success": false})
    }
    var req model.UpdateRoleRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Request tidak valid", "success": false})
    }
    role, err := repository.UpdateRole(db, id, &req)
    if err != nil {
        if err == sql.ErrNoRows {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Role tidak ditemukan", "success": false})
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal mengupdate role", "success": false})
    }
    return c.JSON(fiber.Map{"message": "Berhasil mengupdate role", "success": true, "data": role})
}

func DeleteRoleService(c *fiber.Ctx, db *sql.DB) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ID tidak valid", "success": false})
    }
    if err := repository.DeleteRole(db, id); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal menghapus role", "success": false})
    }
    return c.JSON(fiber.Map{"message": "Berhasil menghapus role", "success": true})
}


