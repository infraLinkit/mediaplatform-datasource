package entity

import (
	"time"

	"gorm.io/gorm"
)

type (
	Campaign struct {
		gorm.Model
		ID                int    `gorm:"primaryKey;autoIncrement" json:"id"`
		CampaignId        string `gorm:"index:idx_campid_unique;not null;size:50" json:"campaign_id"`
		Name              string `gorm:"not null;size:50" json:"name"`
		CampaignObjective string `gorm:"not null;size:50" json:"campaign_objective"`
		Country           string `gorm:"not null;size:50" json:"country"`
		Advertiser        string `gorm:"not null;size:50" json:"advertiser"`
		CreatedAt         time.Time
		UpdatedAt         time.Time
	}

	CampaignDetail struct {
		gorm.Model
		ID                        int       `gorm:"primaryKey;autoIncrement" json:"id"`
		URLServiceKey             string    `gorm:"index:idx_urlservicekey;not null;size:50" json:"urlservicekey"`
		CampaignId                string    `gorm:"index:idx_campdetailid_unique;not null;size:50" json:"campaign_id"`
		Country                   string    `gorm:"not null;size:50" json:"country"`
		Operator                  string    `gorm:"not null;size:50" json:"operator"`
		Partner                   string    `gorm:"not null;size:50" json:"partner"`
		Aggregator                string    `gorm:"not null;size:50" json:"aggregator"`
		Adnet                     string    `gorm:"not null;size:50" json:"adnet"`
		Service                   string    `gorm:"not null;size:50" json:"service"`
		Keyword                   string    `gorm:"not null;size:50" json:"keyword"`
		Subkeyword                string    `gorm:"not null;size:50" json:"subkeyword"`
		IsBillable                bool      `gorm:"not null;size:50" json:"is_billable"`
		Plan                      string    `gorm:"not null;size:50" json:"plan"`
		PO                        string    `gorm:"size:50" json:"po"`
		Cost                      string    `gorm:"not null;size:50" json:"cost"`
		PubId                     string    `gorm:"not null;size:50" json:"pubid"`
		ShortCode                 string    `gorm:"not null;size:50" json:"short_code"`
		DeviceType                string    `gorm:"not null;size:50" json:"device_type"`
		OS                        string    `gorm:"not null;size:50" json:"os"`
		URLType                   string    `gorm:"not null;size:50;default:wap" json:"url_type"`
		ClickType                 int       `gorm:"not null;length:1;default:1" json:"click_type"`
		ClickDelay                int       `gorm:"not null;length:1;default:1" json:"click_delay"`
		ClientType                string    `gorm:"not null;size:50" json:"client_type"`
		TrafficSource             bool      `gorm:"not null;default:false" json:"traffic_source"`
		UniqueClick               bool      `gorm:"not null;default:false" json:"unique_click"`
		URLBanner                 string    `gorm:"size:255;default:NA" json:"url_banner"`
		URLLanding                string    `gorm:"size:255;default:NA" json:"url_landing"`
		URLWarpLanding            string    `gorm:"size:255;default:NA" json:"url_warp_landing"`
		URLService                string    `gorm:"size:255;default:NA" json:"url_service"`
		URLTFCORSmartlink         string    `gorm:"size:255;default:NA" json:"url_tfc_or_smartlink"`
		GlobPost                  bool      `gorm:"not null;default:false" json:"glob_post"`
		URLGlobPost               string    `gorm:"size:255;default:NA" json:"url_glob_post"`
		CustomIntegration         string    `gorm:"size:30;default:NA" json:"custom_integration"`
		IpAddress                 []string  `gorm:"type:text[]"`
		IsActive                  bool      `gorm:"not null;default:true" json:"is_active"`
		MOCapping                 int       `gorm:"not null;length:10;default:0" json:"mo_capping"`
		CounterMOCapping          int       `gorm:"not null;length:10;default:0" json:"counter_mo_capping"`
		StatusCapping             bool      `gorm:"not null;default:false" json:"status_capping"`
		KPIUpperLimitCapping      int       `gorm:"not null;length:20;default:1" json:"kpi_upper_limit_capping"`
		IsMachineLearningCapping  bool      `gorm:"not null;default:false" json:"is_machine_learning_capping"`
		RatioSend                 int       `gorm:"not null;length:10;default:1" json:"ratio_send"`
		RatioReceive              int       `gorm:"not null;length:10;default:4" json:"ratio_receive"`
		CounterMORatio            int       `gorm:"not null;length:10;default:0" json:"counter_mo_ratio"`
		StatusRatio               bool      `gorm:"not null;default:false" json:"status_ratio"`
		KPIUpperLimitRatioSend    int       `gorm:"not null;length:10;default:1" json:"kpi_upper_limit_ratio_send"`
		KPIUpperLimitRatioReceive int       `gorm:"not null;length:10;default:4" json:"kpi_upper_limit_ratio_receive"`
		IsMachineLearningRatio    bool      `gorm:"not null;default:false" json:"is_machine_learning_ratio"`
		APIURL                    string    `gorm:"size:255;default:NA" json:"api_url"`
		LastUpdate                time.Time `json:"last_update"`
		LastUpdateCapping         time.Time `json:"last_update_capping"`
		CostPerConversion         float64   `gorm:"type:double precision" json:"cost_per_conversion"`
		AgencyFee                 float64   `gorm:"type:double precision" json:"agency_fee"`
		TargetDailyBudget         float64   `gorm:"type:double precision" json:"target_daily_budget"`
		URLPostback               string    `gorm:"size:255;default:NA" json:"url_postback"`
		CreatedAt                 time.Time
		UpdatedAt                 time.Time
	}

	ResultCampaign struct {
		Name                      string
		CampaignObjective         string
		Advertiser                string
		ID                        int
		URLServiceKey             string
		CampaignId                string
		Country                   string
		Operator                  string
		Partner                   string
		Aggregator                string
		Adnet                     string
		Service                   string
		Keyword                   string
		Subkeyword                string
		IsBillable                bool
		Plan                      string
		PO                        string
		Cost                      string
		PubId                     string
		ShortCode                 string
		DeviceType                string
		OS                        string
		URLType                   string
		ClickType                 int
		ClickDelay                int
		ClientType                string
		TrafficSource             bool
		UniqueClick               bool
		URLBanner                 string
		URLLanding                string
		URLWarpLanding            string
		URLService                string
		URLTFCORSmartlink         string
		GlobPost                  bool
		URLGlobPost               string
		CustomIntegration         string
		IpAddress                 []string
		IsActive                  bool
		MOCapping                 int
		CounterMOCapping          int
		StatusCapping             bool
		KPIUpperLimitCapping      int
		IsMachineLearningCapping  bool
		RatioSend                 int
		RatioReceive              int
		CounterMORatio            int
		StatusRatio               bool
		KPIUpperLimitRatioSend    int
		KPIUpperLimitRatioReceive int
		IsMachineLearningRatio    bool
		APIURL                    string
		LastUpdate                time.Time
		LastUpdateCapping         time.Time
		CostPerConversion         float64
		AgencyFee                 float64
		TargetDailyBudget         float64
		URLPostback               string
	}

	MO struct {
		gorm.Model
		ID                int       `gorm:"primaryKey;autoIncrement" json:"id"`
		CampaignDetailId  int       `gorm:"index:idx_mocampdetailid_unique;not null;length:20" json:"campaign_detail_id"`
		Pxdate            time.Time `gorm:"not null" json:"pxdate"`
		URLServiceKey     string    `gorm:"index:idx_urlservicekey;not null;size:50" json:"urlservicekey"`
		CampaignId        string    `gorm:"index:idx_campdetailid_unique;not null;size:50" json:"campaign_id"`
		Country           string    `gorm:"not null;size:50" json:"country"`
		Operator          string    `gorm:"not null;size:50" json:"operator"`
		Partner           string    `gorm:"not null;size:50" json:"partner"`
		Aggregator        string    `gorm:"size:50" json:"aggregator"`
		Adnet             string    `gorm:"not null;size:50" json:"adnet"`
		Service           string    `gorm:"not null;size:50" json:"service"`
		Keyword           string    `gorm:"size:50" json:"keyword"`
		Subkeyword        string    `gorm:"size:50" json:"subkeyword"`
		IsBillable        bool      `gorm:"default:false" json:"is_billable"`
		Plan              string    `gorm:"size:50;default:NA" json:"plan"`
		PO                string    `gorm:"size:50;default:NA" json:"po"`
		Cost              string    `gorm:"size:50;default:NA" json:"cost"`
		PubId             string    `gorm:"size:50;default:NA" json:"pubid"`
		ShortCode         string    `gorm:"size:50;default:NA" json:"short_code"`
		URL               string    `gorm:"size:255;default:NA" json:"url"`
		URLType           string    `gorm:"size:50;default:NA" json:"url_type"`
		Pixel             string    `gorm:"size:255;default:NA" json:"pixel"`
		Token             string    `gorm:"size:255;default:NA" json:"token"`
		TrxId             string    `gorm:"size:255;default:NA" json:"trx_id"`
		Msisdn            string    `gorm:"size:255;default:NA" json:"msisdn"`
		IsUsed            bool      `gorm:"not null;default:false" json:"is_used"`
		Browser           string    `gorm:"size:50;default:NA" json:"browser"`
		OS                string    `gorm:"size:50;default:NA" json:"os"`
		Ip                string    `gorm:"size:50;default:NA" json:"ip"`
		ISP               string    `gorm:"size:50;default:NA" json:"isp"`
		ReferralURL       string    `gorm:"size:255;default:NA" json:"referral_url"`
		UserAgent         string    `gorm:"size:255;default:NA" json:"user_agent"`
		TrafficSource     bool      `gorm:"not null;default:false" json:"traffic_source"`
		TrafficSourceData string    `gorm:"size:255;default:NA" json:"traffic_source_data"`
		UserRejected      bool      `gorm:"not null;default:false" json:"user_rejected"`
		UniqueClick       bool      `gorm:"not null;default:false" json:"unique_click"`
		UserDuplicated    bool      `gorm:"not null;default:false" json:"user_duplicated"`
		Handset           string    `gorm:"size:255;default:NA" json:"handset"`
		HandsetCode       string    `gorm:"size:50;default:NA" json:"handset_code"`
		HandsetType       string    `gorm:"size:50;default:NA" json:"handset_type"`
		URLLanding        string    `gorm:"size:255;default:NA" json:"url_landing"`
		URLWarpLanding    string    `gorm:"size:255;default:NA" json:"url_warp_landing"`
		URLService        string    `gorm:"size:255;default:NA" json:"url_service"`
		URLTFCORSmartlink string    `gorm:"size:255;default:NA" json:"url_tfc_or_smartlink"`
		PixelUsedDate     time.Time `gorm:"not null" json:"pixel_used_date"`
		StatusPostback    bool      `gorm:"not null;default:false" json:"status_postback"`
		IsUnique          bool      `gorm:"not null;default:false" json:"is_unique"`
		URLPostback       string    `gorm:"size:255;default:NA" json:"url_postback"`
		StatusURLPostback string    `gorm:"size:50" json:"status_url_postback"`
		ReasonURLPostback string    `gorm:"size:255" json:"reason_url_postback"`
		IsActive          bool      `gorm:"not null;default:false" json:"is_active"`
		CounterMOCapping  int       `gorm:"not null;length:10" json:"counter_mo_capping"`
		MOCapping         int       `gorm:"not null;length:10" json:"mo_capping"`
		StatusCapping     bool      `gorm:"not null;default:false" json:"status_capping"`
		CounterMORatio    int       `gorm:"not null;length:10" json:"counter_mo_ratio"`
		RatioSend         int       `gorm:"not null;length:10;default:1" json:"ratio_send"`
		RatioReceive      int       `gorm:"not null;length:10;default:4" json:"ratio_receive"`
		StatusRatio       bool      `gorm:"not null;default:false" json:"status_ratio"`
		APIURL            string    `gorm:"size:255;default:NA" json:"api_url"`
		CreatedAt         time.Time
		UpdatedAt         time.Time
	}

	PixelStorage struct {
		gorm.Model
		ID                int       `gorm:"primaryKey;autoIncrement" json:"id"`
		CampaignDetailId  int       `gorm:"index:idx_mocampdetailid_unique;not null;length:20" json:"campaign_detail_id"`
		Pxdate            time.Time `gorm:"not null" json:"pxdate"`
		URLServiceKey     string    `gorm:"index:idx_urlservicekey;not null;size:50" json:"urlservicekey"`
		CampaignId        string    `gorm:"index:idx_campdetailid_unique;not null;size:50" json:"campaign_id"`
		Country           string    `gorm:"not null;size:50" json:"country"`
		Operator          string    `gorm:"not null;size:50" json:"operator"`
		Partner           string    `gorm:"not null;size:50" json:"partner"`
		Aggregator        string    `gorm:"size:50" json:"aggregator"`
		Adnet             string    `gorm:"not null;size:50" json:"adnet"`
		Service           string    `gorm:"not null;size:50" json:"service"`
		Keyword           string    `gorm:"size:50" json:"keyword"`
		Subkeyword        string    `gorm:"size:50" json:"subkeyword"`
		IsBillable        bool      `gorm:"default:false" json:"is_billable"`
		Plan              string    `gorm:"size:50;default:NA" json:"plan"`
		PO                string    `gorm:"size:50;default:NA" json:"po"`
		Cost              string    `gorm:"size:50;default:NA" json:"cost"`
		PubId             string    `gorm:"size:50;default:NA" json:"pubid"`
		ShortCode         string    `gorm:"size:50;default:NA" json:"short_code"`
		URL               string    `gorm:"size:255;default:NA" json:"url"`
		URLType           string    `gorm:"size:50;default:NA" json:"url_type"`
		Pixel             string    `gorm:"size:255;default:NA" json:"pixel"`
		Token             string    `gorm:"size:255;default:NA" json:"token"`
		TrxId             string    `gorm:"size:255;default:NA" json:"trx_id"`
		Msisdn            string    `gorm:"size:255;default:NA" json:"msisdn"`
		IsUsed            bool      `gorm:"not null;default:false" json:"is_used"`
		Browser           string    `gorm:"size:50;default:NA" json:"browser"`
		OS                string    `gorm:"size:50;default:NA" json:"os"`
		Ip                string    `gorm:"size:50;default:NA" json:"ip"`
		ISP               string    `gorm:"size:50;default:NA" json:"isp"`
		ReferralURL       string    `gorm:"size:255;default:NA" json:"referral_url"`
		UserAgent         string    `gorm:"size:255;default:NA" json:"user_agent"`
		TrafficSource     bool      `gorm:"not null;default:false" json:"traffic_source"`
		TrafficSourceData string    `gorm:"size:255;default:NA" json:"traffic_source_data"`
		UserRejected      bool      `gorm:"not null;default:false" json:"user_rejected"`
		UniqueClick       bool      `gorm:"not null;default:false" json:"unique_click"`
		UserDuplicated    bool      `gorm:"not null;default:false" json:"user_duplicated"`
		Handset           string    `gorm:"size:255;default:NA" json:"handset"`
		HandsetCode       string    `gorm:"size:50;default:NA" json:"handset_code"`
		HandsetType       string    `gorm:"size:50;default:NA" json:"handset_type"`
		URLLanding        string    `gorm:"size:255;default:NA" json:"url_landing"`
		URLWarpLanding    string    `gorm:"size:255;default:NA" json:"url_warp_landing"`
		URLService        string    `gorm:"size:255;default:NA" json:"url_service"`
		URLTFCORSmartlink string    `gorm:"size:255;default:NA" json:"url_tfc_or_smartlink"`
		PixelUsedDate     time.Time `gorm:"not null" json:"pixel_used_date"`
		StatusPostback    bool      `gorm:"not null;default:false" json:"status_postback"`
		IsUnique          bool      `gorm:"not null;default:false" json:"is_unique"`
		URLPostback       string    `gorm:"size:255;default:NA" json:"url_postback"`
		StatusURLPostback string    `gorm:"size:50" json:"status_url_postback"`
		ReasonURLPostback string    `gorm:"size:255" json:"reason_url_postback"`
		IsActive          bool      `gorm:"not null;default:false" json:"is_active"`
		CounterMOCapping  int       `gorm:"not null;length:10" json:"counter_mo_capping"`
		MOCapping         int       `gorm:"not null;length:10" json:"mo_capping"`
		StatusCapping     bool      `gorm:"not null;default:false" json:"status_capping"`
		CounterMORatio    int       `gorm:"not null;length:10" json:"counter_mo_ratio"`
		RatioSend         int       `gorm:"not null;length:10;default:1" json:"ratio_send"`
		RatioReceive      int       `gorm:"not null;length:10;default:4" json:"ratio_receive"`
		StatusRatio       bool      `gorm:"not null;default:false" json:"status_ratio"`
		APIURL            string    `gorm:"size:255;default:NA" json:"api_url"`
		CreatedAt         time.Time
		UpdatedAt         time.Time
	}

	Postback struct {
		gorm.Model
		ID                int       `gorm:"primaryKey;autoIncrement" json:"id"`
		CampaignDetailId  int       `gorm:"index:idx_mocampdetailid_unique;not null;length:20" json:"campaign_detail_id"`
		Pxdate            time.Time `gorm:"not null" json:"pxdate"`
		URLServiceKey     string    `gorm:"index:idx_urlservicekey;not null;size:50" json:"urlservicekey"`
		CampaignId        string    `gorm:"index:idx_campdetailid_unique;not null;size:50" json:"campaign_id"`
		Country           string    `gorm:"not null;size:50" json:"country"`
		Operator          string    `gorm:"not null;size:50" json:"operator"`
		Partner           string    `gorm:"not null;size:50" json:"partner"`
		Aggregator        string    `gorm:"size:50" json:"aggregator"`
		Adnet             string    `gorm:"not null;size:50" json:"adnet"`
		Service           string    `gorm:"not null;size:50" json:"service"`
		Keyword           string    `gorm:"size:50" json:"keyword"`
		Subkeyword        string    `gorm:"size:50" json:"subkeyword"`
		IsBillable        bool      `gorm:"default:false" json:"is_billable"`
		Plan              string    `gorm:"size:50;default:NA" json:"plan"`
		PO                string    `gorm:"size:50;default:NA" json:"po"`
		Cost              string    `gorm:"size:50;default:NA" json:"cost"`
		PubId             string    `gorm:"size:50;default:NA" json:"pubid"`
		ShortCode         string    `gorm:"size:50;default:NA" json:"short_code"`
		URL               string    `gorm:"size:255;default:NA" json:"url"`
		URLType           string    `gorm:"size:50;default:NA" json:"url_type"`
		Pixel             string    `gorm:"size:255;default:NA" json:"pixel"`
		Token             string    `gorm:"size:255;default:NA" json:"token"`
		TrxId             string    `gorm:"size:255;default:NA" json:"trx_id"`
		Msisdn            string    `gorm:"size:255;default:NA" json:"msisdn"`
		IsUsed            bool      `gorm:"not null;default:false" json:"is_used"`
		Browser           string    `gorm:"size:50;default:NA" json:"browser"`
		OS                string    `gorm:"size:50;default:NA" json:"os"`
		Ip                string    `gorm:"size:50;default:NA" json:"ip"`
		ISP               string    `gorm:"size:50;default:NA" json:"isp"`
		ReferralURL       string    `gorm:"size:255;default:NA" json:"referral_url"`
		UserAgent         string    `gorm:"size:255;default:NA" json:"user_agent"`
		TrafficSource     bool      `gorm:"not null;default:false" json:"traffic_source"`
		TrafficSourceData string    `gorm:"size:255;default:NA" json:"traffic_source_data"`
		UserRejected      bool      `gorm:"not null;default:false" json:"user_rejected"`
		UniqueClick       bool      `gorm:"not null;default:false" json:"unique_click"`
		UserDuplicated    bool      `gorm:"not null;default:false" json:"user_duplicated"`
		Handset           string    `gorm:"size:255;default:NA" json:"handset"`
		HandsetCode       string    `gorm:"size:50;default:NA" json:"handset_code"`
		HandsetType       string    `gorm:"size:50;default:NA" json:"handset_type"`
		URLLanding        string    `gorm:"size:255;default:NA" json:"url_landing"`
		URLWarpLanding    string    `gorm:"size:255;default:NA" json:"url_warp_landing"`
		URLService        string    `gorm:"size:255;default:NA" json:"url_service"`
		URLTFCORSmartlink string    `gorm:"size:255;default:NA" json:"url_tfc_or_smartlink"`
		PixelUsedDate     time.Time `gorm:"not null" json:"pixel_used_date"`
		StatusPostback    bool      `gorm:"not null;default:false" json:"status_postback"`
		IsUnique          bool      `gorm:"not null;default:false" json:"is_unique"`
		URLPostback       string    `gorm:"size:255;default:NA" json:"url_postback"`
		StatusURLPostback string    `gorm:"size:50" json:"status_url_postback"`
		ReasonURLPostback string    `gorm:"size:255" json:"reason_url_postback"`
		IsActive          bool      `gorm:"not null;default:false" json:"is_active"`
		CounterMOCapping  int       `gorm:"not null;length:10" json:"counter_mo_capping"`
		MOCapping         int       `gorm:"not null;length:10" json:"mo_capping"`
		StatusCapping     bool      `gorm:"not null;default:false" json:"status_capping"`
		CounterMORatio    int       `gorm:"not null;length:10" json:"counter_mo_ratio"`
		RatioSend         int       `gorm:"not null;length:10;default:1" json:"ratio_send"`
		RatioReceive      int       `gorm:"not null;length:10;default:4" json:"ratio_receive"`
		StatusRatio       bool      `gorm:"not null;default:false" json:"status_ratio"`
		APIURL            string    `gorm:"size:255;default:NA" json:"api_url"`
		CreatedAt         time.Time
		UpdatedAt         time.Time
	}

	SummaryCampaign struct {
		gorm.Model
		ID                       int       `gorm:"primaryKey;autoIncrement" json:"id"`
		Status                   bool      `gorm:"not null;size:50" json:"status"`
		SummaryDate              time.Time `gorm:"type:date" json:"summary_date"`
		URLServiceKey            string    `gorm:"index:idx_urlservicekey;not null;size:50" json:"urlservicekey"`
		CampaignId               string    `gorm:"index:idx_campdetailid_unique;not null;size:50" json:"campaign_id"`
		CampaignName             string    `gorm:"not null;size:100" json:"campaign_name"`
		Country                  string    `gorm:"not null;size:50" json:"country"`
		Operator                 string    `gorm:"not null;size:50" json:"operator"`
		Partner                  string    `gorm:"not null;size:50" json:"partner"`
		Aggregator               string    `gorm:"not null;size:50" json:"aggregator"`
		Adnet                    string    `gorm:"not null;size:50" json:"adnet"`
		Service                  string    `gorm:"not null;size:50" json:"service"`
		ShortCode                string    `gorm:"not null;size:50" json:"short_code"`
		Traffic                  int       `gorm:"not null;length:20;default:0" json:"traffic"`
		Landing                  int       `gorm:"not null;length:20;default:0" json:"landing"`
		MoReceived               int       `gorm:"not null;length:20;default:0" json:"mo_received"`
		CR                       float64   `gorm:"type:double precision" json:"cr"`
		Postback                 int       `gorm:"not null;length:20;default:0" json:"postback"`
		TotalFP                  int       `gorm:"not null;length:20;default:0" json:"total_fp"`
		SuccessFP                int       `gorm:"not null;length:20;default:0" json:"success_fp"`
		Billrate                 float64   `gorm:"type:double precision" json:"billrate"`
		ROI                      float64   `gorm:"type:double precision" json:"roi"`
		PO                       float64   `gorm:"type:double precision" json:"po"`
		Cost                     float64   `gorm:"type:double precision;not null;length:20;default:0" json:"cost"`
		SBAF                     float64   `gorm:"type:double precision;not null;length:20;default:0" json:"sbaf"`
		SAAF                     float64   `gorm:"type:double precision;not null;length:20;default:0" json:"saaf"`
		CPA                      float64   `gorm:"type:double precision" json:"cpa"`
		Revenue                  float64   `gorm:"type:double precision;not null;length:20;default:0" json:"revenue"`
		URLAfter                 string    `gorm:"size:255;default:NA" json:"url_after"`
		URLBefore                string    `gorm:"size:255;default:NA" json:"url_before"`
		MOLimit                  int       `gorm:"not null;length:10;default:0" json:"mo_limit"`
		RatioSend                int       `gorm:"not null;length:10;default:1" json:"ratio_send"`
		RatioReceive             int       `gorm:"not null;length:10;default:4" json:"ratio_receive"`
		Company                  string    `gorm:"size:255;default:NA" json:"company"`
		ClientType               string    `gorm:"size:30;default:NA" json:"client_type"`
		CostPerConversion        float64   `gorm:"type:double precision" json:"cost_per_conversion"`
		AgencyFee                float64   `gorm:"type:double precision" json:"agency_fee"`
		TargetDailyBudget        float64   `gorm:"type:double precision" json:"target_daily_budget"`
		CrMO                     float64   `gorm:"type:double precision" json:"cr_mo"`
		CrPostback               float64   `gorm:"type:double precision" json:"cr_postback"`
		TotalWakiAgencyFee       float64   `gorm:"type:double precision" json:"total_waki_agency_fee"`
		BudgetUsage              float64   `gorm:"type:double precision" json:"budget_usage"`
		TargetDailyBudgetChanges int       `gorm:"not null;length:12;default:0" json:"target_daily_budget_changes"`
		CreatedAt                time.Time
		UpdatedAt                time.Time
	}

	DataClicked struct {
		gorm.Model
		ID                int       `gorm:"primaryKey;autoIncrement" json:"id"`
		ClickedTime       time.Time `json:"clicked_time"`
		ClickedButtonTime int       `gorm:"not null;length:20" json:"clicked_button_time"`
		HttpStatus        int       `gorm:"not null;length:10" json:"http_status"`
		URLServiceKey     string    `gorm:"index:idx_urlservicekey;not null;size:50" json:"urlservicekey"`
		CampaignId        string    `gorm:"index:idx_campdetailid_unique;not null;size:50" json:"campaign_id"`
		Country           string    `gorm:"not null;size:50" json:"country"`
		Operator          string    `gorm:"not null;size:50" json:"operator"`
		Partner           string    `gorm:"not null;size:50" json:"partner"`
		Aggregator        string    `gorm:"size:50" json:"aggregator"`
		Adnet             string    `gorm:"not null;size:50" json:"adnet"`
		Service           string    `gorm:"not null;size:50" json:"service"`
		ShortCode         string    `gorm:"size:50;default:NA" json:"short_code"`
		Keyword           string    `gorm:"size:50" json:"keyword"`
		Subkeyword        string    `gorm:"size:50" json:"subkeyword"`
		IsBillable        bool      `gorm:"not null;default:false" json:"is_billable"`
		Plan              string    `gorm:"size:50;default:NA" json:"plan"`
		CreatedAt         time.Time
		UpdatedAt         time.Time
	}

	DataLanding struct {
		gorm.Model
		ID            int       `gorm:"primaryKey;autoIncrement" json:"id"`
		LandingTime   time.Time `json:"landing_time"`
		LandedTime    int       `gorm:"not null;length:20" json:"landed_time"`
		HttpStatus    int       `gorm:"not null;length:10" json:"http_status"`
		URLServiceKey string    `gorm:"index:idx_urlservicekey;not null;size:50" json:"urlservicekey"`
		CampaignId    string    `gorm:"index:idx_campdetailid_unique;not null;size:50" json:"campaign_id"`
		Country       string    `gorm:"not null;size:50" json:"country"`
		Operator      string    `gorm:"not null;size:50" json:"operator"`
		Partner       string    `gorm:"not null;size:50" json:"partner"`
		Aggregator    string    `gorm:"size:50" json:"aggregator"`
		Adnet         string    `gorm:"not null;size:50" json:"adnet"`
		Service       string    `gorm:"not null;size:50" json:"service"`
		ShortCode     string    `gorm:"size:50;default:NA" json:"short_code"`
		Keyword       string    `gorm:"size:50" json:"keyword"`
		Subkeyword    string    `gorm:"size:50" json:"subkeyword"`
		IsBillable    bool      `gorm:"not null;default:false" json:"is_billable"`
		Plan          string    `gorm:"size:50;default:NA" json:"plan"`
		CreatedAt     time.Time
		UpdatedAt     time.Time
	}

	DataRedirect struct {
		gorm.Model
		ID                int       `gorm:"primaryKey;autoIncrement" json:"id"`
		RedirectTime      time.Time `json:"redirect_time"`
		RedirectAddedTime int       `gorm:"not null;length:20" json:"redirect_added_time"`
		HttpStatus        int       `gorm:"not null;length:10" json:"http_status"`
		URLServiceKey     string    `gorm:"index:idx_urlservicekey;not null;size:50" json:"urlservicekey"`
		CampaignId        string    `gorm:"index:idx_campdetailid_unique;not null;size:50" json:"campaign_id"`
		Country           string    `gorm:"not null;size:50" json:"country"`
		Operator          string    `gorm:"not null;size:50" json:"operator"`
		Partner           string    `gorm:"not null;size:50" json:"partner"`
		Aggregator        string    `gorm:"size:50" json:"aggregator"`
		Adnet             string    `gorm:"not null;size:50" json:"adnet"`
		Service           string    `gorm:"not null;size:50" json:"service"`
		ShortCode         string    `gorm:"size:50;default:NA" json:"short_code"`
		Keyword           string    `gorm:"size:50" json:"keyword"`
		Subkeyword        string    `gorm:"size:50" json:"subkeyword"`
		IsBillable        bool      `gorm:"not null;default:false" json:"is_billable"`
		Plan              string    `gorm:"size:50;default:NA" json:"plan"`
		CreatedAt         time.Time
		UpdatedAt         time.Time
	}

	DataTraffic struct {
		gorm.Model
		ID               int       `gorm:"primaryKey;autoIncrement" json:"id"`
		TrafficTime      time.Time `json:"traffic_time"`
		TrafficAddedTime int       `gorm:"not null;length:20" json:"traffic_added_time"`
		HttpStatus       int       `gorm:"not null;length:10" json:"http_status"`
		URLServiceKey    string    `gorm:"index:idx_urlservicekey;not null;size:50" json:"urlservicekey"`
		CampaignId       string    `gorm:"index:idx_campdetailid_unique;not null;size:50" json:"campaign_id"`
		Country          string    `gorm:"not null;size:50" json:"country"`
		Operator         string    `gorm:"not null;size:50" json:"operator"`
		Partner          string    `gorm:"not null;size:50" json:"partner"`
		Aggregator       string    `gorm:"size:50" json:"aggregator"`
		Adnet            string    `gorm:"not null;size:50" json:"adnet"`
		Service          string    `gorm:"not null;size:50" json:"service"`
		ShortCode        string    `gorm:"size:50;default:NA" json:"short_code"`
		Keyword          string    `gorm:"size:50" json:"keyword"`
		Subkeyword       string    `gorm:"size:50" json:"subkeyword"`
		IsBillable       bool      `gorm:"not null;default:false" json:"is_billable"`
		Plan             string    `gorm:"size:50;default:NA" json:"plan"`
		CreatedAt        time.Time
		UpdatedAt        time.Time
	}

	ApiPinReport struct {
		gorm.Model
		ID            int       `gorm:"primaryKey;autoIncrement" json:"id"`
		DateSend      time.Time `gorm:"type:date" json:"date_send"`
		Country       string    `gorm:"not null;size:50" json:"country"`
		Company       string    `gorm:"not null;size:50" json:"company"`
		Adnet         string    `gorm:"not null;size:50" json:"adnet"`
		Operator      string    `gorm:"not null;size:50" json:"operator"`
		Service       string    `gorm:"not null;size:50" json:"service"`
		PayoutAdn     string    `gorm:"size:50" json:"payout_adn"`
		PayoutAF      string    `gorm:"size:50" json:"payout_af"`
		TotalMO       int       `gorm:"length:20;default:0" json:"total_mo"`
		TotalPostback int       `gorm:"length:20;default:0" json:"total_postback"`
		SBAF          float64   `gorm:"type:double precision;not null;length:20;default:0" json:"sbaf"`
		SAAF          float64   `gorm:"type:double precision;not null;length:20;default:0" json:"saaf"`
		PricePerMO    float64   `gorm:"type:double precision;not null;length:20;default:0" json:"price_per_mo"`
		WakiRevenue   float64   `gorm:"type:double precision;not null;length:20;default:0" json:"waki_revenue"`
		CreatedAt     time.Time
		UpdatedAt     time.Time
	}

	ApiPinPerformance struct {
		gorm.Model
		ID                  int       `gorm:"primaryKey;autoIncrement" json:"id"`
		DateSend            time.Time `gorm:"type:date" json:"date_send"`
		Country             string    `gorm:"not null;size:50" json:"country"`
		Company             string    `gorm:"not null;size:50" json:"company"`
		Adnet               string    `gorm:"not null;size:50" json:"adnet"`
		Operator            string    `gorm:"not null;size:50" json:"operator"`
		Service             string    `gorm:"not null;size:50" json:"service"`
		PinRequest          int       `gorm:"length:20;default:0" json:"pin_request"`
		UniquePinRequest    int       `gorm:"length:20;default:0" json:"unique_pin_request"`
		PinSent             int       `gorm:"length:20;default:0" json:"pin_sent"`
		PinFailed           int       `gorm:"length:20;default:0" json:"pin_failed"`
		VerifyRequest       int       `gorm:"length:20;default:0" json:"verify_request"`
		VerifyRequestUnique int       `gorm:"length:20;default:0" json:"verify_request_unique"`
		PinOK               int       `gorm:"length:20;default:0" json:"pin_ok"`
		PinNotOK            int       `gorm:"length:20;default:0" json:"pin_not_ok"`
		PinOkSendAdnet      int       `gorm:"length:20;default:0" json:"pin_ok_send_adnet"`
		CPA                 float64   `gorm:"type:double precision;not null;length:20;default:0" json:"cpa"`
		CPAWaki             float64   `gorm:"type:double precision;not null;length:20;default:0" json:"cpa_waki"`
		EstimatedARPU       float64   `gorm:"type:double precision;not null;length:20;default:0" json:"estimated_arpu"`
		SBAF                float64   `gorm:"type:double precision;not null;length:20;default:0" json:"sbaf"`
		SAAF                float64   `gorm:"type:double precision;not null;length:20;default:0" json:"saaf"`
		ChargedMO           float64   `gorm:"type:double precision;not null;length:20;default:0" json:"charged_mo"`
		SubsCR              float64   `gorm:"type:double precision;not null;length:20;default:0" json:"subs_cr"`
		AdnetCR             float64   `gorm:"type:double precision;not null;length:20;default:0" json:"adnet_cr"`
		CAC                 float64   `gorm:"type:double precision;not null;length:20;default:0" json:"cac"`
		PaidCAC             float64   `gorm:"type:double precision;not null;length:20;default:0" json:"paid_cac"`
		CreatedAt           time.Time
		UpdatedAt           time.Time
	}

	Menu struct {
		ID            int    `gorm:"primaryKey;autoIncrement" json:"id"`
		Name          string `gorm:"type:varchar(255);not null" json:"name"`
		Code          string `gorm:"type:varchar(255);not null" json:"code"`
		Url           string `gorm:"type:varchar(255);not null" json:"url"`
		Icon          string `gorm:"type:varchar(255)" json:"icon"`
		Parent        int    `gorm:"type:int" json:"parent"`
		Sort          string `gorm:"type:varchar(255)" json:"sort"`
		ShowOnSidebar bool   `gorm:"type:bool;default:false" json:"show_on_sidebar"`
		Permission    string `gorm:"type:varchar(255)" json:"permission"`
		CreatedAt     time.Time
		UpdatedAt     time.Time
	}

	Role struct {
		ID        int    `gorm:"primaryKey;autoIncrement" json:"id"`
		Name      string `gorm:"type:varchar(255);not null" json:"name"`
		Code      string `gorm:"type:varchar(255);not null" json:"code"`
		CreatedAt time.Time
		UpdatedAt time.Time
		Users     []User `gorm:"foreignKey:RoleID;references:ID"`
	}

	Permission struct {
		ID        int  `gorm:"primaryKey;autoIncrement" json:"id"`
		RoleID    int  `gorm:"type:int" json:"role_id"`
		MenuID    int  `gorm:"type:int" json:"menu_id"`
		Status    bool `gorm:"type:bool;" json:"status"`
		CreatedAt time.Time
		UpdatedAt time.Time

		// Relations
		Role Role `gorm:"foreignKey:RoleID;references:ID"`
		Menu Menu `gorm:"foreignKey:MenuID;references:ID"`
	}

	User struct {
		ID              int            `gorm:"primaryKey;autoIncrement" json:"id"`
		Name            string         `gorm:"type:varchar(255);not null" json:"name"`
		Username        string         `gorm:"type:varchar(255);unique;not null" json:"username"`
		Email           string         `gorm:"type:varchar(255);unique;not null" json:"email"`
		EmailVerifiedAt *time.Time     `json:"email_verified_at,omitempty"`
		Password        string         `gorm:"type:varchar(255);not null" json:"-"`
		RoleID          int            `json:"role_id"`
		LastLogin       *time.Time     `json:"last_login,omitempty"`
		IPAddress       string         `gorm:"type:varchar(255)" json:"ip_address,omitempty"`
		Handset         string         `gorm:"type:varchar(255)" json:"handset,omitempty"`
		Status          bool           `gorm:"default:true" json:"status"`
		TypeUser        string         `gorm:"type:varchar(255)" json:"type_user,omitempty"`
		IsVerify        bool           `gorm:"type:bool" json:"is_verify,omitempty"`
		VerifyBy        string         `gorm:"type:varchar(255)" json:"verify_by,omitempty"`
		VerifyAt        *time.Time     `json:"verify_at,omitempty"`
		CreatedBy       string         `gorm:"type:varchar(255)" json:"created_by,omitempty"`
		UpdatedBy       string         `gorm:"type:varchar(255)" json:"updated_by,omitempty"`
		DeletedAt       gorm.DeletedAt `gorm:"index"`
		CreatedAt       time.Time
		UpdatedAt       time.Time

		// Relations
		Role       Role         `gorm:"foreignKey:RoleID;references:ID"`
		DetailUser []DetailUser `gorm:"foreignKey:UserID;references:ID"`
		UserAdnet  []UserAdnet  `gorm:"foreignKey:UserID;references:ID"`
	}

	DetailUser struct {
		ID         int  `gorm:"primaryKey" json:"id"`
		UserID     int  `json:"user_id"`
		CountryID  int  `json:"country_id"`
		OperatorID int  `json:"operator_id"`
		ServiceID  int  `json:"service_id"`
		Status     bool `json:"status"`
		CreatedAt  time.Time
		UpdatedAt  time.Time

		// Relations
		User User `gorm:"foreignKey:UserID;references:ID"`
	}

	UserAdnet struct {
		ID        int  `gorm:"primaryKey" json:"id"`
		UserID    int  `json:"user_id"`
		AdnetID   int  `json:"adnet_id"`
		Status    bool `json:"status"`
		CreatedAt time.Time
		UpdatedAt time.Time

		// Relations
		User User `gorm:"foreignKey:UserID;references:ID"`
	}

	Country struct {
		ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
		Code        string `gorm:"type:varchar(10)" json:"code" validate:"unique,required,max=10"`
		Code2       string `gorm:"type:varchar(10)" json:"code2" validate:"unique,required,max=10"`
		Name        string `gorm:"type:varchar(80)" json:"name"  validate:"unique,drequired,max=80"`
		NumericCode string `gorm:"type:varchar(10)" json:"numeric_code" form:"numeric_code"  validate:"required,max=10"`
		IsActive    string `gorm:"type:bool;default:true" json:"is_active" form:"is_active" `
	}

	Company struct {
		ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
		Name string `gorm:"type:text" json:"name"`
	}

	Domain struct {
		ID  uint   `gorm:"primaryKey;autoIncrement" json:"id"`
		Url string `gorm:"type:varchar(80)" json:"url"`
	}

	Operator struct {
		ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
		Country    string    `gorm:"type:varchar(10)" json:"country"`
		Name       string    `gorm:"type:varchar(30)" json:"name"`
		IsActive   string    `gorm:"type:bool;default:true" json:"is_active" form:"is_active" `
		Lastupdate time.Time `gorm:"type:timestamp" json:"lastupdate"`
	}
	Partner struct {
		ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
		Country        string    `gorm:"type:varchar(10)" json:"country"`
		Name           string    `gorm:"type:varchar(30)" json:"name"`
		Operator       string    `gorm:"type:varchar(50)" json:"operator"`
		Client         string    `gorm:"type:varchar(50)" json:"client"`
		ClientType     string    `gorm:"type:varchar(50)" json:"client_type"`
		Company        string    `gorm:"type:varchar(75)" json:"company"`
		UrlPostback    string    `gorm:"type:text" json:"url_postback"`
		Postback       string    `gorm:"type:varchar(50)" json:"postback"`
		PostbackMethod string    `gorm:"type:varchar(50)" json:"postback_method"`
		Aggregator     string    `gorm:"type:varchar(50)" json:"aggregator"`
		IsActive       string    `gorm:"type:bool;default:true" json:"is_active" form:"is_active"`
		IsBillable     string    `gorm:"type:bool;default:false" json:"is_billable" form:"is_billable"`
		Lastupdate     time.Time `gorm:"type:timestamp" json:"lastupdate"`
	}

	Service struct {
		ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
		Service  string `gorm:"type:varchar(55)" json:"service"`
		Adn      string `gorm:"type:varchar(20)" json:"adn"`
		Country  string `gorm:"type:varchar(50)" json:"country"`
		Operator string `gorm:"type:varchar(50)" json:"operator"`
	}

	AdnetList struct {
		ID           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
		Code         string `gorm:"type:varchar(30)" json:"code"`
		Name         string `gorm:"type:varchar(30)" json:"name"`
		ApiUrl       string `gorm:"type:text" json:"api_url"`
		ApiUrlBefore string `gorm:"column:api_url" json:"api_url_before"`
		ApiUrlDr     string `gorm:"type:text" json:"api_url_dr"`
		IsActive     string `gorm:"type:bool;default:true" json:"is_active" form:"is_active"`
		IsRemnant    string `gorm:"type:bool;default:false" json:"is_remnant" form:"is_remnant"`
		Lastupdate   string `gorm:"type:timestamp" json:"lastupdate"`
		IsDummyPixel string `gorm:"type:bool;default:false" json:"is_dummy_pixel"`
	}

	Audit struct {
		ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
		UserType      string    `gorm:"type:varchar(255)" json:"user_type"`
		UserID        int       `gorm:"type:int8" json:"user_id"`
		Event         string    `gorm:"type:varchar(255)" json:"event"`
		AuditableType string    `gorm:"type:varchar(255)" json:"auditable_type"`
		AuditableID   string    `gorm:"type:varchar(255)" json:"auditable_id"`
		OldValues     string    `gorm:"type:text" json:"old_values"`
		NewValues     string    `gorm:"type:text" json:"new_values"`
		URL           string    `gorm:"type:text" json:"url"`
		IPAddress     string    `gorm:"type:inet" json:"ip_address"`
		UserAgent     string    `gorm:"type:varchar(1023)" json:"user_agent"`
		Tags          string    `gorm:"type:varchar(255)" json:"tags"`
		CreatedAt     time.Time `gorm:"type:timestamp" json:"created_at"`
		UpdatedAt     time.Time `gorm:"type:timestamp" json:"updated_at"`
		ActionName    string    `gorm:"type:varchar(255)" json:"action_name"`
	}
)
