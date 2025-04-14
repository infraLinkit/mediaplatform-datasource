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

// const PAGESIZE int = 10

func (h *IncomingHandler) DisplayMainstreamReport(c *fiber.Ctx) error {
	c.Set("Content-type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	m := c.Queries()

	page, _ := strconv.Atoi(m["page"])
	pageSize, err := strconv.Atoi(m["page_size"])
	if err != nil {
		pageSize = PAGESIZE
	}
	draw, _ := strconv.Atoi(m["draw"])
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
		Adnet:        m["adnet"],
		Service:      m["service"],
		Draw:         draw,
		Page:         page,
		PageSize:     pageSize,
		Action:       m["action"],
		DateRange:    m["date_range"],
		DateBefore:   m["date_before"],
		DateAfter:    m["date_after"],
	}
	r := h.DisplayMainstreamReportExtra(c, fe)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayMainstreamReportExtra(c *fiber.Ctx, fe entity.DisplayCPAReport) entity.ReturnResponse {
	key := "temp_key_api_mainstream_report" +
		"_" + fe.CampaignId +
		"_" + fe.CampaignName +
		"_" + fe.Country +
		"_" + fe.ClientType +
		"_" + fe.Company +
		"_" + fe.Operator +
		"_" + fe.Partner +
		"_" + fe.Aggregator +
		"_" + fe.Adnet +
		"_" + fe.Service +
		"_" + fe.DateRange +
		"_" + fe.DateBefore +
		"_" + fe.DateAfter +
		"_" + strconv.Itoa(fe.Page) +
		"_" + strconv.Itoa(fe.PageSize) + strings.ReplaceAll(helper.GetIpAddress(c), ".", "_")
	var (
		err                     error
		x                       int
		isempty                 bool
		mainstreamreport        []entity.SummaryCampaign
		displaymainstreamreport []entity.SummaryCampaign
	)

	if mainstreamreport, isempty = h.DS.RGetDisplayMainstreamReport(key, "$"); isempty {
		mainstreamreport, err = h.DS.GetDisplayMainstreamReport(fe)
		s, _ := json.Marshal(mainstreamreport)

		h.DS.SetData(key, "$", string(s))
		h.DS.SetExpireData(key, 60)
	}

	if err != nil {
		return entity.ReturnResponse{
			HttpStatus: fiber.StatusNotFound,
			Rsp: entity.GlobalResponse{
				Code:    fiber.StatusNotFound,
				Message: "empty",
			},
		}
	}

	pagesize := fe.PageSize
	if pagesize == 0 {
		pagesize = PAGESIZE
	}

	if fe.Page >= 2 {
		x = pagesize * (fe.Page - 1)
	} else {
		x = 0
	}

	for i := x; i < len(mainstreamreport) && i < x+pagesize; i++ {
		displaymainstreamreport = append(displaymainstreamreport, mainstreamreport[i])
	}
	if displaymainstreamreport == nil {
		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithDataTable{
				Draw:            fe.Draw,
				Code:            fiber.StatusOK,
				Message:         config.OK_DESC,
				Data:            []entity.SummaryCampaign{},
				RecordsTotal:    len(mainstreamreport),
				RecordsFiltered: len(mainstreamreport),
			},
		}
	}
	return entity.ReturnResponse{
		HttpStatus: fiber.StatusOK,
		Rsp: entity.GlobalResponseWithDataTable{
			Draw:            fe.Draw,
			Code:            fiber.StatusOK,
			Message:         config.OK_DESC,
			Data:            displaymainstreamreport,
			RecordsTotal:    len(mainstreamreport),
			RecordsFiltered: len(mainstreamreport),
		},
	}
}
