package handler

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource/config"
	"github.com/infraLinkit/mediaplatform-datasource/entity"
)


func (h *IncomingHandler) DisplayBudgetIO(c *fiber.Ctx) error {

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
	
	fe := entity.DisplayBudgetIO{
		Continent:    m["continent"],
		Country:      m["country"],
		CampaignType: m["campaign_type"],
		Company:      m["company"],
		Partner:      m["partner"],
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

	r := h.DisplayBudgetIOExtra(c, fe, allowedCompanies)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayBudgetIOExtra(c *fiber.Ctx, fe entity.DisplayBudgetIO, allowedCompanies []string) entity.ReturnResponse {
	// key := "temp_key_api_cpa_report_" + strings.ReplaceAll(helper.GetIpAddress(c), ".", "_")

	var (
		err        error
		total_data int64
		// isempty    bool
		budgetio []entity.BudgetIO
		// displaybudgetio []entity.BudgetIO
	)

	if fe.Action != "" || fe.Reload == "true" {
		budgetio, total_data, err = h.DS.GetDisplayBudgetIO(fe, allowedCompanies)
	} else {
		budgetio, total_data, err = h.DS.GetDisplayBudgetIO(fe, allowedCompanies)
	}

	if err == nil {

		if budgetio == nil {
			budgetio = []entity.BudgetIO{}
		}

		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithDataTable{
				Draw:            fe.Draw,
				Code:            fiber.StatusOK,
				Message:         config.OK_DESC,
				Data:            budgetio,
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

func (h *IncomingHandler) DisplayBudgetIOAll(c *fiber.Ctx) error {

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
	
	fe := entity.DisplayBudgetIO{
		Continent:    m["continent"],
		Country:      m["country"],
		CampaignType: m["campaign_type"],
		Company:      m["company"],
		Partner:      m["partner"],
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


	r := h.DisplayBudgetIOExtraAll(c, fe)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayBudgetIOExtraAll(c *fiber.Ctx, fe entity.DisplayBudgetIO) entity.ReturnResponse {
	// key := "temp_key_api_cpa_report_" + strings.ReplaceAll(helper.GetIpAddress(c), ".", "_")

	var (
		err        error
		total_data int64
		// isempty    bool
		budgetio []entity.BudgetIO
		// displaybudgetio []entity.BudgetIO
	)

	if fe.Action != "" || fe.Reload == "true" {
		budgetio, total_data, err = h.DS.GetDisplayBudgetIOAll(fe)
	} else {
		budgetio, total_data, err = h.DS.GetDisplayBudgetIOAll(fe)
	}

	if err == nil {

		if budgetio == nil {
			budgetio = []entity.BudgetIO{}
		}

		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithDataTable{
				Draw:            fe.Draw,
				Code:            fiber.StatusOK,
				Message:         config.OK_DESC,
				Data:            budgetio,
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

func (h *IncomingHandler) DisplayBudgetIOApproved(c *fiber.Ctx) error {

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
	
	fe := entity.DisplayBudgetIO{
		Continent:    m["continent"],
		Country:      m["country"],
		CampaignType: m["campaign_type"],
		Company:      m["company"],
		Partner:      m["partner"],
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

	r := h.DisplayBudgetIOExtra(c, fe, allowedCompanies)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayBudgetIOApprovedExtra(c *fiber.Ctx, fe entity.DisplayBudgetIO, allowedCompanies []string) entity.ReturnResponse {
	// key := "temp_key_api_cpa_report_" + strings.ReplaceAll(helper.GetIpAddress(c), ".", "_")

	var (
		err        error
		total_data int64
		// isempty    bool
		budgetio []entity.BudgetIO
		// displaybudgetio []entity.BudgetIO
	)

	if fe.Action != "" || fe.Reload == "true" {
		budgetio, total_data, err = h.DS.GetDisplayBudgetIOApproved(fe, allowedCompanies)
	} else {
		budgetio, total_data, err = h.DS.GetDisplayBudgetIOApproved(fe, allowedCompanies)
	}

	if err == nil {

		if budgetio == nil {
			budgetio = []entity.BudgetIO{}
		}

		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithDataTable{
				Draw:            fe.Draw,
				Code:            fiber.StatusOK,
				Message:         config.OK_DESC,
				Data:            budgetio,
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

func (h *IncomingHandler) DisplayBudgetIOApprovedAll(c *fiber.Ctx) error {

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
	
	fe := entity.DisplayBudgetIO{
		Continent:    m["continent"],
		Country:      m["country"],
		CampaignType: m["campaign_type"],
		Company:      m["company"],
		Partner:      m["partner"],
		Draw:         draw,
		Page:         page,
		PageSize:     pageSize,
		Keyword: 	  m["keyword"],
		Action:       m["action"],
		DateRange:    m["date_range"],
		DateBefore:   m["date_before"],
		DateAfter:    m["date_after"],
		Reload:       m["reload"],
		OrderColumn:  m["order_column"],
		OrderDir:     m["order_dir"],
	}


	r := h.DisplayBudgetIOExtraApproved(c, fe)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayBudgetIOExtraApproved(c *fiber.Ctx, fe entity.DisplayBudgetIO) entity.ReturnResponse {
	// key := "temp_key_api_cpa_report_" + strings.ReplaceAll(helper.GetIpAddress(c), ".", "_")

	var (
		err        error
		total_data int64
		// isempty    bool
		budgetio []entity.BudgetIO
		// displaybudgetio []entity.BudgetIO
	)

	if fe.Action != "" || fe.Reload == "true" {
		budgetio, total_data, err = h.DS.GetDisplayBudgetIOApprovedAll(fe)
	} else {
		budgetio, total_data, err = h.DS.GetDisplayBudgetIOApprovedAll(fe)
	}

	if err == nil {

		if budgetio == nil {
			budgetio = []entity.BudgetIO{}
		}

		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithDataTable{
				Draw:            fe.Draw,
				Code:            fiber.StatusOK,
				Message:         config.OK_DESC,
				Data:            budgetio,
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

func (h *IncomingHandler) DisplaySummaryBudgetIO(c *fiber.Ctx) error {
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

	fe := entity.DisplaySummaryBudgetIO{
		Continent:   m["continent"],
		Country:     m["country"],
		Company:     m["company"],
		Partner:     m["partner"],
		Channel:     m["channel"],
		CampaignType: m["campaign_type"],
		Operator:    m["operator"],
		Service:     m["service"],
		ClientType:  m["client_type"],
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

	r := h.DisplaySummaryBudgetIOExtra(c, fe)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplaySummaryBudgetIOExtra(c *fiber.Ctx, fe entity.DisplaySummaryBudgetIO) entity.ReturnResponse {
	var (
		err    error
		result []entity.SummaryBudgetIOAgg
	)

	result, _, err = h.DS.GetDisplaySummaryBudgetIO(fe)

	if err != nil {
		return entity.ReturnResponse{
			HttpStatus: fiber.StatusNotFound,
			Rsp: entity.GlobalResponse{
				Code:    fiber.StatusNotFound,
				Message: "empty",
			},
		}
	}

	if result == nil {
		result = []entity.SummaryBudgetIOAgg{}
	}

	rows := MapSummaryBudgetIOToReportRows(result)

	return entity.ReturnResponse{
		HttpStatus: fiber.StatusOK,
		Rsp: entity.GlobalResponseWithDataTable{
			Draw:            fe.Draw,
			Code:            fiber.StatusOK,
			Message:         config.OK_DESC,
			Data:            rows,
			RecordsTotal:    len(rows),
			RecordsFiltered: len(rows),
		},
	}
}

func MapSummaryBudgetIOToReportRows(data []entity.SummaryBudgetIOAgg) []entity.IOReportRow {
	rows := make([]entity.IOReportRow, 0, len(data))
	for _, d := range data {
		recordedDay := 0
		now := time.Now()
		if d.Month == now.Format("2006-01") {
			recordedDay = now.Day()
		}

		rows = append(rows, entity.IOReportRow{
			BudgetIOID:  d.BudgetIOID,
			Region:      d.Continent,
			Country:     d.Country,
			Company:     d.Company,
			Partner:     d.Partner,
			Operator:    d.Operator,
			Channel:     d.Channel,
			ClientType:  d.ClientType,
			Service:     d.Service,
			Month:       d.Month,
			MOWeek1:     d.MOWeek1,
			MOWeek2:     d.MOWeek2,
			MOWeek3:     d.MOWeek3,
			MOWeek4:     d.MOWeek4,
			CostWeek1:   d.ActualWeek1,
			CostWeek2:   d.ActualWeek2,
			CostWeek3:   d.ActualWeek3,
			CostWeek4:   d.ActualWeek4,
			IOTarget:    d.IOTarget,
			MOTarget:    d.MOTarget,
			TargetCAC:   d.TargetCAC,
			EstLTV:      d.LTV,
			EstROAS:     d.ROAS,
			ROI:         d.ROI,
			RecordedDay: recordedDay,
		})
	}
	return rows
}

func (h *IncomingHandler) UpdateSummaryBudgetIO(c *fiber.Ctx) error {
	var req entity.UpdateSummaryBudgetIORequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "invalid request body",
		})
	}
	if req.ID == 0 && (req.Country == "" || req.Month == "") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "id or (country, month) required",
		})
	}
	if err := h.DS.UpdateSummaryBudgetIO(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "updated",
	})
}

