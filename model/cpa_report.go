package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/infraLinkit/mediaplatform-datasource/entity"
)

func (r *BaseModel) GetDisplayCPAReport(o entity.DisplayCPAReport, allowedCompanies []string) ([]entity.SummaryCampaign, int64, entity.TotalSummaryCampaign, error) {
	var rows *sql.Rows
	var err error
	var total_rows int64
	var TotalSummaryCampaign entity.TotalSummaryCampaign

	query := r.DB.Model(&entity.SummaryCampaign{})
	t_query := r.DB.Model(&entity.SummaryCampaign{})

	if o.CampaignObjective != "" {
		query = query.Where("campaign_objective = ?", o.CampaignObjective)
		t_query = t_query.Where("campaign_objective = ?", o.CampaignObjective)
	} else {
		query = query.Where("campaign_objective = ? OR campaign_objective = ?", "CPA", "UPLOAD SMS")
		t_query = t_query.Where("campaign_objective = ? OR campaign_objective = ?", "CPA", "UPLOAD SMS")
	}

	query = query.Where("mo_received > 0")
	query = query.Where("company IN ?", allowedCompanies)

	t_query = t_query.Where("mo_received > 0")
	t_query = t_query.Where("company IN ?", allowedCompanies)

	if o.Action == "Search" {
		if o.CampaignId != "" {
			query = query.Where("campaign_id = ?", o.CampaignId)
			t_query = t_query.Where("campaign_id = ?", o.CampaignId)
		}
		if o.UrlServiceKey != "" {
			query = query.Where("url_service_key = ?", o.UrlServiceKey)
			t_query = t_query.Where("url_service_key = ?", o.UrlServiceKey)
		}
		if o.Country != "" {
			query = query.Where("country = ?", o.Country)
			t_query = t_query.Where("url_service_key = ?", o.UrlServiceKey)
		}
		if o.Company != "" {
			query = query.Where("company = ?", o.Company)
			t_query = t_query.Where("url_service_key = ?", o.UrlServiceKey)
		}
		if o.ClientType != "" {
			query = query.Where("client_type = ?", o.ClientType)
			t_query = t_query.Where("url_service_key = ?", o.UrlServiceKey)
		}
		if o.Operator != "" {
			query = query.Where("operator = ?", o.Operator)
			t_query = t_query.Where("url_service_key = ?", o.UrlServiceKey)
		}
		if o.CampaignName != "" {
			query = query.Where("campaign_name = ?", o.CampaignName)
			t_query = t_query.Where("campaign_name = ?", o.CampaignName)
		}
		if o.Partner != "" {
			query = query.Where("partner = ?", o.Partner)
			t_query = t_query.Where("partner = ?", o.Partner)
		}
		if o.Adnet != "" { //Publisher
			query = query.Where("adnet = ?", o.Adnet)
			t_query = t_query.Where("partner = ?", o.Partner)
		}
		if o.Service != "" {
			query = query.Where("service = ?", o.Service)
			t_query = t_query.Where("partner = ?", o.Partner)
		}
		if o.DateRange != "" {
			switch strings.ToUpper(o.DateRange) {
			case "TODAY":
				query = query.Where("summary_date = CURRENT_DATE")
				t_query = t_query.Where("summary_date = CURRENT_DATE")
			case "YESTERDAY":
				query = query.Where("summary_date = CURRENT_DATE - INTERVAL '1 DAY'")
				t_query = t_query.Where("summary_date = CURRENT_DATE - INTERVAL '1 DAY'")
			case "LAST7DAY":
				query = query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
				t_query = t_query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
			case "LAST30DAY":
				query = query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '30 DAY' AND CURRENT_DATE")
				t_query = t_query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '30 DAY' AND CURRENT_DATE")
			case "THISMONTH":
				query = query.Where("summary_date >= DATE_TRUNC('month', CURRENT_DATE)")
				t_query = t_query.Where("summary_date >= DATE_TRUNC('month', CURRENT_DATE)")
			case "LASTMONTH":
				query = query.Where("summary_date BETWEEN DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 DAY'")
				t_query = t_query.Where("summary_date BETWEEN DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 DAY'")
			case "CUSTOMRANGE":
				query = query.Where("summary_date BETWEEN ? AND ?", o.DateBefore, o.DateAfter)
				t_query = t_query.Where("summary_date BETWEEN ? AND ?", o.DateBefore, o.DateAfter)
			default:
				query = query.Where("summary_date = ?", o.DateRange)
				t_query = t_query.Where("summary_date = ?", o.DateRange)
			}
		} else {
			query = query.Where("summary_date = CURRENT_DATE")
			t_query = t_query.Where("summary_date = CURRENT_DATE")
		}

		rows, err = query.Order("summary_date DESC").Order("id DESC").Rows()

		if err != nil {
			return []entity.SummaryCampaign{}, 0, TotalSummaryCampaign, err
		}
	} else {

		rows, err = query.Order("summary_date DESC").Order("id DESC").Rows()
		if err != nil {
			return []entity.SummaryCampaign{}, 0, TotalSummaryCampaign, err
		}
	}

	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Rows()
	defer rows.Close()

	var ss []entity.SummaryCampaign

	for rows.Next() {
		var s entity.SummaryCampaign

		r.DB.ScanRows(rows, &s)

		ss = append(ss, s)
	}

	if total_rows > 0 {

		// COUNT THE SUMMARIZE
		_ = t_query.Select("SUM(landing) as landing,"+
			"SUM(mo_received) as mo_received,"+
			"SUM(postback) as postback,"+
			"AVG(po) as price_per_postback,"+
			"SUM(cost_per_conversion) as cost_per_conversion,"+
			"SUM(total_waki_agency_fee) as total_waki_agency_fee,"+
			"SUM(sbaf) as spending_to_adnet,"+
			"SUM(saaf) as spending,"+
			"SUM(technical_fee) as technical_fee,"+
			"SUM(saaf)-SUM(sbaf) as waki_revenue,"+
			"CASE WHEN SUM(landing)>0 THEN ROUND(SUM(mo_received)/SUM(landing)::numeric,5) ELSE 0 END as cr_mo,"+
			"CASE WHEN SUM(landing)>0 THEN ROUND(SUM(postback)/SUM(landing)::numeric,5) ELSE 0 END as cr_postback,"+
			"ROUND(AVG(cpa)::numeric,5) as avg_cpa").Row().Scan(
			&TotalSummaryCampaign.Landing,
			&TotalSummaryCampaign.MoReceived,
			&TotalSummaryCampaign.Postback,
			&TotalSummaryCampaign.PO,
			&TotalSummaryCampaign.CostPerConversion,
			&TotalSummaryCampaign.TotalWakiAgencyFee,
			&TotalSummaryCampaign.SBAF,
			&TotalSummaryCampaign.SAAF,
			&TotalSummaryCampaign.TechnicalFee,
			&TotalSummaryCampaign.WakiRevenue,
			&TotalSummaryCampaign.CrMO,
			&TotalSummaryCampaign.CrPostback,
			&TotalSummaryCampaign.ECPA)
	}

	r.Logs.Debug(fmt.Sprintf("Total data : %d ... \n", len(ss)))

	return ss, total_rows, TotalSummaryCampaign, rows.Err()
}

func (r *BaseModel) CreateCpaReport(s entity.SummaryCampaign) error {
	s.SuccessFP = 0
	s.PO = 0
	s.AgencyFee = 0
	s.TotalWakiAgencyFee = 0
	s.TechnicalFee = 0
	s.CPA = 0
	s.Traffic = 0
	s.CrPostback = 0
	s.CrMO = 0
	return r.DB.Create(&s).Error
}

func (r *BaseModel) UpdateCpaReport(s entity.SummaryCampaign) error {
	if err := r.DB.Model(&entity.SummaryCampaign{}).Where("id = ?", s.ID).Updates(&s).Error; err != nil {
		return err
	}

	// Override field default (misalnya 0) agar ikut terupdate
	resetFields := map[string]interface{}{
		"success_fp":            0,
		"po":                    0,
		"agency_fee":            0,
		"total_waki_agency_fee": 0,
		"technical_fee":         0,
		"cpa":                   0,
		"traffic":               0,
		"cr_postback":           0,
		"cr_mo":                 0,
	}

	return r.DB.Model(&entity.SummaryCampaign{}).Where("id = ?", s.ID).Updates(resetFields).Error
}

func (r *BaseModel) FindSummaryCampaignByUniqueKey(
	summaryDate *time.Time,
	campaignId, country, operator, partner, service, adnet, urlServiceKey string,
) (entity.SummaryCampaign, error) {
	db := r.DB.Model(&entity.SummaryCampaign{})
	if summaryDate != nil {
		db = db.Where("summary_date = ?", *summaryDate)
	}
	if campaignId != "" {
		db = db.Where("LOWER(campaign_id) = LOWER(?)", campaignId)
	}
	if country != "" {
		db = db.Where("LOWER(country) = LOWER(?)", country)
	}
	if operator != "" {
		db = db.Where("LOWER(operator) = LOWER(?)", operator)
	}
	if partner != "" {
		db = db.Where("LOWER(partner) = LOWER(?)", partner)
	}
	if service != "" {
		db = db.Where("LOWER(service) = LOWER(?)", service)
	}
	if adnet != "" {
		db = db.Where("LOWER(adnet) = LOWER(?)", adnet)
	}
	if urlServiceKey != "" {
		db = db.Where("LOWER(url_service_key) = LOWER(?)", urlServiceKey)
	}

	var s entity.SummaryCampaign
	result := db.First(&s)
	return s, result.Error
}

func (r *BaseModel) FindLatestSummaryCampaignByUniqueKey(service, adnet, operator string) (entity.SummaryCampaign, error) {
	var s entity.SummaryCampaign
	result := r.DB.Where("LOWER(service) = LOWER(?) AND LOWER(adnet) = LOWER(?) AND LOWER(operator) = LOWER(?)", service, adnet, operator).
		Order("summary_date DESC").
		First(&s)
	return s, result.Error
}

func (r *BaseModel) GetDisplayMainstreamReport(o entity.DisplayCPAReport, allowedCompanies []string) ([]entity.SummaryCampaign, error) {
	var rows *sql.Rows
	var err error

	query := r.DB.Model(&entity.SummaryCampaign{})
	// fmt.Println(query)
	query = query.Where("campaign_objective = ?", "MAINSTREAM")
	query = query.Where("mo_received > 0")
	query = query.Where("company IN ?", allowedCompanies)
	if o.Action == "Search" {
		if o.CampaignId != "" {
			query = query.Where("LOWER(campaign_id) = LOWER(?)", o.CampaignId)
		}
		if o.Agency != "" {
			query = query.Where("LOWER(adnet) = LOWER(?)", o.Agency)
		}
		if o.UrlServiceKey != "" {
			query = query.Where("LOWER(url_service_key) = LOWER(?)", o.UrlServiceKey)
		}
		if o.Country != "" {
			query = query.Where("LOWER(country) = LOWER(?)", o.Country)
		}
		if o.Company != "" {
			query = query.Where("LOWER(company) = LOWER(?)", o.Company)
		}
		if o.ClientType != "" {
			query = query.Where("LOWER(client_type) = LOWER(?)", o.ClientType)
		}
		if o.Operator != "" {
			query = query.Where("LOWER(operator) = LOWER(?)", o.Operator)
		}
		if o.CampaignName != "" {
			query = query.Where("LOWER(campaign_name) = LOWER(?)", o.CampaignName)
		}
		if o.Channel != "" {
			query = query.Where("LOWER(channel) = LOWER(?)", o.Channel)
		}
		if o.Partner != "" {
			query = query.Where("LOWER(partner) = LOWER(?)", o.Partner)
		}
		if o.Adnet != "" { //Publisher
			query = query.Where("LOWER(adnet) = LOWER(?)", o.Adnet)
		}
		if o.Service != "" {
			query = query.Where("LOWER(service) = LOWER(?0", o.Service)
		}
		if o.DateRange != "" {
			switch strings.ToUpper(o.DateRange) {
			case "TODAY":
				query = query.Where("summary_date = CURRENT_DATE")
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
				query = query.Where("summary_date BETWEEN ? AND ?", o.DateBefore, o.DateAfter)
			case "ALLDATERANGE":
			default:
				query = query.Where("summary_date = ?", o.DateRange)
			}
		} else {
			query = query.Where("summary_date = CURRENT_DATE")
		}

		if o.DataBasedOn != "" {
			switch strings.ToUpper(o.DataBasedOn) {
			case "HIGHEST_PIXEL_RECEIVED":
				query = query.Order("mo_received DESC")
			case "HIGHEST_PIXEL_SEND":
				query = query.Order("postback DESC")
			case "HIGHEST_PRICE_PER_POSTBACK":
				query = query.Order("po DESC")
			case "HIGHEST_COST_PER_CONVERSION":
				query = query.Order("cost_per_conversion DESC")
			case "HIGHEST_AGENCY_FEE":
				query = query.Order("agency_fee DESC")
			case "HIGHEST_SPENDING_TO_ADNETS":
				query = query.Order("sbaf DESC")
			case "HIGHEST_TOTAL_WAKI_AGENCY_FEE":
				query = query.Order("total_waki_agency_fee DESC")
			case "HIGHEST_TOTAL_SPENDING":
				query = query.Order("saaf DESC")
			case "HIGHEST_ECPA":
				query = query.Order("cpa DESC")
			case "HIGHEST_LANDING":
				query = query.Order("traffic DESC")
			case "HIGHEST_POSTBACK":
				query = query.Order("cr_postback DESC")
			case "HIGHEST_MO":
				query = query.Order("cr_mo DESC")
			case "LOWEST_PIXEL_RECEIVED":
				query = query.Order("pixel_received ASC")
			case "LOWEST_PIXEL_SEND":
				query = query.Order("postback ASC")
			case "LOWEST_PRICE_PER_POSTBACK":
				query = query.Order("po ASC")
			case "LOWEST_COST_PER_CONVERSION":
				query = query.Order("cost_per_conversion ASC")
			case "LOWEST_AGENCY_FEE":
				query = query.Order("agency_fee ASC")
			case "LOWEST_SPENDING_TO_ADNETS":
				query = query.Order("sbaf ASC")
			case "LOWEST_TOTAL_WAKI_AGENCY_FEE":
				query = query.Order("total_waki_agency_fee ASC")
			case "LOWEST_TOTAL_SPENDING":
				query = query.Order("saaf ASC")
			case "LOWEST_ECPA":
				query = query.Order("cpa ASC")
			case "LOWEST_LANDING":
				query = query.Order("traffic ASC")
			case "LOWEST_POSTBACK":
				query = query.Order("cr_postback ASC")
			case "LOWEST_MO":
				query = query.Order("cr_mo ASC")
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
	result := r.DB.Exec("UPDATE summary_campaigns SET po = ? WHERE ID = ?", o.Postback, id)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) UpdateAgencyCostModel(o entity.SummaryCampaign) error {
	result := r.DB.Exec("UPDATE summary_campaigns SET agency_fee = COALESCE(?, agency_fee), cost_per_conversion = COALESCE(?, cost_per_conversion) WHERE DATE_TRUNC('month', summary_date) = DATE_TRUNC('month', CURRENT_DATE) AND summary_date >= DATE_TRUNC('month', CURRENT_DATE) + INTERVAL '0 DAY'", o.AgencyFee, o.CostPerConversion)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error

}

func (r *BaseModel) GetDisplayCostReport(o entity.DisplayCostReport, allowedAdnets []string) ([]entity.CostReport, int64, error) {
	var (
		rows       *sql.Rows
		err        error
		total_rows int64
	)

	query := r.DB.Model(&entity.SummaryCampaign{})
	query = query.Where("mo_received > 0")
	query = query.Where("adnet IN ?", allowedAdnets)

	if o.Action == "Search" {

		if o.Adnet != "" {
			query = query.Where("adnet = ?", o.Adnet)
		}

		if o.DateRange != "" {
			switch strings.ToUpper(o.DateRange) {
			case "TODAY":
				query = query.Where("summary_date = CURRENT_DATE")
			case "YESTERDAY":
				query = query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '1 DAY' AND CURRENT_DATE")
			case "LAST 7 DAYS":
				query = query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
			case "LAST 30 DAYS":
				query = query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '30 DAY' AND CURRENT_DATE")
			case "THIS MONTH":
				query = query.Where("summary_date >= DATE_TRUNC('month', CURRENT_DATE)")
			case "LAST MONTH":
				query = query.Where("summary_date BETWEEN DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 DAY'")
			case "CUSTOMRANGE":
				query = query.Where("summary_date BETWEEN ? AND ?", o.DateBefore, o.DateAfter)
			case "ALLDATERANGE":
			default:
				query = query.Where("summary_date = ?", o.DateRange)
			}
		} else {
			// query = query.Where("summary_date = CURRENT_DATE")
		}

		rows, err = query.Select(`
			adnet,
			SUM(postback) as conversion1,
			SUM(sbaf) as cost1,
			NULL as conversion2,
			NULL as cost2
		`).Group("adnet").
			Rows()

		if o.DataBasedOn != "" {
			switch strings.ToUpper(o.DataBasedOn) {
			case "HIGHEST CONVERSION S2S":
				query = query.Order("SUM(postback) DESC")
			case "LOWEST CONVERSION S2S":
				query = query.Order("SUM(postback) ASC")
			case "HIGHEST COST S2S":
				query = query.Order("SUM(sbaf) DESC")
			case "LOWEST COST S2S":
				query = query.Order("SUM(sbaf) ASC")
			default:
				query = query.Order("SUM(postback) DESC")
			}
		}

		if err != nil {
			return nil, 0, err
		}
		if rows == nil {
			return []entity.CostReport{}, 0, nil
		}

	} else {
		rows, err = query.Select(`
			adnet,
			SUM(postback) as conversion1,
			SUM(sbaf) as cost1,
			NULL as conversion2,
			NULL as cost2
		`).Group("adnet").
			Rows()

		if err != nil {
			return nil, 0, err
		}
		if rows == nil {
			return []entity.CostReport{}, 0, nil
		}

	}

	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Rows()

	defer rows.Close()
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

	return results, total_rows, rows.Err()
}
func (r *BaseModel) GetDisplayCostReportDetail(o entity.DisplayCostReport) ([]entity.CostReport, int64, error) {
	var (
		rows       *sql.Rows
		err        error
		total_rows int64
	)

	query := r.DB.Model(&entity.SummaryCampaign{})
	query = query.Where("mo_received > 0")

	if o.Action == "Search" {

		if o.Adnet != "" {
			query = query.Where("adnet = ?", o.Adnet)
		}

		if o.Country != "" {
			query = query.Where("country = ?", o.Country)
		}

		if o.CampaignType != "" {
			switch strings.ToUpper(o.CampaignType) {
			case "S2S":
				query = query.Where("campaign_objective IN ('CPA', 'CPC', 'CPM')")
			case "MAINSTREAM":
				query = query.Where("campaign_objective ? ", "MAINSTREAM")
			case "API":
				query = query.Where("campaign_objective ?", "API")
			default:
				query = query.Where("campaign_objective = ?", o.CampaignType)
			}
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
			adnet,
			summary_date,
			country,
			operator,
			SUM(traffic) as landing,
			CASE WHEN SUM(traffic)=0 THEN 0 ELSE SUM(postback)/SUM(traffic) END as cr_postback,
			CASE WHEN SUM(traffic)=0 THEN 0 ELSE SUM(mo_received)/SUM(traffic) END as cr_mo,
			short_code,
			url_after,
			SUM(postback) as conversion1,
			SUM(mo_received) as mo_conversion1,
			SUM(sbaf) as cost1,
			NULL as conversion2,
			NULL as cost2
		`).Group("summary_date, adnet, country, operator, short_code, url_after").Order("summary_date ASC").
			Rows()

		if o.DataBasedOn != "" {
			switch strings.ToUpper(o.DataBasedOn) {
			case "HIGHEST CONVERSION S2S":
				query = query.Order("SUM(postback) DESC")
			case "LOWEST CONVERSION S2S":
				query = query.Order("SUM(postback) ASC")
			case "HIGHEST COST S2S":
				query = query.Order("SUM(sbaf) DESC")
			case "LOWEST COST S2S":
				query = query.Order("SUM(sbaf) ASC")
			default:
				query = query.Order("SUM(postback) DESC")
			}
		}
		if err != nil {
			return nil, 0, err
		}
		if rows == nil {
			return []entity.CostReport{}, 0, nil
		}

	} else {

		if o.Adnet != "" {
			query = query.Where("adnet = ?", o.Adnet)
		}

		rows, err = query.Select(`
			adnet,
			summary_date,
			country,
			operator,
			SUM(traffic) as landing,
			CASE WHEN SUM(traffic)=0 THEN 0 ELSE SUM(postback)/SUM(traffic) END as cr_postback,
			CASE WHEN SUM(traffic)=0 THEN 0 ELSE SUM(mo_received)/SUM(traffic) END as cr_mo,
			short_code,
			url_after,
			SUM(postback) as conversion1,
			SUM(mo_received) as mo_conversion1,
			SUM(sbaf) as cost1,
			NULL as conversion2,
			NULL as cost2
		`).Group("summary_date, adnet, country, operator, adn, url_after").
			Rows()

		if err != nil {
			return nil, 0, err
		}
	}

	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}
	rows, _ = query_limit.Rows()

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

	return results, total_rows, rows.Err()
}
