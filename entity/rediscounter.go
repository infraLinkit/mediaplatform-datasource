package entity

type RedisCounter struct {
	Key        string `json:"key"`
	Name       string `json:"name"`
	Date       string `json:"date"`
	Time       string `json:"time"`
	HttpStatus int    `json:"http_status"`
	OS         string `json:"os"`
	IsOperator bool   `json:"is_operator"`
}

type DomainServices struct {
	Domain string `json:"domain"`
	Render string `json:"render"`
	Country string `json:"country"`
	Operator string `json:"operator"`
	Service string `json:"service"`
	Company string `json:"company"`
	CompanyLegalName string `json:"company_legal_name"`
	CompanyAddress string `json:"company_address"`
	CompanyEmail string `json:"company_email"`
	CompanyPhone string `json:"company_phone"`
	ServiceCurrency string `json:"service_currency"`
	ServicePrice float64 `json:"service_price"`
	PortalURL string `json:"portal_url"`
}