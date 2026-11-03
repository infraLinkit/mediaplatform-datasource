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

	ServiceOption struct {
		ID   int    `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	}

	ServiceGroupOption struct {
		Group      string          `json:"group"`
		CountryID  int             `json:"country_id"`
		OperatorID int             `json:"operator_id"`
		Services   []ServiceOption `json:"services"`
	}

	UserOptionsData struct {
		Roles     []Role               `json:"roles"`
		Companies []Company            `json:"companies"`
		Agencies  []Agency             `json:"agencies"`
		Adnets    []AdnetList          `json:"adnets"`
		Services  []ServiceGroupOption `json:"services"`
	}

	UserDetailResponse struct {
		User        User          `json:"user"`
		DetailUser  []DetailUser  `json:"detail_user"`
		UserAdnet   []UserAdnet   `json:"user_adnet"`
		UserCompany []UserCompany `json:"user_company"`
		UserAgency  []UserAgency  `json:"user_agency"`
	}
)
