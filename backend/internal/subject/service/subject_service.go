package service

import (
	"cpsu/internal/subject/models"
	"cpsu/internal/subject/repository"
	"strconv"

	authrepo "cpsu/internal/auth/repository"
)

type SubjectService interface {
	GetAllSubjects(param models.SubjectsQueryParam) ([]models.Subjects, error)
	GetSubjectByID(id int) (*models.Subjects, error)
	CreateSubject(req models.SubjectsRequest, userID int, ip string, userAgent string) (*models.Subjects, error)
	UpdateSubject(id int, req models.SubjectsRequest, userID int, ip string, userAgent string) (*models.Subjects, error)
	DeleteSubject(id int, userID int, ip string, userAgent string) error
}

type subjectService struct {
	repo      repository.SubjectRepository
	auditRepo *authrepo.AuditRepository
}

func NewSubjectService(
	repo repository.SubjectRepository,
	auditRepo *authrepo.AuditRepository,
) SubjectService {
	return &subjectService{
		repo:      repo,
		auditRepo: auditRepo,
	}
}

func (s *subjectService) GetAllSubjects(param models.SubjectsQueryParam) ([]models.Subjects, error) {
	return s.repo.GetAllSubjects(param)
}

func (s *subjectService) GetSubjectByID(id int) (*models.Subjects, error) {
	return s.repo.GetSubjectByID(id)
}

func (s *subjectService) CreateSubject(subject models.SubjectsRequest, userID int, ip string, userAgent string) (*models.Subjects, error) {
	created, err := s.repo.CreateSubject(subject)
	if err != nil {
		return nil, err
	}

	err = s.auditRepo.LogAudit(
		userID, "create", "subject", created.SubjectID,
		map[string]interface{}{
			"thai_subject": created.ThaiSubject,
		},
		ip, userAgent,
	)

	return created, nil
}

func (s *subjectService) UpdateSubject(id int, subject models.SubjectsRequest, userID int, ip string, userAgent string) (*models.Subjects, error) {
	updated, err := s.repo.UpdateSubject(id, subject)
	if err != nil {
		return nil, err
	}

	err = s.auditRepo.LogAudit(
		userID, "update", "subject", strconv.Itoa(id),
		map[string]interface{}{
			"thai_subject": subject.ThaiSubject,
		},
		ip, userAgent,
	)

	return updated, nil
}

func (s *subjectService) DeleteSubject(id int, userID int, ip string, userAgent string) error {
	err := s.repo.DeleteSubject(id)
	if err != nil {
		return err
	}

	err = s.auditRepo.LogAudit(
		userID, "delete", "subject", strconv.Itoa(id),
		map[string]interface{}{},
		ip, userAgent,
	)

	return nil
}
