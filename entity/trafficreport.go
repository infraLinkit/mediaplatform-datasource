package entity

type (
	DisplayTrafficReport struct {
		Draw         int    `form:"draw" json:"draw"`
		Date         string `form:"date" json:"date"`
		CampaignID   string `form:"campaign_id" json:"campaign_id"`
		CampaignName string `form:"campaign_name" json:"campaign_name"`
		Country      string `form:"country" json:"country"`
		Operator     string `form:"operator" json:"operator"`
		Service      string `form:"service" json:"service"`
		Adnet        string `form:"adnet" json:"adnet"`

		DataIndicator []string `form:"indicator" json:"indicator"`

		Total []DataIndicator `form:"total" json:"total"`

		Avg []DataIndicator `form:"avarage" json:"avarage"`

		TmoEnd []DataIndicator `form:"tmo_end" json:"tmo_end"`

		TitleDate []Dates `form:"title_dates" json:"title_dates"`

		Page       int    `form:"page" json:"page"`
		PageSize   int    `form:"page_size" json:"page_size"`
		DateRange  string `form:"date_range" json:"date_range"`
		DateBefore string `form:"date_before" json:"date_before"`
		DateAfter  string `form:"date_after" json:"date_after"`
		Action     string `form:"action" json:"action"`
	}

	DataIndicator struct {
		Traffic   int `form:"traffic" json:"traffic"`
		Cr        int `form:"cr" json:"cr"`
		Mo        int `form:"mo_received" json:"mo_received"`
		FirstPush int `form:"first_push" json:"first_push"`
	}

	Dates struct {
		Date  string          `form:"date" json:"date"`
		Value []DataIndicator `form:"value" json:"value"`
	}
)
