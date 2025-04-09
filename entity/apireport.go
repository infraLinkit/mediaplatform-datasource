package entity

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource/config"
	"github.com/infraLinkit/mediaplatform-datasource/helper"
	"github.com/sirupsen/logrus"
)

type (
	DisplayPinReport struct {
		Draw       int    `form:"draw" json:"draw"`
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

	DisplayPinPerformanceReport struct {
		Country    string `form:"country" json:"country"`
		Adnet      string `form:"adnet" json:"adnet"`
		Operator   string `form:"operator" json:"operator"`
		Service    string `form:"service" json:"service"`
		Page       int    `form:"page" json:"page"`
		DateRange  string `form:"date_range" json:"date_range"`
		DateBefore string `form:"date_before" json:"date_before"`
		DateAfter  string `form:"date_after" json:"date_after"`
		Action     string `form:"action" json:"action"`
		Draw       int    `form:"draw" json:"draw"`
		PageSize   int    `form:"page_size" json:"page_size"`
	}

	DisplayConversionLogReport struct {
		Country        string `form:"country" json:"country"`
		Adnet          string `form:"adnet" json:"adnet"`
		Agency         string `form:"agency" json:"agency"`
		Operator       string `form:"operator" json:"operator"`
		CampaignType   string `form:"campaign_type" json:"campaign_type"`
		CampaignId     string `form:"campaign_id" json:"campaign_id"`
		StatusPostback string `form:"status_postback" json:"status_postback"`
		Pixel          string `form:"pixel" json:"pixel"`
		Page           int    `form:"page" json:"page"`
		DateRange      string `form:"date_range" json:"date_range"`
		DateStart      string `form:"date_start" json:"date_start"`
		DateEnd        string `form:"date_end" json:"date_end"`
		Action         string `form:"action" json:"action"`
		Draw           int    `form:"draw" json:"draw"`
		PageSize       int    `form:"page_size" json:"page_size"`
	}

	DisplayCPAReport struct { // cpa
		ID                int       `form:"id" json:"id"`
		SummaryDate       time.Time `form:"summary_date" json:"summary_date"`
		CampaignId        string    `form:"campaign_id" json:"campaign_id"`
		CampaignName      string    `form:"campaign_name" json:"campaign_name"`
		Country           string    `form:"country" json:"country"`
		Operator          string    `form:"operator" json:"operator"`
		Partner           string    `form:"partner" json:"partner"`
		Aggregator        string    `form:"aggregator" json:"aggregator"`
		Adnet             string    `form:"adnet" json:"adnet"`
		Service           string    `form:"service" json:"service"`
		Cost              float64   `form:"cost" json:"cost"`
		SBAF              float64   `form:"sbaf" json:"sbaf"`
		SAAF              float64   `form:"saaf" json:"saaf"`
		RatioSend         int       `form:"ratio_send" json:"ratio_send"`
		RatioReceive      int       `form:"ratio_receive" json:"ratio_receive"`
		Company           string    `form:"company" json:"company"`
		ClientType        string    `form:"client_type" json:"client_type"`
		CostPerConversion float64   `form:"cost_per_conversion" json:"cost_per_conversion"`
		AgencyFee         float64   `form:"agency_fee" json:"agency_fee"`
		PageSize          int       `form:"page_size" json:"page_size"`
		Page              int       `form:"page" json:"page"`
		Action            string    `form:"action" json:"action"`
		DateRange         string    `form:"date_range" json:"date_range"`
		DateBefore        string    `form:"date_before" json:"date_before"`
		DateAfter         string    `form:"date_after" json:"date_after"`
		Draw              int       `form:"draw" json:"draw"`
	}

	CostReport struct {
		Adnet       string  `json:"adnet"`
		Country     string  `json:"country"`
		Operator    string  `json:"operator"`
		Conversion1 float64 `json:"conversion1"`
		Cost1       float64 `json:"cost1"`
		Conversion2 float64 `json:"conversion2"`
		Cost2       float64 `json:"cost2"`
	}

	DisplayCostReport struct {
		Adnet       string   `json:"adnet"`
		Country     string   `json:"country"`
		Operator    string   `json:"operator"`
		Conversion1 *float64 `json:"conversion1"`
		Conversion2 *float64 `json:"conversion2"`
		Cost1       *float64 `json:"revenue1"`
		Cost2       *float64 `json:"revenue2"`
		Action      string   `json:"action"`
		DateRange   string   `json:"date_range"`
		DateBefore  string   `json:"date_before"`
		DateAfter   string   `json:"date_after"`
		Page        int      `json:"page"`
		Draw        int      `json:"draw"`
		DataBasedOn string   `json:"data_based_on"`
	}

	DisplayAlertReport struct {
		Action       string `json:"action"`
		Country      string `json:"country"`
		Operator     string `json:"operator"`
		CampaignName string `json:"campaign_name"`
		Service      string `json:"service"`
		DateRange    string `json:"date_range"`
		DateBefore   string `json:"date_before"`
		DateAfter    string `json:"date_after"`
		Page         int    `json:"page"`
		Draw         int    `json:"draw"`
		PageSize     int    `json:"page_size"`
		ExportData   string `json:"export_data"`
	}

	PerformaceReportParams struct {
		Country      string `form:"country" json:"country"`
		Company      string `form:"company" json:"company"`
		ClientType   string `form:"client_type" json:"client_type"`
		Operator     string `form:"operator" json:"operator"`
		CampaignName string `form:"campaign_name" json:"campaign_name"`
		CampaignType string `form:"campaign_type" json:"campaign_type"`
		Publisher    string `form:"publisher" json:"publisher"`
		Service      string `form:"service" json:"service"`
		CampaignId   string `form:"campaign_id" json:"campaign_id"`
		Partner      string `form:"partner" json:"partner"`
		PageSize     int    `form:"page_size" json:"page_size"`
		Page         int    `form:"page" json:"page"`
		Action       string `form:"action" json:"action"`
		DateStart    string `form:"date_before" json:"date_start"`
		DateEnd      string `form:"date_after" json:"date_end"`
		Draw         int    `form:"draw" json:"draw"`
	}

	PerformanceReport struct {
		Country            string  `json:"country"`
		Company            string  `json:"company"`
		ClientType         string  `json:"client_type"`
		CampaignName       string  `json:"campaign_name"`
		Operator           string  `json:"operator"`
		Service            string  `json:"service"`
		Adnet              string  `json:"adnet"`
		PixelReceived      int     `json:"pixel_received"`
		PixelSend          int     `json:"pixel_send"`
		CRPostback         int     `json:"cr_postback"`
		CRMo               int     `json:"cr_mo"`
		Landing            int     `json:"landing"`
		RatioSend          int     `json:"ratio_send"`
		RatioReceive       int     `json:"ratio_receive"`
		PricePerPostback   float64 `json:"price_per_postback"`
		CostPerConversion  float64 `json:"cost_per_conversion"`
		AgencyFee          float64 `json:"agency_fee"`
		SpendingToAdnets   float64 `json:"spending_to_adnets"`
		TotalWakiAgencyFee float64 `json:"total_waki_agency_fee"`
		TotalSpending      float64 `json:"total_spending"`
		TotalFP            float64 `json:"total_fp"`
		SuccessFP          float64 `json:"success_fp"`
		ECPA               float64 `json:"e_cpa"`
		ARPUROI            float64 `json:"arpu_roi"`
		ARPU90             float64 `json:"arpu_90"`
		BillrateFP         float64 `json:"billrate_fp"`
	}

	ARPUResponse struct {
		Status  int          `json:"status"`
		Message string       `json:"message"`
		Data    *ARPUDataSet `json:"data"`
	}

	ARPUDataSet struct {
		DateHit   string         `json:"date_hit"`
		Country   string         `json:"country"`
		Operator  string         `json:"operator"`
		Service   string         `json:"service"`
		Keyword   string         `json:"keyword"`
		Publisher string         `json:"publisher"`
		Data      []ARPUDataItem `json:"data"`
	}

	ARPUDataItem struct {
		Adnet        string  `json:"adnet"`
		Arpu90       float64 `json:"arpu90"`
		Arpu90Net    float64 `json:"arpu90_net"`
		Arpu90USD    float64 `json:"arpu90_usd"`
		Arpu90USDNet float64 `json:"arpu90_usd_net"`
		Service      string  `json:"service"`
	}
)

func NewInstanceTrxPinReport(c *fiber.Ctx, cfg *config.Cfg) *ApiPinReport {

	m := c.Queries()

	mo, _ := strconv.Atoi(m["mo"])
	postback, _ := strconv.Atoi(m["postback"])

	pin := ApiPinReport{
		Adnet:         m["adnet"],
		Country:       m["country"],
		Service:       m["service"],
		Operator:      m["telco"],
		DateSend:      helper.GetCurrentTime(cfg.TZ, time.RFC3339),
		TotalMO:       mo,
		TotalPostback: postback,
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
