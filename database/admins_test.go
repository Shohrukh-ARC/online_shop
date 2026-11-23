package database_test

import (
	"database/sql"
	"errors"
	"testing"

	"example.com/online-shop/admin"
	"example.com/online-shop/database"
	"github.com/joho/godotenv"
)

func TestCreateAndGetAdmins(t *testing.T) {
	// Connecting to the database
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

	// Administrator creation
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

	admin, err := database.GetAdministratorByID(db, adminID)
	if err != nil {
		t.Fatal("Failed to get administrator by ID:", err)
	}

	if admin.AdminID != adminID {
		t.Errorf("Expected id %d, got %d", adminID, admin.AdminID)
	}

	if admin.Name != a.Name {
		t.Errorf("Expected name %s, got %s", a.Name, admin.Name)
	}

	if admin.Email != a.Email {
		t.Errorf("Expected email %s, got %s", a.Email, admin.Email)
	}

	if admin.RoleID != a.RoleID {
		t.Errorf("Expected role ID %d, got %d", a.RoleID, admin.RoleID)
	}
}

func TestGetAdminByEmail(t *testing.T) {
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
		Email:        "alex.morgan.2@example.com",
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

	admin, err := database.GetAdminByEmail(db, "alex.morgan.2@example.com")
	if err != nil {
		t.Fatal("Failed to get administrator by email:", err)
	}

	if admin.AdminID != adminID {
		t.Errorf("Expected %d, got %d:", adminID, admin.AdminID)
	}

	if admin.Name != a.Name {
		t.Errorf("Expected %s, got %s:", a.Name, admin.Name)
	}

	if admin.Email != a.Email {
		t.Errorf("Expected %s, got %s:", a.Email, admin.Email)
	}

	if admin.RoleID != a.RoleID {
		t.Errorf("Expected %d, got %d:", a.RoleID, admin.RoleID)
	}

	if admin.PasswordHash != passwordHash {
		t.Errorf("Retrieved password hash does not match")
	}

	// Checking nonexistent email
	_, err = database.GetAdminByEmail(db, "nonexistent@example.com")
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("Expected sql.ErrNoRows, got %v", err)
	}
}

func TestGetAdmins(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatal("Error loading .env file")
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

	a1 := admin.Admin{
		Name:         "Alex Morgan",
		Email:        "alex.morgan_3@example.com",
		PasswordHash: passwordHash,
		RoleID:       1,
	}

	a2 := admin.Admin{
		Name:         "Alex Morgans",
		Email:        "alex.morgans_3@example.com",
		PasswordHash: passwordHash,
		RoleID:       1,
	}

	adminID1, err := database.CreateAdmin(db, a1)
	if err != nil {
		t.Fatal("Failed to create administrator:", err)
	}

	t.Cleanup(func() {
		_, err := db.Exec(
			"DELETE FROM admin WHERE adminid = $1",
			adminID1,
		)
		if err != nil {
			t.Errorf("Failed to clean up test adminstrator: %v", err)
		}
	})

	if adminID1 <= 0 {
		t.Fatal("Wrong admin ID")
	}

	adminID2, err := database.CreateAdmin(db, a2)
	if err != nil {
		t.Fatal("Failed to create administrator:", err)
	}

	t.Cleanup(func() {
		_, err := db.Exec(
			"DELETE FROM admin WHERE adminid = $1",
			adminID2,
		)
		if err != nil {
			t.Errorf("Failed to clean up test adminstrator: %v", err)
		}
	})

	if adminID2 <= 0 {
		t.Fatal("Wrong admin ID")
	}

	admins, err := database.GetAdmins(db)
	if err != nil {
		t.Fatal("Failed to get administrators:", err)
	}

	found1 := false
	found2 := false

	for _, retrievedAdmin := range admins {
		if retrievedAdmin.AdminID == adminID1 {
			found1 = true

			if retrievedAdmin.Name != a1.Name {
				t.Errorf("Expected name %s, got %s", a1.Name, retrievedAdmin.Name)
			}

			if retrievedAdmin.Email != a1.Email {
				t.Errorf("Expected email %s, got %s", a1.Email, retrievedAdmin.Email)
			}

			if retrievedAdmin.RoleID != a1.RoleID {
				t.Errorf("Expected role id %d, got %d", a1.RoleID, retrievedAdmin.RoleID)
			}

			if retrievedAdmin.PasswordHash != "" {
				t.Error("Password hash should not be returned")
			}
		}

		if retrievedAdmin.AdminID == adminID2 {
			found2 = true

			if retrievedAdmin.Name != a2.Name {
				t.Errorf("Expected name %s, got %s", a2.Name, retrievedAdmin.Name)
			}

			if retrievedAdmin.Email != a2.Email {
				t.Errorf("Expected email %s, got %s", a2.Email, retrievedAdmin.Email)
			}

			if retrievedAdmin.RoleID != a2.RoleID {
				t.Errorf("Expected role id %d, got %d", a2.RoleID, retrievedAdmin.RoleID)
			}

			if retrievedAdmin.PasswordHash != "" {
				t.Error("Password hash should not be returned")
			}
		}

		if found1 && found2 {
			break
		}
	}

	if !found1 {
		t.Errorf("Administrator with %d id was not found", adminID1)
	}

	if !found2 {
		t.Errorf("Administrator with %d id was not found", adminID2)
	}
}

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

func TestPatchAdmin(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatal("Error loading .env file:", err)
	}

	db, err := database.ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database:%v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	hashedPassword, err := admin.HashPassword("TestPassword123!")
	if err != nil {
		t.Fatalf("Failed to hash password:%v", err)
	}

	a := admin.Admin{
		Name:         "Alex Morgan",
		Email:        "alex.morgan@example.com",
		PasswordHash: hashedPassword,
		RoleID:       1,
		IsActive:     true,
	}

	adminID, err := database.CreateAdmin(db, a)
	if err != nil {
		t.Fatalf("Failed to create administrator information:%v", err)
	}

	t.Cleanup(func() {
		_, err := db.Exec(
			"DELETE FROM admin WHERE adminid = $1",
			adminID,
		)

		if err != nil {
			t.Errorf("Failed to clean up test administrator:%v", err) // Errorf is used instead of Fatalf, becuase program should continue running even if cleaning was unsuccessfull
		}
	})

	if adminID <= 0 {
		t.Fatal("Wrong admin ID")
	}

	updatedName := "Alexander Morgan"
	isActive := false

	toPatch := admin.AdminUpdate{
		Name:     &updatedName,
		IsActive: &isActive,
	}

	err = database.PatchAdmin(db, adminID, toPatch)
	if err != nil {
		t.Fatalf("Could not patch requested information:%v", err)
	}

	patchedAdmin, err := database.GetAdministratorByID(db, adminID)
	if err != nil {
		t.Fatal("Failed to get administrator by ID:", err)
	}

	if patchedAdmin.AdminID != adminID {
		t.Fatalf("Expected id %d, got %d", adminID, patchedAdmin.AdminID)
	}

	if patchedAdmin.Name != *toPatch.Name {
		t.Fatalf("Expected name %s, got %s", *toPatch.Name, patchedAdmin.Name)
	}

	if patchedAdmin.Email != a.Email {
		t.Fatalf("Expected email %s, got %s", a.Email, patchedAdmin.Email)
	}

	if patchedAdmin.RoleID != a.RoleID {
		t.Fatalf("Expected role ID %d, got %d", a.RoleID, patchedAdmin.RoleID)
	}

	if patchedAdmin.IsActive != *toPatch.IsActive {
		t.Fatalf("Expected activeness %v, got %v", *toPatch.IsActive, patchedAdmin.IsActive) // %v - for boolean
	}

	if patchedAdmin.UpdatedAt.Before(patchedAdmin.CreatedAt) {
		t.Error("UpdatedAt should not be earlier than CreatedAt")
	}

	// Empty Check
	emptyCheck := admin.AdminUpdate{}

	err = database.PatchAdmin(db, adminID, emptyCheck)
	if err != nil {
		t.Fatalf("Could not patch requested information:%v", err)
	}

	newPatchedAdmin, err := database.GetAdministratorByID(db, adminID)
	if err != nil {
		t.Fatal("Failed to get administrator by ID:", err)
	}

	if newPatchedAdmin.AdminID != adminID {
		t.Fatalf("Expected id %d, got %d", adminID, newPatchedAdmin.AdminID)
	}

	if newPatchedAdmin.Name != patchedAdmin.Name {
		t.Fatalf("Expected name %s, got %s", patchedAdmin.Name, newPatchedAdmin.Name)
	}

	if newPatchedAdmin.Email != patchedAdmin.Email {
		t.Fatalf("Expected email %s, got %s", patchedAdmin.Email, newPatchedAdmin.Email)
	}

	if newPatchedAdmin.RoleID != patchedAdmin.RoleID {
		t.Fatalf("Expected roleID %d, got %d", patchedAdmin.RoleID, newPatchedAdmin.RoleID)
	}

	if newPatchedAdmin.IsActive != patchedAdmin.IsActive {
		t.Fatalf("Expected activeness %v, got %v", patchedAdmin.IsActive, newPatchedAdmin.IsActive)
	}

	if newPatchedAdmin.UpdatedAt.Before(newPatchedAdmin.CreatedAt) {
		t.Error("UpdatedAt should not be earlier than CreatedAt")
	}

	// Nonexistent ID
	err = database.PatchAdmin(db, -1, emptyCheck)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("Expected sql.ErrNoRows, got %v", err)
	}
}

func TestDeactivateAdmin(t *testing.T) {
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

	passwordHash, err := admin.HashPassword("TestPassword123")
	if err != nil {
		t.Fatal("Failed to hash password", err)
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

	initialAdmin, err := database.GetAdministratorByID(db, adminID)
	if err != nil {
		t.Fatal("Failed to get administrator:", err)
	}

	if !initialAdmin.IsActive {
		t.Fatal("Administrator should initially be active")
	}

	err = database.DeactivateAdmin(db, adminID)
	if err != nil {
		t.Fatal("Failed to deactivate adminstrator:", err)
	}

	deactivatedAdmin, err := database.GetAdministratorByID(db, adminID)
	if err != nil {
		t.Fatal("Failed to get administrator by ID:", err)
	}

	if deactivatedAdmin.AdminID != adminID {
		t.Errorf("Expected id %d, got %d", adminID, deactivatedAdmin.AdminID)
	}

	if deactivatedAdmin.Name != a.Name {
		t.Errorf("Expected name %s, got %s", a.Name, deactivatedAdmin.Name)
	}

	if deactivatedAdmin.Email != a.Email {
		t.Errorf("Expected email %s, got %s", a.Email, deactivatedAdmin.Email)
	}

	if deactivatedAdmin.RoleID != a.RoleID {
		t.Errorf("Expected role id %d, got %d", a.RoleID, deactivatedAdmin.RoleID)
	}

	if deactivatedAdmin.IsActive {
		t.Error("Deactivated administrator should have IsActive = false")
	}

	if deactivatedAdmin.UpdatedAt.Before(deactivatedAdmin.CreatedAt) {
		t.Error("UpdatedAt should not be earlier than CreatedAt")
	}

	// Non-existent admin check
	err = database.DeactivateAdmin(db, -1)

	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("Expected sql.ErrNoRows, got %v", err)
	}
}

func TestLastSuperAdminProtection(t *testing.T) {
	err := godotenv.Overload("../.env")
	if err != nil {
		t.Fatalf("Failed to load test environment: %v", err)
	}

	db, err := database.ConnectDB()
	if err != nil {
		t.Fatal("Failed to connect to database:", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	var dbName string

	err = db.QueryRow("SELECT current_database()").Scan(&dbName)
	if err != nil {
		t.Fatalf("Failed to retrieve database name: %v", err)
	}

	if dbName != "online_shop_test" {
		t.Fatalf("Expected online_shop_test, connected to %s", dbName)
	}

	passwordHash := "$argon2id$v=19$m=19456,t=2,p=1$L3TWzxV82IT16DnW3KVgbA$8UJBpoDkOQFLRhHNkp+cNGeN6WC1sIC4wNPmzZdvG68"

	testAdmin := admin.Admin{
		Name:         "Alex Morgan",
		Email:        "alex.morgan@example.com",
		PasswordHash: passwordHash,
		RoleID:       1,
	}

	superAdminID, err := database.CreateAdmin(db, testAdmin)
	if err != nil {
		t.Fatalf("failed to create super admin: %v", err)
	}

	t.Cleanup(func() {
		_, err := db.Exec(
			"DELETE FROM admin WHERE adminid = $1",
			superAdminID,
		)
		if err != nil {
			t.Errorf("Failed to clean up test administrator: %v", err)
		}
	})

	err = database.DeactivateAdmin(db, superAdminID)
	if !errors.Is(err, database.ErrLastSuperAdmin) {
		t.Fatalf("Expected %v, got %v", database.ErrLastSuperAdmin, err)
	}

	superAdmin, err := database.GetAdministratorByID(db, superAdminID)
	if err != nil {
		t.Fatalf("failed to retrieve super admin: %v", err)
	}

	if !superAdmin.IsActive {
		t.Fatal("Administrator should be active")
	}

	var newRoleID int = 2
	updatedAdmin := admin.AdminUpdate{
		RoleID: &newRoleID,
	}

	err = database.PatchAdmin(db, superAdminID, updatedAdmin)
	if !errors.Is(err, database.ErrLastSuperAdmin) {
		t.Fatalf("Expected %v, got %v", database.ErrLastSuperAdmin, err)
	}

	superAdmin, err = database.GetAdministratorByID(db, superAdminID)
	if err != nil {
		t.Fatalf("failed to get super admin with the given id: %v", err)
	}
	if superAdmin.RoleID != 1 {
		t.Fatal("Super admin ID should not change")
	}

	var newName string = "Alex Williams"
	updatedAdmin = admin.AdminUpdate{
		Name: &newName,
	}

	err = database.PatchAdmin(db, superAdminID, updatedAdmin)
	if err != nil {
		t.Fatalf("failed to patch admin information: %v", err)
	}

	superAdmin, err = database.GetAdministratorByID(db, superAdminID)
	if err != nil {
		t.Fatalf("failed to get super admin with the given id: %v", err)
	}
	if superAdmin.Name != "Alex Williams" {
		t.Fatalf("Expected Alex Williams, got %s", superAdmin.Name)
	}

	if superAdmin.RoleID != 1 || superAdmin.IsActive == false {
		t.Fatal("Either admin is not super admin or is not active")
	}

	passwordHash_2 := "$argon2id$v=19$m=19456,t=2,p=1$3e+kv8zbPeeKXK9IVUOH0A$3JUJ4EPAakfmEXB6Ckuq5a7DEr5/kevohDCfAOhA194"

	testAdmin_2 := admin.Admin{
		Name:         "Michael Brown",
		Email:        "michael.brown@example.com",
		PasswordHash: passwordHash_2,
		RoleID:       1,
	}

	superAdminID_2, err := database.CreateAdmin(db, testAdmin_2)
	if err != nil {
		t.Fatalf("failed to create super admin: %v", err)
	}

	t.Cleanup(func() {
		_, err := db.Exec(
			"DELETE FROM admin WHERE adminid = $1",
			superAdminID_2,
		)
		if err != nil {
			t.Errorf("Failed to clean up test administrator: %v", err)
		}
	})

	err = database.DeactivateAdmin(db, superAdminID_2)
	if err != nil {
		t.Fatalf("Failed to deactivate second Super Admin: %v", err)
	}

	superAdmin_2, err := database.GetAdministratorByID(db, superAdminID_2)
	if err != nil {
		t.Fatalf("failed to retrieve super admin: %v", err)
	}

	if superAdmin_2.IsActive {
		t.Fatal("The second super admin should be inactive")
	}

	superAdmin, err = database.GetAdministratorByID(db, superAdminID)
	if err != nil {
		t.Fatalf("failed to get super admin with the given id: %v", err)
	}

	if superAdmin.RoleID != 1 || !superAdmin.IsActive {
		t.Fatal("Either admin is not super admin or is not active")
	}
}

// err := godotenv.Overload("../.env")
// 	if err != nil {
// 		t.Fatalf("Failed to load test environment: %v", err)
// 	}

// err := godotenv.Overload("../.env.test")
// if err != nil {
//     t.Fatalf("Failed to load test environment: %v", err)
// }
