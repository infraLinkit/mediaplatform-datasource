package handler

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource/config"
	"github.com/infraLinkit/mediaplatform-datasource/entity"
	"github.com/infraLinkit/mediaplatform-datasource/helper"
)

func (h *IncomingHandler) SetData(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	v := c.Params("v")

	if v == "set_target_daily_budget" {
		m := c.Queries()
		target_daily_budget := m["target_daily_budget"]
		country := strings.ToUpper(m["country"])
		operator := strings.ToUpper(m["operator"])

		if target_daily_budget == "" {
			return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{Code: fiber.StatusBadRequest, Message: "target_daily_budget is empty"})
		}

		targetDailyBudget, err := strconv.ParseFloat(target_daily_budget, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{Code: fiber.StatusBadRequest, Message: "target_daily_budget is not a valid number"})
		}

		redisKey := strings.ToLower(helper.Concat("_", "global_setting_tdb", country[0:2], operator[0:3]))

		gset, _ := json.Marshal(entity.GlobalSetting{
			TargetDailyBudget: target_daily_budget,
		})

		h.DS.SetData(redisKey, "$", string(gset))

		h.DS.UpdateCampaignMonitoringBudget(entity.CampaignDetail{
			TargetDailyBudget: targetDailyBudget,
			Country:           country,
			Operator:          operator,
		})

		h.DS.UpdateReportSummaryCampaignMonitoringBudget(entity.SummaryCampaign{
			SummaryDate:       helper.GetCurrentTime(h.Config.TZ, time.RFC3339),
			TargetDailyBudget: targetDailyBudget,
			Country:           country,
			Operator:          operator,
		})

		return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{Code: fiber.StatusBadRequest, Message: "Request parameter unknown"})
	}
}

func (h *IncomingHandler) UpdateAgencyFeeAndCostConversion(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	m := c.Queries()

	v := c.Params("v")

	if v == "update_cpcr_agency" {
		agencyFee := m["agency_fee"]
		costPerConversion := m["cost_per_conversion"]

		h.Logs.Debug(fmt.Sprintf("Received agency_fee: %s, cost_per_conversion: %s", agencyFee, costPerConversion))

		if agencyFee == "" && costPerConversion == "" {
			h.Logs.Error("Missing required fields")
			return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Missing required fields",
			})
		}

		gs, _ := h.DS.GetDataConfig("global_setting", "$")

		if gs == nil {
			h.Logs.Warn("Global setting not found in Redis, initializing new configuration")
			gs = &entity.DataConfig{
				CPCR:              "",
				AgencyFee:         "",
				TargetDailyBudget: "",
			}
		}

		if costPerConversion != gs.CPCR {
			gs.CPCR = costPerConversion
		}

		if agencyFee != gs.AgencyFee {
			gs.AgencyFee = agencyFee
		}

		gset, err := json.Marshal(entity.GlobalSetting{
			CostPerConversion: gs.CPCR,
			AgencyFee:         gs.AgencyFee,
			TargetDailyBudget: gs.TargetDailyBudget,
		})
		if err != nil {
			h.Logs.Error(fmt.Sprintf("Failed to marshal global settings: %v", err))
			return c.Status(fiber.StatusInternalServerError).JSON(entity.GlobalResponse{
				Code:    fiber.StatusInternalServerError,
				Message: "Failed to update global settings",
			})
		}

		h.DS.SetData("global_setting", "$", string(gset))

		h.Logs.Info("Successfully updated agency fee and cost per conversion in Redis")
		return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{
			Code:    fiber.StatusOK,
			Message: "Settings updated successfully",
		})

	} else {
		return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{Code: fiber.StatusBadRequest, Message: "Request parameter unknown"})
	}
}

func (h *IncomingHandler) TrxPinReport(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	pin := entity.NewInstanceTrxPinReport(c, h.Config)
	r := pin.ValidateParams(h.Logs)
	if r.HttpStatus == 200 {
		h.DS.PinReport(*pin)
	}

	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) TrxPerformancePinReport(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	pin := entity.NewInstanceTrxPinPerfonrmanceReport(c, h.Config)
	r := pin.ValidateParams(h.Logs)
	if r.HttpStatus == 200 {
		h.DS.PinPerformanceReport(*pin)
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) UpdateRatio(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	m := c.Queries()
	v := c.Params("v")

	if v == "update_ratio" {
		ratioSend := m["ratio_send"]
		ratioReceived := m["ratio_received"]
		ID := m["id"]

		h.Logs.Debug(fmt.Sprintf("Received ratio_send: %s, ratio_received: %s", ratioSend, ratioReceived))

		if ratioSend == "" || ratioReceived == "" {
			h.Logs.Error("Missing required fields")
			return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Missing Required Fields",
			})
		}

		ratioSendInt, err := strconv.Atoi(ratioSend)
		if err != nil {
			h.Logs.Error(fmt.Sprintf("Failed to parse ratio_send: %v", err))
			return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Invalid ratio_send value",
			})
		}

		ratioReceivedInt, err := strconv.Atoi(ratioReceived)
		if err != nil {
			h.Logs.Error(fmt.Sprintf("Failed to parse ratio_received: %v", err))
			return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Invalid ratio_received value",
			})
		}

		id, err := strconv.Atoi(ID)
		if err != nil {
			h.Logs.Error(fmt.Sprintf("Failed to parse id: %v", err))
			return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Invalid id value",
			})
		}

		_, errFind := h.DS.FindSummaryCampaign(id)

		if errFind != nil {
			h.Logs.Error(fmt.Sprintf("Data Not Found: %v", errFind))
			return c.Status(fiber.StatusNotFound).JSON(entity.GlobalResponse{
				Code:    fiber.StatusNotFound,
				Message: "Data Not Found",
			})
		}

		err = h.DS.UpdateRatioModel(entity.SummaryCampaign{
			RatioSend:    ratioSendInt,
			RatioReceive: ratioReceivedInt,
		}, id)
		h.Logs.Info("Successfully updated ratio in the database.")
		if err != nil {
			h.Logs.Error(fmt.Sprintf("Failed to update ratio: %v", err))
			return c.Status(fiber.StatusInternalServerError).JSON(entity.GlobalResponse{
				Code:    fiber.StatusInternalServerError,
				Message: "Failed to update ratio",
			})
		}

		return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{
			Code:    fiber.StatusOK,
			Message: "Ratio updated successfully",
		})
	}

	return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{
		Code:    fiber.StatusBadRequest,
		Message: "Invalid request parameter",
	})
}

func (h *IncomingHandler) UpdatePostback(c *fiber.Ctx) error {
	c.Set("content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	m := c.Queries()
	v := c.Params("v")

	if v == "update_postback" {
		Postback := m["postback"]
		ID := m["id"]

		h.Logs.Debug(fmt.Sprintf("Received postback: %s", Postback))

		if Postback == "" {
			h.Logs.Error("Postback is empty")
			return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Postback is empty",
			})
		}

		postbackInt, err := strconv.Atoi(Postback)
		if err != nil {
			h.Logs.Error(fmt.Sprintf("Failed to parse postback: %v", err))
			return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Invalid postback value",
			})
		}

		id, err := strconv.Atoi(ID)
		if err != nil {
			h.Logs.Error(fmt.Sprintf("Failed to parse id: %v", err))
			return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Invalid id value",
			})
		}

		_, errFind := h.DS.FindSummaryCampaign(id)

		if errFind != nil {
			h.Logs.Error(fmt.Sprintf("Data Not Found: %v", errFind))
			return c.Status(fiber.StatusNotFound).JSON(entity.GlobalResponse{
				Code:    fiber.StatusNotFound,
				Message: "Data Not Found",
			})
		}

		err = h.DS.UpdatePostbackModel(entity.SummaryCampaign{
			Postback: postbackInt,
		}, id)
		h.Logs.Info("Successfully updated postback in the database.")
		if err != nil {
			h.Logs.Error(fmt.Sprintf("Failed to update postback: %v", err))
			return c.Status(fiber.StatusInternalServerError).JSON(entity.GlobalResponse{
				Code:    fiber.StatusInternalServerError,
				Message: "Failed to update postback",
			})
		}

		return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{
			Code:    fiber.StatusOK,
			Message: "Postback updated successfully",
		})
	}

	return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{
		Code:    fiber.StatusBadRequest,
		Message: "Invalid request parameter",
	})
}

func (h *IncomingHandler) UpdateAgencyCost(c *fiber.Ctx) error {
	c.Set("content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("otf-8", "iso-8859-1")

	m := c.Queries()
	v := c.Params("v")

	if v == "update_agency_cost" {
		AgencyFee := m["agency_fee"]
		CostPerConversion := m["cost_per_conversion"]
		TechnicalFee := m["technical_fee"]

		h.Logs.Debug(fmt.Sprintf("Received agency_fee: %s, cost_per_conversion: %s , technical_fee: %s", AgencyFee, CostPerConversion, TechnicalFee))

		if AgencyFee == "" && CostPerConversion == "" && TechnicalFee == "" {
			h.Logs.Error("Missing required fields")
			return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Missing required fields",
			})
		}

		// var agencyFee float64
		// var err error
		// if AgencyFee != "" {
		// 	if err != nil {
		// 		h.Logs.Error(fmt.Sprintf("Failed to parse float: %v", err))
		// 		return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{
		// 			Code:    fiber.StatusBadRequest,
		// 			Message: "Invalid agency_fee value",
		// 		})
		// 	}
		// }
		// agencyFee, err := strconv.ParseFloat(AgencyFee, 64)

		// var costPerConversion float64
		// if CostPerConversion != "" {
		// 	if err != nil {
		// 		h.Logs.Error(fmt.Sprintf("Failed to parse float : %v", err))
		// 		return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{
		// 			Code:    fiber.StatusBadRequest,
		// 			Message: "Invalid cost_per_conversion value",
		// 		})
		// 	}
		// }
		// costPerConversion, err := strconv.ParseFloat(CostPerConversion, 64)

		// var technicalFee float64
		// if TechnicalFee != "" {
		// 	if err != nil {
		// 		h.Logs.Error(fmt.Sprintf("Failed to parse float : %v", err))
		// 		return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{
		// 			Code:    fiber.StatusBadRequest,
		// 			Message: "Invalid technical_fee value",
		// 		})
		// 	}
		// }
		// technicalFee, err := strconv.ParseFloat(TechnicalFee, 64)

		// err = h.DS.UpdateAgencyCostModel(entity.SummaryCampaign{
		// 	AgencyFee:         agencyFee,
		// 	CostPerConversion: costPerConversion,
		// 	TechnicalFee:      technicalFee,
		// })

		redisKey := strings.ToLower("global_setting")

		gset, err := json.Marshal(entity.GlobalSetting{
			AgencyFee:         AgencyFee,
			CostPerConversion: CostPerConversion,
			TechnicalFee:      TechnicalFee,
		})

		h.DS.SetData(redisKey, "$", string(gset))

		fmt.Printf("Successfully send to Redis AgencyFee: %s, CostPerConversion: %s, TechnicalFee: %s\n", AgencyFee, CostPerConversion, TechnicalFee)

		if err != nil {
			h.Logs.Error(fmt.Sprintf("Failed to update: %v", err))
			return c.Status(fiber.StatusInternalServerError).JSON(entity.GlobalResponse{
				Code:    fiber.StatusInternalServerError,
				Message: "Failed to update",
			})
		}

		return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{
			Code:    fiber.StatusOK,
			Message: fmt.Sprintf("--Successfully updated agency_fee: %s, cost_per_conversion: %s, technical_fee: %s", AgencyFee, CostPerConversion, TechnicalFee),
		})

	}

	return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{
		Code:    fiber.StatusBadRequest,
		Message: "Invalid request parameter",
	})

}
