package entity

import (
	"time"

	"gorm.io/gorm"
)

type (
	SummaryMo struct {
		gorm.Model
		ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
		SummaryDate  time.Time `gorm:"type:timestamptz;index:uidx_summary_mo_check_data_unique,uidx_summary_mo_check_data_perpic" json:"summary_date"`
		CampaignId   string    `gorm:"size:80;index:uidx_summary_mo_check_data_campaign_id,uidx_summary_mo_check_data_unique,uidx_summary_mo_check_data_perpic" json:"campaign_id"`
		CampaignName string    `gorm:"size:50" json:"campaign_name"`
		Partner      string    `gorm:"size:30;index:uidx_summary_mo_check_data_unique,uidx_summary_mo_check_data_perpic" json:"partner"`
		Country      string    `gorm:"size:30;index:uidx_summary_mo_check_data_unique,uidx_summary_mo_check_data_perpic" json:"country"`
		Operator     string    `gorm:"size:30;index:uidx_summary_mo_check_data_unique,uidx_summary_mo_check_data_perpic" json:"operator"`
		Service      string    `gorm:"size:80;index:uidx_summary_mo_check_data_unique,uidx_summary_mo_check_data_perpic" json:"service"`
		Status       bool      `gorm:"default:false" json:"status"`
		Traffic      int       `gorm:"default:0" json:"traffic"`
		Mo           int       `gorm:"default:0" json:"mo"`
		MoLimit      int       `gorm:"default:0" json:"mo_limit"`
		MoCounter1   int       `gorm:"default:0" json:"mo_counter_1"`
		MoCounter2   int       `gorm:"default:0" json:"mo_counter_2"`
		MoCounter3   int       `gorm:"default:0" json:"mo_counter_3"`
		MoCounter4   int       `gorm:"default:0" json:"mo_counter_4"`
		MoCounter5   int       `gorm:"default:0" json:"mo_counter_5"`
		MoCounter6   int       `gorm:"default:0" json:"mo_counter_6"`
		DescMo       string    `gorm:"type:text" json:"desc_mo"`
		PicMoReport  string    `gorm:"size:50;default:''" json:"pic_mo_report"`
		Adnet        string    `gorm:"size:30" json:"adnet"`
		IsLock       bool      `gorm:"default:false" json:"is_lock"`
		SentQueue    bool      `gorm:"default:false" json:"sent_queue"`
		CreatedAt    time.Time
		UpdatedAt    time.Time
	}

	SummaryCr struct {
		gorm.Model
		ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
		SummaryDate  time.Time `gorm:"type:timestamptz;index:uidx_summary_cr_check_data_unique,uidx_summary_cr_check_data_perpic" json:"summary_date"`
		CampaignId   string    `gorm:"size:80;index:uidx_summary_cr_check_data_campaign_id,uidx_summary_cr_check_data_unique,uidx_summary_cr_check_data_perpic" json:"campaign_id"`
		CampaignName string    `gorm:"size:50" json:"campaign_name"`
		Partner      string    `gorm:"size:30;index:uidx_summary_cr_check_data_unique,uidx_summary_cr_check_data_perpic" json:"partner"`
		Country      string    `gorm:"size:30;index:uidx_summary_cr_check_data_unique,uidx_summary_cr_check_data_perpic" json:"country"`
		Operator     string    `gorm:"size:30;index:uidx_summary_cr_check_data_unique,uidx_summary_cr_check_data_perpic" json:"operator"`
		Service      string    `gorm:"size:80;index:uidx_summary_cr_check_data_unique,uidx_summary_cr_check_data_perpic" json:"service"`
		Status       bool      `gorm:"default:false" json:"status"`
		Cr           int       `gorm:"default:0" json:"cr"`
		CrCounter1   int       `gorm:"default:0" json:"cr_counter_1"`
		CrCounter2   int       `gorm:"default:0" json:"cr_counter_2"`
		CrCounter3   int       `gorm:"default:0" json:"cr_counter_3"`
		CrCounter4   int       `gorm:"default:0" json:"cr_counter_4"`
		CrCounter5   int       `gorm:"default:0" json:"cr_counter_5"`
		CrCounter6   int       `gorm:"default:0" json:"cr_counter_6"`
		DescCr       string    `gorm:"type:text" json:"desc_cr"`
		PicCrReport  string    `gorm:"size:50;default:''" json:"pic_cr_report"`
		Adnet        string    `gorm:"size:30" json:"adnet"`
		IsLock       bool      `gorm:"default:false" json:"is_lock"`
		SentQueue    bool      `gorm:"default:false" json:"sent_queue"`
		CreatedAt    time.Time
		UpdatedAt    time.Time
	}

	SummaryCapping struct {
		gorm.Model
		ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
		SummaryDate  time.Time `gorm:"type:timestamptz;index:uidx_summary_capping_check_data_unique,uidx_summary_capping_check_data_perpic" json:"summary_date"`
		CampaignId   string    `gorm:"size:80;index:uidx_summary_capping_check_data_campaign_id,uidx_summary_capping_check_data_unique,uidx_summary_capping_check_data_perpic" json:"campaign_id"`
		CampaignName string    `gorm:"size:50" json:"campaign_name"`
		Partner      string    `gorm:"size:30;index:uidx_summary_capping_check_data_unique,uidx_summary_capping_check_data_perpic" json:"partner"`
		Country      string    `gorm:"size:30;index:uidx_summary_capping_check_data_unique,uidx_summary_capping_check_data_perpic" json:"country"`
		Operator     string    `gorm:"size:30;index:uidx_summary_capping_check_data_unique,uidx_summary_capping_check_data_perpic" json:"operator"`
		Service      string    `gorm:"size:80;index:uidx_summary_capping_check_data_unique,uidx_summary_capping_check_data_perpic" json:"service"`
		Status       bool      `gorm:"default:false" json:"status"`
		Mo           int       `gorm:"default:0" json:"mo"`
		MoLimit      int       `gorm:"default:0" json:"mo_limit"`
		DescMo       string    `gorm:"type:text" json:"desc_mo"`
		PicMoReport  string    `gorm:"size:50;default:''" json:"pic_mo_report"`
		Adnet        string    `gorm:"size:30" json:"adnet"`
		IsLock       bool      `gorm:"default:false" json:"is_lock"`
		SentQueue    bool      `gorm:"default:false" json:"sent_queue"`
		CreatedAt    time.Time
		UpdatedAt    time.Time
	}

	SummaryRatio struct {
		gorm.Model
		ID             int       `gorm:"primaryKey;autoIncrement" json:"id"`
		SummaryDate    time.Time `gorm:"type:timestamptz;index:uidx_summary_ratio_check_data_unique,uidx_summary_ratio_check_data_perpic" json:"summary_date"`
		CampaignId     string    `gorm:"size:80;index:uidx_summary_ratio_check_data_campaign_id,uidx_summary_ratio_check_data_percampaign_id,uidx_summary_ratio_check_data_unique,uidx_summary_ratio_check_data_perpic" json:"campaign_id"`
		CampaignName   string    `gorm:"size:50" json:"campaign_name"`
		Partner        string    `gorm:"size:30;index:uidx_summary_ratio_check_data_percampaign_id,uidx_summary_ratio_check_data_unique,uidx_summary_ratio_check_data_perpic" json:"partner"`
		Country        string    `gorm:"size:30;index:uidx_summary_ratio_check_data_percampaign_id,uidx_summary_ratio_check_data_unique,uidx_summary_ratio_check_data_perpic" json:"country"`
		Operator       string    `gorm:"size:30;index:uidx_summary_ratio_check_data_percampaign_id,uidx_summary_ratio_check_data_unique,uidx_summary_ratio_check_data_perpic" json:"operator"`
		Service        string    `gorm:"size:80;index:uidx_summary_ratio_check_data_percampaign_id,uidx_summary_ratio_check_data_unique,uidx_summary_ratio_check_data_perpic" json:"service"`
		Status         bool      `gorm:"default:false" json:"status"`
		Ratio          string    `gorm:"default:'0'" json:"ratio"`
		ActualRatio    string    `gorm:"default:'0'" json:"actual_ratio"`
		Mo             int       `gorm:"default:0" json:"mo"`
		Postback       int       `gorm:"default:0" json:"postback"`
		DescRatio      string    `gorm:"type:text" json:"desc_ratio"`
		PicRatioReport string    `gorm:"size:50;default:''" json:"pic_ratio_report"`
		Adnet          string    `gorm:"size:30" json:"adnet"`
		IsLock         bool      `gorm:"default:false" json:"is_lock"`
		SentQueue      bool      `gorm:"default:false" json:"sent_queue"`
		CreatedAt      time.Time
		UpdatedAt      time.Time
	}

	SummaryAll struct {
		ID           int       `json:"id"`
		SummaryDate  time.Time `json:"summary_date"`
		CampaignId   string    `json:"campaign_id"`
		CampaignName string    `json:"campaign_name"`
		Partner      string    `json:"partner"`
		Country      string    `json:"country"`
		Operator     string    `json:"operator"`
		Service      string    `json:"service"`
		Status       bool      `json:"status"`
		IsLock       bool      `json:"is_lock"`
		Adnet        string    `json:"adnet"`
		SentQueue    bool      `json:"sent_queue"`

		Error string `json:"error"`
		// summary_mo & summary_capping
		Traffic     int    `json:"traffic"`
		Mo          int    `json:"mo"`       //cap & ratio
		MoLimit     int    `json:"mo_limit"` //cap
		MoCounter1  int    `json:"mo_counter_1"`
		MoCounter2  int    `json:"mo_counter_2"`
		MoCounter3  int    `json:"mo_counter_3"`
		MoCounter4  int    `json:"mo_counter_4"`
		MoCounter5  int    `json:"mo_counter_5"`
		MoCounter6  int    `json:"mo_counter_6"`
		DescMo      string `json:"desc_mo"`       //cap
		PicMoReport string `json:"pic_mo_report"` //cap

		// summary_ratio
		Ratio          string `json:"ratio"`
		ActualRatio    string `json:"actual_ratio"`
		Postback       int    `json:"postback"`
		DescRatio      string `json:"desc_ratio"`
		PicRatioReport string `json:"pic_ratio_report"`

		// summary_cr
		Cr          int    `json:"cr"`
		CrCounter1  int    `json:"cr_counter_1"`
		CrCounter2  int    `json:"cr_counter_2"`
		CrCounter3  int    `json:"cr_counter_3"`
		CrCounter4  int    `json:"cr_counter_4"`
		CrCounter5  int    `json:"cr_counter_5"`
		CrCounter6  int    `json:"cr_counter_6"`
		DescCr      string `json:"desc_cr"`
		PicCrReport string `json:"pic_cr_report"`
	}
)

func (SummaryMo) TableName() string {
	return "summary_mos"
}
func (SummaryCr) TableName() string {
	return "summary_crs"
}
func (SummaryCapping) TableName() string {
	return "summary_cappings"
}
func (SummaryRatio) TableName() string {
	return "summary_ratios"
}
