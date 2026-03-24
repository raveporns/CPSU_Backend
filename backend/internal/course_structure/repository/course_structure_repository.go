package repository

import (
	"database/sql"
	"strconv"
	"strings"

	"cpsu/internal/course_structure/models"
)

type CourseStructureRepository interface {
	GetAllCourseStructure(param models.CourseStructureQueryParam) ([]models.CourseStructure, error)
	GetCourseStructureByID(id int) (*models.CourseStructure, error)
	CreateCourseStructure(req *models.CourseStructureRequest) (*models.CourseStructure, error)
	UpdateCourseStructure(id int, req models.CourseStructureRequest) (*models.CourseStructure, error)
	DeleteCourseStructure(id int) error
}

type courseStructureRepository struct {
	db *sql.DB
}

func NewCourseStructureRepository(db *sql.DB) CourseStructureRepository {
	return &courseStructureRepository{db: db}
}

func (r *courseStructureRepository) GetAllCourseStructure(param models.CourseStructureQueryParam) ([]models.CourseStructure, error) {
	query := `
		SELECT 
			cs.course_structure_id, c.course_id, c.thai_course, cs.detail
		FROM course_structure cs
		LEFT JOIN courses c ON cs.course_id = c.course_id
	`

	conditions := []string{}
	args := []interface{}{}
	argIndex := 1

	if param.CourseID != "" {
		conditions = append(conditions, "c.course_id = $"+strconv.Itoa(argIndex))
		args = append(args, param.CourseID)
		argIndex++
	}

	if param.Search != "" {
		conditions = append(conditions, "cs.detail ILIKE $"+strconv.Itoa(argIndex))
		args = append(args, "%"+param.Search+"%")
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	sort := "cs.course_structure_id"
	switch param.Sort {
	case "course_id", "detail":
		sort = "cs." + param.Sort
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

	var result []models.CourseStructure
	for rows.Next() {
		var cs models.CourseStructure
		if err := rows.Scan(
			&cs.CourseStructureID, &cs.CourseID, &cs.ThaiCourse, &cs.Detail,
		); err != nil {
			return nil, err
		}
		result = append(result, cs)
	}

	return result, nil
}

func (r *courseStructureRepository) GetCourseStructureByID(id int) (*models.CourseStructure, error) {

	query := `
		SELECT 
			cs.course_structure_id, c.course_id,
			c.thai_course, cs.detail
		FROM course_structure cs
		LEFT JOIN courses c ON cs.course_id = c.course_id
		WHERE cs.course_structure_id = $1
	`

	var cs models.CourseStructure
	err := r.db.QueryRow(query, id).Scan(
		&cs.CourseStructureID, &cs.CourseID, &cs.ThaiCourse, &cs.Detail,
	)

	if err != nil {
		return nil, err
	}

	return &cs, nil
}

func (r *courseStructureRepository) CreateCourseStructure(req *models.CourseStructureRequest) (*models.CourseStructure, error) {
	query := `
		INSERT INTO course_structure (course_id, detail)
		VALUES ($1, $2)
		RETURNING course_structure_id
	`

	var cs models.CourseStructure
	err := r.db.QueryRow(query, req.CourseID, req.Detail).
		Scan(&cs.CourseStructureID)
	if err != nil {
		return nil, err
	}

	cs.CourseID = req.CourseID
	cs.Detail = req.Detail

	_ = r.db.QueryRow(
		"SELECT thai_course FROM courses WHERE course_id = $1",
		req.CourseID,
	).Scan(&cs.ThaiCourse)

	return &cs, nil
}

func (r *courseStructureRepository) UpdateCourseStructure(id int, req models.CourseStructureRequest) (*models.CourseStructure, error) {
	query := `
		UPDATE course_structure
		SET course_id = $1, detail = $2
		WHERE course_structure_id = $3
		RETURNING course_structure_id, course_id, detail
	`

	var cs models.CourseStructure
	err := r.db.QueryRow(
		query, req.CourseID, req.Detail, id,
	).Scan(
		&cs.CourseStructureID, &cs.CourseID, &cs.Detail,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	_ = r.db.QueryRow(
		"SELECT thai_course FROM courses WHERE course_id = $1",
		cs.CourseID,
	).Scan(&cs.ThaiCourse)

	return &cs, nil
}

func (r *courseStructureRepository) DeleteCourseStructure(id int) error {
	result, err := r.db.Exec("DELETE FROM course_structure WHERE course_structure_id = $1", id)
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
