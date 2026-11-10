package handler

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource/entity"
	"github.com/lib/pq"
)

// RequireUnique replaces the unique: rules in ServiceManagementController. On update routes the
// row named by :id is excluded. model must be a pointer to the entity, so gorm finds its table.
func (h *IncomingHandler) RequireUnique(model interface{}, columns ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body map[string]interface{}
		if err := json.Unmarshal(c.Body(), &body); err != nil {
			return c.Next()
		}
		for _, column := range columns {
			value, isText := body[column].(string)
			if !isText || value == "" {
				continue
			}
			query := h.DB.Model(model).Where(column+" = ?", value)
			if id := c.Params("id"); id != "" {
				query = query.Where("id <> ?", id)
			}
			var count int64
			if err := query.Count(&count).Error; err != nil {
				h.Logs.Errorf("RequireUnique %s: %v", column, err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to validate input"})
			}
			if count > 0 {
				return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
					"message": "Input Not Valid",
					"errors":  fiber.Map{column: []string{"The " + strings.ReplaceAll(column, "_", " ") + " has already been taken."}},
				})
			}
		}
		return c.Next()
	}
}

// nextMaskedID replaces getMaxId and getMaxIdCompanyGroup, which read the tables directly.
func (h *IncomingHandler) nextMaskedID(model interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var max int64
		if err := h.DB.Model(model).Select("COALESCE(MAX(id), 0)").Scan(&max).Error; err != nil {
			h.Logs.Errorf("nextMaskedID: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to read id"})
		}
		return c.JSON(fiber.Map{"id": max + 1})
	}
}

func (h *IncomingHandler) NextCompanyMaskedID(c *fiber.Ctx) error {
	return h.nextMaskedID(&entity.Company{})(c)
}

func (h *IncomingHandler) NextCompanyGroupMaskedID(c *fiber.Ctx) error {
	return h.nextMaskedID(&entity.CompanyGroup{})(c)
}

type UpdateAdnetRequest struct {
	Name        string   `json:"name"`
	APIURL      string   `json:"api_url"`
	IsActive    bool     `json:"is_active"`
	IsDSP       bool     `json:"is_dsp"`
	CcEmail     []string `json:"cc_email"`
	IsUpdateAll bool     `json:"is_update_all"`
}

// UpdateAdnetAndCampaigns replaces ServiceManagementController::editAdnet, which updated
// adnet_lists and, on request, campaign_details directly.
func (h *IncomingHandler) UpdateAdnetAndCampaigns(c *fiber.Ctx) error {
	var req UpdateAdnetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}
	if req.Name == "" || req.APIURL == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": "Input Not Valid", "errors": fiber.Map{"api_url": []string{"The name and API postback URL are required."}}})
	}

	var adnet struct{ Code string }
	found := h.DB.Table("adnet_lists").Select("code").Where("id = ?", c.Params("id")).Limit(1).Scan(&adnet)
	if found.Error != nil {
		h.Logs.Errorf("UpdateAdnetAndCampaigns: find: %v", found.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update adnet"})
	}
	if found.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Adnet not found"})
	}

	if err := h.DB.Exec("UPDATE adnet_lists SET name = ?, api_url = ?, is_active = ?, cc_email = ?, is_dsp = ?, lastupdate = ? WHERE id = ?",
		strings.ToUpper(req.Name), req.APIURL, req.IsActive, pq.StringArray(req.CcEmail), req.IsDSP, time.Now(), c.Params("id")).Error; err != nil {
		h.Logs.Errorf("UpdateAdnetAndCampaigns: update: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update adnet"})
	}
	if req.IsUpdateAll {
		if err := h.DB.Exec("UPDATE campaign_details SET api_url = ? WHERE adnet = ?", req.APIURL, adnet.Code).Error; err != nil {
			h.Logs.Errorf("UpdateAdnetAndCampaigns: campaigns: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Adnet saved, but campaigns were not updated"})
		}
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Success"})
}
