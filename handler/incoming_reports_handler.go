package handler

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

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
			x = pagesize * (fe.Page - 1)
		} else {
			x = 0
		}

		for i := x; i < len(pinreport) && i < x+pagesize; i++ {

			// h.Logs.Debug(fmt.Sprintf("incr : %d, ID : %d", i, pinreport[i].ID))

			displaypinreport = append(displaypinreport, pinreport[i])
		}

		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithDataTable{
				Code:    fiber.StatusOK,
				Message: config.OK_DESC,
				Data:    displaypinreport,
				Page:    fe.Page,
				Total:   len(pinreport),
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

func (h *IncomingHandler) DisplayConversionLogReport(c *fiber.Ctx) error {

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
	fe := entity.DisplayConversionLogReport{
		Adnet:      m["adnet"],
		Country:    m["country"],
		Operator:   m["operator"],
		Pixel:      m["pixel"],
		CampaignId: m["campaign_id"],
		DateRange:  m["date_range"],
		DateBefore: m["date_before"],
		DateAfter:  m["date_after"],
		Page:       page,
		Action:     m["action"],
		Draw:       draw,
		PageSize:   pageSize,
	}

	r := h.DisplayConversionLogReportExtra(c, fe)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayConversionLogReportExtra(c *fiber.Ctx, fe entity.DisplayConversionLogReport) entity.ReturnResponse {

	key := "temp_key_api_conversion_log_report_" + strings.ReplaceAll(helper.GetIpAddress(c), ".", "_")

	var (
		err                   error
		total_data            int64
		isempty               bool
		conversion_log_report []entity.PixelStorage
	)

	if fe.Action != "" {
		conversion_log_report, total_data, err = h.DS.GetConversionLogReport(fe)
	} else {
		if conversion_log_report, isempty = h.DS.RGetConversionLogReport(key, "$"); isempty {

			conversion_log_report, total_data, err = h.DS.GetConversionLogReport(fe)

			s, _ := json.Marshal(conversion_log_report)

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

func (h *IncomingHandler) DisplayCPAReport(c *fiber.Ctx) error { //dev-cpa

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	m := c.Queries()

	page, _ := strconv.Atoi(m["page"])
	fe := entity.DisplayCPAReport{
		SummaryDate:   time.Time{},
		URLServiceKey: m["urlservicekey"],
		CampaignId:    m["campaign_id"],
		CampaignName:  m["campaign_name"],
		Country:       m["country"],
		Operator:      m["operator"],
		Partner:       m["partner"],
		Aggregator:    m["aggregator"],
		Adnet:         m["adnet"],
		Service:       m["service"],
		ShortCode:     m["short_code"],
		// Traffic:                  0,
		// Landing:                  0,
		// MoReceived:               0,
		// CR:                       0,
		// Postback:                 0,
		// TotalFP:                  0,
		// SuccessFP:                0,
		// Billrate:                 0,
		// ROI:                      0,
		// PO:                       0,
		// Cost:                     0,
		// SBAF:                     0,
		// SAAF:                     0,
		// CPA:                      0,
		// Revenue:                  0,
		// URLAfter:                 "NA",
		// URLBefore:                "NA",
		// MOLimit:                  0,
		// RatioSend:                1,
		// RatioReceive:             4,
		// Company:                  "NA",
		// ClientType:               "NA",
		// CostPerConversion:        0,
		// AgencyFee:                0,
		// TargetDailyBudget:        0,
		// CrMO:                     0,
		// CrPostback:               0,
		// TotalWakiAgencyFee:       0,
		// BudgetUsage:              0,
		// TargetDailyBudgetChanges: 0,
		Page:       page,
		Action:     m["action"],
		DateRange:  m["date_range"],
		DateBefore: m["date_before"],
		DateAfter:  m["date_after"],
	}

	r := h.DisplayCPAReportExtra(c, fe)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayCPAReportExtra(c *fiber.Ctx, fe entity.DisplayCPAReport) entity.ReturnResponse {
	key := "temp_key_api_cpa_report_" + strings.ReplaceAll(helper.GetIpAddress(c), ".", "_")

	var (
		err              error
		x                int
		isempty          bool
		cpareport        []entity.SummaryCampaign
		displaycpareport []entity.SummaryCampaign
	)

	if fe.Action != "" {
		cpareport, err = h.DS.GetDisplayCPAReport(fe)
	} else {

		if cpareport, isempty = h.DS.RGetDisplayCPAReport(key, "$"); isempty {

			cpareport, err = h.DS.GetDisplayCPAReport(fe)

			s, _ := json.Marshal(cpareport)

			h.DS.SetData(key, "$", string(s))
			h.DS.SetExpireData(key, 60)
		}
	}

	if err == nil {

		pagesize := PAGESIZE
		if fe.Page >= 2 {
			x = pagesize * (fe.Page - 1)
		} else {
			x = 0
		}

		for i := x; i < len(cpareport) && i < x+pagesize; i++ {

			// h.Logs.Debug(fmt.Sprintf("incr : %d, ID : %d", i, cpareport[i].ID))

			displaycpareport = append(displaycpareport, cpareport[i])
		}

		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithData{
				Code:    fiber.StatusOK,
				Message: config.OK_DESC,
				Data:    displaycpareport,
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

func (h *IncomingHandler) ExportCpaButton(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	m := c.Queries()

	page, _ := strconv.Atoi(m["page"])
	fe := entity.DisplayCPAReport{
		SummaryDate:              time.Time{},
		URLServiceKey:            m["urlservicekey"],
		CampaignId:               m["campaign_id"],
		CampaignName:             m["campaign_name"],
		Country:                  m["country"],
		Operator:                 m["operator"],
		Partner:                  m["partner"],
		Aggregator:               m["aggregator"],
		Adnet:                    m["adnet"],
		Service:                  m["service"],
		ShortCode:                m["short_code"],
		Traffic:                  0,
		Landing:                  0,
		MoReceived:               0,
		CR:                       0,
		Postback:                 0,
		TotalFP:                  0,
		SuccessFP:                0,
		Billrate:                 0,
		ROI:                      0,
		PO:                       0,
		Cost:                     0,
		SBAF:                     0,
		SAAF:                     0,
		CPA:                      0,
		Revenue:                  0,
		URLAfter:                 "NA",
		URLBefore:                "NA",
		MOLimit:                  0,
		RatioSend:                1,
		RatioReceive:             4,
		Company:                  "NA",
		ClientType:               "NA",
		CostPerConversion:        0,
		AgencyFee:                0,
		TargetDailyBudget:        0,
		CrMO:                     0,
		CrPostback:               0,
		TotalWakiAgencyFee:       0,
		BudgetUsage:              0,
		TargetDailyBudgetChanges: 0,
		Page:                     page,
		Action:                   m["action"],
		DateRange:                m["date_range"],
		DateBefore:               m["date_before"],
		DateAfter:                m["date_after"],
	}

	export_cpa := m["export_cpa"]

	if export_cpa == "true" {

		r := h.ExportCpaReportExtraNoLimit(c, fe)
		return c.Status(r.HttpStatus).JSON(r.Rsp)
	}

	return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{
		Code:    fiber.StatusBadRequest,
		Message: config.BAD_REQUEST_DESC,
	})
}

func (h *IncomingHandler) ExportCpaReportExtraNoLimit(c *fiber.Ctx, fe entity.DisplayCPAReport) entity.ReturnResponse {
	key := "temp_key_api_cpa_report_" + strings.ReplaceAll(helper.GetIpAddress(c), ".", "_")

	var (
		err       error
		cpareport []entity.SummaryCampaign
		isempty   bool
		// displaycpareport []entity.SummaryCampaign
	)

	if fe.Action != "" {
		cpareport, err = h.DS.GetDisplayCPAReport(fe)
	} else {

		if cpareport, isempty = h.DS.RGetDisplayCPAReport(key, "$"); isempty {

			cpareport, err = h.DS.GetDisplayCPAReport(fe)

			s, _ := json.Marshal(cpareport)

			h.DS.SetData(key, "$", string(s))
			h.DS.SetExpireData(key, 60)
		}
	}

	if err == nil {

		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithData{
				Code:    fiber.StatusOK,
				Message: config.OK_DESC,
				Data:    cpareport,
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

func (h *IncomingHandler) DisplayCostReport(c *fiber.Ctx) error {
	c.Set("Content-type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	m := c.Queries()

	page, _ := strconv.Atoi(m["page"])
	draw, _ := strconv.Atoi(m["draw"])
	v := c.Params("v")

	fe := entity.DisplayCostReport{
		Adnet:       m["adnet"],
		Country:     m["country"],
		Operator:    m["operator"],
		Page:        page,
		Action:      m["action"],
		DateRange:   m["date_range"],
		DateBefore:  m["date_before"],
		DateAfter:   m["date_after"],
		DataBasedOn: m["data_based_on"],
		Draw:        draw,
	}

	r := h.DisplayCostReportExtra(c, fe, v)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayCostReportExtra(c *fiber.Ctx, fe entity.DisplayCostReport, v string) entity.ReturnResponse {
	key := "temp_key_api_cost_report_" + strings.ReplaceAll(helper.GetIpAddress(c), ".", "_")
	keydetail := "temp_key_api_cost_report_detail_" + strings.ReplaceAll(helper.GetIpAddress(c), ".", "_")

	var (
		err               error
		isempty           bool
		x                 int
		costreport        []entity.CostReport
		displaycostreport []entity.CostReport
	)
	if v == "list" {
		if fe.Action != "" {
			costreport, err = h.DS.GetDisplayCostReport(fe)
		} else {
			if costreport, isempty = h.DS.RGetDisplayCostReport(key, "$"); isempty {
				costreport, err = h.DS.GetDisplayCostReport(fe)
				s, _ := json.Marshal(costreport)
				h.DS.SetData(key, "$", string(s))
				h.DS.SetExpireData(key, 60)
			}
		}
	} else if v == "detail" {
		if fe.Action != "" {
			costreport, err = h.DS.GetDisplayCostReportDetail(fe)
		} else {
			if costreport, isempty = h.DS.RGetDisplayCostReportDetail(keydetail, "$"); isempty {
				costreport, err = h.DS.GetDisplayCostReportDetail(fe)
				s, _ := json.Marshal(costreport)
				h.DS.SetData(key, "$", string(s))
				h.DS.SetExpireData(key, 60)
			}
		}
	}

	if err == nil {
		pagesize := PAGESIZE
		if fe.Page >= 2 {
			x = pagesize * (fe.Page - 1)
		} else {
			x = 0
		}

		for i := x; i < len(costreport) && i < x+pagesize; i++ {
			displaycostreport = append(displaycostreport, costreport[i])
		}

		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithTable{
				Code:            fiber.StatusOK,
				Message:         config.OK_DESC,
				Data:            displaycostreport,
				Draw:            fe.Draw,
				RecordsTotal:    len(costreport),
				RecordsFiltered: len(costreport),
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

func (h *IncomingHandler) ExportCostButton(c *fiber.Ctx) error {
	c.Set("Content-type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	m := c.Queries()

	page, _ := strconv.Atoi(m["page"])
	draw, _ := strconv.Atoi(m["draw"])
	fe := entity.DisplayCostReport{
		Adnet:       m["adnet"],
		Country:     m["country"],
		Operator:    m["operator"],
		Page:        page,
		Action:      m["action"],
		DateRange:   m["date_range"],
		DateBefore:  m["date_before"],
		DateAfter:   m["date_after"],
		DataBasedOn: m["data_based_on"],
		Draw:        draw,
	}
	export_cost := m["export_cost"]
	if export_cost == "true" {

		r := h.ExportCostReportExtraNoLimit(c, fe)
		return c.Status(r.HttpStatus).JSON(r.Rsp)
	}

	return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{
		Code:    fiber.StatusBadRequest,
		Message: config.BAD_REQUEST_DESC,
	})
}

func (h *IncomingHandler) ExportCostReportExtraNoLimit(c *fiber.Ctx, fe entity.DisplayCostReport) entity.ReturnResponse {
	key := "temp_key_api_cost_report_" + strings.ReplaceAll(helper.GetIpAddress(c), ".", "_")

	var (
		err        error
		costreport []entity.CostReport
		isempty    bool
		// displaycpareport []entity.SummaryCampaign
	)

	if fe.Action != "" {
		costreport, err = h.DS.GetDisplayCostReport(fe)
	} else {
		if costreport, isempty = h.DS.RGetDisplayCostReport(key, "$"); isempty {
			costreport, err = h.DS.GetDisplayCostReport(fe)
			s, _ := json.Marshal(costreport)
			h.DS.SetData(key, "$", string(s))
			h.DS.SetExpireData(key, 60)
		}
	}

	if err == nil {
		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithTable{
				Code:            fiber.StatusOK,
				Message:         config.OK_DESC,
				Data:            costreport,
				Draw:            fe.Draw,
				RecordsTotal:    len(costreport),
				RecordsFiltered: len(costreport),
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
func (h *IncomingHandler) ExportCostDetailButton(c *fiber.Ctx) error {
	c.Set("Content-type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	m := c.Queries()

	page, _ := strconv.Atoi(m["page"])
	draw, _ := strconv.Atoi(m["draw"])
	fe := entity.DisplayCostReport{
		Adnet:       m["adnet"],
		Country:     m["country"],
		Operator:    m["operator"],
		Page:        page,
		Action:      m["action"],
		DateRange:   m["date_range"],
		DateBefore:  m["date_before"],
		DateAfter:   m["date_after"],
		DataBasedOn: m["data_based_on"],
		Draw:        draw,
	}
	export_cost := m["export_cost"]
	if export_cost == "true" {

		r := h.ExportCostReportDetailExtraNoLimit(c, fe)
		return c.Status(r.HttpStatus).JSON(r.Rsp)
	}

	return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{
		Code:    fiber.StatusBadRequest,
		Message: config.BAD_REQUEST_DESC,
	})
}

func (h *IncomingHandler) ExportCostReportDetailExtraNoLimit(c *fiber.Ctx, fe entity.DisplayCostReport) entity.ReturnResponse {
	key := "temp_key_api_cost_report_detail_" + strings.ReplaceAll(helper.GetIpAddress(c), ".", "_")

	var (
		err        error
		costreport []entity.CostReport
		isempty    bool
		// displaycpareport []entity.SummaryCampaign
	)

	if fe.Action != "" {
		costreport, err = h.DS.GetDisplayCostReportDetail(fe)
	} else {
		if costreport, isempty = h.DS.RGetDisplayCostReportDetail(key, "$"); isempty {
			costreport, err = h.DS.GetDisplayCostReportDetail(fe)
			s, _ := json.Marshal(costreport)
			h.DS.SetData(key, "$", string(s))
			h.DS.SetExpireData(key, 60)
		}
	}

	if err == nil {
		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithTable{
				Code:            fiber.StatusOK,
				Message:         config.OK_DESC,
				Data:            costreport,
				Draw:            fe.Draw,
				RecordsTotal:    len(costreport),
				RecordsFiltered: len(costreport),
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
