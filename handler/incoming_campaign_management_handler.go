package handler

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource/config"
	"github.com/infraLinkit/mediaplatform-datasource/entity"
	"github.com/infraLinkit/mediaplatform-datasource/helper"
	"github.com/wiliehidayat87/rmqp"
)

func (h *IncomingHandler) DisplayCampaignManagement(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	m := c.Queries()
	page, _ := strconv.Atoi(m["page"])
	draw, _ := strconv.Atoi(m["draw"])

	fe := entity.DisplayCampaignManagement{
		Adnet:        m["adnet"],
		Country:      m["country"],
		Service:      m["service"],
		Operator:     m["operator"],
		Partner:      m["partner"],
		Status:       m["status"],
		CampaignName: m["campaign_name"],
		CampaignId:   m["campaign_id"],
		Page:         page,
		Action:       m["action"],
	}

	v := c.Params("v")

	if v == "detail" {
		r := h.DisplayCampaignManagementDetail(c, fe)
		return c.Status(r.HttpStatus).JSON(r.Rsp)
	}

	r := h.DisplayCampaignManagementExtra(c, fe, draw)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayCampaignManagementExtra(c *fiber.Ctx, fe entity.DisplayCampaignManagement, draw int) entity.ReturnResponse {
	key := "temp_key_api_campaign_management_" + strings.ReplaceAll(helper.GetIpAddress(c), ".", "_")

	var (
		err                       error
		isempty                   bool
		campaignmanagement        []entity.CampaignManagementData
		displaycampaignmanagement []entity.CampaignManagementData
	)

	if fe.Action != "" {
		campaignmanagement, err = h.DS.GetCampaignManagement(fe)
	} else {
		if campaignmanagement, isempty = h.DS.RGetCampaignManagement(key, "$"); isempty {
			campaignmanagement, err = h.DS.GetCampaignManagement(fe)
			s, _ := json.Marshal(campaignmanagement)
			h.DS.SetData(key, "$", string(s))
			h.DS.SetExpireData(key, 60)
		}
	}

	if err == nil {
		totalRecords := len(campaignmanagement)
		pagesize := PAGESIZE
		start := (fe.Page - 1) * pagesize
		end := start + pagesize

		if start < totalRecords {
			if end > totalRecords {
				end = totalRecords
			}
			displaycampaignmanagement = campaignmanagement[start:end]
		}

		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithDataTable{
				Draw:            fe.Page,
				Code:            fiber.StatusOK,
				Message:         config.OK_DESC,
				Data:            displaycampaignmanagement,
				RecordsTotal:    totalRecords,
				RecordsFiltered: totalRecords,
			},
		}
	}

	return entity.ReturnResponse{
		HttpStatus: fiber.StatusNotFound,
		Rsp: entity.GlobalResponse{
			Code:    fiber.StatusNotFound,
			Message: "empty",
		},
	}
}

func (h *IncomingHandler) DisplayCampaignManagementDetail(c *fiber.Ctx, fe entity.DisplayCampaignManagement) entity.ReturnResponse {
	var (
		err                             error
		campaignmanagementdetail        []entity.CampaignManagementDataDetail
		displaycampaignmanagementdetail []entity.CampaignManagementDataDetail
	)

	campaignmanagementdetail, err = h.DS.GetCampaignManagementDetail(fe)

	if err == nil {
		displaycampaignmanagementdetail = campaignmanagementdetail
		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithData{
				Code:    fiber.StatusOK,
				Message: config.OK_DESC,
				Data:    displaycampaignmanagementdetail,
			},
		}
	}

	return entity.ReturnResponse{
		HttpStatus: fiber.StatusNotFound,
		Rsp: entity.GlobalResponseWithData{
			Code:    fiber.StatusNotFound,
			Message: "empty",
		},
	}
}

func (h *IncomingHandler) SendCampaignHandler(c *fiber.Ctx) error {
	var campaignData entity.DataCampaignAction

	if err := c.BodyParser(&campaignData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	bodyReq, err := json.Marshal(campaignData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to serialize data"})
	}

	published := h.Rmqp.PublishMsg(rmqp.PublishItems{
		ExchangeName: h.Config.RabbitMQCampaignManagementExchangeName,
		QueueName:    h.Config.RabbitMQCampaignManagementQueueName,
		ContentType:  h.Config.RabbitMQDataType,
		Payload:      string(bodyReq),
		Priority:     0,
	})

	if !published {
		h.Logs.Debug(fmt.Sprintf("[x] Failed published: Data: %s ...", string(bodyReq)))
	} else {
		h.Logs.Debug(fmt.Sprintf("[v] Published: Data: %s ...", string(bodyReq)))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Campaign data sent to RabbitMQ"})
}
