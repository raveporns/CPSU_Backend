package repository

import (
	"database/sql"
	"strconv"
	"strings"

	"cpsu/internal/subject/models"
)

type SubjectRepository interface {
	GetAllSubjects(param models.SubjectsQueryParam) ([]models.Subjects, error)
	GetSubjectByID(id int) (*models.Subjects, error)
	CreateSubject(req models.SubjectsRequest) (*models.Subjects, error)
	UpdateSubject(id int, req models.SubjectsRequest) (*models.Subjects, error)
	DeleteSubject(id int) error
}

type subjectRepository struct {
	db *sql.DB
}

func NewSubjectRepository(db *sql.DB) SubjectRepository {
	return &subjectRepository{db: db}
}

func (r *subjectRepository) GetAllSubjects(param models.SubjectsQueryParam) ([]models.Subjects, error) {
	query := `
		SELECT 
			s.id, s.subject_id, c.course_id, c.thai_course,s.plan_type, 
			s.semester, s.thai_subject, s.eng_subject, s.credits, 
			s.compulsory_subject, s.condition, d.description_id, 
			d.description_thai, d.description_eng, cl.clo_id, cl.clo
		FROM subjects s
		LEFT JOIN courses c ON s.course_id = c.course_id
		LEFT JOIN description d ON s.description_id = d.description_id
		LEFT JOIN clo cl ON s.clo_id = cl.clo_id
	`

	conditions := []string{}
	args := []interface{}{}
	argIndex := 1

	if param.SubjectID != "" {
		conditions = append(conditions, "s.subject_id = $"+strconv.Itoa(argIndex))
		args = append(args, param.SubjectID)
		argIndex++
	}

	if param.CourseID != "" {
		conditions = append(conditions, "c.course_id = $"+strconv.Itoa(argIndex))
		args = append(args, param.CourseID)
		argIndex++
	}

	if param.PlanType != "" {
		conditions = append(conditions, "s.plan_type = $"+strconv.Itoa(argIndex))
		args = append(args, param.PlanType)
		argIndex++
	}

	if param.Semester != "" {
		conditions = append(conditions, "s.semester = $"+strconv.Itoa(argIndex))
		args = append(args, param.Semester)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	sort := "s.id"
	if param.Sort != "" {
		sort = "s." + param.Sort
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

	var subjects []models.Subjects
	for rows.Next() {
		var subject models.Subjects
		err := rows.Scan(
			&subject.ID, &subject.SubjectID, &subject.CourseID, &subject.ThaiCourse,
			&subject.PlanType, &subject.Semester, &subject.ThaiSubject,
			&subject.EngSubject, &subject.Credits, &subject.CompulsorySubject,
			&subject.Condition, &subject.DescriptionID, &subject.DescriptionThai,
			&subject.DescriptionEng, &subject.CloID, &subject.CLO,
		)
		if err != nil {
			return nil, err
		}
		subjects = append(subjects, subject)
	}

	return subjects, nil
}

func (r *subjectRepository) GetSubjectByID(id int) (*models.Subjects, error) {
	query := `
		SELECT 
			s.id, s.subject_id, c.course_id, c.thai_course,s.plan_type, 
			s.semester, s.thai_subject, s.eng_subject, s.credits, 
			s.compulsory_subject, s.condition, d.description_id, 
			d.description_thai, d.description_eng, cl.clo_id, cl.clo
		FROM subjects s
		LEFT JOIN courses c ON s.course_id = c.course_id
		LEFT JOIN description d ON s.description_id = d.description_id
		LEFT JOIN clo cl ON s.clo_id = cl.clo_id
		WHERE s.id = $1
	`

	row := r.db.QueryRow(query, id)

	var subject models.Subjects
	err := row.Scan(
		&subject.ID, &subject.SubjectID, &subject.CourseID, &subject.ThaiCourse,
		&subject.PlanType, &subject.Semester, &subject.ThaiSubject,
		&subject.EngSubject, &subject.Credits, &subject.CompulsorySubject,
		&subject.Condition, &subject.DescriptionID, &subject.DescriptionThai,
		&subject.DescriptionEng, &subject.CloID, &subject.CLO,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &subject, nil
}

func (r *subjectRepository) CreateSubject(req models.SubjectsRequest) (*models.Subjects, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if req.DescriptionThai != nil || req.DescriptionEng != nil {
		_, err = tx.Exec(`
			INSERT INTO description (description_id, description_thai, description_eng)
			VALUES ($1,$2,$3)
			ON CONFLICT (description_id) DO UPDATE 
			SET description_thai=EXCLUDED.description_thai, description_eng=EXCLUDED.description_eng
		`, req.SubjectID, req.DescriptionThai, req.DescriptionEng)
		if err != nil {
			return nil, err
		}
		req.DescriptionID = &req.SubjectID
	}

	if req.CLO != nil {
		_, err = tx.Exec(`
			INSERT INTO clo (clo_id, clo)
			VALUES ($1,$2)
			ON CONFLICT (clo_id) DO UPDATE 
			SET clo=EXCLUDED.clo
		`, req.SubjectID, req.CLO)
		if err != nil {
			return nil, err
		}
		req.CloID = &req.SubjectID
	}

	var subject models.Subjects
	err = tx.QueryRow(`
		INSERT INTO subjects (
			subject_id, course_id, plan_type, semester, thai_subject, eng_subject,
			credits, compulsory_subject, condition, description_id, clo_id
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		RETURNING id
	`,
		req.SubjectID, req.CourseID, req.PlanType, req.Semester,
		req.ThaiSubject, req.EngSubject, req.Credits, req.CompulsorySubject,
		req.Condition, req.DescriptionID, req.CloID,
	).Scan(&subject.ID)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return r.GetSubjectByID(subject.ID)
}

func (r *subjectRepository) UpdateSubject(id int, req models.SubjectsRequest) (*models.Subjects, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if req.DescriptionThai != nil || req.DescriptionEng != nil {
		_, err = tx.Exec(`
			INSERT INTO description (description_id, description_thai, description_eng)
			VALUES ($1,$2,$3)
			ON CONFLICT (description_id) DO UPDATE 
			SET description_thai=EXCLUDED.description_thai, description_eng=EXCLUDED.description_eng
		`, req.SubjectID, req.DescriptionThai, req.DescriptionEng)
		if err != nil {
			return nil, err
		}
		req.DescriptionID = &req.SubjectID
	}

	if req.CLO != nil {
		_, err = tx.Exec(`
			INSERT INTO clo (clo_id, clo)
			VALUES ($1,$2)
			ON CONFLICT (clo_id) DO UPDATE 
			SET clo=EXCLUDED.clo
		`, req.SubjectID, req.CLO)
		if err != nil {
			return nil, err
		}
		req.CloID = &req.SubjectID
	}

	_, err = tx.Exec(`
		UPDATE subjects
		SET subject_id=$1, course_id=$2, plan_type=$3, semester=$4, 
		    thai_subject=$5, eng_subject=$6, credits=$7, compulsory_subject=$8, 
		    condition=$9, description_id=$10, clo_id=$11
		WHERE id=$12
	`, req.SubjectID, req.CourseID, req.PlanType, req.Semester,
		req.ThaiSubject, req.EngSubject, req.Credits, req.CompulsorySubject,
		req.Condition, req.DescriptionID, req.CloID, id,
	)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return r.GetSubjectByID(id)
}

func (r *subjectRepository) DeleteSubject(id int) error {
	result, err := r.db.Exec("DELETE FROM subjects WHERE id = $1", id)
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
