package repository

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"cpsu/internal/auth/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetAllUser(param models.UserQueryParam) ([]models.UserResponse, error) {
	query := `
		SELECT u.user_id, u.username, u.email, r.role_id, r.name
		FROM users u
		LEFT JOIN user_roles ur ON u.user_id = ur.user_id
		LEFT JOIN roles r ON ur.role_id = r.role_id
	`

	conditions := []string{}
	args := []interface{}{}
	argIndex := 1

	if param.UserID > 0 {
		conditions = append(conditions, "u.user_id = $"+strconv.Itoa(argIndex))
		args = append(args, param.UserID)
		argIndex++
	}

	if param.RoleID > 0 {
		conditions = append(conditions, "r.role_id = $"+strconv.Itoa(argIndex))
		args = append(args, param.RoleID)
		argIndex++
	}

	if param.Search != "" {
		conditions = append(conditions, "(u.username ILIKE $"+strconv.Itoa(argIndex)+" OR u.email ILIKE $"+strconv.Itoa(argIndex)+")")
		args = append(args, "%"+param.Search+"%")
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	sortColumn := "u.user_id"
	switch param.Sort {
	case "username":
		sortColumn = "u.username"
	case "email":
		sortColumn = "u.email"
	case "role":
		sortColumn = "r.name"
	}

	order := "ASC"
	if strings.ToUpper(param.Order) == "DESC" {
		order = "DESC"
	}

	query += " ORDER BY " + sortColumn + " " + order

	if param.Limit > 0 {
		query += " LIMIT " + strconv.Itoa(param.Limit)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.UserResponse
	for rows.Next() {
		var u models.UserResponse

		var roleID sql.NullInt64
		var roleName sql.NullString

		if err := rows.Scan(&u.UserID, &u.Username, &u.Email, &roleID, &roleName); err != nil {
			return nil, err
		}

		if roleID.Valid {
			u.RoleID = int(roleID.Int64)
		}

		if roleName.Valid {
			u.Name = roleName.String
		}

		result = append(result, u)
	}
	return result, nil
}

func (r *UserRepository) CreateUser(username, email string) (int, error) {
	query := `
		INSERT INTO users (username, email, is_active)
		VALUES ($1, $2, true)
		RETURNING user_id
	`

	var userID int
	err := r.db.QueryRow(
		query, username, email,
	).Scan(&userID)

	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *UserRepository) DeleteUser(id int) error {
	query := `
		UPDATE users
		SET deleted_at = NOW()
		WHERE user_id = $1
		  AND deleted_at IS NULL
	`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	query := `
		SELECT user_id, username, email, password_hash, is_active,
		       created_at, last_login, deleted_at
		FROM users
		WHERE email = $1
		  AND deleted_at IS NULL
	`

	var user models.User
	var lastLogin sql.NullTime

	err := r.db.QueryRow(query, email).Scan(
		&user.UserID, &user.Username, &user.Email, &user.PasswordHash,
		&user.IsActive, &user.CreatedAt, &lastLogin, &user.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}

	return &user, nil
}

func (r *UserRepository) FindByID(userID int) (*models.User, error) {
	query := `
		SELECT user_id, username, email, password_hash, is_active, created_at, last_login
		FROM users
		WHERE user_id = $1
	`

	var user models.User
	var lastLogin sql.NullTime

	err := r.db.QueryRow(query, userID).Scan(
		&user.UserID, &user.Username, &user.Email, &user.PasswordHash,
		&user.IsActive, &user.CreatedAt, &lastLogin,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}

	return &user, nil
}

func (r *UserRepository) UpdateLastLogin(userID int) error {
	query := `
		UPDATE users
		SET last_login = NOW(),
		    updated_at = NOW()
		WHERE user_id = $1
	`
	_, err := r.db.Exec(query, userID)
	return err
}
