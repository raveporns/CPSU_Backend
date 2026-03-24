package repository

import (
	"context"
	"cpsu/internal/auth/models"
	"database/sql"
	"encoding/json"
)

type AuditRepository struct {
	db *sql.DB
}

func NewAuditRepository(db *sql.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) GetAllAuditLog(ctx context.Context) ([]models.AuditLog, error) {

	rows, err := r.db.QueryContext(ctx, `
		SELECT
			a.id, u.user_id, u.username, u.email, a.action, 
			a.resource, a.resource_id, a.details, a.created_at
		FROM audit_logs a 
		LEFT JOIN users u ON a.user_id = u.user_id
		ORDER BY a.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var audits []models.AuditLog
	var username, email sql.NullString

	for rows.Next() {
		var audit models.AuditLog
		var detailsBytes []byte

		err := rows.Scan(
			&audit.ID, &audit.UserID, &username, &email,
			&audit.Action, &audit.Resource, &audit.ResourceID,
			&detailsBytes, &audit.CreatedAt,
		)

		if username.Valid {
			audit.Username = username.String
		}
		if email.Valid {
			audit.Email = email.String
		}

		if err != nil {
			return nil, err
		}

		if len(detailsBytes) > 0 {
			if err := json.Unmarshal(detailsBytes, &audit.Details); err != nil {
				return nil, err
			}
		}

		audits = append(audits, audit)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return audits, nil
}

func (r *AuditRepository) LogAudit(userID int, action string, resource string, resourceID string, details map[string]interface{}, ipAddress string, userAgent string) error {
	var detailsJSON []byte
	if details != nil {
		detailsJSON, _ = json.Marshal(details)
	}

	query := `
		INSERT INTO audit_logs(user_id, action, resource, resource_id, details, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(
		query, userID, action, resource, resourceID,
		detailsJSON, ipAddress, userAgent,
	)

	return err
}
