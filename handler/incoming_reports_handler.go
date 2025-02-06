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
	}

	r := h.DisplayPinPerformanceReportExtra(c, fe)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayPinPerformanceReportExtra(c *fiber.Ctx, fe entity.DisplayPinPerformanceReport) entity.ReturnResponse {

	key := "temp_key_api_pin_performance_report_" + strings.ReplaceAll(helper.GetIpAddress(c), ".", "_")

	var (
		err                         error
		x                           int
		isempty                     bool
		pinperformancereport        []entity.ApiPinPerformance
		displaypinperformancereport []entity.ApiPinPerformance
	)

	if fe.Action != "" {
		pinperformancereport, err = h.DS.GetApiPinPerformanceReport(fe)
	} else {
		if pinperformancereport, isempty = h.DS.RGetApiPinPerformanceReport(key, "$"); isempty {

			pinperformancereport, err = h.DS.GetApiPinPerformanceReport(fe)

			s, _ := json.Marshal(pinperformancereport)

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

		for i := x; i < len(pinperformancereport); i++ {

			//h.Logs.Debug(fmt.Sprintf("incr : %d, ID : %d", i, pinreport[i].ID))

			displaypinperformancereport = append(displaypinperformancereport, pinperformancereport[i])
			if i == pagesize {
				break
			}
		}

		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithData{
				Code:    fiber.StatusOK,
				Message: config.OK_DESC,
				Data:    displaypinperformancereport,
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

	r := h.DisplayCPAReportExtra(c, fe)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayCPAReportExtra(c *fiber.Ctx, fe entity.DisplayCPAReport) entity.ReturnResponse {
	// key := "temp_key_api_cpa_report_" + strings.ReplaceAll(helper.GetIpAddress(c), ".", "_")

	var (
		err error
		x   int
		// isempty          bool
		cpareport        []entity.SummaryCampaign
		displaycpareport []entity.SummaryCampaign
	)

	// if fe.Action != "" {
	cpareport, err = h.DS.GetDisplayCPAReport(fe)
	// }
	// else {

	// if cpareport, isempty = h.DS.RGetDisplayCPAReport(key, "$"); isempty {

	// 	cpareport, err = h.DS.RGetDisplayCPAReport(fe)

	// 	s, _ := json.Marshal(cpareport)

	// 	h.DS.SetData(key, "$", string(s))
	// 	h.DS.SetExpireData(key, 60)
	// }
	// }

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

		for i := x; i < len(cpareport); i++ {

			// h.Logs.Debug(fmt.Sprintf("incr : %d, ID : %d", i, cpareport[i].ID))

			displaycpareport = append(displaycpareport, cpareport[i])
			if i == pagesize {
				break
			}
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
