package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/infraLinkit/mediaplatform-datasource/entity"
)

func (r *BaseModel) GetDisplayCPAReport(o entity.DisplayCPAReport) ([]entity.SummaryCampaign, error) {
	var rows *sql.Rows
	var err error

	query := r.DB.Model(&entity.SummaryCampaign{})
	// fmt.Println(query)
	if o.Action == "Search" {
		if o.Country != "" {
			query = query.Where("country = ?", o.Country)
		}
		if o.Company != "" {
			query = query.Where("company = ?", o.Company)
		}
		if o.ClientType != "" {
			query = query.Where("client_type = ?", o.ClientType)
		}
		if o.Operator != "" {
			query = query.Where("operator = ?", o.Operator)
		}
		if o.CampaignName != "" {
			query = query.Where("campaign_name = ?", o.CampaignName)
		}
		if o.Partner != "" {
			query = query.Where("partner = ?", o.Partner)
		}
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
			case "LAST 7 DAY":
				query = query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
			case "LAST 30 DAY":
				query = query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '30 DAY' AND CURRENT_DATE")
			case "THIS MONTH":
				query = query.Where("summary_date >= DATE_TRUNC('month', CURRENT_DATE)")
			case "LAST MONTH":
				query = query.Where("summary_date BETWEEN DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 DAY'")
			case "CUSTOM RANGE":
				query = query.Where("summary_date BETWEEN ? AND ?", o.DateBefore, o.DateAfter)
			default:
				query = query.Where("summary_date = ?", o.DateRange)
			}
		}

		rows, err = query.Order("created_at DESC").Order("id DESC").Rows()
		if err != nil {
			return nil, err
		}
	} else {
		rows, err = query.Order("created_at DESC").Order("id DESC").Rows()
		if err != nil {
			return nil, err
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

func (r *BaseModel) FindSummaryCampaign(id int) (entity.SummaryCampaign, error) {
	var s entity.SummaryCampaign
	result := r.DB.First(&s, id)
	r.Logs.Debug(fmt.Sprintf("Total data : %d ... \n", result.RowsAffected))
	return s, result.Error
}

func (r *BaseModel) UpdateRatioModel(o entity.SummaryCampaign, id int) error {
	result := r.DB.Exec("UPDATE summary_campaigns SET ratio_send = ?, ratio_receive = ? WHERE ID = ?", o.RatioSend, o.RatioReceive, id)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) UpdatePostbackModel(o entity.SummaryCampaign, id int) error {
	result := r.DB.Exec("UPDATE summary_campaigns SET postback = ? WHERE ID = ?", o.Postback, id)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) UpdateAgencyCostModel(o entity.SummaryCampaign) error {
	result := r.DB.Exec("UPDATE summary_campaigns SET agency_fee = COALESCE(?, agency_fee), cost_per_conversion = COALESCE(?, cost_per_conversion) WHERE DATE_TRUNC('month', summary_date) = DATE_TRUNC('month', CURRENT_DATE) AND summary_date >= DATE_TRUNC('month', CURRENT_DATE) + INTERVAL '0 DAY'", o.AgencyFee, o.CostPerConversion)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error

}

func (r *BaseModel) GetDisplayCostReport(o entity.DisplayCostReport) ([]entity.CostReport, error) {
	var (
		rows *sql.Rows
		err  error
	)

	query := r.DB.Model(&entity.SummaryCampaign{})

	if o.Action == "Search" {

		if o.Adnet != "" {
			query = query.Where("adnet = ?", o.Adnet)
		}

		if o.DataBasedOn != "" {
			switch strings.ToUpper(o.DataBasedOn) {
			case "HIGHEST CONVERSION S2S":
				query = query.Order("conversion1 DESC")
			case "LOWEST CONVERSION S2S":
				query = query.Order("conversion1 ASC")
			case "HIGHEST CONVERSION API":
				query = query.Order("conversion2 DESC")
			case "LOWEST CONVERSION API":
				query = query.Order("conversion2 ASC")
			case "HIGHEST COST S2S":
				query = query.Order("cost1 DESC")
			case "LOWEST COST S2S":
				query = query.Order("cost1 ASC")
			case "HIGHEST COST API":
				query = query.Order("cost2 DESC")
			case "LOWEST COST API":
				query = query.Order("cost2 ASC")
			default:
				query = query.Order("conversion1 DESC")
			}
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

		rows, err = query.Select(`
			adnet,
			SUM(postback) as conversion1,
			SUM(sbaf) as cost1,
			NULL as conversion2,
			NULL as cost2
		`).Group("adnet").
			Rows()

		if err != nil {
			return nil, err
		}
		if rows == nil {
			return []entity.CostReport{}, nil
		}
		defer rows.Close()
	} else {
		rows, err = query.Order("summary_date DESC").Order("id DESC").Rows()
		if err != nil {
			return nil, err
		}
	}

	var results []entity.CostReport

	for rows.Next() {
		var item entity.CostReport
		if err := r.DB.ScanRows(rows, &item); err != nil {
			r.Logs.Error(fmt.Sprintf("Error scanning row: %v", err))
			continue
		}
		results = append(results, item)
	}

	r.Logs.Debug(fmt.Sprintf("Total data : %d ...\n", len(results)))

	return results, rows.Err()
}
func (r *BaseModel) GetDisplayCostReportDetail(o entity.DisplayCostReport) ([]entity.CostReport, error) {
	var (
		rows *sql.Rows
		err  error
	)

	query := r.DB.Model(&entity.SummaryCampaign{})

	if o.Action == "Search" {

		if o.Adnet != "" {
			query = query.Where("adnet = ?", o.Adnet)
		}

		if o.DateRange != "" {
			today := time.Now()
			switch o.DateRange {
			case "Today":
				query = query.Where("summary_date = CURRENT_DATE")
			case "Yesterday":
				query = query.Where("summary_date = CURRENT_DATE - INTERVAL '1 DAY'")
			case "Last 7 Days":
				query = query.Where("summary_date >= CURRENT_DATE - INTERVAL '6 DAY'")
			case "Last 30 Days":
				query = query.Where("summary_date >= CURRENT_DATE - INTERVAL '29 DAY'")
			case "This Month":
				query = query.Where("EXTRACT(YEAR FROM summary_date) = ? AND EXTRACT(MONTH FROM summary_date) = ?", today.Year(), int(today.Month()))
			case "Last Month":
				lastMonth := today.AddDate(0, -1, 0)
				query = query.Where("EXTRACT(YEAR FROM summary_date) = ? AND EXTRACT(MONTH FROM summary_date) = ?", lastMonth.Year(), int(lastMonth.Month()))
			case "Custom Range":
				if o.DateBefore != "" && o.DateAfter != "" {
					query = query.Where("summary_date BETWEEN ? AND ?", o.DateBefore, o.DateAfter)
				}
			}
		}

		rows, err = query.Select(`
			country,
			operator,
			SUM(postback) as conversion1,
			SUM(sbaf) as cost1,
			NULL as conversion2,
			NULL as cost2
		`).Group("country, operator").
			Rows()

		if err != nil {
			return nil, err
		}
		if rows == nil {
			return []entity.CostReport{}, nil
		}
		defer rows.Close()
	} else {
		rows, err = query.Order("summary_date DESC").Order("id DESC").Rows()
		if err != nil {
			return nil, err
		}
	}

	var results []entity.CostReport

	for rows.Next() {
		var item entity.CostReport
		if err := r.DB.ScanRows(rows, &item); err != nil {
			r.Logs.Error(fmt.Sprintf("Error scanning row: %v", err))
			continue
		}
		results = append(results, item)
	}

	r.Logs.Debug(fmt.Sprintf("Total data : %d ...\n", len(results)))

	return results, rows.Err()
}
