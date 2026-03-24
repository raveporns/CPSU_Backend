package models

type AssignRoleRequest struct {
	RoleID int `json:"role_id" binding:"required"`
}
