package entity

import "time"

// ParamsCampaignSpendingChannel holds all filter parameters for the spending channel report.
type ParamsCampaignSpendingChannel struct {
	DataType        string
	DataIndicators  []string
	ReportType      string
	Country         string
	Operator        string
	PartnerName     string
	CampaignName    string
	Adnet           string
	Service         string
	ChannelCampaign string
	DataBasedOn     string
	DateRange       string
	DateStart       string
	DateEnd         string
	DateCustomRange string
	All             string
	ChartType       string
	UrlServiceKey   string
	ViewType        string
}

// CampaignSpendingChannelMonitoring is the row structure returned by GetSpendingChannelMonitoring.
// It includes the Channel column from summary_campaigns so we can group by canonical channel name.
type CampaignSpendingChannelMonitoring struct {
	Country      string    `gorm:"column:country"`
	Operator     string    `gorm:"column:operator"`
	Partner      string    `gorm:"column:partner"`
	Service      string    `gorm:"column:service"`
	Adnet        string    `gorm:"column:adnet"`
	// Channel is the raw channel value from summary_campaigns.channel
	// (e.g. "cpa", "google traffic", "tiktok", "fbmeta", "dsp", "s2s", "telco_channel")
	// For api_pin_reports rows this is set to "API" in the SQL.
	Channel      string    `gorm:"column:channel"`
	UrlServiceKey string   `gorm:"column:url_service_key"`
	CampaignName string    `gorm:"column:campaign_name"`
	CampaignId   string    `gorm:"column:campaign_id"`
	SummaryDate  time.Time `gorm:"column:summary_date"`
	SBAF         float64   `gorm:"column:sbaf"`
}