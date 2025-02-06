package handler

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource/config"
	"github.com/infraLinkit/mediaplatform-datasource/entity"
	"github.com/infraLinkit/mediaplatform-datasource/helper"
)

const PAGESIZE int = 10

func (h *IncomingHandler) DisplayPinReport(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	m := c.Queries()

	page, _ := strconv.Atoi(m["page"])
	fe := entity.DisplayPinReport{
		Adnet:      m["adnet"],
		Country:    m["country"],
		Service:    m["service"],
		Operator:   m["operator"],
		DateRange:  m["date_range"],
		DateBefore: m["date_before"],
		DateAfter:  m["date_after"],
		Page:       page,
		Action:     m["action"],
	}

	r := h.DisplayPinReportExtra(c, fe)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayPinReportExtra(c *fiber.Ctx, fe entity.DisplayPinReport) entity.ReturnResponse {

	key := "temp_key_api_pin_report_" + strings.ReplaceAll(helper.GetIpAddress(c), ".", "_")

	var (
		err              error
		x                int
		isempty          bool
		pinreport        []entity.ApiPinReport
		displaypinreport []entity.ApiPinReport
	)

	if fe.Action != "" {
		pinreport, err = h.DS.GetApiPinReport(fe)
	} else {
		if pinreport, isempty = h.DS.RGetApiPinReport(key, "$"); isempty {

			pinreport, err = h.DS.GetApiPinReport(fe)

			s, _ := json.Marshal(pinreport)

			h.DS.SetData(key, "$", string(s))
			h.DS.SetExpireData(key, 60)
		}
	}

	if err == nil {

		pagesize := PAGESIZE
		if fe.Page >= 2 {
			x = PAGESIZE * (fe.Page - 1)
		} else if fe.Page == 1 {
			x = 0
			pagesize = pagesize - 1
		} else {
			x = 0
			pagesize = pagesize - 1
		}

		for i := x; i < len(pinreport); i++ {

			//h.Logs.Debug(fmt.Sprintf("incr : %d, ID : %d", i, pinreport[i].ID))

			displaypinreport = append(displaypinreport, pinreport[i])
			if i == pagesize {
				break
			}
		}

		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithData{
				Code:    fiber.StatusOK,
				Message: config.OK_DESC,
				Data:    displaypinreport,
			},
		}

	} else {

		return entity.ReturnResponse{
			HttpStatus: fiber.StatusNotFound,
			Rsp: entity.GlobalResponse{
				Code:    fiber.StatusNotFound,
				Message: "empty",
			},
		}
	}
}

func (h *IncomingHandler) DisplayPinPerformanceReport(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	m := c.Queries()

	page, _ := strconv.Atoi(m["page"])
	pageSize, err := strconv.Atoi(m["page_size"])
	if err != nil {
		pageSize = 10
	}

	draw, _ := strconv.Atoi(m["draw"])
	fe := entity.DisplayPinPerformanceReport{
		Adnet:      m["adnet"],
		Country:    m["country"],
		Service:    m["service"],
		Operator:   m["operator"],
		DateRange:  m["date_range"],
		DateBefore: m["date_before"],
		DateAfter:  m["date_after"],
		Page:       page,
		Action:     m["action"],
		Draw:       draw,
		PageSize:   pageSize,
	}

	r := h.DisplayPinPerformanceReportExtra(c, fe)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayPinPerformanceReportExtra(c *fiber.Ctx, fe entity.DisplayPinPerformanceReport) entity.ReturnResponse {

	key := "temp_key_api_pin_performance_report_" + strings.ReplaceAll(helper.GetIpAddress(c), ".", "_")

	var (
		err                  error
		total_data           int64
		isempty              bool
		pinperformancereport []entity.ApiPinPerformance
	)

	if fe.Action != "" {
		pinperformancereport, total_data, err = h.DS.GetApiPinPerformanceReport(fe)
	} else {
		if pinperformancereport, isempty = h.DS.RGetApiPinPerformanceReport(key, "$"); isempty {

			pinperformancereport, total_data, err = h.DS.GetApiPinPerformanceReport(fe)

			s, _ := json.Marshal(pinperformancereport)

			h.DS.SetData(key, "$", string(s))
			h.DS.SetExpireData(key, 60)
		}
	}

	if err == nil {

		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithDataTable{
				Code:            fiber.StatusOK,
				Message:         config.OK_DESC,
				Data:            pinperformancereport,
				Draw:            fe.Draw,
				RecordsTotal:    int(total_data),
				RecordsFiltered: int(total_data),
			},
		}

	} else {

		return entity.ReturnResponse{
			HttpStatus: fiber.StatusNotFound,
			Rsp: entity.GlobalResponse{
				Code:    fiber.StatusNotFound,
				Message: "empty",
			},
		}
	}
}

func (h *IncomingHandler) DisplayConversionLog(c *fiber.Ctx, fe entity.DisplayConversionLog) entity.ReturnResponse {

	// key := "temp_key_api_conversion_log_report_" + strings.ReplaceAll(helper.GetIpAddress(c), ".", "_")

	var (
		err        error
		total_data int64
		// isempty               bool
		conversion_log_report []entity.PixelStorage
	)

	// if fe.Action != "" {
	conversion_log_report, total_data, err = h.DS.GetConversionLog(fe)
	// } else {
	// 	if conversion_log_report, isempty = h.DS.RGetApiPinPerformanceReport(key, "$"); isempty {

	// 		pinperformancereport, total_data, err = h.DS.GetApiPinPerformanceReport(fe)

	// 		s, _ := json.Marshal(pinperformancereport)

	// 		h.DS.SetData(key, "$", string(s))
	// 		h.DS.SetExpireData(key, 60)

	// 	}
	// }

	if err == nil {

		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithDataTable{
				Code:            fiber.StatusOK,
				Message:         config.OK_DESC,
				Data:            conversion_log_report,
				Draw:            fe.Draw,
				RecordsTotal:    int(total_data),
				RecordsFiltered: int(total_data),
			},
		}

	} else {

		return entity.ReturnResponse{
			HttpStatus: fiber.StatusNotFound,
			Rsp: entity.GlobalResponse{
				Code:    fiber.StatusNotFound,
				Message: "empty",
			},
		}
	}
}
