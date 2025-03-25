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

func (r *BaseModel) GetConversionLogReport(o entity.DisplayConversionLogReport) ([]entity.MO, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.MO{})
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
				query = query.Where("pxdate BETWEEN ? AND ?", o.DateBefore, o.DateAfter)
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

	var ss []entity.MO
	for rows.Next() {
		var s entity.MO
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	r.Logs.Debug(fmt.Sprintf("Total data : %d ...\n", len(ss)))

	return ss, total_rows, rows.Err()
}

func (r *BaseModel) GetPeformanceReport(o entity.PerformaceReportParams) ([]entity.ApiPinPerformance, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.ApiPinPerformance{})

	query.Select("country, company, client_type, campaign_name, operator, service, adnet, publisher, date_send, SUM(pin_request) AS pin_request, SUM(unique_pin_request) AS unique_pin_request, SUM(pin_sent) AS pin_sent, SUM(pin_failed) AS pin_failed, SUM(verify_request) AS verify_request, SUM(verify_request_unique) AS verify_request_unique, SUM(pin_ok) AS pin_ok, SUM(pin_not_ok) AS pin_not_ok, SUM(pin_ok_send_adnet) AS pin_ok_send_adnet, SUM(cpa) AS cpa, SUM(cpa_waki) AS cpa_waki, SUM(estimated_arpu) AS estimated_arpu, SUM(sbaf) AS sbaf, SUM(saaf) AS saaf, SUM(charged_mo) AS charged_mo, SUM(subs_cr) AS subs_cr, SUM(adnet_cr) AS adnet_cr, SUM(cac) AS cac, SUM(paid_cac) AS paid_cac, SUM(cr_mo) AS cr_mo, SUM(cr_postback) AS cr_postback, SUM(landing) AS landing, SUM(roi) AS roi, SUM(arpu90) AS arpu90, SUM(billing_rate_fp) AS billing_rate_fp, SUM(ratio) AS ratio, SUM(price_per_postback) AS price_per_postback, SUM(cost_per_conversion) AS cost_per_conversion, SUM(agency_fee) AS agency_fee, SUM(total_waki_agency_fee) AS total_waki_agency_fee, SUM(total_spending) AS total_spending")

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
			query = query.Where("publisher = ?", o.ClientType)
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

	query = query.Where("date_send BETWEEN ? AND ?", dateStart, dateEnd)

	query.Group("country, company, client_type, campaign_name, operator, service, adnet, publisher, date_send")

	// Get the total count after applying filters
	query.Unscoped().Count(&total_rows)

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
