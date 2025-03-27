package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/infraLinkit/mediaplatform-datasource/entity"
)

func (r *BaseModel) PinReport(o entity.ApiPinReport) int {

	result := r.DB.Create(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return int(o.ID)
}

func (r *BaseModel) GetApiPinReport(o entity.DisplayPinReport) ([]entity.ApiPinReport, error) {

	var (
		rows *sql.Rows
	)

	query := r.DB.Model(&entity.ApiPinReport{})

	if o.Action == "Search" {
		if o.Country != "" {
			query = query.Where("country = ?", o.Country)
		}
		if o.Operator != "" {
			query = query.Where("operator = ?", o.Operator)
		}
		if o.Service != "" {
			query = query.Where("service = ?", o.Service)
		}
		if o.DateRange != "" {
			switch strings.ToUpper(o.DateRange) {
			case "TODAY":
				query = query.Where("date_send = CURRENT_DATE")
			case "YESTERDAY":
				query = query.Where("date_send BETWEEN CURRENT_DATE - INTERVAL '1 DAY' AND CURRENT_DATE")
			case "LAST7DAY":
				query = query.Where("date_send BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
			case "LAST30DAY":
				query = query.Where("date_send BETWEEN CURRENT_DATE - INTERVAL '30 DAY' AND CURRENT_DATE")
			case "THISMONTH":
				query = query.Where("date_send >= DATE_TRUNC('month', CURRENT_DATE)")
			case "LASTMONTH":
				query = query.Where("date_send BETWEEN DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 DAY'")
			case "CUSTOMRANGE":
				query = query.Where("date_send BETWEEN ? AND ?", o.DateBefore, o.DateAfter)
			default:
				query = query.Where("date_send = ?", o.DateRange)
			}
		}

		rows, _ = query.Order("date_send").Rows()
	} else {
		rows, _ = query.Rows()
	}

	defer rows.Close()

	var (
		ss []entity.ApiPinReport
	)

	for rows.Next() {

		var s entity.ApiPinReport

		// ScanRows scans a row into a struct
		r.DB.ScanRows(rows, &s)

		ss = append(ss, s)
	}

	r.Logs.Debug(fmt.Sprintf("Total data : %d ...\n", len(ss)))

	return ss, rows.Err()
}

func (r *BaseModel) PinPerformanceReport(o entity.ApiPinPerformance) int {

	result := r.DB.Create(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return int(o.ID)
}

func (r *BaseModel) GetApiPinPerformanceReport(o entity.DisplayPinPerformanceReport) ([]entity.ApiPinPerformance, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.ApiPinPerformance{})
	if o.Action == "Search" {
		if o.Country != "" {
			query = query.Where("country = ?", o.Country)
		}
		if o.Operator != "" {
			query = query.Where("operator = ?", o.Operator)
		}
		if o.Service != "" {
			query = query.Where("service = ?", o.Service)
		}
		if o.DateRange != "" {
			switch strings.ToUpper(o.DateRange) {
			case "TODAY":
				query = query.Where("date_send = CURRENT_DATE")
			case "YESTERDAY":
				query = query.Where("date_send BETWEEN CURRENT_DATE - INTERVAL '1 DAY' AND CURRENT_DATE")
			case "LAST7DAY":
				query = query.Where("date_send BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
			case "LAST30DAY":
				query = query.Where("date_send BETWEEN CURRENT_DATE - INTERVAL '30 DAY' AND CURRENT_DATE")
			case "THISMONTH":
				query = query.Where("date_send >= DATE_TRUNC('month', CURRENT_DATE)")
			case "LASTMONTH":
				query = query.Where("date_send BETWEEN DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 DAY'")
			case "CUSTOMRANGE":
				query = query.Where("date_send BETWEEN ? AND ?", o.DateBefore, o.DateAfter)
			default:
				query = query.Where("date_send = ?", o.DateRange)
			}
		}
	}

	// Get the total count after applying filters
	query.Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Order("date_send").Rows()
	defer rows.Close()

	var ss []entity.ApiPinPerformance
	for rows.Next() {
		var s entity.ApiPinPerformance
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	r.Logs.Debug(fmt.Sprintf("Total data : %d ...\n", len(ss)))

	return ss, total_rows, rows.Err()
}

func (r *BaseModel) GetConversionLogReport(o entity.DisplayConversionLogReport) ([]entity.PixelStorage, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.PixelStorage{})
	query = query.Where("is_used = ?", "true")
	if o.CampaignType == "mainstream" {
		query = query.Where("campaign_objective = ?", "MAINSTREAM").Where("status_postback = ? ", "true")
	} else {
		query = query.Where("campaign_objective IN ?", []string{"CPA", "CPC", "CPI", "CPM"})
	}
	if o.Action == "Search" {
		if o.Country != "" {
			query = query.Where("country = ?", o.Country)
		}
		if o.Operator != "" {
			query = query.Where("operator = ?", o.Operator)
		}
		if o.Pixel != "" {
			query = query.Where("pixel = ?", o.Pixel)
		}
		if o.CampaignId != "" {
			query = query.Where("campaign_id = ?", o.CampaignId)
		}
		if o.CampaignType == "mainstream" {
			if o.Agency != "" {
				//compare to adnet maybe will change in the future
				query = query.Where("adnet = ?", o.Agency)
			}
		}
		if o.CampaignType == "s2s" {
			if o.StatusPostback != "" {
				query = query.Where("status_postback = ?", o.StatusPostback)
			}
			if o.Adnet != "" {
				query = query.Where("adnet = ?", o.Adnet)
			}
		}
		if o.DateRange != "" {
			switch strings.ToUpper(o.DateRange) {
			case "TODAY":
				query = query.Where("pxdate = CURRENT_DATE")
			case "YESTERDAY":
				query = query.Where("pxdate BETWEEN CURRENT_DATE - INTERVAL '1 DAY' AND CURRENT_DATE")
			case "LAST7DAY":
				query = query.Where("pxdate BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
			case "LAST30DAY":
				query = query.Where("pxdate BETWEEN CURRENT_DATE - INTERVAL '30 DAY' AND CURRENT_DATE")
			case "THISMONTH":
				query = query.Where("pxdate >= DATE_TRUNC('month', CURRENT_DATE)")
			case "LASTMONTH":
				query = query.Where("pxdate BETWEEN DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 DAY'")
			case "CUSTOMRANGE":
				dateEnd, _ := time.Parse("2006-01-02", o.DateEnd)
				query = query.Where("pxdate BETWEEN ? AND ?", o.DateStart, dateEnd.AddDate(0, 0, 1))
			default:
				query = query.Where("pxdate = ?", o.DateRange)
			}
		}
	}

	// Get the total count after applying filters
	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Order("pxdate").Rows()
	defer rows.Close()

	var ss []entity.PixelStorage
	for rows.Next() {
		var s entity.PixelStorage
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	return ss, total_rows, rows.Err()
}

func (r *BaseModel) GetPerformanceReport(o entity.PerformaceReportParams) ([]entity.PerformanceReport, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.SummaryCampaign{})

	query.Select(`country, company, client_type, campaign_name, operator, service, adnet, SUM(mo_received) AS pixel_received, SUM(postback) as pixel_send, SUM(cr_postback) as cr_postback,
SUM(cr_mo) as cr_mo, SUM(landing) as landing, SUM(ratio_send) as ratio_send, SUM(ratio_receive) as ratio_receive,SUM(po) as price_per_postback,SUM(cost_per_conversion) as cost_per_conversion,
SUM(agency_fee) as agency_fee, SUM(postback*po) as spending_to_adnets, SUM(total_waki_agency_fee), SUM(total_waki_agency_fee + po*postback) as total_spending,sum(cpa) as e_cpa`)

	if o.Action == "Search" {
		if o.Country != "" {
			query = query.Where("country = ?", o.Country)
		}
		if o.Operator != "" {
			query = query.Where("operator = ?", o.Operator)
		}
		if o.Partner != "" {
			query = query.Where("country = ?", o.Partner)
		}
		if o.CampaignType != "" {
			query = query.Where("campaign_type = ?", o.CampaignType)
		}
		if o.CampaignId != "" {
			query = query.Where("campaign_id = ?", o.CampaignId)
		}
		if o.CampaignName != "" {
			query = query.Where("campaign_name = ?", o.CampaignName)
		}
		if o.ClientType != "" {
			query = query.Where("client_type = ?", o.ClientType)
		}
		if o.Publisher != "" {
			query = query.Where("adnet = ?", o.Publisher)
		}
	}
	now := time.Now()
	dateStart, errStart := time.Parse("2006-01-02", o.DateStart)
	dateEnd, errEnd := time.Parse("2006-01-02", o.DateEnd)
	if errStart != nil {
		dateStart = now.AddDate(0, 0, -30)
	}
	if errEnd != nil {
		dateEnd = now
	}

	query = query.Where("summary_date BETWEEN ? AND ?", dateStart, dateEnd)

	query.Group("country, company, client_type, campaign_name, operator, service, adnet")

	// Get the total count after applying filters
	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Order("country").Rows()
	defer rows.Close()

	var ss []entity.PerformanceReport
	for rows.Next() {
		var s entity.PerformanceReport
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	r.Logs.Debug(fmt.Sprintf("Total data : %d ...\n", len(ss)))

	return ss, total_rows, rows.Err()
}

func (r *BaseModel) GetDataDistinctPerformanceReport(o entity.ApiPinPerformance) ([]entity.ApiPinPerformance, error) {

	var (
		rows *sql.Rows
	)

	query := r.DB.Model(&entity.ApiPinPerformance{})

	rows, _ = query.Distinct("date_send", "country", "operator", "service").Where("date_send = CURRENT_DATE()").Order("date_send").Rows()

	defer rows.Close()

	var (
		ss []entity.ApiPinPerformance
	)

	for rows.Next() {

		var s entity.ApiPinPerformance

		// ScanRows scans a row into a struct
		r.DB.ScanRows(rows, &s)

		ss = append(ss, s)
	}

	r.Logs.Debug(fmt.Sprintf("Total data : %d ...\n", len(ss)))

	return ss, rows.Err()
}
