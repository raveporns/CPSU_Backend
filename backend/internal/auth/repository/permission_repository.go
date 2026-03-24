package repository

import (
	"database/sql"
)

type PermissionRepository struct {
	db *sql.DB
}

func NewPermissionRepository(db *sql.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) CheckUserPermission(userID int, permissionName string) (bool, error) {
	query := `
		SELECT p.permission_id
		FROM permissions p
		JOIN role_permissions rp ON p.permission_id = rp.permission_id
		JOIN user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = $1 AND p.name = $2
		LIMIT 1
	`

	var permissionID int
	err := r.db.QueryRow(query, userID, permissionName).Scan(&permissionID)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
