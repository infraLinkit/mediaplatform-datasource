package handler

import (
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// alertSource mirrors AlertDashboardRepository. The tables are the singular ones the Laravel
// models read, not the plural names on the entity.SummaryMo family.
type alertSource struct {
	table     string
	errorName string
	columns   []string
}

var (
	alertNoMo = alertSource{"summary_mo", "No MO", []string{
		"summary_date", "campaign_id", "campaign_name", "country", "operator", "service", "status", "pic_mo_report", "desc_mo",
		"mo_counter_6", "mo_counter_5", "mo_counter_4", "mo_counter_3", "mo_counter_2", "mo_counter_1", "mo_limit", "mo", "traffic", "is_lock",
	}}
	alertCapReached = alertSource{"summary_capping", "Cap Reached", []string{
		"summary_date", "campaign_id", "campaign_name", "partner", "country", "operator", "service", "status", "mo", "mo_limit", "desc_mo", "pic_mo_report", "is_lock",
	}}
	alertCrDrop = alertSource{"summary_cr", "CR Drop", []string{
		"summary_date", "campaign_id", "campaign_name", "country", "operator", "service", "status", "pic_cr_report", "desc_cr",
		"cr_counter_6", "cr_counter_5", "cr_counter_4", "cr_counter_3", "cr_counter_2", "cr_counter_1", "cr", "is_lock",
	}}
	alertRatio = alertSource{"summary_ratio", "Ratio Alert", []string{
		"summary_date", "campaign_id", "campaign_name", "country", "operator", "service", "status", "pic_ratio_report", "desc_ratio", "postback", "mo", "actual_ratio", "ratio", "is_lock",
	}}
	alertSourcesByType = map[string][]alertSource{
		"all":        {alertNoMo, alertCapReached, alertCrDrop, alertRatio},
		"nomo":       {alertNoMo},
		"crdrop":     {alertCrDrop},
		"ratio":      {alertRatio},
		"capreached": {alertCapReached},
	}
	alertSourceByError = map[string]alertSource{
		alertNoMo.errorName:       alertNoMo,
		alertCapReached.errorName: alertCapReached,
		alertCrDrop.errorName:     alertCrDrop,
		alertRatio.errorName:      alertRatio,
	}
)

const alertTimeLayout = "2006-01-02 15:04:05.999999"

func (h *IncomingHandler) alertNow() time.Time {
	if h.Config != nil && h.Config.TZ != nil {
		return time.Now().In(h.Config.TZ)
	}
	return time.Now()
}

// alertDateBounds follows AlertDashboardRepository::applyDateFilter. ok is false for "All Date Range".
func alertDateBounds(dateRange, custom string, now time.Time) (start, end time.Time, ok bool) {
	day := func(t time.Time) time.Time { return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()) }
	endOfDay := func(t time.Time) time.Time { return day(t).AddDate(0, 0, 1).Add(-time.Microsecond) }
	month := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	switch dateRange {
	case "Today":
		return day(now), endOfDay(now), true
	case "Yesterday":
		y := now.AddDate(0, 0, -1)
		return day(y), endOfDay(y), true
	case "Last 7 Days":
		return now.AddDate(0, 0, -7), now, true
	case "Last 30 Days":
		return now.AddDate(0, 0, -30), now, true
	case "This Month":
		return month, month.AddDate(0, 1, 0).Add(-time.Microsecond), true
	case "Last Month":
		return month.AddDate(0, -1, 0), month.Add(-time.Microsecond), true
	case "Costum Range":
		parts := strings.Split(custom, " to ")
		if len(parts) != 2 {
			return start, end, false
		}
		from, errFrom := time.ParseInLocation("2006-01-02", strings.TrimSpace(parts[0]), now.Location())
		to, errTo := time.ParseInLocation("2006-01-02", strings.TrimSpace(parts[1]), now.Location())
		if errFrom != nil || errTo != nil {
			return start, end, false
		}
		return from, endOfDay(to), true
	}
	return start, end, false
}

// DisplayAlertDashboard replaces the direct queries in AlertDashboardRepository (list, filter,
// nomoList, crdropList, ratioalert, cappingList).
func (h *IncomingHandler) DisplayAlertDashboard(c *fiber.Ctx) error {
	sources, known := alertSourcesByType[c.Query("type", "all")]
	if !known {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "unknown alert type"})
	}

	start, end, bounded := alertDateBounds(c.Query("date-range"), c.Query("date-custom-range"), h.alertNow())

	rows := []map[string]interface{}{}
	for _, source := range sources {
		query := h.DB.Table(source.table).Select(source.columns)
		for param, column := range map[string]string{"id": "campaign_id", "country": "country", "operator": "operator", "campaign_name": "campaign_name", "service": "service"} {
			if value := c.Query(param); value != "" {
				query = query.Where(column+" = ?", value)
			}
		}
		if bounded {
			query = query.Where("summary_date BETWEEN ? AND ?", start.Format(alertTimeLayout), end.Format(alertTimeLayout))
		}
		var found []map[string]interface{}
		if err := query.Find(&found).Error; err != nil {
			h.Logs.Errorf("DisplayAlertDashboard: %s: %v", source.table, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load alerts"})
		}
		for _, row := range found {
			rows = append(rows, normalizeAlertRow(row, source.errorName))
		}
	}

	sort.SliceStable(rows, func(i, j int) bool {
		a, _ := rows[i]["time"].(string)
		b, _ := rows[j]["time"].(string)
		return a > b
	})

	return c.JSON(fiber.Map{"code": fiber.StatusOK, "data": rows})
}

// normalizeAlertRow keeps the database wall clock, because change-status matches summary_date
// exactly, and turns numeric columns (scanned as bytes) into text.
func normalizeAlertRow(row map[string]interface{}, errorName string) map[string]interface{} {
	row["error"] = errorName
	switch value := row["summary_date"].(type) {
	case time.Time:
		row["time"] = value.Format(alertTimeLayout)
	default:
		row["time"] = value
	}
	delete(row, "summary_date")
	for key, value := range row {
		if raw, isBytes := value.([]byte); isBytes {
			row[key] = string(raw)
		}
	}
	return row
}

type ChangeAlertStatusRequest struct {
	ID     string `json:"id"`
	Status *bool  `json:"status"`
	Error  string `json:"error"`
	Time   string `json:"time"`
}

// ChangeAlertStatus replaces AlertDashboardRepository::changeStatus.
func (h *IncomingHandler) ChangeAlertStatus(c *fiber.Ctx) error {
	var req ChangeAlertStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	source, known := alertSourceByError[req.Error]
	if !known || req.ID == "" || req.Time == "" || req.Status == nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "id, status, error and time are required"})
	}

	var row struct {
		ID           int
		Status       bool
		IsLock       bool
		CampaignName string
	}
	result := h.DB.Table(source.table).Select("id, status, is_lock, campaign_name").
		Where("campaign_id = ? AND summary_date = ?", req.ID, req.Time).Limit(1).Scan(&row)
	if result.Error != nil {
		h.Logs.Errorf("ChangeAlertStatus: find: %v", result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update status"})
	}
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": errors.New("alert not found for given id and time").Error()})
	}
	if row.IsLock {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "This alert is locked"})
	}
	if err := h.DB.Table(source.table).Where("id = ?", row.ID).Update("status", *req.Status).Error; err != nil {
		h.Logs.Errorf("ChangeAlertStatus: update: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update status"})
	}

	return c.JSON(fiber.Map{
		"success":       true,
		"message":       "Status updated successfully",
		"status_before": row.Status,
		"campaign_name": row.CampaignName,
	})
}
