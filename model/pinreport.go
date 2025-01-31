package model

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/infraLinkit/mediaplatform-datasource/entity"
	"gorm.io/gorm"
)

func (r *BaseModel) PinReport(o entity.ApiPinReport) int {

	result := r.DB.Create(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return int(o.ID)
}

func (r *BaseModel) GetApiPinReport(o entity.DisplayPinReport) ([]entity.ApiPinReport, error) {

	var (
		q    string
		rows *sql.Rows
	)

	if o.Action == "Search" {
		if o.Country != "" {
			q = q + gorm.Expr("country = ?", o.Country).SQL
		}
		if o.Operator != "" {
			q = q + gorm.Expr(" AND operator = ?", o.Operator).SQL
		}
		if o.Service != "" {
			q = q + gorm.Expr(" AND service = ?", o.Service).SQL
		}
		if o.DateRange != "" {
			switch strings.ToUpper(o.DateRange) {
			case "TODAY":
				q = q + " AND date_send = CURRENT_DATE"
			case "YESTERDAY":
				q = q + gorm.Expr(" AND date_send BETWEEN CURRENT_DATE - INTERVAL ?", "1 DAY").SQL
			case "LAST7DAY":
				q = q + gorm.Expr(" AND date_send = CURRENT_DATE - INTERVAL ?", "7 DAY").SQL
			case "LAST30DAY":
				q = q + gorm.Expr(" AND date_send = CURRENT_DATE - INTERVAL ?", "30 DAY").SQL
			case "THISMONTH":
				q = q + gorm.Expr(" AND date_send BETWEEN CURRENT_DATE - INTERVAL ? AND CURRENT_DATE", "30 DAY").SQL
			case "LASTMONTH":
				q = q + gorm.Expr(" AND date_send BETWEEN CURRENT_DATE - INTERVAL ? AND CURRENT_DATE - INTERVAL ?", "60 DAY", "30 DAY").SQL
			case "CUSTOMRANGE":
				q = q + gorm.Expr(" AND date_send BETWEEN ? AND ?", o.DateBefore, o.DateAfter).SQL
			default:
				q = q + gorm.Expr(" AND date_send = ?", o.DateRange).SQL
			}
		}

		q = strings.TrimLeft(q, " AND")
		rows, _ = r.DB.Model(&entity.ApiPinReport{}).Where(q).Order("date_send").Rows()

	} else {

		rows, _ = r.DB.Model(&entity.ApiPinReport{}).Rows()
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
