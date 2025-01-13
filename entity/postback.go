package entity

import (
	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource/helper"
	"github.com/sirupsen/logrus"
)

type (
	Postback struct {
		CookieKey     string
		URLServiceKey string `json:"urlservicekey"`
		AffSub        string `json:"aff_sub"`
	}

	PostbackData struct {
		CmpDetail DataConfig
		Pxs       PixelStorage
	}
)

func NewDataPostback(c *fiber.Ctx) *Postback {

	m := c.Queries()

	CookieKey := helper.Concat("-", helper.GetIpAddress(c), m["urlservicekey"], m["aff_sub"])

	return &Postback{
		CookieKey:     CookieKey,
		URLServiceKey: m["urlservicekey"],
		AffSub:        m["aff_sub"],
	}
}

func (p *Postback) ValidateParams(Logs *logrus.Logger) GlobalResponse {

	if p.URLServiceKey == "" {

		return GlobalResponse{Code: fiber.StatusBadRequest, Message: "urlservicekey empty or not found"}

	} else if p.AffSub == "" {

		return GlobalResponse{Code: fiber.StatusBadRequest, Message: "pixel empty"}

	} else {
		Logs.Debug("All traffic service is valid ...\n")

		return GlobalResponse{Code: fiber.StatusOK, Message: ""}
	}
}
