package handler

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource/config"
	"github.com/infraLinkit/mediaplatform-datasource/entity"
	"github.com/infraLinkit/mediaplatform-datasource/helper"
	"github.com/wiliehidayat87/rmqp"
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
			return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{Code: fiber.StatusBadRequest, Message: "v is empty"})
		} else {

			gs, _ := h.DS.GetDataConfig("global_setting", "$")

			tdb := gs.TargetDailyBudget
			if target_daily_budget != tdb {
				tdb = target_daily_budget
			}
			var dc []entity.DataConfig
			dc = append(dc, entity.DataConfig{
				TargetDailyBudget: tdb,
				Country:           country,
				Operator:          operator,
			})

			bodyReq, _ := json.Marshal(entity.DataCampaignAction{
				Action: "UPDATE_CAMP_MONITOR_BUDGET", DataConfig: dc})

			corId := "CMN" + helper.GetUniqId(h.Config.TZ)

			published := h.Rmqp.PublishMsg(rmqp.PublishItems{
				ExchangeName: h.Config.RabbitMQCampaignManagementExchangeName,
				QueueName:    h.Config.RabbitMQCampaignManagementQueueName,
				ContentType:  h.Config.RabbitMQDataType,
				CorId:        corId,
				Payload:      string(bodyReq),
				Priority:     0,
			})

			if !published {

				h.Logs.Debug(fmt.Sprintf("[x] Failed published: %s, Data: %s ...", corId, string(bodyReq)))

			} else {

				h.Logs.Debug(fmt.Sprintf("[v] Published: %s, Data: %s ...", corId, string(bodyReq)))
			}

			/* gs, _ := h.DS.GetDataConfig("global_setting", "$")

			tdb := gs.TargetDailyBudget
			if target_daily_budget != tdb {
				tdb = target_daily_budget
			}

			gset, _ := json.Marshal(entity.GlobalSetting{
				CostPerConversion: gs.CPCR,
				AgencyFee:         gs.AgencyFee,
				TargetDailyBudget: tdb,
			})

			h.DS.SetData("global_setting", "$", string(gset))

			h.DS.UpdateCampaignMonitoringBudget(entity.DataConfig{
				TargetDailyBudget: tdb,
				Country:           country,
				Operator:          operator,
			})

			h.DS.UpdateReportSummaryCampaignMonitoringBudget(helper.GetFormatTime(h.Config.TZ, "2006-01-02"),
				entity.DataConfig{
					TargetDailyBudget: tdb,
					Country:           country,
					Operator:          operator,
				}) */

			return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
		}
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

	// Parse Traffic Data
	//m := c.Queries()

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}
