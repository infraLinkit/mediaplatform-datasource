package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/infraLinkit/mediaplatform-datasource/entity"
)

func uniqueStrings(arr []string) []string {
	encountered := make(map[string]bool)
	result := []string{}

	for _, v := range arr {
		if !encountered[v] {
			encountered[v] = true
			result = append(result, v)
		}
	}
	return result
}

func GetDaysInMonth(year int, month time.Month) int {
	// 1. Find the first day of the *next* month.
	// We create a time.Time object for the 1st day of the specified month/year.
	firstOfCurrentMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)

	// Add one month to get the first day of the next month.
	firstOfNextMonth := firstOfCurrentMonth.AddDate(0, 1, 0)

	// 2. Subtract one day from the first day of the next month.
	// This gives us the last day of the current month.
	lastOfCurrentMonth := firstOfNextMonth.AddDate(0, 0, -1)

	// 3. The Day component of the last day is the number of days in the month.
	return lastOfCurrentMonth.Day()
}

func (r *BaseModel) GetReport(country string, operator string, client_type string, partner string, service string, campaign_objective string, date_range string, date_before string, date_after string, allowedAdnets []string, allowedCompanies []string) (entity.SummaryDashboardReport, error) {
	query := r.DB.Model(&entity.SummaryCampaign{})

	// GROUP BY BASED ON DATE_RANGE TYPE
	select_date := " DATE(summary_date) as date, "

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
		query = query.Where("summary_date BETWEEN TO_DATE(?, 'YYYY-MM') AND TO_DATE(?, 'YYYY-MM') + INTERVAL '1 month' - INTERVAL '1 day' ", date_before, date_after)
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
		} else {
			start_date = date_list[0]
			end_date = date_list[len(date_list)-1]
		}

		format := "2006-01-02"

		if date_range == "MONTHLY" {
			format = "2006-01"
		}

		start, _ := time.Parse(format, start_date)
		end, _ := time.Parse(format, end_date)

		current := start
		for !current.After(end) {
			date_list_result = append(date_list_result, current.Format(format))
			current = current.AddDate(0, 0, 1)
		}

		date_list_result = uniqueStrings(date_list_result)

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

func (r *BaseModel) GetDisplayDashboard(date_range string, date_before string, date_after string, allowedAdnets []string, allowedCompanies []string) (entity.SummaryDashboardData, error) {

	query := r.DB.Model(&entity.SummaryCampaign{})
	query_last_month := r.DB.Model(&entity.SummaryCampaign{})

	var SummaryDashboard entity.SummaryDashboardData

	var date_list []string

	currentTime := time.Now()

	switch date_range {
	case "TODAY":
		date_list = append(date_list, currentTime.Format("2006-01-02"))
	case "YESTERDAY":
		date_list = append(date_list, currentTime.AddDate(0, 0, -1).Format("2006-01-02"))
	case "LAST7DAY":
		for i := 0; i < 7; i++ {
			date := currentTime.AddDate(0, 0, -i).Format("2006-01-02")
			date_list = append(date_list, date)
		}
	case "LAST30DAY":
		current_day := currentTime.Day()
		for i := current_day; i <= 30+current_day; i++ {
			newDate := time.Date(currentTime.Year(), currentTime.Month(), i, 0, 0, 0, 0, time.Local)
			date_list = append(date_list, newDate.Format("2006-01-02"))
		}
	case "THISMONTH":
		current_day := currentTime.Day()
		last_day := GetDaysInMonth(currentTime.Year(), currentTime.Month())
		for i := current_day; i <= last_day; i++ {
			newDate := time.Date(currentTime.Year(), currentTime.Month(), i, 0, 0, 0, 0, time.Local)
			date_list = append(date_list, newDate.Format("2006-01-02"))
		}
	case "LASTMONTH":
		lastMonth := currentTime.AddDate(0, -1, 0)
		current_day := 1
		last_day := GetDaysInMonth(lastMonth.Year(), lastMonth.Month())
		for i := current_day; i <= last_day; i++ {
			newDate := time.Date(lastMonth.Year(), lastMonth.Month(), i, 0, 0, 0, 0, time.Local)
			date_list = append(date_list, newDate.Format("2006-01-02"))
		}
	case "CUSTOMRANGE":
		start, _ := time.Parse("2006-01-02", date_before)
		end, _ := time.Parse("2006-01-02", date_after)
		current := start
		for !current.After(end) {
			date_list = append(date_list, current.Format("2006-01-02"))
			current = current.AddDate(0, 0, 1)
		}
	}

	query = query.Where("summary_date IN ? ", date_list)
	query.Where("adnet IN ?", allowedAdnets)
	query.Where("company IN ?", allowedCompanies)

	where := ""

	for _, date := range date_list {
		where +=
			`CASE WHEN EXTRACT(DAY FROM (DATE_TRUNC('month', DATE '` + date + `') - INTERVAL '1 day')) < EXTRACT(DAY FROM DATE '` + date + `')
        THEN '0001-01-01' ELSE DATE '` + date + `' - INTERVAL '1 month'
    	END ,`
	}

	query_last_month.Where("summary_date IN(" + strings.TrimSuffix(where, ",") + ")")
	query_last_month.Where("adnet IN ?", allowedAdnets)
	query_last_month.Where("company IN ?", allowedCompanies)

	var dsp struct {
		Code  string "code"
		IsDsp bool   "is_dsp"
	}

	where_dsp := ""
	where_non_dsp := ""

	query_adnet := r.DB.Model(&entity.AdnetList{})
	rows, err := query_adnet.Select("code,is_dsp").Rows()

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			r.DB.ScanRows(rows, &dsp)
			if dsp.IsDsp == true {
				where_dsp += "'" + dsp.Code + "',"
			} else {
				where_non_dsp += "'" + dsp.Code + "',"
			}
		}
	}

	if where_dsp == "" {
		where_dsp = " AND false "
	} else {
		where_dsp = " AND adnet IN (" + strings.TrimSuffix(where_dsp, ",") + ")"
	}

	if where_non_dsp == "" {
		where_non_dsp = " AND false "
	} else {
		where_non_dsp = " AND adnet IN (" + strings.TrimSuffix(where_non_dsp, ",") + ")"
	}

	fmt.Println("where_non_dsp ", where_non_dsp)

	rows, err = query.Select(
		`summary_date as date,
		 SUM(mo_received) as total_mo,
		 COUNT(DISTINCT adnet) as total_active_adnet,
		 SUM(sbaf) as total_spending,
		 SUM(CASE WHEN campaign_objective IN('CPA','UPLOAD SMS') ` + where_non_dsp + ` THEN sbaf ELSE 0 END) as total_s2s_spending,
		 0 as total_api_spending,
		 SUM(CASE WHEN campaign_objective IN('MAINSTREAM') THEN sbaf ELSE 0 END) as total_mainstream_spending,
		 SUM(CASE WHEN campaign_objective IN('CPA','UPLOAD SMS') ` + where_dsp + ` THEN sbaf ELSE 0 END) as total_dsp_spending
		`).Group("summary_date").Order("summary_date ASC").Rows()

	if err == nil {

		defer rows.Close()
		var ss []entity.SummaryDashboardDetail

		for rows.Next() {

			var s entity.SummaryDashboardDetail
			r.DB.ScanRows(rows, &s)
			s.Date = strings.TrimSuffix(s.Date, "T00:00:00Z")

			SummaryDashboard.TotalActiveAdnet += s.TotalActiveAdnet
			SummaryDashboard.TotalMO += s.TotalMO
			SummaryDashboard.TotalSpending += s.TotalSpending
			SummaryDashboard.TotalS2SSpending += s.TotalS2sSpending
			SummaryDashboard.TotalMainstreamSpending += s.TotalMainstreamSpending
			SummaryDashboard.TotalDSPSpending += s.TotalDSPSpending

			total_spending_api := 0.0
			total_mo := 0

			api_query := r.DB.Model(&entity.ApiPinReport{})
			api_query.Where("date_send = ? ", s.Date)
			_ = api_query.Select("SUM(sbaf),SUM(total_mo)").Limit(1).Row().Scan(&total_spending_api, &total_mo)
			SummaryDashboard.TotalAPISpending += total_spending_api
			SummaryDashboard.TotalSpending += total_spending_api
			ss = append(ss, s)

		}

		// GET 1 MONTH PRIOR DATA
		rows, _ = query_last_month.Select(
			`summary_date as date,
		 SUM(mo_received) as total_mo,
		 COUNT(DISTINCT adnet) as total_active_adnet,
		 SUM(sbaf) as total_spending,
		 SUM(CASE WHEN campaign_objective IN('CPA','UPLOAD SMS') ` + where_non_dsp + ` THEN sbaf ELSE 0 END) as total_s2s_spending,
		 0 as total_api_spending,
		 SUM(CASE WHEN campaign_objective IN('MAINSTREAM') THEN sbaf ELSE 0 END) as total_mainstream_spending,
		 SUM(CASE WHEN campaign_objective IN('CPA','UPLOAD SMS') ` + where_dsp + ` THEN sbaf ELSE 0 END) as total_dsp_spending
		`).Group("summary_date").Order("summary_date ASC").Rows()

		defer rows.Close()
		var sl []entity.SummaryDashboardDetail

		for rows.Next() {
			var s entity.SummaryDashboardDetail
			r.DB.ScanRows(rows, &s)
			s.Date = strings.TrimSuffix(s.Date, "T00:00:00Z")
			sl = append(sl, s)
		}

		SummaryDashboard.DateList = date_list
		var DetailChartData []entity.DetailChartData
		var DetailChart entity.DetailChartData

		for _, date := range SummaryDashboard.DateList {

			DetailChart.Date = date
			DetailChart.TotalMO = 0
			DetailChart.TotalSpending = 0
			DetailChart.LastMonthTotalMO = 0
			DetailChart.LastMonthTotalSpending = 0
			DetailChart.LastMonthDate = ""

			for _, detail := range ss {

				if date == detail.Date {
					DetailChart.TotalMO = detail.TotalMO
					DetailChart.TotalSpending = detail.TotalSpending
				}

				for _, last_detail := range sl {
					prev_date, _ := time.Parse("2006-01-02", last_detail.Date)
					prev_date = prev_date.AddDate(0, 1, 0)
					if date == prev_date.Format("2006-01-02") {
						DetailChart.LastMonthTotalMO = last_detail.TotalMO
						DetailChart.LastMonthTotalSpending = last_detail.TotalSpending
						DetailChart.LastMonthDate = prev_date.Format("2006-01-02")
						// ADD WITH API
						total_spending_api := 0.0
						total_mo := 0
						api_query := r.DB.Model(&entity.ApiPinReport{})
						api_query.Where("date_send = ? ", date)
						_ = api_query.Select("SUM(sbaf),SUM(total_mo)").Limit(1).Row().Scan(&total_spending_api, &total_mo)

						DetailChart.LastMonthTotalSpending += total_spending_api
						DetailChart.LastMonthTotalMO += total_mo
						break
					}
				}
			}

			DetailChartData = append(DetailChartData, DetailChart)
		}

		SummaryDashboard.DetailChartData = DetailChartData
		return SummaryDashboard, nil
	}

	SummaryDashboard.DateList = date_list
	return SummaryDashboard, nil
}
