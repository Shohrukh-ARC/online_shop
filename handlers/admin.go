package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"example.com/online-shop/admin"
	"example.com/online-shop/auth"
	"example.com/online-shop/database"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/pgconn"
)

type CreateAdminRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleID   int    `json:"role_id"`
}

type CreateAdminResponse struct {
	AdminID int    `json:"admin_id"`
	Message string `json:"message"`
}

func CreateAdminHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Wrong Method", http.StatusMethodNotAllowed)
			return
		}

		var req CreateAdminRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
			return
		}

		if len(strings.TrimSpace(req.Name)) == 0 || len(strings.TrimSpace(req.Email)) == 0 || len(strings.TrimSpace(req.Password)) == 0 {
			http.Error(w, "Name, email, and password cannot be empty", http.StatusBadRequest)
			return
		}

		email, err := mail.ParseAddress(req.Email)
		if err != nil {
			http.Error(w, "Invalid email format", http.StatusBadRequest)
			return
		}

		if email.Address != req.Email {
			http.Error(w, "Invalid email format", http.StatusBadRequest)
			return
		}

		if req.RoleID <= 0 || req.RoleID > 3 {
			http.Error(w, "Invalid role ID", http.StatusBadRequest)
			return
		}

		hashedPassword, err := admin.HashPassword(req.Password)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		newAdmin := admin.Admin{
			Name:         req.Name,
			Email:        req.Email,
			PasswordHash: hashedPassword,
			RoleID:       req.RoleID,
		}

		adminID, err := database.CreateAdmin(db, newAdmin)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // 23505 - unique constraint violation
			http.Error(w, "Cannot use existing email", http.StatusConflict)
			return
		}

		if err != nil {
			http.Error(w, "failed to created administrator", http.StatusInternalServerError)
			return
		}

		var response CreateAdminResponse
		response.AdminID = adminID
		response.Message = "Administrator created successfully"

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			return
		}
	}
}

func GetAdminHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Wrong Method", http.StatusMethodNotAllowed)
			return
		}

		admins, err := database.GetAdmins(db)
		if err != nil {
			http.Error(w, "Failed to get administrators", http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(admins)
		if err != nil {
			http.Error(w, "Failed to encode administrators", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func GetAdminByIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
			return
		}

		idString := r.PathValue("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(w, "Failed to convert id to integer", http.StatusBadRequest)
			return
		}

		if id <= 0 {
			http.Error(w, "Invalid id value", http.StatusBadRequest)
			return
		}

		admin, err := database.GetAdministratorByID(db, id)
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Couldn't find administrator with given ID", http.StatusNotFound)
			return
		}

		if err != nil {
			http.Error(w, "Failed to get admin", http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(admin)
		if err != nil {
			http.Error(w, "Failed to encode administrator", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func PatchAdminHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			http.Error(w, "Wrong Method", http.StatusMethodNotAllowed)
			return
		}

		idString := r.PathValue("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(w, "Failed to convert id to integer", http.StatusBadRequest)
			return
		}

		if id <= 0 {
			http.Error(w, "Invalid administrator ID", http.StatusBadRequest)
			return
		}

		newAdmin := admin.AdminUpdate{}
		err = json.NewDecoder(r.Body).Decode(&newAdmin)
		if err != nil {
			http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
			return
		}

		if newAdmin.Name != nil {
			if len(strings.TrimSpace(*newAdmin.Name)) == 0 {
				http.Error(w, "Name should be provided", http.StatusBadRequest)
				return
			}
		}

		if newAdmin.Email != nil {
			if len(strings.TrimSpace(*newAdmin.Email)) == 0 {
				http.Error(w, "Email should be provided", http.StatusBadRequest)
				return
			}

			email, err := mail.ParseAddress(*newAdmin.Email)
			if err != nil {
				http.Error(w, "Invalid email format", http.StatusBadRequest)
				return
			}

			if email.Address != *newAdmin.Email {
				http.Error(w, "Invalid email format", http.StatusBadRequest)
				return
			}
		}

		if newAdmin.RoleID != nil {
			if *newAdmin.RoleID <= 0 || *newAdmin.RoleID > 3 {
				http.Error(w, "Invalid Role ID", http.StatusBadRequest)
				return
			}
		}

		if newAdmin.Name == nil && newAdmin.Email == nil && newAdmin.RoleID == nil && newAdmin.IsActive == nil {
			http.Error(w, "Wrong input", http.StatusBadRequest)
			return
		}

		authAdminID, ok := auth.GetAuthenticatedAdminID(r.Context())
		if !ok {
			http.Error(w, "failed to retrieve authenticated administor ID", http.StatusUnauthorized)
			return
		}

		if authAdminID == id {
			if newAdmin.RoleID != nil && *newAdmin.RoleID != 1 {
				http.Error(w, "You cannot change your role", http.StatusForbidden)
				return
			}

			if newAdmin.IsActive != nil && !*newAdmin.IsActive {
				http.Error(w, "You cannot change your activeness", http.StatusForbidden)
				return
			}
		}

		err = database.PatchAdmin(db, id, newAdmin)

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}

		if errors.Is(err, database.ErrLastSuperAdmin) {
			http.Error(w, "Cannot remove the last active Super Admin", http.StatusConflict)
			return
		}

		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No admin found with given id", http.StatusNotFound)
			return
		}

		if err != nil {
			http.Error(w, "Failed to patch admin", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func DeactivateAdminHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			http.Error(w, "Wrong Method", http.StatusMethodNotAllowed)
			return
		}

		idString := r.PathValue("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(w, "Failed to convert id to integer", http.StatusBadRequest)
			return
		}

		if id <= 0 {
			http.Error(w, "Invalid administrator ID", http.StatusBadRequest)
			return
		}

		adminID, ok := auth.GetAuthenticatedAdminID(r.Context())
		if !ok {
			http.Error(w, "failed to deactivate admin", http.StatusUnauthorized)
			return
		}

		if id == adminID {
			http.Error(w, "You cannot deactivate your own account", http.StatusForbidden)
			return
		}

		err = database.DeactivateAdmin(db, id)
		if errors.Is(err, database.ErrLastSuperAdmin) {
			http.Error(w, "cannot deactivate the last superadmin", http.StatusConflict)
			return
		}

		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No admin found with given id", http.StatusNotFound)
			return
		}

		if err != nil {
			http.Error(w, "Failed to deactivate admin", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJvbmxpbmUtc2hvcCIsInN1YiI6IjY2IiwiZXhwIjoxNzk1MTY4MTg4LCJpYXQiOjE3OTUxNjcyODh9.VcvCqfDdkSnIKISBbnm0Tl1y5FVdAzHaSbSl_YL7yAA
// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJvbmxpbmUtc2hvcCIsInN1YiI6IjY3IiwiZXhwIjoxNzk1MTQ2MjExLCJpYXQiOjE3OTUxNDUzMTF9.A0wVssJvDqa5GX0E295l183mGbQtCGhcK1MecuOB9lU
