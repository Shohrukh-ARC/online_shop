package main

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/joho/godotenv"
	"example.com/online-shop/admin"
	"example.com/online-shop/database"
)

func TestUpdateAdmin(t *testing.T) {
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
		IsActive:     true,
	}

	adminID, err := database.CreateAdmin(db, a)
	if err != nil {
		t.Fatal("Failed to create administrator information:", err)
	}

	t.Cleanup(func() {
		_, err := db.Exec(
			"DELETE FROM admin WHERE adminid = $1",
			adminID,
		)

		if err != nil {
			t.Errorf("Failed to clean up test administrator:%v", err)
		}
	})

	if adminID <= 0 {
		t.Fatal("Wrong admin ID")
	}

	updatedA := admin.Admin{
		AdminID:      adminID,
		Name:         "Alexander Morgan",
		Email:        "alexander.morgan@example.com",
		PasswordHash: passwordHash,
		RoleID:       2,
		IsActive:     false,
	}

	err = database.UpdateAdmin(db, updatedA)
	if err != nil {
		t.Fatal("Failed to update administrator:", err)
	}

	updatedAdmin, err := database.GetAdministratorByID(db, adminID)
	if err != nil {
		t.Fatal("Failed to get administrator by ID:", err)
	}

	if updatedAdmin.AdminID != adminID {
		t.Errorf("Expected id %d, got %d", adminID, updatedAdmin.AdminID)
	}

	if updatedAdmin.Name != updatedA.Name {
		t.Errorf("Expected name %s, got %s", updatedA.Name, updatedAdmin.Name)
	}

	if updatedAdmin.Email != updatedA.Email {
		t.Errorf("Expected email %s, got %s", updatedA.Email, updatedAdmin.Email)
	}

	if updatedAdmin.RoleID != updatedA.RoleID {
		t.Errorf("Expected role id %d, got %d", updatedA.RoleID, updatedAdmin.RoleID)
	}

	if updatedAdmin.IsActive != updatedA.IsActive {
		t.Errorf("Expected is_active %v, got %v", updatedA.IsActive, updatedAdmin.IsActive)
	}

	if updatedAdmin.UpdatedAt.Before(updatedAdmin.CreatedAt) {
		t.Error("UpdatedAt should not be earlier than CreatedAt")
	}

	nonexistentAdmin := admin.Admin{
		AdminID:  -1,
		Name:     "Nonexistent Admin",
		Email:    "nonexistent@example.com",
		RoleID:   1,
		IsActive: true,
	}

	err = database.UpdateAdmin(db, nonexistentAdmin)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("Expected sql.ErrNoRows, got %v", err)
	}
}
