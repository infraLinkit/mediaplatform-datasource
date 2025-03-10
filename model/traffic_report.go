package model

import (
	"database/sql"
	"fmt"
	"strings"

	// "encoding/json"
	// "errors"

	"github.com/infraLinkit/mediaplatform-datasource/entity"
)

func (r *BaseModel) GetDisplayTrafficReport(o entity.DisplayTrafficReport) ([]entity.SummaryCampaign, error) {
	var rows *sql.Rows
	var err error

	query := r.DB.Model(&entity.SummaryCampaign{})
	// fmt.Println(query)
	if o.Action == "Search" {
		if o.Country != "" {
			query = query.Where("country = ?", o.Country)
		}
		// if o.Company != "" {
		// 	query = query.Where("company = ?", o.Company)
		// }
		// if o.ClientType != "" {
		// 	query = query.Where("client_type = ?", o.ClientType)
		// }
		if o.Operator != "" {
			query = query.Where("operator = ?", o.Operator)
		}
		if o.CampaignName != "" {
			query = query.Where("campaign_name = ?", o.CampaignName)
		}
		// if o.Partner != "" {
		// 	query = query.Where("partner = ?", o.Partner)
		// }
		if o.Adnet != "" { //Publisher
			query = query.Where("adnet = ?", o.Adnet)
		}
		if o.Service != "" {
			query = query.Where("service = ?", o.Service)
		}
		if o.DateRange != "" {
			switch strings.ToUpper(o.DateRange) {
			case "TODAY":
				query = query.Where("summary_date = CURRENT_DATE")
			case "YESTERDAY":
				query = query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '1 DAY' AND CURRENT_DATE")
			case "LAST7DAY":
				query = query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
			case "LAST30DAY":
				query = query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '30 DAY' AND CURRENT_DATE")
			case "THISMONTH":
				query = query.Where("summary_date >= DATE_TRUNC('month', CURRENT_DATE)")
			case "LASTMONTH":
				query = query.Where("summary_date BETWEEN DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 DAY'")
			case "CUSTOMRANGE":
				query = query.Where("summary_date BETWEEN ? AND ?", o.DateBefore, o.DateAfter)
			default:
				query = query.Where("summary_date = ?", o.DateRange)
			}
		}

		rows, err = query.Order("created_at DESC").Order("id DESC").Rows()
		if err != nil {
			return []entity.SummaryCampaign{}, err
		}
	} else {
		rows, err = query.Order("summary_date DESC").Order("id DESC").Rows()
		if err != nil {
			return []entity.SummaryCampaign{}, err
		}
	}

	if rows == nil {
		return []entity.SummaryCampaign{}, nil
	}
	defer rows.Close()

	var ss []entity.SummaryCampaign

	for rows.Next() {
		var s entity.SummaryCampaign

		r.DB.ScanRows(rows, &s)

		ss = append(ss, s)
	}

	r.Logs.Debug(fmt.Sprintf("Total data : %d ... \n", len(ss)))

	return ss, rows.Err()
}
