package auth

import (
	"database/sql"
	"net/http"
)
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(db *sql.DB) http.Hand