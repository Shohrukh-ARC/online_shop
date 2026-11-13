package auth

import (
	"database/sql"
	"errors"

	"example.com/online-shop/admin"
	"example.com/online-shop/database"
)

func AuthenticateAdmin(
	db *sql.DB,
	email string,
	password string,
) (admin.Admin, error) {
	newAdmin, err := database.GetAdminByEmail(db, email)
	if errors.Is(err, sql.ErrNoRows) {
		return admin.Admin{}, errors.New("invalid email or password")
	}
	if err != nil {
		return admin.Admin{}, err
	}

	passwordValid, err := admin.VerifyPassword(password, newAdmin.PasswordHash)
	if err != nil {
		return admin.Admin{}, err
	}

	if !passwordValid {
		return admin.Admin{}, errors.New("invalid email or password")
	}

	if !newAdmin.IsActive {
		return admin.Admin{}, errors.New("Admin status should be active")
	}

	newAdmin.PasswordHash = ""

	return newAdmin, nil
}
