package entity

import "time"

type (
	DisplayUserLogList struct {
		ID            uint      `json:"id"`
		UserType      string    `json:"user_type"`
		UserID        int       `json:"user_id"`
		Event         string    `json:"event"`
		AuditableType string    `json:"auditable_type"`
		AuditableID   string    `json:"auditable_id"`
		OldValues     string    `json:"old_values"`
		NewValues     string    `json:"new_values"`
		URL           string    `json:"url"`
		IPAddress     string    `json:"ip_address"`
		UserAgent     string    `json:"user_agent"`
		Tags          *string   `json:"tags"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
		ActionName    string    `json:"action_name"`
		Email         string    `json:"email"`
		Name          string    `json:"name"`
		Username      string    `json:"username"`
		Role          string    `json:"role"`
		ActionDate    string    `json:"action_date"`
	}
)
