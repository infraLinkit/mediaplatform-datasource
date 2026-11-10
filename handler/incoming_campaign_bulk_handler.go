package handler

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource/entity"
	"github.com/infraLinkit/mediaplatform-datasource/helper"
)

// campaignConfig reads the Redis config; when the key is missing it falls back to the database row, so
// one expired key cannot crash a bulk action half way through. The second result reports a Redis hit.
func (h *IncomingHandler) campaignConfig(redisKey, urlServiceKey string) (*entity.DataConfig, bool) {
	cfg, err := h.DS.GetDataConfig(redisKey, "$")
	if err == nil && cfg != nil {
		return cfg, true
	}
	fallback := &entity.DataConfig{}
	var detail entity.CampaignDetail
	if err := h.DB.Where("url_service_key = ?", urlServiceKey).First(&detail).Error; err == nil {
		fallback.Country = detail.Country
		fallback.Operator = detail.Operator
		fallback.Partner = detail.Partner
		fallback.Adnet = detail.Adnet
		fallback.Service = detail.Service
		fallback.CampaignId = detail.CampaignId
	}
	return fallback, false
}

// applyCampaignStatus switches one campaign detail in Redis, its row and today's summary.
func (h *IncomingHandler) applyCampaignStatus(urlServiceKey string, active bool) {
	cfgRediskey := helper.Concat("-", urlServiceKey, "configIdx")
	cfgCmp, cached := h.campaignConfig(cfgRediskey, urlServiceKey)
	cfgCmp.IsActive = active
	if active {
		cfgCmp.StatusCapping = false
	}
	if cached {
		cfgDataConfig, _ := json.Marshal(cfgCmp)
		h.DS.SetData(cfgRediskey, "$", string(cfgDataConfig))
	}

	h.DS.UpdateStatusCounterMOCampaignDetail(entity.CampaignDetail{
		URLServiceKey: urlServiceKey,
		IsActive:      active,
	})
	h.DS.UpdateSummaryCampaign(entity.SummaryCampaign{
		SummaryDate:   helper.GetCurrentTime(h.Config.TZ, time.RFC3339),
		Status:        active,
		URLServiceKey: urlServiceKey,
		Country:       cfgCmp.Country,
		Operator:      cfgCmp.Operator,
		Partner:       cfgCmp.Partner,
		Adnet:         cfgCmp.Adnet,
		Service:       cfgCmp.Service,
		CampaignId:    cfgCmp.CampaignId,
	})
}

// applyCampaignDelete removes one campaign detail from Redis and the database,
// and the campaign itself once its last detail is gone.
func (h *IncomingHandler) applyCampaignDelete(urlServiceKey, campaignID string) error {
	cfgRediskey := helper.Concat("-", urlServiceKey, "configIdx")
	cfgCmp, _ := h.campaignConfig(cfgRediskey, urlServiceKey)
	h.DS.DelData(cfgRediskey, "$")

	ctrRedisKey := helper.Concat("-", urlServiceKey, "counterIdx")
	h.DS.DelData(ctrRedisKey, "$")

	sumRedisKey := helper.Concat("-", urlServiceKey, "summary")
	h.DS.DelData(sumRedisKey, "$")

	h.DS.DelCampaignDetail(entity.CampaignDetail{
		URLServiceKey: urlServiceKey,
		Country:       cfgCmp.Country,
		Operator:      cfgCmp.Operator,
		Partner:       cfgCmp.Partner,
		Adnet:         cfgCmp.Adnet,
		Service:       cfgCmp.Service,
		CampaignId:    campaignID,
	})
	h.DS.DelSummaryCampaign(entity.SummaryCampaign{
		SummaryDate:   helper.GetCurrentTime(h.Config.TZ, time.RFC3339),
		URLServiceKey: urlServiceKey,
		Country:       cfgCmp.Country,
		Operator:      cfgCmp.Operator,
		Partner:       cfgCmp.Partner,
		Adnet:         cfgCmp.Adnet,
		Service:       cfgCmp.Service,
		CampaignId:    campaignID,
	})

	count, err := h.DS.CountCampaignDetailByCampaignID(entity.CampaignDetail{CampaignId: campaignID})
	if err != nil {
		return errors.New("Failed to count campaign details")
	}
	if count == 0 {
		if err := h.DS.DelCampaign(entity.Campaign{CampaignId: campaignID}); err != nil {
			return errors.New("Failed to delete campaign")
		}
	}
	return nil
}

type campaignBulkRequest struct {
	CampaignIDs []string `json:"campaign_ids"`
	SelectedIDs []int    `json:"selected_ids"`
	Status      bool     `json:"status"`
}

type campaignBulkTarget struct {
	URLServiceKey string
	CampaignID    string
}

// campaignBulkTargets resolves the selection the way CampaignRepository's bulk methods do: an unknown
// campaign stops everything, and a campaign with none of the selected details is skipped.
func (h *IncomingHandler) campaignBulkTargets(req campaignBulkRequest) ([]campaignBulkTarget, error) {
	if len(req.CampaignIDs) == 0 {
		return nil, errors.New("No campaign IDs provided.")
	}
	var targets []campaignBulkTarget
	for _, campaignID := range req.CampaignIDs {
		var known int64
		if err := h.DB.Table("campaign_details").Where("campaign_id = ?", campaignID).Count(&known).Error; err != nil {
			return nil, err
		}
		if known == 0 {
			return nil, errors.New("No campaign found for ID: " + campaignID)
		}
		var keys []string
		err := h.DB.Table("campaign_details").Where("campaign_id = ? AND id IN ?", campaignID, req.SelectedIDs).Pluck("url_service_key", &keys).Error
		if err != nil {
			return nil, err
		}
		for _, key := range keys {
			targets = append(targets, campaignBulkTarget{URLServiceKey: key, CampaignID: campaignID})
		}
	}
	return targets, nil
}

func (h *IncomingHandler) BulkUpdateCampaignStatus(c *fiber.Ctx) error {
	const failed = "An error occurred while preparing the campaigns for update. Please try again."
	var req campaignBulkRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}
	targets, err := h.campaignBulkTargets(req)
	if err != nil {
		h.Logs.Errorf("BulkUpdateCampaignStatus: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": failed})
	}
	for _, target := range targets {
		h.applyCampaignStatus(target.URLServiceKey, req.Status)
	}
	action := "disabled"
	if req.Status {
		action = "enabled"
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "message": "Selected campaigns have been " + action + "."})
}

func (h *IncomingHandler) BulkDeleteCampaigns(c *fiber.Ctx) error {
	const failed = "An error occurred while deleting the campaigns. Please try again."
	var req campaignBulkRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}
	targets, err := h.campaignBulkTargets(req)
	if err != nil {
		h.Logs.Errorf("BulkDeleteCampaigns: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": failed})
	}
	for _, target := range targets {
		if err := h.applyCampaignDelete(target.URLServiceKey, target.CampaignID); err != nil {
			h.Logs.Errorf("BulkDeleteCampaigns: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": failed})
		}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "message": "Selected campaigns have been deleted."})
}
