package handler

import (
	"net/http"
	"strconv"

	"cpsu/internal/auth/service"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	RoleService *service.RoleService
}

func NewRoleHandler(roleService *service.RoleService) *RoleHandler {
	return &RoleHandler{RoleService: roleService}
}

func (h *RoleHandler) AssignRole(c *gin.Context) {
	targetUserID, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		RoleID int `json:"role_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	actorUserID := c.GetInt("user_id")

	if err := h.RoleService.AssignRole(
		targetUserID, req.RoleID, actorUserID, c.ClientIP(), c.GetHeader("User-Agent"),
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "role assigned successfully"})
}
