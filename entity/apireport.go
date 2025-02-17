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
		Draw       int    `form:"draw" json:"draw"`
		PageSize   int    `form:"page_size" json:"page_size"`
	}

	DisplayConversionLogReport struct {
		Country    string `form:"country" json:"country"`
		Adnet      string `form:"adnet" json:"adnet"`
		Operator   string `form:"operator" json:"operator"`
		CampaignId string `form:"campaign_id" json:"campaign_id"`
		Pixel      string `form:"pixel" json:"pixel"`
		Page       int    `form:"page" json:"page"`
		DateRange  string `form:"date_range" json:"date_range"`
		DateBefore string `form:"date_before" json:"date_before"`
		DateAfter  string `form:"date_after" json:"date_after"`
		Action     string `form:"action" json:"action"`
		Draw       int    `form:"draw" json:"draw"`
		PageSize   int    `form:"page_size" json:"page_size"`
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
