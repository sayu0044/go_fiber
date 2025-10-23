package postgre

import (
	"database/sql"
	model "go-fiber/app/model/postgre"
	utils "go-fiber/utils/postgre"

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
		return c.Status(500).JSON(model.LoginResponse{
			Success: false,
			Message: "Gagal generate token",
			Data:    model.LoginData{},
		})
	}

	return c.JSON(model.LoginResponse{
		Success: true,
		Message: "Login berhasil",
		Data: model.LoginData{
			User:  user,
			Token: token,
		},
	})
}

// Handler untuk melihat profile user yang sedang login
func GetProfileService(c *fiber.Ctx, db *sql.DB) error {
	userID := c.Locals("user_id").(int)
	username := c.Locals("username").(string)
	role := c.Locals("role").(string)

	return c.JSON(model.GetProfileResponse{
		Success: true,
		Message: "Profile berhasil diambil",
		Data: model.ProfileData{
			UserID:   userID,
			Username: username,
			Role:     role,
		},
	})
}
