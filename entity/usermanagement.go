package entity

import "time"

type (
	UserManagementData struct {
		ID             int       `json:"id"`
		Username       string    `json:"username"`
		Name           string    `json:"name"`
		Email          string    `json:"email"`
		Role           string    `json:"role"`
		Status         bool      `json:"status"`
		LastLogin      *time.Time `json:"last_login"`
		IPAddress      string    `json:"ip_address"`
		Handset        string    `json:"handset"`
		TotalCountries int       `json:"total_countries"`
		TotalAdnets    int       `json:"total_adnets"`
	}

	UserApprovalRequestData struct {
		ID       int `json:"id"`
		Username string `json:"username"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		TypeUser string `json:"type_user"`
	}	

	UserCounts struct {
		TotalUsers    int `json:"totalUsers"`
		ActiveUsers   int `json:"activeUsers"`
		NonActiveUsers int `json:"nonActiveUsers"`
	}
	
)
