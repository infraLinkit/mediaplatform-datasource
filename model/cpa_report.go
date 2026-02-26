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
	//fmt.Println("Company: ", o.Company)
	query := r.DB.Model(&entity.SummaryCampaign{}).Select(`
		summary_campaigns.*,
		CASE WHEN campaign_objective='UPLOAD SMS' THEN saaf
			 ELSE
			 (
				CASE 
					WHEN LOWER(client_type) = 'external' 
						THEN (mo_received * poaf) 
					ELSE (total_waki_agency_fee + (po * postback) + technical_fee) 
				END
			)	 
		END AS saaf,
		CASE WHEN campaign_objective='UPLOAD SMS' THEN sbaf
			 ELSE (po * postback)
		END AS sbaf,
		CASE WHEN campaign_objective='UPLOAD SMS' THEN saaf-sbaf
			 ELSE
			 (
				CASE 
					WHEN LOWER(client_type) = 'external' 
						THEN (mo_received * poaf) 
					ELSE (total_waki_agency_fee + (po * postback) + technical_fee) 
				END - (po * postback)
			 )
		END AS revenue
	`).Where("mo_received > 0").Where("company IN ?", allowedCompanies)
	t_query.Where("mo_received > 0").Where("company IN ?", allowedCompanies)

	if o.CampaignObjective != "" {
		query.Where("campaign_objective = ? ", o.CampaignObjective)
		t_query.Where("campaign_objective = ? ", o.CampaignObjective)
	} else {
		query.Where("campaign_objective IN ?", []string{"CPA", "UPLOAD SMS", "SINGLE URL S2S"})
		t_query.Where("campaign_objective IN ?", []string{"CPA", "UPLOAD SMS", "SINGLE URL S2S"})
	}

	//fmt.Println("o.CampaignObjective: ", o.CampaignObjective)

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
			t_query = t_query.Where("company = ?", o.Company)
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
	if total_rows > 0 {
		// COUNT THE SUMMARIZE
		_ = t_query.Select(
			`SUM(landing) as landing,
			 SUM(mo_received) as mo_received,
			 SUM(postback) as postback,
			 AVG(po) as price_per_postback,
			 SUM(cost_per_conversion) as cost_per_conversion,
			 SUM(total_waki_agency_fee) as total_waki_agency_fee,
			 SUM(CASE WHEN campaign_objective='UPLOAD SMS' THEN sbaf 
			 		  ELSE po * postback
				 END) as spending_to_adnet, -- sbaf
			 SUM(CASE WHEN campaign_objective='UPLOAD SMS' THEN saaf
			 		  ELSE
						CASE
							WHEN LOWER(client_type) = 'external'
								THEN (mo_received * poaf)
							ELSE (total_waki_agency_fee + (po * postback) + technical_fee)
						END
				 END) as spending, --saaf
			 SUM(technical_fee) as technical_fee,
			 SUM(
			     CASE WHEN campaign_objective='UPLOAD SMS' THEN saaf-sbaf
			     ELSE
					CASE
						WHEN LOWER(client_type) = 'external'
							THEN (mo_received * poaf)
						ELSE (total_waki_agency_fee + (po * postback) + technical_fee)
					END
				 END) - 
				 SUM(CASE WHEN campaign_objective='UPLOAD SMS' THEN 0
				    	 ELSE po * postback
					 END) as waki_revenue,
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
	campaignId, country, operator, partner, service, adnet, urlServiceKey, campaignObjective string,
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
	if campaignObjective != "" {
		db = db.Where("LOWER(campaign_objective) = LOWER(?)", campaignObjective)
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
	`).Where("campaign_objective LIKE ?", "%MAINSTREAM%").
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

	var apiAdnets []string

	_ = r.DB.Model(&entity.ApiPinReport{}).
		Distinct("adnet").Pluck("adnet", &apiAdnets)

	allowedAdnets = append(allowedAdnets, apiAdnets...)

	//query := r.DB.Model(&entity.SummaryCampaign{})
	query := r.DB.Table(`
	   (select date_send as summary_date,adnet,
		SUM(total_mo) as mo_received,
		sum(total_postback) as conversion,
		sum(sbaf) as "cost",
		sum(saaf) as saaf,'api' as type
		from api_pin_reports WHERE total_mo > 0 group by adnet,date_send
		UNION
		select summary_date as summary_date,
		adnet,
		SUM(mo_received) as mo_received,
		SUM(postback) as conversion,
		SUM(CASE WHEN campaign_objective='UPLOAD SMS' THEN sbaf
			ELSE po * postback END) as cost,
		SUM(CASE WHEN campaign_objective='UPLOAD SMS' THEN saaf
			ELSE CASE WHEN LOWER(client_type) = 'external' THEN (mo_received * poaf)
			ELSE
				(total_waki_agency_fee + (po * postback) + technical_fee)
		    END
		END) as saaf,'s2s' as "type"
		FROM "summary_campaigns" WHERE mo_received > 0
		AND "summary_campaigns"."deleted_at" IS NULL GROUP BY summary_date,"adnet")
		as t`)

	if len(o.Adnets) > 0 {
		query = query.Where("adnet IN ?", o.Adnets)
	} else {
		query = query.Where("adnet IN ?", allowedAdnets)
	}

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
            SUM(CASE WHEN type = 's2s' THEN conversion ELSE 0 END) AS conversion1,
			SUM(CASE WHEN type = 's2s' THEN cost ELSE 0 END) AS cost1,
            SUM(CASE WHEN type = 'api' THEN conversion ELSE 0 END) AS conversion2,
            SUM(CASE WHEN type = 'api' THEN cost ELSE 0 END) AS cost2,
			SUM(CASE WHEN type = 's2s' THEN saaf ELSE 0 END) AS saaf1,
            SUM(CASE WHEN type = 'api' THEN saaf ELSE 0 END) AS saaf2
		`).Group("adnet").
			Rows()

		if o.DataBasedOn != "" {
			switch strings.ToUpper(o.DataBasedOn) {
			case "HIGHEST CONVERSION S2S":
				query = query.Order("SUM(CASE WHEN type = 's2s' THEN conversion ELSE 0 END) DESC")
			case "LOWEST CONVERSION S2S":
				query = query.Order("SUM(CASE WHEN type = 's2s' THEN conversion ELSE 0 END) ASC")
			case "HIGHEST COST S2S":
				query = query.Order("SUM(CASE WHEN type = 's2s' THEN cost ELSE 0 END) DESC")
			case "LOWEST COST S2S":
				query = query.Order("SUM(CASE WHEN type = 's2s' THEN cost ELSE 0 END) ASC")

			case "HIGHEST CONVERSION API":
				query = query.Order("SUM(CASE WHEN type = 'api' THEN conversion ELSE 0 END) DESC")
			case "LOWEST CONVERSION API":
				query = query.Order("SUM(CASE WHEN type = 'api' THEN conversion ELSE 0 END) ASC")
			case "HIGHEST COST API":
				query = query.Order("SUM(CASE WHEN type = 'api' THEN cost ELSE 0 END) DESC")
			case "LOWEST COST API":
				query = query.Order("SUM(CASE WHEN type = 'api' THEN cost ELSE 0 END) ASC")
			default:
				query = query.Order("SUM(CASE WHEN type = 's2s' THEN conversion ELSE 0 END) DESC")
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
            SUM(CASE WHEN type = 's2s' THEN conversion ELSE 0 END) AS conversion1,
			SUM(CASE WHEN type = 's2s' THEN cost ELSE 0 END) AS cost1,
            SUM(CASE WHEN type = 'api' THEN conversion ELSE 0 END) AS conversion2,
            SUM(CASE WHEN type = 'api' THEN cost ELSE 0 END) AS cost2,
			SUM(CASE WHEN type = 's2s' THEN saaf ELSE 0 END) AS saaf1,
            SUM(CASE WHEN type = 'api' THEN saaf ELSE 0 END) AS saaf2
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
				query = query.Where("campaign_objective IN ('CPA', 'UPLOAD SMS', 'SINGLE URL S2S)")
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
			SUM(CASE WHEN campaign_objective='UPLOAD SMS' THEN sbaf
				ELSE po * postback END) as cost1,
			NULL as conversion2,
			NULL as cost2,
			SUM(CASE WHEN campaign_objective='UPLOAD SMS' THEN saaf
				ELSE CASE WHEN LOWER(client_type) = 'external' THEN (mo_received * poaf)
					 ELSE
					 	(total_waki_agency_fee + (po * postback) + technical_fee)
					 END 
				END) as saaf1,
			NULL as saaf2
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
			SUM(CASE WHEN campaign_objective='UPLOAD SMS' THEN sbaf
				ELSE po * postback END) as cost1,
			NULL as conversion2,
			NULL as cost2,
			SUM(CASE WHEN campaign_objective='UPLOAD SMS' THEN saaf
				ELSE CASE WHEN LOWER(client_type) = 'external' THEN (mo_received * poaf)
					 ELSE
					 	(total_waki_agency_fee + (po * postback) + technical_fee)
					 END 
				END) as saaf1,
			NULL as saaf2
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

func (r *BaseModel) AddSMSReport(s entity.SummaryCampaign) error {
	SQL := `
	INSERT INTO summary_campaigns 
	(created_at, updated_at, status, summary_date, url_service_key, campaign_id, campaign_name, country, operator,
	partner, aggregator, adnet, service, short_code, traffic, landing, mo_received, cr, postback,
	total_fp, success_fp, billrate, roi, po, first_push, cost, sbaf, saaf, cpa, revenue,
	url_after, url_before, mo_limit, ratio_send, ratio_receive, company, client_type,
	cost_per_conversion, agency_fee, target_daily_budget, cr_mo, cr_postback, total_waki_agency_fee,
	budget_usage, target_daily_budget_changes, technical_fee, campaign_objective, 
	channel, price_per_mo, target_monthly_budget, poaf)
	SELECT NOW(), NOW(), true, ?,
	cd.url_service_key, cd.campaign_id, cp.name, cd.country, 
	?, ?, cd.aggregator, ?, ?, ?, 0, 0, ?, -- mo_received
	0, ?, -- postback
	0, 0, 0, 0, ?, -- po
	0, -- first_push
	0, -- cost
	?, -- sbaf
	?, -- saaf
	0, 0, cd.url_landing, cd.url_landing, ?, ?,
	?, -- ratio_receive 
	pt.company, pt.client_type, 0, 0, 0, 0, 0, 0, 0, 0, 0, 'UPLOAD SMS',
	'NA', 0, 0, 0
	from campaign_details cd 
	left join partners as pt on pt.name=cd.partner 
	left join campaigns as cp on cp.id = cd.campaign_id::INTEGER where 
	cd.url_service_key = ?
	ON CONFLICT (summary_date,url_service_key,campaign_id,country,operator,partner,adnet,service,campaign_objective) 
	DO UPDATE SET 
		updated_at=NOW(),
		po = EXCLUDED.po,
		sbaf = EXCLUDED.sbaf,
		saaf = EXCLUDED.saaf,
		ratio_send = EXCLUDED.ratio_send,
		ratio_receive = EXCLUDED.ratio_receive,
		mo_received = EXCLUDED.mo_received,
		postback = EXCLUDED.postback `

	q := r.DB.Exec(SQL, s.SummaryDate, s.Operator, s.Partner, s.Adnet,
		s.Service, s.ShortCode, s.MoReceived, s.Postback,
		s.PO, s.SBAF, s.SAAF, 500, s.RatioSend, s.RatioReceive,
		s.URLServiceKey)
	return q.Error
}
