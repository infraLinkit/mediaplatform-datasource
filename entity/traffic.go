package entity

import (
	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource/config"
	"github.com/infraLinkit/mediaplatform-datasource/helper"
	"github.com/sirupsen/logrus"
)

type (
	Traffic struct {
		Date          string
		Key           string
		KeyCfg        string
		KeyCounter    string
		KeyDataMining string
		DataTraffic   DataTraffic
	}

	DataTraffic struct {
		URLServiceKey string `form:"urlservicekey" json:"urlservicekey" xml:"urlservicekey"`
		Aff_Sub       string `form:"aff_sub" json:"aff_sub" xml:"aff_sub"`
		Partner       string `form:"p" json:"p" xml:"p"`
		Service       string `form:"srv" json:"srv" xml:"srv"`
		Adnet         string `form:"ad" json:"ad" xml:"ad"`
		PubId         string `form:"pubid" json:"pubid" xml:"pubid"`
	}
)

func NewInstanceTraffic(cfg *config.Cfg, o DataTraffic) *Traffic {

	date := helper.GetFormatTime(cfg.TZ, "20060102")
	//key := helper.Concat("-", o.URLServiceKey)

	return &Traffic{
		Date:       date,
		Key:        o.URLServiceKey,
		KeyCfg:     helper.Concat("-", o.URLServiceKey, "configIdx"),
		KeyCounter: helper.Concat("-", o.URLServiceKey, "counterIdx"),
		//KeyDataMining: helper.Concat("-", date, key, "dataminingIdx"),
		DataTraffic: o,
	}
}

func (t *Traffic) ValidateParams(Logs *logrus.Logger, traffic *Traffic) GlobalResponse {

	if traffic.DataTraffic.URLServiceKey == "" {
		Logs.Debug("Receive traffic keyaccess param is empty ...\n")

		return GlobalResponse{Code: fiber.StatusBadRequest, Message: "parameters is not complete"}
	} else if traffic.DataTraffic.Aff_Sub == "" {
		Logs.Debug("Receive traffic aff_sub param is empty ...\n")

		return GlobalResponse{Code: fiber.StatusBadRequest, Message: "parameters is not complete"}
	} else if traffic.DataTraffic.Adnet == "" {
		Logs.Debug("Receive traffic adnet param is empty ...\n")

		return GlobalResponse{Code: fiber.StatusBadRequest, Message: "parameters is not complete"}
	} else if traffic.DataTraffic.Partner == "" {
		Logs.Debug("Receive traffic partner param is empty ...\n")

		return GlobalResponse{Code: fiber.StatusBadRequest, Message: "parameters is not complete"}
	} else if traffic.DataTraffic.Service == "" {
		Logs.Debug("Receive traffic service param is empty ...\n")

		return GlobalResponse{Code: fiber.StatusBadRequest, Message: "parameters is not complete"}
	} else {
		Logs.Debug("All traffic service is valid ...\n")

		return GlobalResponse{Code: fiber.StatusOK, Message: ""}
	}
}
