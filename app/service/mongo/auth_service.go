package mongo

import (
	"go-fiber/app/model/mongo"
	repository "go-fiber/app/repository/mongo"
	utils "go-fiber/utils/mongo"

	"github.com/gofiber/fiber/v2"
	mongoDB "go.mongodb.org/mongo-driver/mongo"
)

// Handler untuk login
func LoginService(c *fiber.Ctx, db *mongoDB.Database) error {
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

	alumni, err := repository.GetAlumniByEmail(db, req.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error database"})
	}
	if alumni == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Email atau password salah"})
	}

	if !utils.CheckPassword(req.Password, alumni.Password) {
		return c.Status(401).JSON(fiber.Map{"error": "Email atau password salah"})
	}

	// Get role name from roles collection
	role, err := repository.GetRoleByObjectID(db, alumni.RoleID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error fetching role"})
	}

	user := mongo.User{
		ID:       alumni.ID,
		Username: alumni.Email,
		Email:    alumni.Email,
		Role:     role.Name,
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		return c.Status(500).JSON(mongo.LoginResponse{
			Success: false,
			Message: "Gagal generate token",
			Data:    mongo.LoginData{},
		})
	}

	return c.JSON(mongo.LoginResponse{
		Success: true,
		Message: "Login berhasil",
		Data: mongo.LoginData{
			User:  user,
			Token: token,
		},
	})
}

// Handler untuk melihat profile user yang sedang login
func GetProfileService(c *fiber.Ctx, db *mongoDB.Database) error {
	userID := c.Locals("user_id").(string)
	username := c.Locals("username").(string)
	role := c.Locals("role").(string)

	return c.JSON(mongo.GetProfileResponse{
		Success: true,
		Message: "Profile berhasil diambil",
		Data: mongo.ProfileData{
			UserID:   userID,
			Username: username,
			Role:     role,
		},
	})
}
