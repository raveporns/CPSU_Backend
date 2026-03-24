package middlewares

import (
	"net/http"

	"cpsu/internal/auth/repository"

	"github.com/gin-gonic/gin"
)

type PermissionMiddleware struct {
	PermissionRepo *repository.PermissionRepository
}

func NewPermissionMiddleware(permissionRepo *repository.PermissionRepository) *PermissionMiddleware {
	return &PermissionMiddleware{PermissionRepo: permissionRepo}
}

func (m *PermissionMiddleware) RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("user_id")
		if userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		ok, err := m.PermissionRepo.CheckUserPermission(userID, permission)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "permission check failed"})
			c.Abort()
			return
		}

		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"error":    "insufficient permissions",
				"required": permission,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
