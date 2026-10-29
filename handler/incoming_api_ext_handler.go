package handler

import (
	"crypto/subtle"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource/config"
	"github.com/infraLinkit/mediaplatform-datasource/entity"
)

const (
	extDefaultLimit = 1000
	extMaxLimit     = 5000
)

type extCampaignDetail struct {
	Date       string  `json:"date"`
	CampaignId string  `json:"campaign_id"`
	Service    string  `json:"service"`
	ShortCode  string  `json:"shortcode"`
	Adnet      string  `json:"adnet"`
	Landing    int     `json:"landing"`
	MO         int     `json:"mo"`
	Postback   int     `json:"postback"`
	CrMO       float64 `json:"cr_mo"`
	CrPostback float64 `json:"cr_postback"`
	Cost       float64 `json:"cost"`
}

// ExtAPIKeyMiddleware guards /v1/ext with the static key from EXTAPIKEY (header X-API-Key).
func (h *IncomingHandler) ExtAPIKeyMiddleware(c *fiber.Ctx) error {
	key := h.Config.ExtApiKey
	if key == "" {
		h.Logs.Error("ExtAPIKeyMiddleware: EXTAPIKEY is not configured")
		return c.Status(fiber.StatusServiceUnavailable).JSON(entity.GlobalResponse{
			Code:    fiber.StatusServiceUnavailable,
			Message: "external api disabled",
		})
	}

	if subtle.ConstantTimeCompare([]byte(c.Get("X-API-Key")), []byte(key)) != 1 {
		return c.Status(fiber.StatusUnauthorized).JSON(entity.GlobalResponse{
			Code:    fiber.StatusUnauthorized,
			Message: "invalid api key",
		})
	}

	return c.Next()
}

// GetExtCampaignDetails returns date, campaign_id, service, shortcode, adnet, landing, mo, postback,
// cr_mo, cr_postback, cost (saaf) from summary_campaigns.
// Optional query params: service, shortcode, start_date, end_date (YYYY-MM-DD), limit, page.
// campaign_id in the response is summary_campaigns.url_service_key. The range defaults to today and
// rows are summed per date, url_service_key, service, short_code, adnet. cr_mo = mo/landing and
// cr_postback = postback/landing (ratios, same as the summary cron), recomputed from the sums.
func (h *IncomingHandler) GetExtCampaignDetails(c *fiber.Ctx) error {
	service := c.Query("service")
	shortCode := c.Query("shortcode", c.Query("short_code"))

	today := time.Now().Format("2006-01-02")
	startDate := c.Query("start_date", today)
	endDate := c.Query("end_date", startDate)

	start, errStart := time.Parse("2006-01-02", startDate)
	end, errEnd := time.Parse("2006-01-02", endDate)
	if errStart != nil || errEnd != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{
			Code:    fiber.StatusBadRequest,
			Message: "start_date/end_date must be YYYY-MM-DD",
		})
	}
	if end.Before(start) {
		return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{
			Code:    fiber.StatusBadRequest,
			Message: "end_date must be on/after start_date",
		})
	}

	limit := c.QueryInt("limit", extDefaultLimit)
	if limit < 1 || limit > extMaxLimit {
		limit = extDefaultLimit
	}
	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}

	q := h.DB.Model(&entity.SummaryCampaign{}).
		Select(`TO_CHAR(summary_date, 'YYYY-MM-DD') AS date,
			url_service_key AS campaign_id,
			service,
			short_code,
			adnet,
			COALESCE(SUM(landing), 0) AS landing,
			COALESCE(SUM(mo_received), 0) AS mo,
			COALESCE(SUM(postback), 0) AS postback,
			COALESCE(SUM(mo_received)::float8 / NULLIF(SUM(landing), 0), 0) AS cr_mo,
			COALESCE(SUM(postback)::float8 / NULLIF(SUM(landing), 0), 0) AS cr_postback,
			COALESCE(SUM(saaf), 0) AS cost`).
		Where("summary_date BETWEEN ? AND ?", startDate, endDate).
		Group("summary_date, url_service_key, service, short_code, adnet")

	if service != "" {
		q = q.Where("LOWER(service) = LOWER(?)", service)
	}
	if shortCode != "" {
		q = q.Where("LOWER(short_code) = LOWER(?)", shortCode)
	}

	var rows []extCampaignDetail
	if err := q.Order("summary_date ASC, url_service_key ASC, adnet ASC, service ASC, short_code ASC").
		Limit(limit).Offset((page - 1) * limit).
		Scan(&rows).Error; err != nil {
		h.Logs.WithError(err).Error("GetExtCampaignDetails: database error")
		return c.Status(fiber.StatusInternalServerError).JSON(entity.GlobalResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "database error",
		})
	}

	if rows == nil {
		rows = []extCampaignDetail{}
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithData{
		Code:    fiber.StatusOK,
		Message: config.OK_DESC,
		Data:    rows,
	})
}
