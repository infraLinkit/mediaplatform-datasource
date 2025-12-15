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

	t_query := r.DB.Model(&entity.SummaryCampaign{})

	query := r.DB.Model(&entity.SummaryCampaign{}).Select(`
		summary_campaigns.*,
		CASE 
			WHEN LOWER(client_type) = 'external' 
				THEN (mo_received * poaf) 
			ELSE (total_waki_agency_fee + (po * postback) + technical_fee) 
		END AS saaf,
		(po * postback) AS sbaf,
		(
			CASE 
				WHEN LOWER(client_type) = 'external' 
					THEN (mo_received * poaf) 
				ELSE (total_waki_agency_fee + (po * postback) + technical_fee) 
			END - (po * postback)
		) AS revenue
	`).Where("mo_received > 0").Where("company IN ?", allowedCompanies)
	t_query.Where("mo_received > 0").Where("company IN ?", allowedCompanies)

	if o.CampaignObjective != "" {
		query.Where("campaign_objective = ? ", o.CampaignObjective)
		t_query.Where("campaign_objective = ? ", o.CampaignObjective)
	} else {
		query.Where("campaign_objective = ? OR campaign_objective = ?", "CPA", "UPLOAD SMS")
		t_query.Where("campaign_objective = ? OR campaign_objective = ?", "CPA", "UPLOAD SMS")
	}

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
			t_query = t_query.Where("country = ?", o.Country)
		}
		if o.Company != "" {
			query = query.Where("company = ?", o.Company)
			t_query = t_query.Where("country = ?", o.Country)
		}
		if o.ClientType != "" {
			query = query.Where("client_type = ?", o.ClientType)
			t_query = t_query.Where("client_type = ?", o.ClientType)
		}
		if o.Operator != "" {
			query = query.Where("operator = ?", o.Operator)
			t_query = t_query.Where("operator = ?", o.Operator)
		}
		if o.CampaignName != "" {
			query = query.Where("campaign_name = ?", o.CampaignName)
			t_query = t_query.Where("campaign_name = ?", o.CampaignName)
		}
		if o.Partner != "" {
			query = query.Where("partner = ?", o.Partner)
			t_query = t_query.Where("partner = ?", o.Partner)
		}
		if len(o.Adnets) > 0 {
			query = query.Where("adnet IN ?", o.Adnets)
			t_query = t_query.Where("adnet IN ?", o.Adnets)
		}
		if o.Service != "" {
			query = query.Where("service = ?", o.Service)
			t_query = t_query.Where("service = ?", o.Service)
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
	} else {
		t_query = t_query.Where("summary_date = CURRENT_DATE")
	}

	if o.OrderColumn != "" {
		dir := "ASC"
		if strings.ToUpper(o.OrderDir) == "DESC" {
			dir = "DESC"
		}

		switch o.OrderColumn {
		case "saaf":
			query = query.Order(fmt.Sprintf(`
				CASE 
					WHEN LOWER(client_type) = 'external' 
						THEN (mo_received * poaf) 
					ELSE (total_waki_agency_fee + (po * postback) + technical_fee) 
				END %s
			`, dir))
		case "sbaf":
			query = query.Order(fmt.Sprintf("(po * postback) %s", dir))
		case "revenue":
			query = query.Order(fmt.Sprintf(`
				(
					CASE 
						WHEN LOWER(client_type) = 'external' 
							THEN (mo_received * poaf) 
						ELSE (total_waki_agency_fee + (po * postback) + technical_fee) 
					END - (po * postback)
				) %s
			`, dir))
		default:
			query = query.Order(fmt.Sprintf("%s %s", o.OrderColumn, dir))
		}
	} else {
		query = query.Order("summary_date DESC").Order("id DESC")
	}

	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, err = query_limit.Rows()
	if err != nil {
		return []entity.SummaryCampaign{}, 0, entity.TotalSummaryCampaign{}, err
	}
	defer rows.Close()

	var ss []entity.SummaryCampaign
	for rows.Next() {
		var s entity.SummaryCampaign
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	// GET TOTAL HEADER
	/*summary_campaigns.*,
	CASE
		WHEN LOWER(client_type) = 'external'
			THEN (mo_received * poaf)
		ELSE (total_waki_agency_fee + (po * postback) + technical_fee)
	END AS saaf,
	(po * postback) AS sbaf,
	(
		CASE
			WHEN LOWER(client_type) = 'external'
				THEN (mo_received * poaf)
			ELSE (total_waki_agency_fee + (po * postback) + technical_fee)
		END - (po * postback)
	) AS revenue*/
	if total_rows > 0 {
		// COUNT THE SUMMARIZE
		_ = t_query.Select(
			`SUM(landing) as landing,
			 SUM(mo_received) as mo_received,
			 SUM(postback) as postback,
			 AVG(po) as price_per_postback,
			 SUM(cost_per_conversion) as cost_per_conversion,
			 SUM(total_waki_agency_fee) as total_waki_agency_fee,
			 SUM(po * postback) as spending_to_adnet, --SBAF
			 SUM(CASE
					WHEN LOWER(client_type) = 'external'
						THEN (mo_received * poaf)
					ELSE (total_waki_agency_fee + (po * postback) + technical_fee)
				END) as spending, --SAAF
			 SUM(technical_fee) as technical_fee,
			 SUM(CASE
					WHEN LOWER(client_type) = 'external'
						THEN (mo_received * poaf)
					ELSE (total_waki_agency_fee + (po * postback) + technical_fee)
				END) - SUM(po * postback) as waki_revenue, -- SAAF - SBAF
			 CASE WHEN SUM(landing)>0 THEN ROUND(SUM(mo_received)/SUM(landing)::numeric,5) ELSE 0 END as cr_mo,
			 CASE WHEN SUM(landing)>0 THEN ROUND(SUM(postback)/SUM(landing)::numeric,5) ELSE 0 END as cr_postback,
			 ROUND(AVG(cpa)::numeric,5) as avg_cpa`).Row().Scan(
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

func (r *BaseModel) FindLatestSummaryCampaignByUniqueKey(service, adnet, operator string, partner string) (entity.SummaryCampaign, error) {
	var s entity.SummaryCampaign
	result := r.DB.Where("LOWER(service) = LOWER(?) AND LOWER(adnet) = LOWER(?) AND LOWER(operator) = LOWER(?) AND LOWER(partner)=LOWER(?)", service, adnet, operator, partner).
		Order("summary_date DESC").
		First(&s)
	return s, result.Error
}

func (r *BaseModel) GetDisplayMainstreamReport(o entity.DisplayCPAReport, allowedCompanies []string) ([]entity.SummaryCampaign, int64, error) {
	var rows *sql.Rows
	var err error
	var total_rows int64

	query := r.DB.Model(&entity.SummaryCampaign{}).Select(`
		summary_campaigns.*,
		(poaf * postback) AS saaf,
		(po * postback) AS sbaf,
		(CASE WHEN mo_received > 0 THEN (poaf * postback) / mo_received ELSE 0 END) AS price_per_mo,
		((poaf * postback) - (po * postback)) AS revenue
	`).Where("campaign_objective = ?", "MAINSTREAM").
		Where("mo_received > 0").
		Where("company IN ?", allowedCompanies)

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
		if len(o.Adnets) > 0 {
			query = query.Where("adnet IN ?", o.Adnets)
		}
		if o.Service != "" {
			query = query.Where("LOWER(service) = LOWER(?)", o.Service)
		}

		// Date range filter
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
				// no filter
			default:
				query = query.Where("summary_date = ?", o.DateRange)
			}
		} else {
			query = query.Where("summary_date = CURRENT_DATE")
		}
	}

	if o.OrderColumn != "" {
		dir := "ASC"
		if strings.ToUpper(o.OrderDir) == "DESC" {
			dir = "DESC"
		}

		switch o.OrderColumn {
		case "saaf":
			query = query.Order(fmt.Sprintf("(poaf * postback) %s", dir))
		case "sbaf":
			query = query.Order(fmt.Sprintf("(po * postback) %s", dir))
		case "price_per_mo":
			query = query.Order(fmt.Sprintf("(CASE WHEN mo_received > 0 THEN (poaf * postback) / mo_received ELSE 0 END) %s", dir))
		case "revenue":
			query = query.Order(fmt.Sprintf("((poaf * postback) - (po * postback)) %s", dir))
		default:
			query = query.Order(fmt.Sprintf("%s %s", o.OrderColumn, dir))
		}
	} else {
		query = query.Order("summary_date DESC").Order("id DESC")
	}

	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, err = query_limit.Rows()
	if err != nil {
		return []entity.SummaryCampaign{}, 0, err
	}
	defer rows.Close()

	var ss []entity.SummaryCampaign
	for rows.Next() {
		var s entity.SummaryCampaign
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	r.Logs.Debug(fmt.Sprintf("Total data : %d ... \n", len(ss)))
	return ss, total_rows, rows.Err()
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
func (r *BaseModel) GetDisplayCostReportDetail(o entity.DisplayCostReport, allowedAdnets []string) ([]entity.CostReport, int64, error) {
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
			//today := time.Now()
			switch o.DateRange {
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
			case "ALLDATERANGE":
			default:
				query = query.Where("summary_date = ?", o.DateRange)
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

func (r *BaseModel) GetSummaryReportById(id []string) ([]entity.SummaryCampaign, error) {
	query := r.DB.Model(&entity.SummaryCampaign{}).Where("mo_received > 0 AND id IN ?", id)
	rows, err := query.Rows()

	if err != nil {
		return []entity.SummaryCampaign{}, err
	}

	var results []entity.SummaryCampaign
	for rows.Next() {
		var s entity.SummaryCampaign
		r.DB.ScanRows(rows, &s)
		results = append(results, s)
	}

	return results, nil
}
