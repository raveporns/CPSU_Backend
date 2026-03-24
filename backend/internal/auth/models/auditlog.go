package models

import "time"

type AuditLog struct {
	ID         int                    `json:"id"`
	UserID     int                    `json:"user_id"`
	Username   string                 `json:"username"`
	Email      string                 `json:"email"`
	Action     string                 `json:"action"`
	Resource   string                 `json:"resource"`
	ResourceID string                 `json:"resource_id"`
	Details    map[string]interface{} `json:"details"`
	CreatedAt  time.Time              `json:"created_at"`
}

type AuditLogResponse struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Action      string    `json:"action"`
	Resource    string    `json:"resource"`
	ResourceID  string    `json:"resource_id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
