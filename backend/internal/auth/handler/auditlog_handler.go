// internal/audit/handler/audit_handler.go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"cpsu/internal/auth/service"
)

type AuditHandler struct {
	AuditService *service.AuditService
}

func NewAuditHandler(auditService *service.AuditService) *AuditHandler {
	return &AuditHandler{AuditService: auditService}
}

func (h *AuditHandler) GetAllAuditLog(c *gin.Context) {
	audits, err := h.AuditService.GetAllAuditLog(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get audit logs",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": audits,
	})
}
