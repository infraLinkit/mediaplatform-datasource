package entity


// RedirectionTimeParams defines the query parameters for redirection time reports
type RedirectionTimeParams struct {
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

// RedirectionTimeSummary represents the aggregated redirection time data
type RedirectionTimeSummary struct {
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

// RedirectionTimeChartData represents data for chart visualization
type RedirectionTimeChartData struct {
	Labels   []string       `json:"labels"`
	Datasets []ChartDataset `json:"datasets"`
}

type ChartDataset struct {
	Label           string    `json:"label"`
	Data           []float64 `json:"data"`
	BackgroundColor string    `json:"backgroundColor"`
	BorderColor     string    `json:"borderColor"`
}

// RedirectionTimeReport represents the complete report structure
type RedirectionTimeReport struct {
	Summary      []RedirectionTimeSummary `json:"summary"`
	ChartData    RedirectionTimeChartData `json:"chart_data"`
	KPIMetrics   map[string]float64      `json:"kpi_metrics"`
	ExceededKPIs []string                `json:"exceeded_kpis"`
}

// KPIThresholds defines the KPI thresholds for redirection metrics
type RedirectionKPIStats struct {
	ExceedLoadTimeCount     int `json:"exceed_load_time_count"`
	ExceedResponseTimeCount int `json:"exceed_response_time_count"`
	BelowSuccessRateCount   int `json:"below_success_rate_count"`
}


type RedirectionChartData struct {
	Hour          string  `json:"hour"`             // "01" to "23"
	TotalLoadTime float64 `json:"total_load_time"`  // averaged value
	ResponseTime  float64 `json:"response_time"`    // averaged value
}
