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

	SummaryDashboard struct {
		TotalMO                 int     `json:"total_mo"`
		TotalActiveAdnet        int     `json:"total_active_adnet"`
		TotalSpending           float64 `json:"total_spending"`
		TotalS2SSpending        float64 `json:"total_s2s_spending"`
		TotalAPISpending        float64 `json:"total_api_spending"`
		TotalMainstreamSpending float64 `json:"total_mainstream_spending"`
		TotalDSPSpending        float64 `json:"total_dsp_spending"`
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

/*
func NewInstanceTrxPinReport(c *fiber.Ctx, cfg *config.Cfg) *ApiPinReport {

	m := c.Queries()

	mo, _ := strconv.Atoi(m["mo"])
	postback, _ := strconv.Atoi(m["postback"])
	sbaf, _ := strconv.ParseFloat(m["sbaf"], 64)
	saaf, _ := strconv.ParseFloat(m["saaf"], 64)
	price_per_mo, _ := strconv.ParseFloat(m["price_per_mo"], 64)
	waki_revenue, _ := strconv.ParseFloat(m["waki_revenue"], 64)

	pin := ApiPinReport{
		CampaignId:    m["campaign_id"],
		Country:       m["country"],
		Company:       m["company"],
		Adnet:         m["adnet"],
		Service:       m["service"],
		Operator:      m["telco"],
		DateSend:      helper.GetCurrentTime(cfg.TZ, time.RFC3339),
		PayoutAdn:     m["payout_adn"],
		PayoutAF:      m["payout_af"],
		TotalMO:       mo,
		TotalPostback: postback,
		SBAF:          sbaf,
		SAAF:          saaf,
		PricePerMO:    price_per_mo,
		WakiRevenue:   waki_revenue,
	}

	return &pin
}

func (t *ApiPinReport) ValidateParams(Logs *logrus.Logger) ReturnResponse {

	if t.Adnet == "" {

		return ReturnResponse{HttpStatus: fiber.StatusBadRequest, Rsp: GlobalResponse{Code: fiber.StatusBadRequest, Message: "Parameter Adnet is mandatory"}}

	} else if t.Country == "" {

		return ReturnResponse{HttpStatus: fiber.StatusBadRequest, Rsp: GlobalResponse{Code: fiber.StatusBadRequest, Message: "Parameter Country is mandatory"}}

	} else if t.Service == "" {

		return ReturnResponse{HttpStatus: fiber.StatusBadRequest, Rsp: GlobalResponse{Code: fiber.StatusBadRequest, Message: "Parameter Service is mandatory"}}

	} else if t.Operator == "" {

		return ReturnResponse{HttpStatus: fiber.StatusBadRequest, Rsp: GlobalResponse{Code: fiber.StatusBadRequest, Message: "Parameter Operator is mandatory"}}

	} else {

		return ReturnResponse{HttpStatus: fiber.StatusOK, Rsp: GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC}}

	}
}

func NewInstanceTrxPinPerfonrmanceReport(c *fiber.Ctx, cfg *config.Cfg) *ApiPinPerformance {
	m := c.Queries()

	pinRequest, _ := strconv.Atoi(m["pin_request"])
	uniquePinRequest, _ := strconv.Atoi(m["unique_pin_request"])
	pinSent, _ := strconv.Atoi(m["pin_sent"])
	pinFailed, _ := strconv.Atoi(m["pin_failed"])
	verifyRequest, _ := strconv.Atoi(m["verify_request"])
	verifyRequestUnique, _ := strconv.Atoi(m["verify_request_unique"])
	pinOK, _ := strconv.Atoi(m["pin_ok"])
	pinNotOK, _ := strconv.Atoi(m["pin_not_ok"])
	pinOkSendAdnet, _ := strconv.Atoi(m["pin_ok_send_adnet"])

	pin := ApiPinPerformance{
		Adnet:               m["adnet"],
		Country:             m["country"],
		Service:             m["service"],
		Operator:            m["telco"],
		DateSend:            helper.GetCurrentTime(cfg.TZ, time.RFC3339),
		PinRequest:          pinRequest,
		UniquePinRequest:    uniquePinRequest,
		PinSent:             pinSent,
		PinFailed:           pinFailed,
		VerifyRequest:       verifyRequest,
		VerifyRequestUnique: verifyRequestUnique,
		PinOK:               pinOK,
		PinNotOK:            pinNotOK,
		PinOkSendAdnet:      pinOkSendAdnet,
	}
	return &pin
}

func (t *ApiPinPerformance) ValidateParams(Logs *logrus.Logger) ReturnResponse {

	if t.Adnet == "" {

		return ReturnResponse{HttpStatus: fiber.StatusBadRequest, Rsp: GlobalResponse{Code: fiber.StatusBadRequest, Message: "Parameter Adnet is mandatory"}}

	} else if t.Country == "" {

		return ReturnResponse{HttpStatus: fiber.StatusBadRequest, Rsp: GlobalResponse{Code: fiber.StatusBadRequest, Message: "Parameter Country is mandatory"}}

	} else if t.Service == "" {

		return ReturnResponse{HttpStatus: fiber.StatusBadRequest, Rsp: GlobalResponse{Code: fiber.StatusBadRequest, Message: "Parameter Service is mandatory"}}

	} else if t.Operator == "" {

		return ReturnResponse{HttpStatus: fiber.StatusBadRequest, Rsp: GlobalResponse{Code: fiber.StatusBadRequest, Message: "Parameter Operator is mandatory"}}

	} else {

		return ReturnResponse{HttpStatus: fiber.StatusOK, Rsp: GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC}}

	}
}

func NewInstancePinPerformance(c *fiber.Ctx, cfg *config.Cfg) *ApiPinPerformance {
	m := c.Queries()

	toInt := func(key string) int {
		val, _ := strconv.Atoi(m[key])
		return val
	}

	toFloat := func(key string) float64 {
		val, _ := strconv.ParseFloat(m[key], 64)
		return val
	}

	pin := ApiPinPerformance{
		DateSend:            helper.GetCurrentTime(cfg.TZ, time.RFC3339),
		Country:             m["country"],
		Company:             m["company"],
		Adnet:               m["adnet"],
		Operator:            m["operator"],
		Service:             m["service"],
		PinRequest:          toInt("pin_request"),
		UniquePinRequest:    toInt("unique_pin_request"),
		PinSent:             toInt("pin_sent"),
		PinFailed:           toInt("pin_failed"),
		VerifyRequest:       toInt("verify_request"),
		VerifyRequestUnique: toInt("verify_request_unique"),
		PinOK:               toInt("pin_ok"),
		PinNotOK:            toInt("pin_not_ok"),
		PinOkSendAdnet:      toInt("pin_ok_send_adnet"),
		CPA:                 toFloat("cpa"),
		CPAWaki:             toFloat("cpa_waki"),
		EstimatedARPU:       toFloat("estimated_arpu"),
		SBAF:                toFloat("sbaf"),
		SAAF:                toFloat("saaf"),
		ChargedMO:           toFloat("charged_mo"),
		SubsCR:              toFloat("subs_cr"),
		AdnetCR:             toFloat("adnet_cr"),
		CAC:                 toFloat("cac"),
		PaidCAC:             toFloat("paid_cac"),
		CrMO:                toFloat("cr_mo"),
		CrPostback:          toFloat("cr_postback"),
		Landing:             toInt("landing"),
		ROI:                 toFloat("roi"),
		Arpu90:              toFloat("arpu90"),
		BillingRateFP:       toFloat("billing_rate_fp"),
		Ratio:               toFloat("ratio"),
		PricePerPostback:    toFloat("price_per_postback"),
		CostPerConversion:   toFloat("cost_per_conversion"),
		AgencyFee:           toFloat("agency_fee"),
		TotalWakiAgencyFee:  toFloat("total_waki_agency_fee"),
		TotalSpending:       toFloat("total_spending"),
		// ClientType:          m["client_type"],
		// CampaignName: m["campaign_name"],
	}

	return &pin
}
*/
