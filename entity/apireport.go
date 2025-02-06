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
		Country    string `form:"country" json:"country"`
		Adnet      string `form:"adnet" json:"adnet"`
		Operator   string `form:"operator" json:"operator"`
		Service    string `form:"service" json:"service"`
		Page       int    `form:"page" json:"page"`
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
	}

	DisplayCPAReport struct { // cpa
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
		Page                     int       `form:"page" json:"page"`
		Action                   string    `form:"action" json:"action"`
		DateRange                string    `form:"date_range" json:"date_range"`
		DateBefore               string    `form:"date_before" json:"date_before"`
		DateAfter                string    `form:"date_after" json:"date_after"`
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
