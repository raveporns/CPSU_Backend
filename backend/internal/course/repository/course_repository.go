package repository

import (
	"database/sql"
	"strconv"
	"strings"

	"cpsu/internal/course/models"
)

type CourseRepository interface {
	GetAllCourses(param models.CoursesQueryParam) ([]models.Courses, error)
	GetCourseByID(id string) (*models.Courses, error)
	CreateCourse(req models.CoursesRequest) (*models.Courses, error)
	UpdateCourse(id string, req models.CoursesRequest) (*models.Courses, error)
	DeleteCourse(id string) error
}

type courseRepository struct {
	db *sql.DB
}

func NewCourseRepository(db *sql.DB) CourseRepository {
	return &courseRepository{db: db}
}

func (r *courseRepository) GetAllCourses(param models.CoursesQueryParam) ([]models.Courses, error) {
	query := `
		SELECT 
			c.course_id, c.degree, c.major, c.year, c.thai_course, 
			c.eng_course, c.thai_degree, c.eng_degree, c.admission_req, 
			c.graduation_req, c.philosophy, c.objective, c.tuition, c.credits, cp.career_paths_id, 
			cp.career_paths, p.plo_id, p.plo, c.detail_url, c.status
		FROM courses c
		LEFT JOIN career_paths cp ON c.career_paths_id = cp.career_paths_id
		LEFT JOIN plo p ON c.plo_id = p.plo_id
	`

	conditions := []string{}
	args := []interface{}{}
	argIndex := 1

	if param.Degree != "" {
		conditions = append(conditions, "c.degree = $"+strconv.Itoa(argIndex))
		args = append(args, param.Degree)
		argIndex++
	}

	if param.Major != "" {
		conditions = append(conditions, "c.major = $"+strconv.Itoa(argIndex))
		args = append(args, param.Major)
		argIndex++
	}

	if param.Year > 0 {
		conditions = append(conditions, "c.year = $"+strconv.Itoa(argIndex))
		args = append(args, param.Year)
		argIndex++
	}
	if param.Status != "" {
		conditions = append(conditions, "c.status = $"+strconv.Itoa(argIndex))
		args = append(args, param.Status)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	sort := "c.year"
	if param.Sort != "" {
		sort = "c." + param.Sort
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

	var courses []models.Courses
	for rows.Next() {
		var course models.Courses
		err := rows.Scan(
			&course.CourseID, &course.Degree, &course.Major, &course.Year,
			&course.ThaiCourse, &course.EngCourse, &course.ThaiDegree,
			&course.EngDegree, &course.AdmissionReq, &course.GraduationReq,
			&course.Philosophy, &course.Objective, &course.Tuition, &course.Credits,
			&course.CareerPathsID, &course.CareerPaths, &course.PloID,
			&course.PLO, &course.DetailURL, &course.Status,
		)
		if err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}

	return courses, nil
}

func (r *courseRepository) GetCourseByID(id string) (*models.Courses, error) {
	query := `
		SELECT 
			c.course_id, c.degree, c.major, c.year, c.thai_course, 
			c.eng_course, c.thai_degree, c.eng_degree, c.admission_req, 
			c.graduation_req, c.philosophy, c.objective, c.tuition, c.credits, cp.career_paths_id, 
			cp.career_paths, p.plo_id, p.plo, c.detail_url, c.status
		FROM courses c
		LEFT JOIN career_paths cp ON c.career_paths_id = cp.career_paths_id
		LEFT JOIN plo p ON c.plo_id = p.plo_id
		WHERE c.course_id = $1
	`

	row := r.db.QueryRow(query, id)

	var course models.Courses
	err := row.Scan(
		&course.CourseID, &course.Degree, &course.Major, &course.Year,
		&course.ThaiCourse, &course.EngCourse, &course.ThaiDegree,
		&course.EngDegree, &course.AdmissionReq, &course.GraduationReq,
		&course.Philosophy, &course.Objective, &course.Tuition, &course.Credits,
		&course.CareerPathsID, &course.CareerPaths, &course.PloID,
		&course.PLO, &course.DetailURL, &course.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &course, nil
}

func (r *courseRepository) CreateCourse(req models.CoursesRequest) (*models.Courses, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var careerPathsID int
	err = tx.QueryRow(`SELECT career_paths_id FROM career_paths WHERE career_paths = $1`, req.CareerPaths).Scan(&careerPathsID)
	if err == sql.ErrNoRows {
		err = tx.QueryRow(`INSERT INTO career_paths (career_paths) VALUES($1) RETURNING career_paths_id`, req.CareerPaths).Scan(&careerPathsID)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	var ploID int
	err = tx.QueryRow(`SELECT plo_id FROM plo WHERE plo = $1`, req.PLO).Scan(&ploID)
	if err == sql.ErrNoRows {
		err = tx.QueryRow(`INSERT INTO plo (plo) VALUES($1) RETURNING plo_id`, req.PLO).Scan(&ploID)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	_, err = tx.Exec(`
        INSERT INTO courses (
            course_id, degree, major, year, thai_course, eng_course, thai_degree,
			eng_degree, admission_req, graduation_req, philosophy, objective, 
			tuition, credits,career_paths_id, plo_id, detail_url, status
        )
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)
    `,
		req.CourseID, req.Degree, req.Major, req.Year, req.ThaiCourse, req.EngCourse, req.ThaiDegree,
		req.EngDegree, req.AdmissionReq, req.GraduationReq, req.Philosophy, req.Objective,
		req.Tuition, req.Credits, careerPathsID, ploID, req.DetailURL, req.Status,
	)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return r.GetCourseByID(req.CourseID)
}

func (r *courseRepository) UpdateCourse(id string, req models.CoursesRequest) (*models.Courses, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var careerPathsID int
	err = tx.QueryRow(`SELECT career_paths_id FROM career_paths WHERE career_paths = $1`, req.CareerPaths).Scan(&careerPathsID)
	if err == sql.ErrNoRows {
		err = tx.QueryRow(`INSERT INTO career_paths (career_paths) VALUES($1) RETURNING career_paths_id`, req.CareerPaths).Scan(&careerPathsID)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	var ploID int
	err = tx.QueryRow(`SELECT plo_id FROM plo WHERE plo = $1`, req.PLO).Scan(&ploID)
	if err == sql.ErrNoRows {
		err = tx.QueryRow(`INSERT INTO plo (plo) VALUES($1) RETURNING plo_id`, req.PLO).Scan(&ploID)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	_, err = tx.Exec(`
		UPDATE courses
		SET degree=$1, major=$2, year=$3, thai_course=$4, eng_course=$5, thai_degree=$6, eng_degree=$7,
		    admission_req=$8, graduation_req=$9, philosophy=$10, objective=$11, tuition=$12,
		    credits=$13, career_paths_id=$14, plo_id=$15, detail_url=$16, status=$17
		WHERE course_id=$18
	`,
		req.Degree, req.Major, req.Year, req.ThaiCourse, req.EngCourse, req.ThaiDegree, req.EngDegree,
		req.AdmissionReq, req.GraduationReq, req.Philosophy, req.Objective, req.Tuition,
		req.Credits, careerPathsID, ploID, req.DetailURL, req.Status, id,
	)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return r.GetCourseByID(id)
}

func (r *courseRepository) DeleteCourse(id string) error {
	result, err := r.db.Exec("DELETE FROM courses WHERE course_id = $1", id)
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
