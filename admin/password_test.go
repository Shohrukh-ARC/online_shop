package admin

import "testing"

func TestPasswordHashing(t *testing.T) {
	password := "MySecurePassword123!"

	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	// Test correct password
	valid, err := VerifyPassword(password, hashedPassword)
	if err != nil {
		t.Fatalf("failed to verify password: %v", err)
	}

	if !valid {
		t.Error("correct password was rejected")
	}

	// Test incorrect password
	valid, err = VerifyPassword("WrongPassword", hashedPassword)
	if err != nil {
		t.Fatalf("failed to verify password: %v", err)
	}

	if valid {
		t.Error("incorrect password was accepted")
	}
}
