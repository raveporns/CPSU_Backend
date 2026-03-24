package repository

import (
	"database/sql"
	"strconv"
	"strings"

	"cpsu/internal/roadmap/models"
)

type RoadmapRepository interface {
	GetAllRoadmap(param models.RoadmapQueryParam) ([]models.Roadmap, error)
	GetRoadmapByID(id int) (*models.Roadmap, error)
	CreateRoadmap(req *models.RoadmapRequest) (*models.Roadmap, error)
	DeleteRoadmap(id int) error
}

type roadmapRepository struct {
	db *sql.DB
}

func NewRoadmapRepository(db *sql.DB) RoadmapRepository {
	return &roadmapRepository{db: db}
}

func (r *roadmapRepository) GetAllRoadmap(param models.RoadmapQueryParam) ([]models.Roadmap, error) {
	query := `
		SELECT r.roadmap_id, c.course_id, c.thai_course, r.roadmap_url
		FROM roadmap r
		LEFT JOIN courses c ON r.course_id = c.course_id
	`

	conditions := []string{}
	args := []interface{}{}
	argIndex := 1

	if param.CourseID != "" {
		conditions = append(conditions, "c.course_id = $"+strconv.Itoa(argIndex))
		args = append(args, param.CourseID)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	sort := "r.roadmap_id"
	if param.Sort != "" {
		sort = "r." + param.Sort
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

	var roadmaps []models.Roadmap
	for rows.Next() {
		var roadmap models.Roadmap
		err := rows.Scan(&roadmap.RoadmapID, &roadmap.CourseID, &roadmap.ThaiCourse, &roadmap.RoadmapURL)
		if err != nil {
			return nil, err
		}
		roadmaps = append(roadmaps, roadmap)
	}

	return roadmaps, nil
}

func (r *roadmapRepository) GetRoadmapByID(id int) (*models.Roadmap, error) {
	query := `
		SELECT r.roadmap_id, c.course_id, c.thai_course, r.roadmap_url
		FROM roadmap r
		LEFT JOIN courses c ON r.course_id = c.course_id
		WHERE r.roadmap_id = $1
	`
	row := r.db.QueryRow(query, id)

	var roadmap models.Roadmap
	err := row.Scan(&roadmap.RoadmapID, &roadmap.CourseID, &roadmap.ThaiCourse, &roadmap.RoadmapURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &roadmap, nil
}

func (r *roadmapRepository) CreateRoadmap(req *models.RoadmapRequest) (*models.Roadmap, error) {
	query := `
		INSERT INTO roadmap (course_id, roadmap_url)
		VALUES ($1, $2)
		RETURNING roadmap_id
	`

	var roadmap models.Roadmap
	err := r.db.QueryRow(query, req.CourseID, req.RoadmapURL).Scan(&roadmap.RoadmapID)
	if err != nil {
		return nil, err
	}

	roadmap.CourseID = req.CourseID
	roadmap.RoadmapURL = req.RoadmapURL

	row := r.db.QueryRow("SELECT thai_course FROM courses WHERE course_id = $1", req.CourseID)
	_ = row.Scan(&roadmap.ThaiCourse)

	return &roadmap, nil
}

func (r *roadmapRepository) DeleteRoadmap(id int) error {
	result, err := r.db.Exec("DELETE FROM roadmap WHERE roadmap_id = $1", id)
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
