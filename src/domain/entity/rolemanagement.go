package entity

type (
	RoleManagementData struct {
		ID          int     `json:"id"`
		Code        string   `json:"code"`
		Name        string   `json:"name"`
		Permissions []string `json:"permissions"`
	}
)
