package handler

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource/config"
	"github.com/infraLinkit/mediaplatform-datasource/entity"
	"gorm.io/gorm"
)

// phpFormatNumber is the CMS format_number helper: thousands separators, up to ten decimals, no trailing zeros.
func phpFormatNumber(value float64) string {
	text := strconv.FormatFloat(value, 'f', 10, 64)
	whole, fraction, _ := strings.Cut(text, ".")
	negative := strings.HasPrefix(whole, "-")
	whole = strings.TrimPrefix(whole, "-")
	var grouped strings.Builder
	for i, digit := range whole {
		if i > 0 && (len(whole)-i)%3 == 0 {
			grouped.WriteByte(',')
		}
		grouped.WriteRune(digit)
	}
	result := grouped.String()
	if fraction = strings.TrimRight(fraction, "0"); fraction != "" {
		result += "." + fraction
	}
	if negative && result != "0" {
		result = "-" + result
	}
	return result
}

// phpFloat is PHP's (float) cast on the po column, whose empty value is the text NA.
func phpFloat(value string) float64 {
	number, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil {
		return 0
	}
	return number
}

type campaignSettingRow struct {
	ID                        int        `json:"id"`
	URLServiceKey             string     `json:"url_service_key"`
	CampaignName              string     `json:"campaign_name"`
	CampaignObjective         string     `json:"campaign_objective"`
	CampaignID                string     `json:"campaign_id"`
	Country                   string     `json:"country"`
	Operator                  string     `json:"operator"`
	Partner                   string     `json:"partner"`
	Publisher                 string     `json:"publisher"`
	Service                   string     `json:"service"`
	Channel                   string     `json:"channel"`
	URLService                string     `json:"url_service"`
	URLLanding                string     `json:"url_landing"`
	URLWarpLanding            string     `json:"url_warp_landing"`
	StatusSubmitKeyMainstream bool       `json:"status_submit_key_mainstream"`
	KeyMainstream             string     `json:"key_mainstream"`
	GoogleSheet               string     `json:"google_sheet"`
	GoogleSheetBillable       string     `json:"google_sheet_billable"`
	MOCapping                 int        `json:"mo_capping"`
	PO                        string     `json:"po"`
	MOLimit                   string     `gorm:"-" json:"mo_limit"`
	Payout                    string     `gorm:"-" json:"payout"`
	RatioSend                 int        `json:"ratio_send"`
	RatioReceive              int        `json:"ratio_receive"`
	CreatedAt                 *time.Time `json:"created_at"`
}

const campaignSettingColumns = `d.id, d.url_service_key, c.name AS campaign_name, c.campaign_objective, d.campaign_id, d.country, d.operator,
	d.partner, d.adnet AS publisher, d.service, d.channel, d.url_service, d.url_landing, d.url_warp_landing,
	d.status_submit_key_mainstream, d.key_mainstream, d.google_sheet, d.google_sheet_billable,
	d.mo_capping, d.po, d.ratio_send, d.ratio_receive, c.created_at`

func campaignSettingPage(c *fiber.Ctx) (page, size int) {
	size, err := strconv.Atoi(c.Query("page_size"))
	if err != nil || size <= 0 {
		size = 10
	}
	page, _ = strconv.Atoi(c.Query("page"))
	if page < 1 {
		page = 1
	}
	return page, size
}

func applyCampaignDetailFilters(c *fiber.Ctx, query *gorm.DB, filters map[string]string) *gorm.DB {
	for param, column := range filters {
		if value := c.Query(param); value != "" {
			query = query.Where(column+" = ?", value)
		}
	}
	return query
}

// campaignSettingSearch is the DataTables search box over the visible columns.
func campaignSettingSearch(query *gorm.DB, keyword string, columns ...string) *gorm.DB {
	if keyword == "" {
		return query
	}
	conditions := make([]string, len(columns))
	values := make([]interface{}, len(columns))
	for i, column := range columns {
		conditions[i] = "CAST(" + column + " AS TEXT) ILIKE ?"
		values[i] = "%" + keyword + "%"
	}
	return query.Where("("+strings.Join(conditions, " OR ")+")", values...)
}

func (h *IncomingHandler) campaignSettingRows(c *fiber.Ctx, objective string, filters map[string]string, searchColumns []string) error {
	query := h.DB.Table("campaign_details AS d").
		Joins("JOIN campaigns c ON c.campaign_id = d.campaign_id").
		Where(objective)
	query = applyCampaignDetailFilters(c, query, filters)
	if before, after := c.Query("created_date_before"), c.Query("created_date_after"); before != "" && after != "" {
		query = query.Where("c.created_at BETWEEN ? AND ?", before, after)
	}
	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		h.Logs.Errorf("campaignSettingRows: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch campaign setting"})
	}
	query = campaignSettingSearch(query, c.Query("search"), searchColumns...)
	var filtered int64
	if err := query.Session(&gorm.Session{}).Count(&filtered).Error; err != nil {
		h.Logs.Errorf("campaignSettingRows: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch campaign setting"})
	}

	order := "d.created_at DESC"
	if c.Query("order_column") == "created_at" {
		order = "c.created_at DESC"
		if strings.EqualFold(c.Query("order_dir"), "asc") {
			order = "c.created_at ASC"
		}
	}
	page, size := campaignSettingPage(c)
	rows := []campaignSettingRow{}
	err := query.Select(campaignSettingColumns).Order(order).Limit(size).Offset((page - 1) * size).Scan(&rows).Error
	if err != nil {
		h.Logs.Errorf("campaignSettingRows: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch campaign setting"})
	}
	for i := range rows {
		rows[i].MOLimit = phpFormatNumber(float64(rows[i].MOCapping))
		rows[i].Payout = phpFormatNumber(phpFloat(rows[i].PO))
	}
	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithDataTable{
		Code: fiber.StatusOK, Message: config.OK_DESC, Data: rows, RecordsTotal: int(total), RecordsFiltered: int(filtered),
	})
}

// CampaignSettingAdnetLevel is CampaignRepository::getCampaignSetting.
func (h *IncomingHandler) CampaignSettingAdnetLevel(c *fiber.Ctx) error {
	return h.campaignSettingRows(c, "c.campaign_objective NOT LIKE '%MAINSTREAM%'",
		map[string]string{"country": "d.country", "operator": "d.operator", "publisher": "d.adnet", "service": "d.service"},
		[]string{"d.url_service_key", "c.name", "d.country", "d.operator", "d.partner", "d.adnet", "d.service", "d.mo_capping", "d.po", "d.url_service", "d.url_warp_landing", "c.created_at"})
}

// CampaignSettingMainstream is CampaignRepository::getCampaignMainstreamSetting.
func (h *IncomingHandler) CampaignSettingMainstream(c *fiber.Ctx) error {
	return h.campaignSettingRows(c, "c.campaign_objective LIKE '%MAINSTREAM%'",
		map[string]string{"campaign_name": "c.name", "country": "d.country", "operator": "d.operator", "service": "d.service"},
		[]string{"d.url_service_key", "c.name", "d.country", "d.operator", "d.adnet", "d.service", "d.channel", "d.url_service", "d.url_warp_landing", "d.mo_capping", "d.po", "c.created_at"})
}

// CampaignSettingServiceLevel is CampaignRepository::getCampaignSettingService: one row per distinct service limit.
func (h *IncomingHandler) CampaignSettingServiceLevel(c *fiber.Ctx) error {
	base := h.DB.Table("campaign_details AS d").
		Joins("JOIN campaigns c ON c.campaign_id = d.campaign_id").
		Where("c.campaign_objective NOT LIKE '%MAINSTREAM%' AND c.campaign_objective <> 'UPLOAD SMS'")
	base = applyCampaignDetailFilters(c, base, map[string]string{"country": "d.country", "operator": "d.operator", "service": "d.service"})
	const columns = "d.country, d.operator, d.partner, d.service, d.mo_capping_service"

	count := func(query *gorm.DB) (int64, error) {
		var total int64
		err := h.DB.Table("(?) AS combos", query.Session(&gorm.Session{}).Distinct(columns)).Count(&total).Error
		return total, err
	}
	total, err := count(base)
	if err != nil {
		h.Logs.Errorf("CampaignSettingServiceLevel: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch campaign setting"})
	}
	searched := campaignSettingSearch(base, c.Query("search"), "d.country", "d.operator", "d.partner", "d.service", "d.mo_capping_service")
	filtered, err := count(searched)
	if err != nil {
		h.Logs.Errorf("CampaignSettingServiceLevel: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch campaign setting"})
	}

	type serviceRow struct {
		Country          string `json:"country"`
		Operator         string `json:"operator"`
		Partner          string `json:"partner"`
		Service          string `json:"service"`
		MOCappingService int    `json:"mo_capping_service"`
	}
	page, size := campaignSettingPage(c)
	rows := []serviceRow{}
	err = searched.Distinct(columns).Order("d.country, d.operator, d.partner, d.service").Limit(size).Offset((page - 1) * size).Scan(&rows).Error
	if err != nil {
		h.Logs.Errorf("CampaignSettingServiceLevel: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch campaign setting"})
	}
	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithDataTable{
		Code: fiber.StatusOK, Message: config.OK_DESC, Data: rows, RecordsTotal: int(total), RecordsFiltered: int(filtered),
	})
}

// RequireRoleCode is for what Blade gates on the role itself rather than on a menu permission.
func (h *IncomingHandler) RequireRoleCode(code string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, identified := c.Locals("user_id").(int)
		if !identified {
			return c.Next()
		}
		var granted int64
		err := h.DB.Table("users").
			Joins("JOIN roles ON roles.id = users.role_id").
			Where("users.id = ? AND users.deleted_at IS NULL AND roles.code = ?", userID, code).
			Count(&granted).Error
		if err != nil {
			h.Logs.Errorf("RequireRoleCode: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to check permission"})
		}
		if granted == 0 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "You do not have permission.", "status": "error"})
		}
		return c.Next()
	}
}

// lpBannerExtension allows the Laravel image types except SVG: /storage is public, so an SVG would be stored XSS.
func lpBannerExtension(detected string) string {
	switch detected {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	}
	return ""
}

func lpDesignErrors(lpType, description string) []string {
	var errs []string
	if strings.TrimSpace(lpType) == "" {
		errs = append(errs, "The mainstream lp type field is required.")
	} else if len([]rune(lpType)) > 255 {
		errs = append(errs, fmt.Sprintf("The mainstream lp type field must not be greater than %d characters.", 255))
	}
	if strings.TrimSpace(description) == "" {
		errs = append(errs, "The description field is required.")
	}
	return errs
}
