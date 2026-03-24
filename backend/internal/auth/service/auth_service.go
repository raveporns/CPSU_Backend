package service

import (
	"errors"
	"time"

	"cpsu/internal/auth/models"
	"cpsu/internal/auth/repository"
	"cpsu/internal/auth/utils"
)

type AuthService struct {
	UserRepo  *repository.UserRepository
	RoleRepo  *repository.RoleRepository
	TokenRepo *repository.TokenRepository
	AuditRepo *repository.AuditRepository
}

func NewAuthService(
	userRepo *repository.UserRepository,
	roleRepo *repository.RoleRepository,
	tokenRepo *repository.TokenRepository,
	auditRepo *repository.AuditRepository,
) *AuthService {
	return &AuthService{
		UserRepo:  userRepo,
		RoleRepo:  roleRepo,
		TokenRepo: tokenRepo,
		AuditRepo: auditRepo,
	}
}

func (s *AuthService) Login(req models.LoginRequest, ipAddress string, userAgent string) (*models.LoginResponse, error) {
	user, err := s.UserRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	if err := utils.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	roles, err := s.RoleRepo.GetUserRoles(user.UserID)
	if err != nil {
		roles = []string{}
	}

	accessToken, err := utils.GenerateAccessToken(user.UserID, user.Username, roles)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.UserID, user.Username)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	_ = s.TokenRepo.StoreRefreshToken(user.UserID, refreshToken, expiresAt)
	_ = s.UserRepo.UpdateLastLogin(user.UserID)

	_ = s.AuditRepo.LogAudit(
		user.UserID, "login", "auth", "",
		map[string]interface{}{
			"email": user.Email,
		},
		ipAddress,
		userAgent,
	)

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: models.UserInfo{
			UserID:   user.UserID,
			Username: user.Username,
			Email:    user.Email,
			Roles:    roles,
		},
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (string, error) {
	userID, ok, err := s.TokenRepo.IsRefreshTokenValid(refreshToken)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errors.New("invalid refresh token")
	}

	user, err := s.UserRepo.FindByID(userID)
	if err != nil || user == nil {
		return "", errors.New("user not found")
	}

	roles, err := s.RoleRepo.GetUserRoles(user.UserID)
	if err != nil {
		roles = []string{}
	}

	accessToken, err := utils.GenerateAccessToken(user.UserID, user.Username, roles)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *AuthService) Logout(refreshToken string, userID int, ipAddress string, userAgent string) error {
	if refreshToken != "" {
		_ = s.TokenRepo.RevokeRefreshToken(refreshToken)
	}
	_ = s.AuditRepo.LogAudit(userID, "logout", "auth", "", nil, ipAddress, userAgent)
	return nil
}
