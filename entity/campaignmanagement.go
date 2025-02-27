package entity

type (
	DisplayCampaignManagement struct {
		Country      string `form:"country" json:"country"`
		Adnet        string `form:"adnet" json:"adnet"`
		Operator     string `form:"operator" json:"operator"`
		Service      string `form:"service" json:"service"`
		Status       string `form:"status" json:"status"`
		Partner      string `form:"partner" json:"partner"`
		CampaignName string `form:"campaign_name" json:"campaign_name"`
		CampaignId   string `form:"campaign_id" json:"campaign_id"`
		Page         int    `form:"page" json:"page"`
		Draw         int    `form:"draw" json:"draw"`
		Action       string `form:"action" json:"action"`
	}

	CampaignManagementData struct {
		ID            []int  `json:"id"`
		CampaignID    string `json:"campaign_id"`
		CampaignName  string `json:"campaign_name"`
		Country       string `json:"country"`
		Partner       string `json:"partner"`
		TotalOperator int    `json:"total_operator"`
		Service       string `json:"service"`
		TotalAdnet    int    `json:"total_adnet"`
		ShortCode     string `json:"short_code"`
		IsActive      bool   `json:"is_active"`
	}

	CampaignManagementDetail struct {
		ID             int    `json:"id"`
		CampaignID     string `json:"campaign_id"`
		CampaignName   string `json:"campaign_name"`
		Country        string `json:"country"`
		Operator       string `json:"operator"`
		Service        string `json:"service"`
		Adnet          string `json:"adnet"`
		Partner        string `json:"partner"`
		ShortCode      string `json:"short_code"`
		MOLimit        int    `json:"mo_limit"`
		Payout         string `json:"po"`
		RatioSend      int    `json:"ratio_send"`
		RatioReceive   int    `json:"ratio_receive"`
		URLPostback    string `json:"url_postback"`
		URLService     string `json:"url_service"`
		URLanding      string `json:"url_landing"`
		URLWarpLanding string `json:"url_warp_landing"`
		APIURL         string `json:"api_url"`
		IsActive       bool   `json:"is_active"`
		UrlServiceKey  string `json:"url_service_key"`
	}

	CampaignManagementDataDetail struct {
		Operator string                     `json:"operator"`
		Service  string                     `json:"service"`
		Details  []CampaignManagementDetail `json:"details"`
	}

	CampaignCounts struct {
		TotalCampaigns          int `json:"total_campaign"`
		TotalActiveCampaigns    int `json:"total_active_campaign"`
		TotalNonActiveCampaigns int `json:"total_inactive_campaign"`
	}
)
