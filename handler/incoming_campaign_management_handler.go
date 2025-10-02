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
		CampaignType: m["campaign_type"],
		CampaignId:   m["campaign_id"],
		Page:         page,
		Draw:         draw,
		Action:       m["action"],
		URLServiceKey:  m["url_service_key"],
		OrderColumn:    m["order_column"],
		OrderDir:       m["order_dir"],
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

func (h *IncomingHandler) EditCampaign(c *fiber.Ctx) error {
	o := new(entity.CampaignDetail)
	if err := c.BodyParser(&o); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	s := new(entity.SummaryCampaign)
	_ = c.BodyParser(&s)

	h.Logs.Debug(fmt.Sprintf("campaign detail : %#v ...", o))
	h.Logs.Debug(fmt.Sprintf("summary campaign : %#v ...", s))

	cfgRediskey := helper.Concat("-", o.URLServiceKey, "configIdx")
	cfgCmp, _ := h.DS.GetDataConfig(cfgRediskey, "$")
	mocappingChanged := o.MOCapping != cfgCmp.MOCapping

	cfgCmp.PO = o.PO
	cfgCmp.RatioSend = o.RatioSend
	cfgCmp.RatioReceive = o.RatioReceive
	cfgCmp.MOCapping = o.MOCapping
	cfgCmp.LastUpdate = helper.GetFormatTime(h.Config.TZ, time.RFC3339)
	if mocappingChanged {
		cfgCmp.StatusCapping = false
	}
	cfgDataConfig, _ := json.Marshal(cfgCmp)
	h.DS.SetData(cfgRediskey, "$", string(cfgDataConfig))

	now := time.Now().In(h.Config.TZ)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, h.Config.TZ)

	summaryLocal := s.SummaryDate.In(h.Config.TZ)
	summaryLocal = time.Date(
		summaryLocal.Year(),
		summaryLocal.Month(),
		summaryLocal.Day(),
		0, 0, 0, 0,
		h.Config.TZ,
	)
	
	if s.SummaryDate.IsZero() || summaryLocal.Equal(today) {
		h.DS.EditSettingCampaignDetail(entity.CampaignDetail{
			PO:            o.PO,
			MOCapping:     o.MOCapping,
			RatioSend:     o.RatioSend,
			RatioReceive:  o.RatioReceive,
			LastUpdate:    helper.GetCurrentTime(h.Config.TZ, time.RFC3339),
			URLServiceKey: o.URLServiceKey,
			Country:       cfgCmp.Country,
			Operator:      cfgCmp.Operator,
			Partner:       cfgCmp.Partner,
			Adnet:         cfgCmp.Adnet,
			Service:       cfgCmp.Service,
			CampaignId:    o.CampaignId,
			StatusCapping:  bool(mocappingChanged),
		})
	}

	pos, _ := strconv.ParseFloat(strings.TrimSpace(o.PO), 64)
	h.DS.EditSettingSummaryCampaign(entity.SummaryCampaign{
		SummaryDate:   s.SummaryDate,
		PO:            pos,
		MOLimit:       o.MOCapping,
		RatioSend:     o.RatioSend,
		RatioReceive:  o.RatioReceive,
		URLServiceKey: o.URLServiceKey,
		Country:       cfgCmp.Country,
		Operator:      cfgCmp.Operator,
		Partner:       cfgCmp.Partner,
		Adnet:         cfgCmp.Adnet,
		Service:       cfgCmp.Service,
		CampaignId:    o.CampaignId,
	})

	return c.Status(fiber.StatusOK).Send([]byte("OK"))
}

func (h *IncomingHandler) DelCampaign(c *fiber.Ctx) error {

	o := new(entity.CampaignDetail)

	if err := c.BodyParser(&o); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	} else {

		h.Logs.Debug(fmt.Sprintf("data : %#v ...", o))

		// DELETE REDIS KEY
		cfgRediskey := helper.Concat("-", o.URLServiceKey, "configIdx")
		cfgCmp, _ := h.DS.GetDataConfig(cfgRediskey, "$")
		h.DS.DelData(cfgRediskey, "$")

		// DROP Index redis
		h.DS.R.Conn().B().FtDropindex().Index(cfgRediskey).Build()

		ctrRedisKey := helper.Concat("-", o.URLServiceKey, "counterIdx")

		h.DS.DelData(ctrRedisKey, "$")

		// DROP Index redis
		h.DS.R.Conn().B().FtDropindex().Index(ctrRedisKey).Build()

		sumRedisKey := helper.Concat("-", o.URLServiceKey, "summary")

		h.DS.DelData(sumRedisKey, "$")

		// DROP Index redis
		h.DS.R.Conn().B().FtDropindex().Index(sumRedisKey).Build()

		h.DS.DelCampaignDetail(entity.CampaignDetail{
			URLServiceKey: o.URLServiceKey,
			Country:       cfgCmp.Country,
			Operator:      cfgCmp.Operator,
			Partner:       cfgCmp.Partner,
			Adnet:         cfgCmp.Adnet,
			Service:       cfgCmp.Service,
			CampaignId:    o.CampaignId,
		})

		h.DS.DelSummaryCampaign(entity.SummaryCampaign{
			SummaryDate:   helper.GetCurrentTime(h.Config.TZ, time.RFC3339),
			URLServiceKey: o.URLServiceKey,
			Country:       cfgCmp.Country,
			Operator:      cfgCmp.Operator,
			Partner:       cfgCmp.Partner,
			Adnet:         cfgCmp.Adnet,
			Service:       cfgCmp.Service,
			CampaignId:    o.CampaignId,
		})

		count, err := h.DS.CountCampaignDetailByCampaignID(entity.CampaignDetail{CampaignId: o.CampaignId})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to count campaign details"})
		}

		if count == 0 {
			err = h.DS.DelCampaign(entity.Campaign{CampaignId: o.CampaignId})
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete campaign"})
			}
		}

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

func (h *IncomingHandler) UpdateKeyMainstream(c *fiber.Ctx) error {

	o := new(entity.CampaignDetail)

	if err := c.BodyParser(&o); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	} else {

		h.Logs.Debug(fmt.Sprintf("data : %#v ...", o))

		// Update to redis with key
		cfgRediskey := helper.Concat("-", o.URLServiceKey, "configIdx")
		cfgCmp, _ := h.DS.GetDataConfig(cfgRediskey, "$")
		cfgCmp.StatusSubmitKeyMainstream = o.StatusSubmitKeyMainstream
		cfgCmp.KeyMainstream = o.KeyMainstream

		cfgDataConfig, _ := json.Marshal(cfgCmp)

		h.DS.SetData(cfgRediskey, "$", string(cfgDataConfig))

		err = h.DS.UpdateKeyMainstreamCampaignDetail(entity.CampaignDetail{
			StatusSubmitKeyMainstream: o.StatusSubmitKeyMainstream,
			KeyMainstream:             o.KeyMainstream,
			URLServiceKey:             o.URLServiceKey,
			CampaignId:                o.CampaignId,
		})

		if err != nil {
			h.Logs.Error(fmt.Sprintf("failed updating db: %v", err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed save to db"})
		}

		return c.Status(fiber.StatusOK).Send([]byte("OK"))
	}
}

func (h *IncomingHandler) UpdateGoogleSheet(c *fiber.Ctx) error {

	o := new(entity.CampaignDetail)

	if err := c.BodyParser(&o); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	} else {

		h.Logs.Debug(fmt.Sprintf("data : %#v ...", o))

		// Update to redis with key
		cfgRediskey := helper.Concat("-", o.URLServiceKey, "configIdx")
		cfgCmp, _ := h.DS.GetDataConfig(cfgRediskey, "$")
		cfgCmp.GoogleSheet = o.GoogleSheet

		cfgDataConfig, _ := json.Marshal(cfgCmp)

		h.DS.SetData(cfgRediskey, "$", string(cfgDataConfig))

		err = h.DS.UpdateGoogleSheetCampaignDetail(entity.CampaignDetail{
			GoogleSheet:               o.GoogleSheet,
			URLServiceKey:             o.URLServiceKey,
			CampaignId:                o.CampaignId,
		})

		if err != nil {
			h.Logs.Error(fmt.Sprintf("failed updating db: %v", err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed save to db"})
		}

		return c.Status(fiber.StatusOK).Send([]byte("OK"))
	}
}

func (h *IncomingHandler) EditMOCappingServiceS2S(c *fiber.Ctx) error {
	o := new(entity.CampaignDetail)
	if err := c.BodyParser(o); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	h.Logs.Debug(fmt.Sprintf("data : %#v ...", o))

	keys, err := h.DS.ScanKeys("*-configIdx")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal ambil key Redis"})
	}

	updated := false
	updatedCount := 0

	for _, key := range keys {
		data, err := h.DS.GetDataJSON(key)
		if err != nil {
			continue
		}

		country := data["country"]
		operator := data["operator"]
		partner := data["partner"]
		service := data["service"]

		if country == o.Country && operator == o.Operator && partner == o.Partner && service == o.Service {
			h.DS.SetData(key, "$.mo_capping_service", strconv.Itoa(o.MOCappingService))
			h.DS.SetData(key, "$.status_capping", strconv.FormatBool(false))
			updated = true
			updatedCount++
			h.Logs.Debug(fmt.Sprintf("Updated config key: %s", key))
		}
	}

	if !updated {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Data not found"})
	}

	h.Logs.Debug(fmt.Sprintf("Total configs updated: %d", updatedCount))

	err = h.DS.UpdateMOCappingS2S(entity.CampaignDetail{
		MOCappingService: o.MOCappingService,
		LastUpdate:       helper.GetCurrentTime(h.Config.TZ, time.RFC3339),
		Country:          o.Country,
		Operator:         o.Operator,
		Partner:          o.Partner,
		Service:          o.Service,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update DB"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       fmt.Sprintf("Success update %d configs", updatedCount),
		"updated_count": updatedCount,
	})
}

func (h *IncomingHandler) EditPOAF(c *fiber.Ctx) error {

	o := new(entity.SummaryCampaign)

	if err := c.BodyParser(&o); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	} else {

		h.Logs.Debug(fmt.Sprintf("data : %#v ...", o))

		h.DS.EditPOAFIncSummaryCampaign(entity.IncSummaryCampaign{
			SummaryDate:   o.SummaryDate,
			POAF: 		   o.POAF,
			URLServiceKey: o.URLServiceKey,
		})

		h.DS.EditPOAFSummaryCampaign(entity.SummaryCampaign{
			SummaryDate:   o.SummaryDate,
			POAF:          o.POAF,
			URLServiceKey: o.URLServiceKey,
		})

		return c.Status(fiber.StatusOK).Send([]byte("OK"))
	}
}


func (h *IncomingHandler) EditCampaignManagementDetail(c *fiber.Ctx) error {

	o := new(entity.CampaignDetail)

	if err := c.BodyParser(&o); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	} else {

		h.Logs.Debug(fmt.Sprintf("data : %#v ...", o))

		// Update to redis with key
		cfgRediskey := helper.Concat("-", o.URLServiceKey, "configIdx")
		cfgCmp, _ := h.DS.GetDataConfig(cfgRediskey, "$")

		cfgCmp.APIURL = o.APIURL

		cfgDataConfig, _ := json.Marshal(cfgCmp)

		h.DS.SetData(cfgRediskey, "$", string(cfgDataConfig))

		h.DS.EditCampaignManagementDetail(entity.CampaignDetail{
			APIURL: cfgCmp.APIURL,
			URLServiceKey: o.URLServiceKey,
			CampaignId:    o.CampaignId,
		})

		return c.Status(fiber.StatusOK).Send([]byte("OK"))
	}
}