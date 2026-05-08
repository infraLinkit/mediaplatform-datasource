package entity

type DisplayGoogleTrafficReport struct {
	CampaignId    string `json:"campaign_id"`
	CampaignName  string `json:"campaign_name"`
	UrlServiceKey string `json:"url_service_key"`
	Country       string `json:"country"`
	Operator      string `json:"operator"`
	Partner       string `json:"partner"`
	Service       string `json:"service"`
	Company       string `json:"company"`
	Adnet         string `json:"adnet"`
	AdgroupID     string `json:"adgroup_id"`

	PeriodType string `json:"period_type"` // "weekly" | "monthly" | "custom"
	Month      string `json:"month"`       // YYYY-MM  (weekly & monthly)
	Year       string `json:"year"`        // YYYY     (monthly)
	Week       string `json:"week"`        // "1"–"4" atau "all" (weekly)
	DateFrom   string `json:"date_from"`   // YYYY-MM-DD (custom)
	DateTo     string `json:"date_to"`     // YYYY-MM-DD (custom)

	Action   string `json:"action"`
	Draw     int    `json:"draw"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

// GoogleTrafficReportRow – satu baris data dalam response
type GoogleTrafficReportRow struct {
	URLServiceKey string  `json:"url_service_key"`
	CampaignId    string  `json:"campaign_id"`
	CampaignName  string  `json:"campaign_name"`
	Country       string  `json:"country"`
	Operator      string  `json:"operator"`
	Partner       string  `json:"partner"`
	Adnet         string  `json:"adnet"`
	Service       string  `json:"service"`
	Company       string  `json:"company"`
	AdgroupID     string  `json:"adgroup_id"`
	Placement     string  `json:"placement"`

	// Period
	PeriodLabel string `json:"period_label"` // "Week 3 (2026-04)" / "2026-04" / "2026-04-15"
	SummaryDate string `json:"summary_date"`

	// Billing metrics
	StatusSuccess int     `json:"status_success"` // TRUE
	StatusFailed  int     `json:"status_failed"`  // FALSE
	TotalBill     int     `json:"total_bill"`
	BillRate      float64 `json:"bill_rate"` // (success / total) * 100
}

// GoogleTrafficTotalSummary – aggregat grand total, diisi ke field TotalSummary
type GoogleTrafficTotalSummary struct {
	TotalSuccess   int     `json:"total_success"`
	TotalFailed    int     `json:"total_failed"`
	TotalBill      int     `json:"total_bill"`
	AvgBillRate    float64 `json:"avg_bill_rate"`
	TotalCampaigns int     `json:"total_campaigns"`
	TotalAdgroups  int     `json:"total_adgroups"`
}