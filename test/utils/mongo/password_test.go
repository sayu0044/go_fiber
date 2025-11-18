package mongo_test

import (
	"testing"

	utils "go-fiber/utils/mongo"
)

func TestHashAndCheckPassword(t *testing.T) {
	plain := "SuperSecretPassword123!"

	hash, err := utils.HashPassword(plain)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if hash == "" || hash == plain {
		t.Fatalf("Hash should not be empty or equal to plain text")
	}

	if ok := utils.CheckPassword(plain, hash); !ok {
		t.Fatalf("CheckPassword should return true for correct password")
	}
	if ok := utils.CheckPassword("wrong-password", hash); ok {
		t.Fatalf("CheckPassword should return false for incorrect password")
	}
}


