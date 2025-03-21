package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

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
		Draw:         draw,
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
		campaignmanagement, _, err = h.DS.GetCampaignManagement(fe)
	} else {
		if campaignmanagement, isempty = h.DS.RGetCampaignManagement(key, "$"); isempty {
			campaignmanagement, _, err = h.DS.GetCampaignManagement(fe)
			s, _ := json.Marshal(campaignmanagement)
			h.DS.SetData(key, "$", string(s))
			h.DS.SetExpireData(key, 60)
		}
	}

	if err == nil {
		totalRecords := len(campaignmanagement)
		pagesize := PAGESIZE
		page := 1
		if fe.Page > 0 {
			page = fe.Page
		}
		start := (page - 1) * pagesize
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
				Draw:            fe.Draw,
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
			Rsp: entity.GlobalResponseWithDataTable{
				Draw:    fe.Draw,
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
	var campaignData map[string]interface{}

	if err := c.BodyParser(&campaignData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var bodyReq bytes.Buffer
	encoder := json.NewEncoder(&bodyReq)
	encoder.SetEscapeHTML(false) // Prevents encoding & as \u0026

	if err := encoder.Encode(campaignData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to serialize data"})
	}

	published := h.Rmqp.PublishMsg(rmqp.PublishItems{
		ExchangeName: h.Config.RabbitMQCampaignManagementExchangeName,
		QueueName:    h.Config.RabbitMQCampaignManagementQueueName,
		ContentType:  h.Config.RabbitMQDataType,
		Payload:      bodyReq.String(), // Send the properly formatted JSON
		Priority:     0,
	})

	if !published {
		h.Logs.Debug(fmt.Sprintf("[x] Failed published: Data: %s ...", bodyReq.String()))
	} else {
		h.Logs.Debug(fmt.Sprintf("[v] Published: Data: %s ...", bodyReq.String()))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Campaign data sent to RabbitMQ"})
}

func (h *IncomingHandler) UpdateStatusCampaign(c *fiber.Ctx) error {

	o := new(entity.CampaignDetail)

	if err := c.BodyParser(&o); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	} else {

		h.Logs.Debug(fmt.Sprintf("data : %#v ...", o))

		// Update to redis with key
		cfgRediskey := helper.Concat("-", o.URLServiceKey, "configIdx")
		cfgCmp, _ := h.DS.GetDataConfig(cfgRediskey, "$")
		cfgCmp.IsActive = o.IsActive

		cfgDataConfig, _ := json.Marshal(cfgCmp)

		h.DS.SetData(cfgRediskey, "$", string(cfgDataConfig))

		// Update to database
		h.DS.UpdateCampaignDetail(entity.CampaignDetail{
			ID:       o.ID,
			IsActive: o.IsActive,
		})

		h.DS.UpdateSummaryCampaign(entity.SummaryCampaign{
			SummaryDate:   helper.GetCurrentTime(h.Config.TZ, time.RFC3339),
			Status:        o.IsActive,
			URLServiceKey: o.URLServiceKey,
			Country:       cfgCmp.Country,
			Operator:      cfgCmp.Operator,
			Partner:       cfgCmp.Partner,
			Adnet:         cfgCmp.Adnet,
			Service:       cfgCmp.Service,
			CampaignId:    cfgCmp.CampaignId,
		})

		return c.Status(fiber.StatusOK).Send([]byte("OK"))
	}
}

func (h *IncomingHandler) GetCampaignCounts(c *fiber.Ctx) error {
	var input entity.DisplayCampaignManagement

	if err := c.QueryParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Ambil hanya counts dari GetCampaignManagement
	_, counts, err := h.DS.GetCampaignManagement(input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"total_campaign":          counts.TotalCampaigns,
		"total_active_campaign":   counts.TotalActiveCampaigns,
		"total_inactive_campaign": counts.TotalNonActiveCampaigns,
	})
}
