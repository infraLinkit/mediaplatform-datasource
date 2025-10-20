package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/infraLinkit/mediaplatform-datasource/entity"
)

func (r *BaseModel) GetReport(country string, operator string, client_type string, partner string, service string, campaign_objective string, date_range string, date_before string, date_after string, allowedAdnets []string, allowedCompanies []string) (entity.SummaryDashboardReport, error) {
	query := r.DB.Model(&entity.SummaryCampaign{})

	// GROUP BY BASED ON DATE_RANGE TYPE
	select_date := " DATE(summary_date) as date, "

	fmt.Println("DATE RANGE: ", date_range)

	switch date_range {
	case "TODAY":
		query.Where("summary_date = CURRENT_DATE")
	case "YESTERDAY":
		query = query.Where("summary_date = CURRENT_DATE - INTERVAL '1 DAY'")
	case "LAST7DAY":
		query = query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
	case "LAST30DAY":
		query = query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '30 DAY' AND CURRENT_DATE")
	case "THISMONTH":
		query = query.Where("summary_date >= DATE_TRUNC('month', CURRENT_DATE)")
	case "LASTMONTH":
		query = query.Where("summary_date BETWEEN DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 DAY'")
	case "CUSTOMRANGE":
		query = query.Where("summary_date BETWEEN ? AND ?", date_before, date_after)
	case "MONTHLY":
		query = query.Where("summary_date BETWEEN ? AND ?", date_before, date_after)
		select_date = "TO_CHAR(summary_date,'YYYY-MM') as date, "
	}

	rows, err := query.Select(select_date +
		"SUM(mo_received) as mo_received, " +
		"SUM(postback) as mo_sent, " +
		"SUM(sbaf) as spending_to_adnets, " +
		"SUM(saaf) as spending, " +
		"SUM(saaf)-SUM(sbaf) as waki_revenue").Group("date").Order("date ASC").Rows()

	if err == nil {
		defer rows.Close()

		var ss []entity.SummaryDashboardReportDetail
		var date_list []string
		var date_list_result []string
		var start_date string
		var end_date string

		for rows.Next() {
			var s entity.SummaryDashboardReportDetail

			r.DB.ScanRows(rows, &s)

			s.Date = strings.TrimSuffix(s.Date, "T00:00:00Z")

			ss = append(ss, s)
			date_list = append(date_list, s.Date)
		}

		var SummaryDashboardReport entity.SummaryDashboardReport
		SummaryDashboardReport.DateRange = date_range
		SummaryDashboardReport.Detail = ss

		if len(date_list) == 0 {
			start_date = ""
			end_date = ""
		} else if len(date_list) == 1 {
			start_date = date_list[0]
			end_date = date_list[1]
		} else {
			start_date = date_list[0]
			end_date = date_list[len(date_list)-1]
		}

		start, _ := time.Parse("2006-01-02", start_date)
		end, _ := time.Parse("2006-01-02", end_date)

		//fmt.Println("START ", start)
		//fmt.Println("END ", end)

		current := start
		for !current.After(end) {
			date_list_result = append(date_list_result, current.Format("2006-01-02"))
			current = current.AddDate(0, 0, 1)
		}

		//fmt.Println("DATE LIST ", date_list_result)
		SummaryDashboardReport.DateList = date_list_result

		return SummaryDashboardReport, nil
	}

	return entity.SummaryDashboardReport{DateRange: date_range, Detail: []entity.SummaryDashboardReportDetail{}}, nil
}

func (r *BaseModel) GetCampaign(order_type string, order_by string, offset string, date_range string, date_before string, date_after string, allowedAdnets []string, allowedCompanies []string) ([]entity.TopCampaign, error) {
	/*
			"MO_RECEIVED"
		    "SPENDING"
		    "CR_MO"
		    "CR_POSTBACK"
		    "E_CPA"
	*/
	query := r.DB.Model(&entity.SummaryCampaign{})
	field_order := ""
	desc := "DESC"

	if order_type == "WORST" {
		desc = "ASC"
	}

	fmt.Println("ORDER BY ", order_by)

	switch order_by {
	case "MO_RECEIVED":
		field_order = "mo_received"
	case "SPENDING":
		field_order = "sbaf"
	case "CR_MO":
		field_order = "cr_mo"
	case "CR_POSTBACK":
		field_order = "cr_postback"
	case "E_CPA":
		field_order = "cpa"
	}

	query.Order(field_order + " " + desc)
	limit, e := strconv.Atoi(offset)

	if e != nil {
		limit = 5
	}

	switch date_range {
	case "TODAY":
		query.Where("summary_date = CURRENT_DATE")
	case "YESTERDAY":
		query = query.Where("summary_date = CURRENT_DATE - INTERVAL '1 DAY'")
	case "LAST7DAY":
		query = query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
	case "LAST30DAY":
		query = query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '30 DAY' AND CURRENT_DATE")
	case "THISMONTH":
		query = query.Where("summary_date >= DATE_TRUNC('month', CURRENT_DATE)")
	case "LASTMONTH":
		query = query.Where("summary_date BETWEEN DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 DAY'")
	case "CUSTOMRANGE":
		query = query.Where("summary_date BETWEEN ? AND ?", date_before, date_after)
	}

	/*CampaignID string  `json:"campaign_id"`
	Country    string  `json:"country"`
	Landing    int     `json:"landing"`
	MO         int     `json:"mo_received"`
	Postback   int     `json:"postback"`
	CRMO       float64 `json:"cr_mo"`
	CRPostback float64 `json:"cr_postback"`
	URL        string  `json:"url"`
	ECPA*/

	rows, err := query.Select(
		`campaign_id,
		country,
		landing,
		mo_received,
		postback,
		cr_mo,
		cr_postback,
		url_after,
		cpa`).Limit(limit).Rows()

	if err == nil {
		defer rows.Close()

		var ss []entity.TopCampaign

		for rows.Next() {
			var s entity.TopCampaign

			r.DB.ScanRows(rows, &s)

			c := r.DB.Model(&entity.Country{})
			_ = c.Select("name").Where("code=?", s.Country).Row().Scan(&s.CountryName)

			ss = append(ss, s)
		}

		return ss, nil
	}

	return []entity.TopCampaign{}, err
}

func (r *BaseModel) GetDisplayDashboard(date_range string, date_before string, date_after string, allowedAdnets []string, allowedCompanies []string) (entity.SummaryDashboard, error) {

	query := r.DB.Model(&entity.SummaryCampaign{})
	api_query := r.DB.Model(&entity.ApiPinReport{})

	var SummaryDashboard entity.SummaryDashboard

	switch date_range {
	case "TODAY":
		query.Where("summary_date = CURRENT_DATE")
		api_query.Where("date_send = CURRENT_DATE")
	case "YESTERDAY":
		query = query.Where("summary_date = CURRENT_DATE - INTERVAL '1 DAY'")
		api_query.Where("date_send = CURRENT_DATE - INTERVAL '1 DAY'")
	case "LAST7DAY":
		query = query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
		api_query.Where("date_send BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
	case "LAST30DAY":
		query = query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '30 DAY' AND CURRENT_DATE")
		api_query.Where("date_send BETWEEN CURRENT_DATE - INTERVAL '30 DAY' AND CURRENT_DATE")
	case "THISMONTH":
		query = query.Where("summary_date >= DATE_TRUNC('month', CURRENT_DATE)")
		api_query.Where("date_send >= DATE_TRUNC('month', CURRENT_DATE)")
	case "LASTMONTH":
		query = query.Where("summary_date BETWEEN DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 DAY'")
		api_query.Where("date_send BETWEEN DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 DAY'")
	case "CUSTOMRANGE":
		query = query.Where("summary_date BETWEEN ? AND ?", date_before, date_after)
		api_query.Where("date_send BETWEEN ? AND ?", date_before, date_after)
	}

	query.Where("adnet IN ?", allowedAdnets)
	query.Where("company IN ?", allowedCompanies)

	_ = query.Select(
		`COUNT(DISTINCT adnet),
		 SUM(mo_received),
		 SUM(sbaf),
		 SUM(CASE WHEN campaign_objective IN('CPA','UPLOAD SMS') THEN sbaf ELSE 0 END),
		 SUM(CASE WHEN campaign_objective IN('MAINSTREAM') THEN sbaf ELSE 0 END),
		 0
		 `).Row().Scan(
		&SummaryDashboard.TotalActiveAdnet,
		&SummaryDashboard.TotalMO,
		&SummaryDashboard.TotalSpending,
		&SummaryDashboard.TotalS2SSpending,
		&SummaryDashboard.TotalMainstreamSpending,
		&SummaryDashboard.TotalDSPSpending)

	_ = api_query.Select("SUM(sbaf)").Row().Scan(&SummaryDashboard.TotalAPISpending)

	SummaryDashboard.TotalSpending = SummaryDashboard.TotalSpending + SummaryDashboard.TotalAPISpending

	return SummaryDashboard, nil
}
