package service

import (
    "database/sql"
    "go-fiber/app/model"
    "go-fiber/utils"
    "github.com/gofiber/fiber/v2"
)

// Handler untuk login
func LoginService(c *fiber.Ctx, db *sql.DB) error {
    type loginRequest struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    var req loginRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Request body tidak valid"})
    }
    if req.Email == "" || req.Password == "" {
        return c.Status(400).JSON(fiber.Map{"error": "Email dan password harus diisi"})
    }
    var user model.User
    var passwordHash string
    err := db.QueryRow(`
        SELECT a.id, a.email, a.password, r.name
        FROM alumni a
        JOIN roles r ON r.id = a.role_id
        WHERE a.email = $1
    `, req.Email).Scan(&user.ID, &user.Email, &passwordHash, &user.Role)
    if err != nil {
        if err == sql.ErrNoRows {
            return c.Status(401).JSON(fiber.Map{"error": "Email atau password salah"})
        }
        return c.Status(500).JSON(fiber.Map{"error": "Error database"})
    }
    if !utils.CheckPassword(req.Password, passwordHash) {
        return c.Status(401).JSON(fiber.Map{"error": "Email atau password salah"})
    }
    user.Username = user.Email
    token, err := utils.GenerateToken(user)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal generate token"})
    }
    response := model.LoginResponse{User: user, Token: token}
    return c.JSON(fiber.Map{"success": true, "message": "Login berhasil", "data": response})
}

// Handler untuk melihat profile user yang sedang login
func GetProfileService(c *fiber.Ctx, db *sql.DB) error {
	userID := c.Locals("user_id").(int)
	username := c.Locals("username").(string)
	role := c.Locals("role").(string)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Profile berhasil diambil",
		"data": fiber.Map{
			"user_id":  userID,
			"username": username,
			"role":     role,
		},
	})
}
