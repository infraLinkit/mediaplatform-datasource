package entity

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource/helper"
	"github.com/sirupsen/logrus"
)

type (
	PostbackReceive struct {
		CookieKey     string
		URLServiceKey string `json:"urlservicekey"`
		AffSub        string `json:"aff_sub"`
		Method        string `json:"method"`
		Msisdn        string `json:"msisdn"`
		Trxid         string `json:"trxid"`
	}

	PostbackData struct {
		CmpDetail DataConfig
		Pxs       PixelStorage
	}
)

func NewDataPostback(c *fiber.Ctx) *PostbackReceive {

	m := c.Queries()

	CookieKey := helper.Concat("-", helper.GetIpAddress(c), m["urlservicekey"], m["aff_sub"])

	return &PostbackReceive{
		CookieKey:     CookieKey,
		URLServiceKey: m["urlservicekey"],
		AffSub:        m["aff_sub"],
		Msisdn:        m["msisdn"],
		Trxid:         m["trxid"],
	}
}

func (p *PostbackReceive) ValidateParams(Logs *logrus.Logger) GlobalResponse {

	if p.URLServiceKey == "" {

		return GlobalResponse{Code: fiber.StatusBadRequest, Message: "urlservicekey empty or not found"}

	} else if p.AffSub == "" {

		return GlobalResponse{Code: fiber.StatusBadRequest, Message: "pixel empty"}

	} else {
		Logs.Debug("All traffic service is valid ...\n")

		return GlobalResponse{Code: fiber.StatusOK, Message: ""}
	}
}

func NewDataPostbackV2(c *fiber.Ctx) *PostbackReceive {

	m := c.Queries()

	CookieKey := helper.Concat("-", helper.GetIpAddress(c), m["aff_sub"])

	return &PostbackReceive{
		CookieKey: CookieKey,
		AffSub:    m["aff_sub"],
		Method:    strings.ToUpper(m["method"]),
		Msisdn:    m["msisdn"],
		Trxid:     m["trxid"],
	}
}

func (p *PostbackReceive) ValidateParamsV2(Logs *logrus.Logger) GlobalResponse {

	if p.AffSub == "" {

		return GlobalResponse{Code: fiber.StatusBadRequest, Message: "pixel empty"}

	} else {
		Logs.Debug("All traffic service is valid ...\n")

		return GlobalResponse{Code: fiber.StatusOK, Message: ""}
	}
}
