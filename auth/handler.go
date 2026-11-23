package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Wrong Method", http.StatusMethodNotAllowed)
			return
		}

		var admin LoginRequest

		err := json.NewDecoder(r.Body).Decode(&admin)
		if err != nil {
			http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
			return
		}

		if len(strings.TrimSpace(admin.Email)) == 0 || len(strings.TrimSpace(admin.Password)) == 0 {
			http.Error(w, "Email and password are required", http.StatusBadRequest)
			return
		}

		adminAuth, err := AuthenticateAdmin(db, admin.Email, admin.Password)
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		accessToken, err := GenerateJWT(adminAuth.AdminID)
		if err != nil {
			http.Error(w, "Failed to generate JWT", http.StatusInternalServerError)
			return
		}

		var response LoginResponse
		response.AccessToken = accessToken
		response.TokenType = "Bearer"
		response.ExpiresIn = 900

		w.Header().Set("Content-Type", "application/json")

		data, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}
}

func AdminProfileHandler(w http.ResponseWriter, r *http.Request) {
	adminID, ok := r.Context().Value(adminIDKey).(int)
	// 	r.Context() retrieves the request's context.
	// .Value(adminIDKey) retrieves the value associated with the key.
	// .(int) checks whether that value is an integer.
	// adminID receives the integer value.
	// ok is true if the value exists and has type int; otherwise, it's false.

	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "Authenticated administrator ID: %d", adminID)
}
