package entity

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource/helper"
	"github.com/sirupsen/logrus"
)

type (
	Postback struct {
		CookieKey  string
		Country    string `json:"country"`
		Operator   string `json:"operator"`
		Partner    string `json:"partner"`
		ServiceId  string `json:"serv_id"`
		Keyword    string `json:"keyword"`
		TrxId      string `json:"trxid"`
		Msisdn     string `json:"msisdn"`
		Px         string `json:"px"`
		IsBillable bool   `json:"is_billable"`
	}

	PostbackData struct {
		CmpDetail DataConfig
		Pxs       PixelStorage
	}
)

func NewDataPostback(c *fiber.Ctx) *Postback {

	m := c.Queries()

	isBillable, _ := strconv.ParseBool(m["is_billable"])

	CookieKey := helper.Concat("-", helper.GetIpAddress(c), m["country"], m["operator"], m["partner"], m["serv_id"], m["keyword"], m["msisdn"], m["px"], m["is_billable"])

	return &Postback{
		CookieKey:  CookieKey,
		Country:    m["country"],
		Operator:   m["operator"],
		Partner:    m["partner"],
		ServiceId:  m["serv_id"],
		Keyword:    m["keyword"],
		Msisdn:     m["msisdn"],
		Px:         m["px"],
		TrxId:      m["trxid"],
		IsBillable: isBillable,
	}
}

func (p *Postback) ValidateParams(Logs *logrus.Logger) GlobalResponse {

	if p.Country == "" {

		return GlobalResponse{Code: fiber.StatusBadRequest, Message: "country empty"}

	} else if p.Operator == "" {

		return GlobalResponse{Code: fiber.StatusBadRequest, Message: "operator empty"}

	} else if p.Partner == "" {

		return GlobalResponse{Code: fiber.StatusBadRequest, Message: "partner empty"}

	} else if p.ServiceId == "" {

		return GlobalResponse{Code: fiber.StatusBadRequest, Message: "serv_id empty"}

	} else if p.Keyword == "" {

		return GlobalResponse{Code: fiber.StatusBadRequest, Message: "keyword empty"}

	} else if p.Msisdn == "" {

		return GlobalResponse{Code: fiber.StatusBadRequest, Message: "msisdn empty"}

	} else if p.Px == "" {

		return GlobalResponse{Code: fiber.StatusBadRequest, Message: "pixel empty"}

	} else {
		Logs.Debug("All traffic service is valid ...\n")

		return GlobalResponse{Code: fiber.StatusOK, Message: ""}
	}
}
