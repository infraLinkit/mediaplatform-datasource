package entity

import (
	"time"

	"gorm.io/gorm"
)

type (
	DisplayCampaignSummary struct {
		DataType       string   `form:"data_type" json:"data_type"`
		ReportType     string   `form:"report_type" json:"report_type"`
		Country        string   `form:"country" json:"country"`
		Operator       string   `form:"operator" json:"operator"`
		PartnerName    string   `form:"partner_name" json:"partner_name"`
		Adnet          string   `form:"adnet" json:"adnet"`
		Service        string   `form:"service" json:"service"`
		DataIndicators []string `form:"data_indicators" json:"data_indicators"`
		DataBasedOn    string   `form:"data_based_on" json:"data_based_on"`
		DateRange      string   `form:"date_range" json:"date_range"`
		CustomRange    string   `form:"custom_range" json:"custom_range"`
	}

	CampaignSummaryMonitoring struct {
		gorm.Model
		ID                 int       `gorm:"primaryKey;autoIncrement" json:"id"`
		Status             bool      `gorm:"not null;size:50" json:"status"`
		SummaryDate        time.Time `gorm:"type:date" json:"summary_date"`
		CampaignId         string    `gorm:"index:idx_campdetailid_unique;not null;size:50" json:"campaign_id"`
		CampaignName       string    `gorm:"not null;size:100" json:"campaign_name"`
		Country            string    `gorm:"not null;size:50" json:"country"`
		Operator           string    `gorm:"not null;size:50" json:"operator"`
		Partner            string    `gorm:"not null;size:50" json:"partner"`
		Adnet              string    `gorm:"not null;size:50" json:"adnet"`
		Service            string    `gorm:"not null;size:50" json:"service"`
		Traffic            int       `gorm:"not null;length:20;default:0" json:"traffic"`
		MoReceived         int       `gorm:"not null;length:20;default:0" json:"mo_received"`
		CPA                float64   `gorm:"type:double precision" json:"cpa"`
		AgencyFee          float64   `gorm:"type:double precision" json:"agency_fee"`
		TargetDailyBudget  float64   `gorm:"type:double precision" json:"target_daily_budget"`
		CrMO               float64   `gorm:"type:double precision" json:"cr_mo"`
		CrPostback         float64   `gorm:"type:double precision" json:"cr_postback"`
		BudgetUsage        float64   `gorm:"type:double precision" json:"budget_usage"`
		WakiRevenue        float64   `gorm:"type:double precision;not null;length:20;default:0" json:"waki_revenue"`
		Fp                 float64   `gorm:"type:double precision;not null;length:20;default:0" json:"fp"`
		MoSent             float64   `gorm:"type:double precision;not null;length:20;default:0" json:"mo_sent"`
		SpendingToAdnets   float64   `gorm:"type:double precision;not null;length:20;default:0" json:"spending_to_adnets"`
		TotalSpending      float64   `gorm:"type:double precision;not null;length:20;default:0" json:"total_spending"`
		TotalWakiAgencyFee float64   `gorm:"type:double precision" json:"total_waki_agency_fee"`
	}
)
type Tabler interface {
	TableName() string
}

func (CampaignSummaryMonitoring) TableName() string {
	return "summary_campaigns"
}
