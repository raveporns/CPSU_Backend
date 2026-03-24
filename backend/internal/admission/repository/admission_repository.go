package repository

import (
	"database/sql"
	"strconv"
	"strings"

	"cpsu/internal/admission/models"
)

type AdmissionRepository interface {
	GetAllAdmission(param models.AdmissionQueryParam) ([]models.Admission, error)
	GetAdmissionByID(id int) (*models.Admission, error)
	CreateAdmission(req models.AdmissionRequest) (*models.Admission, error)
	UpdateAdmission(id int, req models.AdmissionRequest) (*models.Admission, error)
	DeleteAdmission(id int) error
}

type admissionRepository struct {
	db *sql.DB
}

func NewAdmissionRepository(db *sql.DB) AdmissionRepository {
	return &admissionRepository{db: db}
}

func (r *admissionRepository) GetAllAdmission(param models.AdmissionQueryParam) ([]models.Admission, error) {
	query := `
		SELECT admission_id, round, detail, file_image
		FROM admission
	`

	conditions := []string{}
	args := []interface{}{}
	argIndex := 1

	if param.Round != "" {
		conditions = append(conditions, "round = $"+strconv.Itoa(argIndex))
		args = append(args, param.Round)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	sort := "admission_id"
	if param.Sort != "" {
		sort = param.Sort
	}

	order := "ASC"
	if strings.ToUpper(param.Order) == "DESC" {
		order = "DESC"
	}

	query += " ORDER BY " + sort + " " + order

	if param.Limit > 0 {
		query += " LIMIT " + strconv.Itoa(param.Limit)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var admission []models.Admission
	for rows.Next() {
		var a models.Admission
		if err := rows.Scan(
			&a.AdmissionID,
			&a.Round,
			&a.Detail,
			&a.FileImage,
		); err != nil {
			return nil, err
		}
		admission = append(admission, a)
	}

	return admission, nil
}

func (r *admissionRepository) GetAdmissionByID(id int) (*models.Admission, error) {
	query := `
		SELECT admission_id, round, detail, file_image
		FROM admission
		WHERE admission_id = $1
	`

	row := r.db.QueryRow(query, id)

	var a models.Admission
	err := row.Scan(
		&a.AdmissionID,
		&a.Round,
		&a.Detail,
		&a.FileImage,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &a, nil
}

func (r *admissionRepository) CreateAdmission(req models.AdmissionRequest) (*models.Admission, error) {
	var newID int

	err := r.db.QueryRow(`
		INSERT INTO admission (round, detail, file_image)
		VALUES ($1, $2, $3)
		RETURNING admission_id
	`,
		req.Round,
		req.Detail,
		req.FileImage,
	).Scan(&newID)

	if err != nil {
		return nil, err
	}

	return r.GetAdmissionByID(newID)
}

func (r *admissionRepository) UpdateAdmission(id int, req models.AdmissionRequest) (*models.Admission, error) {
	var updatedID int

	err := r.db.QueryRow(`
		UPDATE admission
		SET round = $1,
		    detail = $2,
		    file_image = $3
		WHERE admission_id = $4
		RETURNING admission_id
	`,
		req.Round,
		req.Detail,
		req.FileImage,
		id,
	).Scan(&updatedID)

	if err != nil {
		return nil, err
	}

	return r.GetAdmissionByID(updatedID)
}

func (r *admissionRepository) DeleteAdmission(id int) error {
	result, err := r.db.Exec(
		"DELETE FROM admission WHERE admission_id = $1",
		id,
	)
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
