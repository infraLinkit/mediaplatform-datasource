package model

import (
	"database/sql"
	"fmt"
	"strings"

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

func (r *BaseModel) GetConversionLog(o entity.DisplayConversionLog) ([]entity.PixelStorage, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.PixelStorage{})
	if o.Action == "Search" {
		if o.Country != "" {
			query = query.Where("country = ?", o.Country)
		}
		if o.Operator != "" {
			query = query.Where("operator = ?", o.Operator)
		}
		if o.Pixel != "" {
			query = query.Where("service = ?", o.Pixel)
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

	var ss []entity.PixelStorage
	for rows.Next() {
		var s entity.PixelStorage
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	r.Logs.Debug(fmt.Sprintf("Total data : %d ...\n", len(ss)))

	return ss, total_rows, rows.Err()
}
