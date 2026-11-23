package auth

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"example.com/online-shop/database"
)

type contextKey string

const adminIDKey contextKey = "adminID"
const adminRoleKey contextKey = "adminRole"

func AdminAuthMiddleware(db *sql.DB, next http.Handler) http.Handler {
	// next -> HTTP handler that should run after successfull authentication
	// returned http.Handler is the wrapper that checks authentication before executing "next"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if len(authHeader) == 0 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		if !strings.EqualFold(parts[0], "Bearer") {
			http.Error(w, "Wrong token type", http.StatusUnauthorized)
			return
		}

		if len(parts[1]) == 0 {
			http.Error(w, "JWT cannot be empty", http.StatusUnauthorized)
			return
		}

		jwt := parts[1]
		claims, err := VerifyJWT(jwt)
		if err != nil {
			http.Error(w, "failed to verify JWT", http.StatusUnauthorized)
			return
		}

		adminID, err := strconv.Atoi(claims.Subject)
		if err != nil {
			http.Error(w, "failed to get adminID", http.StatusUnauthorized)
			return
		}
		if adminID <= 0 {
			http.Error(w, "adminID should be positive", http.StatusUnauthorized)
			return
		}

		admin, err := database.GetAdministratorByID(db, adminID)
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Administrator with this ID does not exist", http.StatusUnauthorized)
			return
		}

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !admin.IsActive {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx1 := context.WithValue(r.Context(), adminIDKey, adminID)
		ctx2 := context.WithValue(ctx1, adminRoleKey, admin.RoleID)
		// ctx -> new context containing ID
		// r.Context() -> Gets the current HTTP request's context.
		// adminID -> Verified administrator ID

		next.ServeHTTP(w, r.WithContext(ctx2))
	})
}

func RequireAdminRole(allowedRoles ...int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(adminRoleKey).(int)

			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "Forbidden", http.StatusForbidden)
		})
	}
}

func GetAuthenticatedAdminID(ctx context.Context) (int, bool) {
	adminID, ok := ctx.Value(adminIDKey).(int)
	if !ok || adminID <= 0 {
		return 0, false
	}

	return adminID, true
}
