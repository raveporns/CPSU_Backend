// internal/audit/service/audit_service.go
package service

import (
	"context"

	"cpsu/internal/auth/models"
	"cpsu/internal/auth/repository"
)

type AuditService struct {
	AuditRepo *repository.AuditRepository
}

func NewAuditService(auditRepo *repository.AuditRepository) *AuditService {
	return &AuditService{AuditRepo: auditRepo}
}

func (s *AuditService) GetAllAuditLog(ctx context.Context) ([]models.AuditLogResponse, error) {
	audits, err := s.AuditRepo.GetAllAuditLog(ctx)
	if err != nil {
		return nil, err
	}

	var res []models.AuditLogResponse

	for _, a := range audits {
		res = append(res, models.AuditLogResponse{
			ID:          a.ID,
			UserID:      a.UserID,
			Username:    a.Username,
			Email:       a.Email,
			Action:      a.Action,
			Resource:    a.Resource,
			ResourceID:  a.ResourceID,
			Description: Description(a),
			CreatedAt:   a.CreatedAt,
		})
	}

	return res, nil
}

func Description(a models.AuditLog) string {
	switch a.Action {
	case "login":
		return "เข้าสู่ระบบ"
	case "logout":
		return "ออกจากระบบ"
	case "create":
		return "เพิ่มข้อมูลใหม่"
	case "update":
		return "แก้ไขข้อมูล"
	case "delete":
		return "ลบข้อมูล"
	case "assign_role":
		return "ให้สิทธิ์ผู้ใช้งาน"
	default:
		return "มีการดำเนินการในระบบ"
	}
}
