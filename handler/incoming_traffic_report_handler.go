package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource/config"
	"github.com/infraLinkit/mediaplatform-datasource/entity"
	// "github.com/infraLinkit/mediaplatform-datasource/helper"
)

func (h *IncomingHandler) DisplayTrafficReport(c *fiber.Ctx) error {
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
	fe := entity.DisplayTrafficReport{
		Draw:       draw,
		Page:       page,
		PageSize:   pageSize,
		Action:     m["action"],
		DateRange:  m["date_range"],
		DateBefore: m["date_before"],
		DateAfter:  m["date_after"],
	}

	r := h.DisplayTrafficReportExtra(c, fe)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayTrafficReportExtra(c *fiber.Ctx, fe entity.DisplayTrafficReport) entity.ReturnResponse {
	// key := "tempt_key_api_traffic_report_" + strings.ReplaceAll(helper, GetIpAddress(c), ".", "_")

	var (
		err error
		x   int
		// isempty     bool
		trafficreport        []entity.SummaryCampaign
		displaytrafficreport []entity.SummaryCampaign
	)

	if fe.Action != "" {
		trafficreport, err = h.DS.GetDisplayTrafficReport(fe)
	} else {
		trafficreport, err = h.DS.GetDisplayTrafficReport(fe)
	}

	if err == nil {

		pagesize := fe.PageSize
		if pagesize == 0 {
			pagesize = PAGESIZE
		}
		if fe.Page >= 2 {
			x = pagesize * (fe.Page - 1)
		} else {
			x = 0
		}

		for i := x; i < len(trafficreport) && i < x+pagesize; i++ {

			displaytrafficreport = append(displaytrafficreport, trafficreport[i])
		}
		if displaytrafficreport == nil {
			return entity.ReturnResponse{
				HttpStatus: fiber.StatusOK,
				Rsp: entity.GlobalResponseWithDataTable{
					Draw:            fe.Draw,
					Code:            fiber.StatusOK,
					Message:         config.OK_DESC,
					Data:            []entity.DisplayTrafficReport{},
					RecordsTotal:    len(trafficreport),
					RecordsFiltered: len(trafficreport),
				},
			}
		}

		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithDataTable{
				Draw:            fe.Draw,
				Code:            fiber.StatusOK,
				Message:         config.OK_DESC,
				Data:            displaytrafficreport,
				RecordsTotal:    len(trafficreport),
				RecordsFiltered: len(trafficreport),
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
