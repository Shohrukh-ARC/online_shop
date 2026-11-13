package admin

import "time"

type Admin struct {
	AdminID      int       `json:"adminID"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	RoleID       int       `json:"roleID"`
	IsActive     bool      `json:"isActive"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type AdminUpdate struct {
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	RoleID   *int    `json:"role_id"`
	IsActive *bool   `json:"is_active"`
}
