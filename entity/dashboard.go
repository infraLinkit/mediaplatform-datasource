package entity

type (
	DisplayDashboard struct {
		Period  string `form:"period" json:"period"`
		Metrics string `form:"metrics" json:"metrics"`

		Country    string `form:"country" json:"country"`
		Adnet      string `form:"adnet" json:"adnet"`
		Operator   string `form:"operator" json:"operator"`
		Service    string `form:"service" json:"service"`
		Page       int    `form:"page" json:"page"`
		PageSize   int    `form:"page_size" json:"page_size"`
		DateRange  string `form:"date_range" json:"date_range"`
		DateBefore string `form:"date_before" json:"date_before"`
		DateAfter  string `form:"date_after" json:"date_after"`
		Action     string `form:"action" json:"action"`
	}

	SummaryDashboardDetail struct {
		Date                    string  `json:"date"`
		TotalMO                 int     `json:"total_mo"`
		TotalActiveAdnet        int     `json:"total_active_adnet"`
		TotalSpending           float64 `json:"total_spending"`
		TotalS2sSpending        float64 `json:"total_s2s_spending"`
		TotalApiSpending        float64 `json:"total_api_spending"`
		TotalMainstreamSpending float64 `json:"total_mainstream_spending"`
		TotalDSPSpending        float64 `json:"total_dsp_spending"`
	}

	DetailChartData struct {
		Date                   string  `json:"date"`
		TotalMO                int     `json:"total_mo"`
		TotalSpending          float64 `json:"total_spending"`
		LastMonthTotalMO       int     `json:"lastmonth_total_mo"`
		LastMonthTotalSpending float64 `json:"lastmonth_total_spending"`
		LastMonthDate          string  `json:"last_month_date"`
	}

	SummaryDashboardData struct {
		TotalMO                 int     `json:"total_mo"`
		TotalActiveAdnet        int     `json:"total_active_adnet"`
		TotalSpending           float64 `json:"total_spending"`
		TotalS2SSpending        float64 `json:"total_s2s_spending"`
		TotalAPISpending        float64 `json:"total_api_spending"`
		TotalMainstreamSpending float64 `json:"total_mainstream_spending"`
		TotalDSPSpending        float64 `json:"total_dsp_spending"`
		//	Detail                  []SummaryDashboardDetail `json:"detail"`
		DateList        []string          `json:"date_list"`
		DetailChartData []DetailChartData `json:"detail_chart_data"`
	}

	SummaryDashboardReportDetail struct {
		Date             string  `json:"date"`
		MOReceived       int     `json:"mo_received"`
		MOSent           int     `json:"mo_sent"`
		SpendingToAdnets float64 `json:"spending_to_adnets"`
		Spending         float64 `json:"spending"`
		WAKIRevenue      float64 `json:"waki_revenue"`
	}

	SummaryDashboardReport struct {
		DateRange string                         `json:"date_range"`
		DateList  []string                       `json:"date_list"`
		Detail    []SummaryDashboardReportDetail `json:"detail"`
	}

	SummaryMODetail struct {
		DateNow string `json:"date_now"`
	}

	TopCampaign struct {
		CampaignID  string  `json:"campaign_id"`
		Country     string  `json:"country"`
		CountryName string  `json:"country_name"`
		Landing     int     `json:"landing"`
		MO          int     `json:"mo_received"`
		Postback    int     `json:"postback"`
		CRMO        float64 `json:"cr_mo"`
		CRPostback  float64 `json:"cr_postback"`
		URL         string  `json:"url"`
		ECPA        string  `json:"e_cpa"`
	}

	SummaryTopBestCampaign struct {
		Campaign []TopCampaign `json:"campaign"`
	}

	SummaryTopWorstCampaign struct {
		Campaign []TopCampaign `json:"campaign"`
	}
)
