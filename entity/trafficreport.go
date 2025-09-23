package entity


// TrafficReportParams defines the query parameters for TrafficReport time reports
type TrafficReportParams struct {
	DataType              string   `json:"data_type"`
	ReportType           string   `json:"report_type"`
	Country              string   `json:"country"`
	Operator             string   `json:"operator"`
	PartnerName          string   `json:"partner_name"`
	Service              string   `json:"service"`
	Adnet                string   `json:"adnet"`
	TypeData             string   `json:"type_data"`
	CampaignId           string   `json:"campaign_id"`
	URLServiceKey        string   `json:"url_service_key"`
	DataIndicators       []string `json:"data_indicators"`
	DataBasedOn          string   `json:"data_based_on"`
	DataBasedOnIndicator string   `json:"data_based_on_indicator"`
	DateRange           string   `json:"date_range"`
	DateStart           string   `json:"date_start"`
	DateEnd             string   `json:"date_end"`
	DateCustomRange     string   `json:"date_custom_range"`
	All                 string   `json:"all"`
	ChartType           string   `json:"chart_type"`
	CampaignName        string   `json:"campaign_name"`
}

// TrafficReportSummary represents the aggregated TrafficReport time data
type TrafficReportSummary struct {
	DateTime     string                          `json:"date_time"`
	CampaignId   string                          `json:"campaign_id"`
	CampaignName string                          `json:"campaign_name"`
	Country      string                          `json:"country"`
	Operator     string                          `json:"operator"`
	Partner      string                          `json:"partner"`
	Adnet        string                          `json:"adnet"`
	Service      string                          `json:"service"`
	URL          string                          `json:"url"`
	Metrics      map[string]float64              `json:"metrics"`
	DailyMetrics map[string]map[string]float64   `json:"daily_metrics"`
}