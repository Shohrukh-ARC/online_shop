package auth

import (
	"testing"

	"example.com/online-shop/admin"
	"example.com/online-shop/database"
	"github.com/joho/godotenv"
)

func TestAuthenticationAdmin(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatal("Error loading .env file:", err)
	}

	db, err := database.ConnectDB()
	if err != nil {
		t.Fatal("Failed to connect to database:", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	passwordHash, err := admin.HashPassword("TestPassword123!")
	if err != nil {
		t.Fatal("Failed to hash password:", err)
	}

	a := admin.Admin{
		Name:         "Alex Morgan",
		Email:        "alex.morgan@example.com",
		PasswordHash: passwordHash,
		RoleID:       1,
	}

	adminID, err := database.CreateAdmin(db, a)
	if err != nil {
		t.Fatal("Failed to create administrator:", err)
	}

	t.Cleanup(func() {
		_, err := db.Exec(
			"DELETE FROM admin WHERE adminid = $1",
			adminID,
		)
		if err != nil {
			t.Errorf("Failed to clean up test administrator: %v", err)
		}
	})

	if adminID <= 0 {
		t.Fatal("Wrong admin ID")
	}

	newAdmin, err := AuthenticateAdmin(
		db,
		"alex.morgan@example.com",
		"TestPassword123!",
	)
	if err != nil {
		t.Fatal("Failed to authenticate administrator:", err)
	}

	if newAdmin.AdminID != adminID {
		t.Errorf("Expected ID %d, got %d", adminID, newAdmin.AdminID)
	}

	if newAdmin.Email != a.Email {
		t.Errorf("Expected email %s, got %s", a.Email, newAdmin.Email)
	}

	if newAdmin.PasswordHash != "" {
		t.Error("PasswordHash should be empty after authentication")
	}

	// Wrong password check
	// In this case administrator exists, but wrong password in given
	_, err = AuthenticateAdmin(db, a.Email, "WrongPassword123!")
	if err == nil {
		t.Fatal("Authentication had to get rejected")
	}

	if err.Error() != "invalid email or password" {
		t.Errorf("Expected invalid email or password, got %v", err)
	}

	//  Nonexistent password check
	_, err = AuthenticateAdmin(db, "nonexistent@example.com", "TestPassword123!")
	if err == nil {
		t.Fatal("Authenticaiton had to fail")
	}

	if err.Error() != "invalid email or password" {
		t.Errorf("Expected invalid email or password, got %v", err)
	}

	// Testing inactive administrator
	err = database.DeactivateAdmin(db, adminID)
	if err != nil {
		t.Fatal("Failed to deactivate administrator:", err)
	}

	deactivatedAdmin, err := database.GetAdministratorByID(db, adminID)
	if err != nil {
		t.Fatal("Failed to retrieve administrator:", err)
	}

	if deactivatedAdmin.IsActive {
		t.Fatal("Deactivation did not work")
	}

	_, err = AuthenticateAdmin(
		db,
		"alex.morgan@example.com",
		"TestPassword123!",
	)
	if err == nil {
		t.Fatal("Authenticaiton failed")
	}

	if err.Error() != "Admin status should be active" {
		t.Errorf("Expected inactive administrator error, got %v", err)
	}
}
