package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"example.com/online-shop/admin"
)

var ErrLastSuperAdmin = errors.New("cannot deactivate the last active super admin")

func CreateAdmin(db *sql.DB, a admin.Admin) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO admin (name, email, password_hash, role_id)
		VALUES ($1, $2, $3, $4)
		RETURNING adminid;
	`

	err := db.QueryRowContext(
		ctx,
		query,
		a.Name,
		a.Email,
		a.PasswordHash,
		a.RoleID,
	).Scan(&a.AdminID)

	if err != nil {
		return 0, fmt.Errorf("failed to create admin: %w", err)
	}

	return a.AdminID, nil
}

func GetAdministratorByID(db *sql.DB, id int) (admin.Admin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT adminID, name, email, role_id, is_active, created_at, updated_at
		FROM admin
		WHERE adminID = $1
	`

	var a admin.Admin

	err := db.QueryRowContext(ctx, query, id).Scan(
		&a.AdminID,
		&a.Name,
		&a.Email,
		&a.RoleID,
		&a.IsActive,
		&a.CreatedAt,
		&a.UpdatedAt,
	)

	if err != nil {
		return admin.Admin{}, fmt.Errorf("failed to retrieve admin information: %w", err)
	}

	return a, nil
}

func GetAdminByEmail(db *sql.DB, email string) (admin.Admin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var a admin.Admin

	// SELECT * is not used to match QueryRowContext
	query := ` 
		SELECT adminid, name, email, password_hash, role_id,    
       					is_active, created_at, updated_at
		FROM admin
		WHERE email = $1
	`
	err := db.QueryRowContext(ctx, query, email).Scan(
		&a.AdminID,
		&a.Name,
		&a.Email,
		&a.PasswordHash,
		&a.RoleID,
		&a.IsActive,
		&a.CreatedAt,
		&a.UpdatedAt,
	)

	if err != nil {
		return admin.Admin{}, fmt.Errorf("failed to retrieve admin information: %w", err)
	}

	return a, nil
}

func GetAdmins(db *sql.DB) ([]admin.Admin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	adm := make([]admin.Admin, 0)

	query := `
	SELECT adminid, name, email, role_id,
       					is_active, created_at, updated_at
	FROM admin
	`

	rows, err := db.QueryContext(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve admin information: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var a admin.Admin

		err := rows.Scan(
			&a.AdminID,
			&a.Name,
			&a.Email,
			&a.RoleID,
			&a.IsActive,
			&a.CreatedAt,
			&a.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to retrieve admin information: %w", err)
		}

		adm = append(adm, a)
	}

	err = rows.Err()

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve admin information: %w", err)
	}

	return adm, nil
}

func UpdateAdmin(db *sql.DB, a admin.Admin) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE admin
		SET name = $1,
			email = $2,
			role_id = $3,
			is_active = $4,
			updated_at = CURRENT_TIMESTAMP
		WHERE adminID = $5;
	`

	result, err := db.ExecContext(
		ctx,
		query,
		a.Name,
		a.Email,
		a.RoleID,
		a.IsActive,
		a.AdminID,
	)

	if err != nil {
		return fmt.Errorf("failed to update admin: %w", err)
	}

	rowAffected, err := result.RowsAffected()

	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func PatchAdmin(db *sql.DB, adminID int, updates admin.AdminUpdate) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start trancation: %w", err)
	}
	defer tx.Rollback()

	lock_query := `
		SELECT pg_advisory_xact_lock(1001);
	`
	_, err = tx.ExecContext(ctx, lock_query)
	if err != nil {
		return fmt.Errorf("failed to activate lock query: %w", err)
	}

	var roleID int
	var isActive bool

	deactivateQuery := `
		SELECT role_id, is_active
		FROM admin
		WHERE adminid = $1;
	`

	err = tx.QueryRowContext(ctx, deactivateQuery, adminID).Scan(
		&roleID,
		&isActive,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return sql.ErrNoRows
	}
	if err != nil {
		return fmt.Errorf("failed to retrieve admin with given role id and activeness: %w", err)
	}

	if roleID == 1 && isActive && ((updates.RoleID != nil && *updates.RoleID != 1) || (updates.IsActive != nil && !*updates.IsActive)) {

		var activeSuperAdmins int
		countQuery := `
			SELECT COUNT(*)
			FROM admin
			WHERE role_id = 1 AND is_active = TRUE;
		`
		if err := tx.QueryRowContext(ctx, countQuery).Scan(&activeSuperAdmins); err != nil {
			return fmt.Errorf("failed to count active super admins: %w", err)
		}

		if activeSuperAdmins <= 1 {
			return ErrLastSuperAdmin
		}
	}

	query := `
		UPDATE admin
		SET
			name = COALESCE($1, name),
			email = COALESCE($2, email),
			role_id = COALESCE($3, role_id),
			is_active = COALESCE($4, is_active),
			updated_at = CURRENT_TIMESTAMP
		WHERE adminid = $5;
	`

	result, err := tx.ExecContext(
		ctx,
		query,
		updates.Name,
		updates.Email,
		updates.RoleID,
		updates.IsActive,
		adminID,
	)
	if err != nil {
		return fmt.Errorf("failed to patch administrator: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows:%w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit administrator update: %w", err)
	}

	return nil
}

func DeactivateAdmin(db *sql.DB, id int) error {
	ctx2, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tx, err := db.BeginTx(ctx2, nil)
	if err != nil {
		return fmt.Errorf("failed to start trancation: %w", err)
	}
	defer tx.Rollback()

	lock_query := `
		SELECT pg_advisory_xact_lock(1001);
	`
	_, err = tx.ExecContext(ctx2, lock_query)
	if err != nil {
		return fmt.Errorf("failed to activate lock query: %w", err)
	}

	var roleID int
	var isActive bool

	deactivateQuery := `
		SELECT role_id, is_active
		FROM admin
		WHERE adminid = $1;
	`

	err = tx.QueryRowContext(ctx2, deactivateQuery, id).Scan(
		&roleID,
		&isActive,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return sql.ErrNoRows
	}
	if err != nil {
		return fmt.Errorf("failed to retrieve admin with given role id and activeness: %w", err)
	}

	if roleID == 1 && isActive {
		var activeSuperAdmins int
		countQuery := `
			SELECT COUNT(*)
			FROM admin
			WHERE role_id = 1 AND is_active = TRUE;
		`
		if err := tx.QueryRowContext(ctx2, countQuery).Scan(&activeSuperAdmins); err != nil {
			return fmt.Errorf("failed to count active super admins: %w", err)
		}

		if activeSuperAdmins <= 1 {
			return ErrLastSuperAdmin
		}
	}

	query := `
		UPDATE admin
		SET is_active = FALSE,
			updated_at = CURRENT_TIMESTAMP
		WHERE adminID = $1
	`

	result, err := tx.ExecContext(ctx2, query, id)
	if err != nil {
		return fmt.Errorf("failed to update admin status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Failed to commit changes: %w", err)
	}
	return nil
}
