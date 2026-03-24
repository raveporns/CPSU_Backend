package service

import (
	"cpsu/internal/calendar/models"
	"cpsu/internal/calendar/repository"
	"strconv"

	authrepo "cpsu/internal/auth/repository"
)

type CalendarService interface {
	GetAllCalendars(param models.CalendarQueryParam) ([]models.Calendar, error)
	GetCalendarByID(id int) (*models.Calendar, error)
	CreateCalendar(req models.CalendarRequest, userID int, ip string, userAgent string) (*models.Calendar, error)
	UpdateCalendar(id int, req models.CalendarRequest, userID int, ip string, userAgent string) (*models.Calendar, error)
	DeleteCalendar(id int, userID int, ip string, userAgent string) error
}

type calendarService struct {
	repo      repository.CalendarRepository
	auditRepo *authrepo.AuditRepository
}

func NewCalendarService(repo repository.CalendarRepository, auditRepo *authrepo.AuditRepository) CalendarService {
	return &calendarService{
		repo:      repo,
		auditRepo: auditRepo,
	}
}

func (s *calendarService) GetAllCalendars(param models.CalendarQueryParam) ([]models.Calendar, error) {
	return s.repo.GetAllCalendars(param)
}

func (s *calendarService) GetCalendarByID(id int) (*models.Calendar, error) {
	return s.repo.GetCalendarByID(id)
}

func (s *calendarService) CreateCalendar(req models.CalendarRequest, userID int, ip string, userAgent string) (*models.Calendar, error) {
	created, err := s.repo.CreateCalendar(&req)
	if err != nil {
		return nil, err
	}

	err = s.auditRepo.LogAudit(
		userID, "create", "calendar",
		strconv.Itoa(created.CalenderID),
		map[string]interface{}{
			"title": created.Title,
		},
		ip, userAgent,
	)

	return created, nil
}

func (s *calendarService) UpdateCalendar(id int, req models.CalendarRequest, userID int, ip string, userAgent string) (*models.Calendar, error) {
	req.CalenderID = id

	updated, err := s.repo.UpdateCalendar(&req)
	if err != nil {
		return nil, err
	}

	err = s.auditRepo.LogAudit(
		userID, "update", "calendar", strconv.Itoa(id),
		map[string]interface{}{
			"title": req.Title,
		},
		ip, userAgent,
	)

	return updated, nil
}

func (s *calendarService) DeleteCalendar(id int, userID int, ip string, userAgent string) error {
	err := s.repo.DeleteCalendar(id)
	if err != nil {
		return err
	}

	err = s.auditRepo.LogAudit(
		userID, "delete", "calendar", strconv.Itoa(id),
		map[string]interface{}{},
		ip, userAgent,
	)

	return nil
}
