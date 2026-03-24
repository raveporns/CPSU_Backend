package service

import (
	"cpsu/internal/course/models"
	"cpsu/internal/course/repository"

	authrepo "cpsu/internal/auth/repository"
)

type CourseService interface {
	GetAllCourses(param models.CoursesQueryParam) ([]models.Courses, error)
	GetCourseByID(id string) (*models.Courses, error)
	CreateCourse(course models.CoursesRequest, userID int, ip string, userAgent string) (*models.Courses, error)
	UpdateCourse(id string, course models.CoursesRequest, userID int, ip string, userAgent string) (*models.Courses, error)
	DeleteCourse(id string, userID int, ip string, userAgent string) error
}

type courseService struct {
	repo      repository.CourseRepository
	auditRepo *authrepo.AuditRepository
}

func NewCourseService(
	repo repository.CourseRepository,
	auditRepo *authrepo.AuditRepository,
) CourseService {
	return &courseService{
		repo:      repo,
		auditRepo: auditRepo,
	}
}

func (s *courseService) GetAllCourses(param models.CoursesQueryParam) ([]models.Courses, error) {
	return s.repo.GetAllCourses(param)
}

func (s *courseService) GetCourseByID(id string) (*models.Courses, error) {
	return s.repo.GetCourseByID(id)
}

func (s *courseService) CreateCourse(course models.CoursesRequest, userID int, ip string, userAgent string) (*models.Courses, error) {

	created, err := s.repo.CreateCourse(course)
	if err != nil {
		return nil, err
	}

	err = s.auditRepo.LogAudit(
		userID, "create", "course", created.CourseID,
		map[string]interface{}{
			"thai_course": created.ThaiCourse,
		},
		ip, userAgent,
	)

	return created, nil
}

func (s *courseService) UpdateCourse(id string, course models.CoursesRequest, userID int, ip string, userAgent string) (*models.Courses, error) {

	updated, err := s.repo.UpdateCourse(id, course)
	if err != nil {
		return nil, err
	}

	err = s.auditRepo.LogAudit(
		userID, "update", "course", id,
		map[string]interface{}{
			"thai_course": course.ThaiCourse,
		},
		ip, userAgent,
	)

	return updated, nil
}

func (s *courseService) DeleteCourse(id string, userID int, ip string, userAgent string) error {

	err := s.repo.DeleteCourse(id)
	if err != nil {
		return err
	}

	err = s.auditRepo.LogAudit(
		userID, "delete", "course", id,
		map[string]interface{}{},
		ip, userAgent,
	)

	return nil
}
