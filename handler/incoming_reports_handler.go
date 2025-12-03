package handler

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource/config"
	"github.com/infraLinkit/mediaplatform-datasource/entity"
	"github.com/infraLinkit/mediaplatform-datasource/helper"
	"github.com/wiliehidayat87/rmqp"
)

const PAGESIZE int = 10

func (h *IncomingHandler) DisplayPinReport(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	m := c.Queries()

	page, _ := strconv.Atoi(m["page"])
	pageSize, err := strconv.Atoi(m["page_size"])
	if err != nil {
		pageSize = PAGESIZE
	}
	draw, _ := strconv.Atoi(m["draw"])
	var adnets []string
	for k, v := range m {
		if strings.HasPrefix(k, "adnet[") {
			adnets = append(adnets, v)
		}
	}
	fe := entity.DisplayPinReport{
		DateSend:    time.Time{},
		CampaignId:  m["campaign_id"],
		Country:     m["country"],
		Company:     m["company"],
		Operator:    m["operator"],
		Partner:     m["partner"],
		Aggregator:  m["aggregator"],
		Adnets:      adnets,
		Service:     m["service"],
		Draw:        draw,
		Page:        page,
		PageSize:    pageSize,
		Action:      m["action"],
		DateRange:   m["date_range"],
		DateBefore:  m["date_before"],
		DateAfter:   m["date_after"],
		Reload:      m["reload"],
		OrderColumn: m["order_column"],
		OrderDir:    m["order_dir"],
	}

	allowedCompanies, _ := c.Locals("companies").([]string)

	r := h.DisplayPinReportExtra(c, fe, allowedCompanies)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayPinReportExtra(c *fiber.Ctx, fe entity.DisplayPinReport, allowedCompanies []string) entity.ReturnResponse {
	var (
		err        error
		total_data int64
		apireport  []entity.ApiPinReport
	)

	if fe.Action != "" || fe.Reload == "true" {
		fmt.Println("-----", fe.Reload, "-----")
		apireport, total_data, err = h.DS.GetDisplayPinReport(fe, allowedCompanies)
	} else {

		apireport, total_data, err = h.DS.GetDisplayPinReport(fe, allowedCompanies)
	}

	if err == nil {

		if apireport == nil {
			apireport = []entity.ApiPinReport{}
		}

		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithDataTable{
				Draw:            fe.Draw,
				Code:            fiber.StatusOK,
				Message:         config.OK_DESC,
				Data:            apireport,
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

func (h *IncomingHandler) EditPOAFAPIReport(c *fiber.Ctx) error {

	o := new(entity.ApiPinReport)

	if err := c.BodyParser(&o); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	} else {

		h.Logs.Debug(fmt.Sprintf("data : %#v ...", o))

		h.DS.EditPOAFAPIReport(entity.ApiPinReport{
			DateSend:   o.DateSend,
			PayoutAF:   o.PayoutAF,
			CampaignId: o.CampaignId,
		})

		return c.Status(fiber.StatusOK).Send([]byte("OK"))
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

	page, errPage := strconv.Atoi(m["page"])
	pageSize, err := strconv.Atoi(m["page_size"])
	if err != nil {
		pageSize = 10
	}
	if errPage != nil {
		page = 1
	}

	draw, _ := strconv.Atoi(m["draw"])
	fe := entity.DisplayConversionLogReport{
		Adnet:          m["adnet"],
		Agency:         m["agency"],
		Country:        m["country"],
		Operator:       m["operator"],
		Pixel:          m["pixel"],
		CampaignType:   m["campaign_type"],
		StatusPostback: m["status_postback"],
		CampaignId:     m["campaign_id"],
		CampaignName:   m["campaign_name"],
		DateRange:      m["date_range"],
		DateStart:      m["date_start"],
		DateEnd:        m["date_end"],
		Page:           page,
		Action:         m["action"],
		Draw:           draw,
		PageSize:       pageSize,
		Order:          m["order"],
	}

	r := h.DisplayConversionLogReportExtra(c, fe)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayPerformanceReport(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	m := c.Queries()

	page, _ := strconv.Atoi(m["page"])
	pageSize, errRequest := strconv.Atoi(m["page_size"])
	if errRequest != nil {
		pageSize = 10
	}
	draw, _ := strconv.Atoi(m["draw"])
	params := entity.PerformaceReportParams{
		Page:         page,
		Action:       m["action"],
		Draw:         draw,
		PageSize:     pageSize,
		Country:      m["country"],
		Company:      m["company"],
		ClientType:   m["client_type"],
		Operator:     m["operator"],
		CampaignName: m["campaign_name"],
		CampaignType: m["campaign_type"],
		Publisher:    m["publisher"],
		Service:      m["service"],
		DateStart:    m["date_start"],
		DateEnd:      m["date_end"],
	}

	var (
		errResponse             error
		total_data              int64
		performance_report_list []entity.PerformanceReport
	)

	// key := "temp_key_api_company_" + strings.ReplaceAll(helper.GetIpAddress(c), ".", "_")

	// need to add redis mechanism here

	performance_report_list, total_data, errResponse = h.DS.GetPerformanceReport(params)

	r := entity.ReturnResponse{
		HttpStatus: fiber.StatusNotFound,
		Rsp: entity.GlobalResponse{
			Code:    fiber.StatusNotFound,
			Message: "empty",
		},
	}

	if errResponse == nil {
		r = entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithDataTable{
				Code:            fiber.StatusOK,
				Message:         config.OK_DESC,
				Data:            performance_report_list,
				Draw:            params.Draw,
				RecordsTotal:    int(total_data),
				RecordsFiltered: int(total_data),
			},
		}

	}

	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayConversionLogReportExtra(c *fiber.Ctx, fe entity.DisplayConversionLogReport) entity.ReturnResponse {

	// key := "temp_key_api_conversion_log_report_" + strings.ReplaceAll(helper.GetIpAddress(c), ".", "_")

	var (
		err        error
		total_data int64
		// isempty               bool
		conversion_log_report []entity.PixelStorage
	)

	// if fe.Action != "" {
	conversion_log_report, total_data, err = h.DS.GetConversionLogReport(fe)
	// } else {
	// 	if conversion_log_report, isempty = h.DS.RGetConversionLogReport(key, "$"); isempty {

	// 		conversion_log_report, total_data, err = h.DS.GetConversionLogReport(fe)

	// 		s, _ := json.Marshal(conversion_log_report)

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

func (h *IncomingHandler) DisplayCPAReport(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	m := c.Queries()

	page, _ := strconv.Atoi(m["page"])
	pageSize, err := strconv.Atoi(m["page_size"])
	if err != nil {
		pageSize = PAGESIZE
	}
	draw, _ := strconv.Atoi(m["draw"])
	var adnets []string
	for k, v := range m {
		if strings.HasPrefix(k, "adnet[") {
			adnets = append(adnets, v)
		}
	}
	fe := entity.DisplayCPAReport{
		SummaryDate:  time.Time{},
		CampaignId:   m["campaign_id"],
		CampaignName: m["campaign_name"],
		Country:      m["country"],
		ClientType:   m["client_type"],
		Company:      m["company"],
		Operator:     m["operator"],
		Partner:      m["partner"],
		Aggregator:   m["aggregator"],
		Adnets:       adnets,
		Service:      m["service"],
		Draw:         draw,
		Page:         page,
		PageSize:     pageSize,
		Action:       m["action"],
		DateRange:    m["date_range"],
		DateBefore:   m["date_before"],
		DateAfter:    m["date_after"],
		Reload:       m["reload"],
		OrderColumn:  m["order_column"],
		OrderDir:     m["order_dir"],
	}

	allowedCompanies, _ := c.Locals("companies").([]string)

	r := h.DisplayCPAReportExtra(c, fe, allowedCompanies)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayCPAReportExtra(c *fiber.Ctx, fe entity.DisplayCPAReport, allowedCompanies []string) entity.ReturnResponse {
	// key := "temp_key_api_cpa_report_" + strings.ReplaceAll(helper.GetIpAddress(c), ".", "_")

	var (
		err        error
		total_data int64
		// isempty    bool
		cpareport            []entity.SummaryCampaign
		TotalSummaryCampaign entity.TotalSummaryCampaign
		// displaycpareport []entity.SummaryCampaign
	)

	if fe.Action != "" || fe.Reload == "true" {
		cpareport, total_data, TotalSummaryCampaign, err = h.DS.GetDisplayCPAReport(fe, allowedCompanies)
	} else {
		cpareport, total_data, TotalSummaryCampaign, err = h.DS.GetDisplayCPAReport(fe, allowedCompanies)
	}

	if err == nil {

		if cpareport == nil {
			cpareport = []entity.SummaryCampaign{}
		}

		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithDataTable{
				Draw:            fe.Draw,
				Code:            fiber.StatusOK,
				Message:         config.OK_DESC,
				Data:            cpareport,
				TotalSummary:    TotalSummaryCampaign,
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

func (h *IncomingHandler) ExportCpaButton(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	m := c.Queries()

	page, _ := strconv.Atoi(m["page"])
	var adnets []string
	for k, v := range m {
		if strings.HasPrefix(k, "adnet[") {
			adnets = append(adnets, v)
		}
	}
	fe := entity.DisplayCPAReport{
		SummaryDate:  time.Time{},
		CampaignId:   m["campaign_id"],
		CampaignName: m["campaign_name"],
		Country:      m["country"],
		Operator:     m["operator"],
		Partner:      m["partner"],
		Aggregator:   m["aggregator"],
		Adnets:       adnets,
		Service:      m["service"],
		Page:         page,
		Action:       m["action"],
		DateRange:    m["date_range"],
		DateBefore:   m["date_before"],
		DateAfter:    m["date_after"],
		OrderColumn:  m["order_column"],
		OrderDir:     m["order_dir"],
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
		err                  error
		cpareport            []entity.SummaryCampaign
		isempty              bool
		TotalSummaryCampaign entity.TotalSummaryCampaign
		// displaycpareport []entity.SummaryCampaign
	)

	allowedCompanies, _ := c.Locals("companies").([]string)

	if fe.Action != "" {
		cpareport, _, TotalSummaryCampaign, err = h.DS.GetDisplayCPAReport(fe, allowedCompanies)
	} else {

		if cpareport, isempty = h.DS.RGetDisplayCPAReport(key, "$"); isempty {

			cpareport, _, TotalSummaryCampaign, err = h.DS.GetDisplayCPAReport(fe, allowedCompanies)

			s, _ := json.Marshal(cpareport)

			h.DS.SetData(key, "$", string(s))
			h.DS.SetExpireData(key, 60)
		}
	}

	if err == nil {
		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithData{
				Code:         fiber.StatusOK,
				Message:      config.OK_DESC,
				Data:         cpareport,
				TotalSummary: TotalSummaryCampaign,
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

	page, errPage := strconv.Atoi(m["page"])
	pageSize, err := strconv.Atoi(m["page_size"])
	if err != nil {
		pageSize = 10
	}
	if errPage != nil {
		page = 10
	}

	draw, _ := strconv.Atoi(m["draw"])
	v := c.Params("v")

	fe := entity.DisplayCostReport{
		Adnet:        m["adnet"],
		Country:      m["country"],
		Operator:     m["operator"],
		CampaignType: m["campaign_type"],
		Page:         page,
		Action:       m["action"],
		DateRange:    m["date_range"],
		DateBefore:   m["date_before"],
		DateAfter:    m["date_after"],
		DataBasedOn:  m["data_based_on"],
		PageSize:     pageSize,
		Draw:         draw,
	}

	allowedAdnets, _ := c.Locals("adnets").([]string)

	r := h.DisplayCostReportExtra(c, fe, v, allowedAdnets)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayCostReportExtra(c *fiber.Ctx, fe entity.DisplayCostReport, v string, allowedAdnets []string) entity.ReturnResponse {
	key := "temp_key_api_cost_report_" + strings.ReplaceAll(helper.GetIpAddress(c), ".", "_")
	keydetail := "temp_key_api_cost_report_detail_" + strings.ReplaceAll(helper.GetIpAddress(c), ".", "_")

	var (
		err        error
		isempty    bool
		total_data int64
		costreport []entity.CostReport
		// displaycostreport []entity.CostReport
	)
	if v == "list" {
		if fe.Action != "" {
			costreport, total_data, err = h.DS.GetDisplayCostReport(fe, allowedAdnets)
		} else {
			if costreport, isempty = h.DS.RGetDisplayCostReport(key, "$"); isempty {
				costreport, total_data, err = h.DS.GetDisplayCostReport(fe, allowedAdnets)
				s, _ := json.Marshal(costreport)
				h.DS.SetData(key, "$", string(s))
				h.DS.SetExpireData(key, 60)
			}
		}
	} else if v == "detail" {
		if fe.Action != "" {
			costreport, total_data, err = h.DS.GetDisplayCostReportDetail(fe)
		} else {
			if costreport, isempty = h.DS.RGetDisplayCostReportDetail(keydetail, "$"); isempty {
				costreport, total_data, err = h.DS.GetDisplayCostReportDetail(fe)
				s, _ := json.Marshal(costreport)
				h.DS.SetData(key, "$", string(s))
				h.DS.SetExpireData(key, 60)
			}
		}
	}

	if err == nil {
		// pagesize := PAGESIZE
		// if fe.Page >= 2 {
		// 	x = pagesize * (fe.Page - 1)
		// } else {
		// 	x = 0
		// }

		// for i := x; i < len(costreport) && i < x+pagesize; i++ {
		// 	displaycostreport = append(displaycostreport, costreport[i])
		// }

		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithDataTable{
				Draw:            fe.Draw,
				Code:            fiber.StatusOK,
				Message:         config.OK_DESC,
				Data:            costreport,
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

		allowedAdnets, _ := c.Locals("adnets").([]string)

		r := h.ExportCostReportExtraNoLimit(c, fe, allowedAdnets)
		return c.Status(r.HttpStatus).JSON(r.Rsp)
	}

	return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{
		Code:    fiber.StatusBadRequest,
		Message: config.BAD_REQUEST_DESC,
	})
}

func (h *IncomingHandler) ExportCostReportExtraNoLimit(c *fiber.Ctx, fe entity.DisplayCostReport, allowedAdnets []string) entity.ReturnResponse {
	key := "temp_key_api_cost_report_" + strings.ReplaceAll(helper.GetIpAddress(c), ".", "_")

	var (
		err        error
		costreport []entity.CostReport
		isempty    bool
		total_data int64
		// displaycpareport []entity.SummaryCampaign
	)

	if fe.Action != "" {
		costreport, total_data, err = h.DS.GetDisplayCostReport(fe, allowedAdnets)
	} else {
		if costreport, isempty = h.DS.RGetDisplayCostReport(key, "$"); isempty {
			costreport, total_data, err = h.DS.GetDisplayCostReport(fe, allowedAdnets)
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
		total_data int64
		// displaycpareport []entity.SummaryCampaign
	)

	if fe.Action != "" {
		costreport, total_data, err = h.DS.GetDisplayCostReportDetail(fe)
	} else {
		if costreport, isempty = h.DS.RGetDisplayCostReportDetail(key, "$"); isempty {
			costreport, total_data, err = h.DS.GetDisplayCostReportDetail(fe)
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

func (h *IncomingHandler) DisplayDefaultInput(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	gs, err := h.DS.GetDataConfig("global_setting", "$")
	if err != nil {
		h.Logs.Error(fmt.Sprintf("Failed to get global settings: %v", err))
		return c.Status(fiber.StatusInternalServerError).JSON(entity.GlobalResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get global settings",
		})
	}

	// Convert string values to float64
	costPerConversion, _ := strconv.ParseFloat(gs.CPCR, 64)
	agencyFee, _ := strconv.ParseFloat(gs.AgencyFee, 64)
	technicalFee, _ := strconv.ParseFloat(gs.TechnicalFee, 64)
	targetDailyBudget, _ := strconv.ParseFloat(gs.TargetDailyBudget, 64)

	// Validasi jika error, jadi default agency 5, cost 0.06, tech fee 5
	if err != nil {
		costPerConversion = 0.06
		agencyFee = 5
		technicalFee = 5
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithData{
		Code:    fiber.StatusOK,
		Message: config.OK_DESC,
		Data: entity.DefaultInput{
			CostPerConversion: costPerConversion,
			AgencyFee:         agencyFee,
			TechnicalFee:      technicalFee,
			TargetDailyBudget: targetDailyBudget,
		},
	})
}

func (h *IncomingHandler) ResendData(c *fiber.Ctx) error {

	var ids []string
	total, _ := strconv.Atoi(c.FormValue("total"))

	for i := 0; i < total; i++ {
		ids = append(ids, c.FormValue("id["+strconv.Itoa(i)+"]"))
		//fmt.Printf("Current iteration: %d\n", i)
	}

	baseURL := h.Config.APILINKITDashboard

	reports, _ := h.DS.GetSummaryReportById(ids)
	errorCounter := 0

	for _, sc := range reports {
		q := url.Values{
			"date":           {sc.SummaryDate.Format("2006-10-02")},
			"campaign_id":    {sc.URLServiceKey},
			"publisher":      {sc.Adnet},
			"adnet":          {sc.Adnet},
			"operator":       {sc.Partner},
			"adn":            {sc.ShortCode},
			"client":         {sc.Partner},
			"aggregator":     {sc.Aggregator},
			"country":        {sc.Country},
			"service":        {sc.Service},
			"mo_received":    {strconv.Itoa(sc.MoReceived)},
			"mo_postback":    {strconv.Itoa(sc.Postback)},
			"total_mo":       {strconv.Itoa(sc.MoReceived)},
			"total_postback": {strconv.Itoa(sc.Postback)},
			"landing":        {strconv.Itoa(sc.Traffic)},
			"cr_mo_received": {strconv.FormatFloat(sc.CrMO, 'f', 2, 64)},
			"cr_mo_postback": {strconv.FormatFloat(sc.CrPostback, 'f', 2, 64)},
			"url_campaign":   {sc.URLAfter},
			"url_service":    {sc.URLBefore},
			"sbaf":           {strconv.FormatFloat(sc.SBAF, 'f', 2, 64)},
			"saaf":           {strconv.FormatFloat(sc.SAAF, 'f', 2, 64)},
			"spending":       {strconv.FormatFloat(sc.SAAF, 'f', 2, 64)},
			"campaign":       {sc.CampaignObjective},
			"payout":         {strconv.FormatFloat(sc.PO, 'f', 2, 64)},
			"price_per_mo":   {strconv.FormatFloat(sc.PricePerMO, 'f', 2, 64)},
		}

		fullURL := fmt.Sprintf("%s?%s", baseURL, q.Encode())
		message := `{"url":"` + fullURL + `"}`
		published := h.Rmqp.PublishMsg(rmqp.PublishItems{
			ExchangeName: "E_RESENDCAMPAIGNDATA",
			QueueName:    "Q_RESENDCAMPAIGNDATA",
			ContentType:  "application/json",
			Payload:      message, // Send the properly formatted JSON
			Priority:     0,
		})

		if !published {
			errorCounter++
			h.Logs.Debug(fmt.Sprintf("[x] Failed published: Data: %s ...", message))
			//fmt.Println(fmt.Sprintf("[x] Failed published: Data: %s ...", message))
		} else {
			h.Logs.Debug(fmt.Sprintf("[v] Published: Data: %s ...", message))
			//fmt.Println(fmt.Sprintf("[v] Published: Data: %s ...", message))
		}
		//fmt.Println(message)

	}

	if errorCounter > 0 {
		return c.Status(fiber.StatusOK).SendString(`{"status":"NOK","error":"Some data not published"}`)
	}

	return c.Status(fiber.StatusOK).SendString(`{"status":"OK","error":""}`)
}
