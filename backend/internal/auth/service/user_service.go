package service

import (
	"cpsu/internal/auth/models"
	"cpsu/internal/auth/repository"
	"strconv"
)

type UserService struct {
	UserRepo  *repository.UserRepository
	AuditRepo *repository.AuditRepository
}

func NewUserService(
	userRepo *repository.UserRepository,
	auditRepo *repository.AuditRepository,
) *UserService {
	return &UserService{
		UserRepo:  userRepo,
		AuditRepo: auditRepo,
	}
}

func (s *UserService) GetAllUser(param models.UserQueryParam) ([]models.UserResponse, error) {
	return s.UserRepo.GetAllUser(param)
}

func (s *UserService) CreateUser(req models.UserRequest, ipAddress string, userAgent string) error {

	userID, err := s.UserRepo.CreateUser(req.Username, req.Email)
	if err != nil {
		return err
	}

	_ = s.AuditRepo.LogAudit(
		userID, "create", "user", "",
		map[string]interface{}{
			"email":    req.Email,
			"username": req.Username,
		},
		ipAddress,
		userAgent,
	)

	return nil
}

func (s *UserService) DeleteUser(targetUserID int, actorUserID int, ipAddress string, userAgent string) error {
	if err := s.UserRepo.DeleteUser(targetUserID); err != nil {
		return err
	}

	_ = s.AuditRepo.LogAudit(
		actorUserID, "delete", "user",
		strconv.Itoa(targetUserID),
		map[string]interface{}{},
		ipAddress, userAgent,
	)

	return nil
}
