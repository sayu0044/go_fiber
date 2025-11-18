package mongo_test

import (
	"testing"

	model "go-fiber/app/model/mongo"
	utils "go-fiber/utils/mongo"
)

func TestGenerateAndValidateToken(t *testing.T) {
	user := model.User{
		Username: "alice",
		Email:    "alice@example.com",
		Role:     "admin",
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken error: %v", err)
	}
	if token == "" {
		t.Fatalf("token should not be empty")
	}

	claims, err := utils.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken error: %v", err)
	}
	if claims.Username != user.Username || claims.Role != user.Role {
		t.Fatalf("claims mismatch, got username=%s role=%s", claims.Username, claims.Role)
	}
}


