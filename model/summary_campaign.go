package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/infraLinkit/mediaplatform-datasource/entity"
	"gorm.io/gorm/clause"
)

func (r *BaseModel) DelSummaryCampaign(o entity.SummaryCampaign) error {

	result := r.DB.
		Where("summary_date = ? AND url_service_key = ? AND country = ? AND operator = ? AND partner = ? AND service = ? AND adnet = ? AND campaign_id = ?", o.SummaryDate, o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId).
		Delete(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) EditSettingSummaryCampaign(o entity.SummaryCampaign) error {

	result := r.DB.Model(&o).
		Where("summary_date = ? AND url_service_key = ? AND country = ? AND operator = ? AND partner = ? AND service = ? AND adnet = ? AND campaign_id = ?", o.SummaryDate, o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId).
		Updates(entity.SummaryCampaign{PO: o.PO, MOLimit: o.MOLimit, RatioSend: o.RatioSend, RatioReceive: o.RatioReceive})

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) UpdateSummaryCampaign(o entity.SummaryCampaign) error {

	result := r.DB.Model(&o).
		Where("summary_date = ? AND url_service_key = ? AND country = ? AND operator = ? AND partner = ? AND service = ? AND adnet = ? AND campaign_id = ?", o.SummaryDate, o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId).
		Updates(entity.SummaryCampaign{Status: o.Status})

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) SummaryCampaign(o entity.SummaryCampaign) int {

	result := r.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "summary_date"},
			{Name: "campaign_id"},
			{Name: "country"},
			{Name: "partner"},
			{Name: "operator"},
			{Name: "url_service_key"},
			{Name: "service"},
			{Name: "adnet"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"traffic":               o.Traffic,
			"landing":               o.Landing,
			"mo_received":           o.MoReceived,
			"cr_mo":                 o.CrMO,
			"cr_postback":           o.CrPostback,
			"postback":              o.Postback,
			"total_fp":              o.TotalFP,
			"success_fp":            o.SuccessFP,
			"billrate":              o.Billrate,
			"po":                    o.PO,
			"sbaf":                  o.SBAF,
			"saaf":                  o.SAAF,
			"cpa":                   o.CPA,
			"revenue":               o.Revenue,
			"url_after":             o.URLAfter,
			"url_before":            o.URLBefore,
			"mo_limit":              o.MOLimit,
			"ratio_send":            o.RatioSend,
			"ratio_receive":         o.RatioReceive,
			"client_type":           o.ClientType,
			"cost_per_conversion":   o.CostPerConversion,
			"agency_fee":            o.AgencyFee,
			"total_waki_agency_fee": o.TotalWakiAgencyFee,
			"target_daily_budget":   o.TargetDailyBudget,
			"budget_usage":          o.BudgetUsage,
			"campaign_name":         o.CampaignName}),
	}).Create(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return o.ID
}

func (r *BaseModel) DataTraffic(o entity.DataTraffic) int {

	result := r.DB.Create(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return int(o.ID)
}

func (r *BaseModel) DataLanding(o entity.DataLanding) int {

	result := r.DB.Create(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return int(o.ID)
}

func (r *BaseModel) DataClicked(o entity.DataClicked) int {

	result := r.DB.Create(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return int(o.ID)
}

func (r *BaseModel) DataRedirect(o entity.DataRedirect) int {

	result := r.DB.Create(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return int(o.ID)
}

func (r *BaseModel) UpdateCPAReportSummaryCampaign(o entity.SummaryCampaign) error {

	result := r.DB.Model(&o).
		Where("summary_date = ? AND url_service_key = ? AND country = ? AND operator = ? AND partner = ? AND service = ? AND adnet = ? AND campaign_id = ?", o.SummaryDate, o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId).
		Updates(entity.SummaryCampaign{CostPerConversion: o.CostPerConversion, AgencyFee: o.AgencyFee})

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error

}

func (r *BaseModel) UpdateReportSummaryCampaignMonitoringBudget(o entity.SummaryCampaign) error {

	result := r.DB.Model(&o).
		Where("summary_date = ? AND country = ? AND operator = ?", o.SummaryDate, o.Country, o.Operator).
		Updates(entity.SummaryCampaign{TargetDailyBudget: o.TargetDailyBudget})

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) GetSummaryCampaignMonitoring(filter entity.DisplayCampaignSummary) ([]entity.CampaignSummaryMonitoring, time.Time, time.Time, error) {
	var (
		rows *sql.Rows
	)
	query := r.DB.Model(&entity.CampaignSummaryMonitoring{})

	// Apply Indicator Selection
	selectedFields := []string{"summary_date", "country", "partner", "operator", "service", "adnet"}
	formattedIndicators := formatQueryIndicators(filter.DataIndicators, filter.DataType)
	selectedFields = append(selectedFields, formattedIndicators...)

	query.Select(selectedFields)

	// Set default values
	if filter.DateRange == "" {
		filter.DateRange = "this_month"
		query.Where("EXTRACT(MONTH FROM summary_date) = ?", int(time.Now().Month())).
			Where("EXTRACT(YEAR FROM summary_date) = ?", time.Now().Year())
	}
	if filter.DataType == "" {
		filter.DataType = "daily_report"
	}

	// Apply filters
	if filter.Country != "" {
		query.Where("country = ?", filter.Country)
	}
	if filter.Operator != "" {
		query.Where("operator = ?", filter.Operator)
	}
	if filter.Adnet != "" {
		query.Where("adnet = ?", filter.Adnet)
	}
	if filter.PartnerName != "" {
		query.Where("partner = ?", filter.PartnerName)
	}
	if filter.Service != "" {
		query.Where("service = ?", filter.Service)
	}

	// Handle Date Range
	var startDate, endDate time.Time
	now := time.Now()

	switch strings.ToLower(filter.DateRange) {
	case "today":
		startDate, endDate = now, now
	case "yesterday":
		startDate, endDate = now.AddDate(0, 0, -1), now.AddDate(0, 0, -1)
	case "last_7_days":
		startDate, endDate = now.AddDate(0, 0, -6), now
	case "last_30_days":
		startDate, endDate = now.AddDate(0, 0, -30), now
	case "this_month":
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endDate = time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 999999999, now.Location())
	case "last_month":
		startDate = time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location())
		endDate = time.Date(now.Year(), now.Month(), 0, 23, 59, 59, 999999999, now.Location())
	case "custom_range":
		if filter.CustomRange != "" {
			dates := strings.Split(filter.CustomRange, " - ")
			if len(dates) == 2 {
				startDate, _ = time.Parse("01/02/2006", strings.TrimSpace(dates[0]))
				endDate, _ = time.Parse("01/02/2006", strings.TrimSpace(dates[1]))
			}
		}
	}

	// Ensure end date is not in the future
	if endDate.After(now) {
		endDate = now
	}

	// Apply date range filter
	query.Where("summary_date BETWEEN ? AND ?", startDate, endDate)

	// Grouping for monthly reports
	if filter.DataType == "monthly_report" {
		query.Group("EXTRACT(YEAR FROM summary_date), EXTRACT(MONTH FROM summary_date), country, partner, operator, service, adnet")
	}

	rows, _ = query.Unscoped().Rows()

	defer rows.Close()

	var (
		ss []entity.CampaignSummaryMonitoring
	)

	for rows.Next() {

		var s entity.CampaignSummaryMonitoring

		// ScanRows scans a row into a struct
		r.DB.ScanRows(rows, &s)

		ss = append(ss, s)
	}
	return ss, startDate, endDate, rows.Err()

}

// helper
func formatQueryIndicators(selects []string, dataType string) []string {
	var formattedSelects []string

	for _, value := range selects {
		var formattedValue string

		if dataType == "monthly_report" {
			switch value {
			case "waki_revenue":
				formattedValue = "SUM(saaf - sbaf) AS waki_revenue"
			case "budget_usage":
				formattedValue = "SUM(CASE WHEN target_daily_budget = 0 THEN 0 ELSE (sbaf / target_daily_budget * 100) END) AS budget_usage"
			case "spending_to_adnets", "total_spending":
				formattedValue = fmt.Sprintf("SUM(%s) AS %s", value, value)
			case "fp":
				formattedValue = "SUM(first_push) AS fp"
			case "mo_sent":
				formattedValue = "SUM(postback) AS mo_sent"
			case "traffic":
				formattedValue = "SUM(landing) AS traffic"
			default:
				formattedValue = fmt.Sprintf("SUM(%s) AS %s", value, value)
			}
		} else { // Daily Report
			switch value {
			case "waki_revenue":
				formattedValue = "saaf - sbaf AS waki_revenue"
			case "budget_usage":
				formattedValue = "CASE WHEN target_daily_budget = 0 THEN NULL ELSE (sbaf / target_daily_budget * 100) END AS budget_usage, sbaf AS sbaf_t, target_daily_budget AS target_daily_budget_t"
			case "fp":
				formattedValue = "first_push AS fp"
			case "mo_sent":
				formattedValue = "postback AS mo_sent"
			case "spending_to_adnets":
				formattedValue = "sbaf AS spending_to_adnets"
			case "total_spending":
				formattedValue = "saaf AS total_spending"
			case "traffic":
				formattedValue = "landing AS traffic"
			default:
				formattedValue = fmt.Sprintf("%s AS %s", value, value)
			}
		}

		formattedSelects = append(formattedSelects, formattedValue)
	}

	return formattedSelects
}

/* package model

import (
	"context"
	"fmt"
	"strconv"

	"github.com/infraLinkit/mediaplatform-datasource/entity"

	_ "github.com/lib/pq"
)

const (
	DELSUMMARYCAMPAIGN         = "DELETE FROM summary_campaign WHERE DATE(summary_date) = '%s' AND urlservicekey = '%s' AND country = '%s' AND operator = '%s' AND partner = '%s' AND service = '%s' AND adnet = '%s' AND campaign_id = '%s'"
	EDITSETTINGSUMMARYCAMPAIGN = "UPDATE summary_campaign SET po = '%s', mo_limit = %d, ratio_send = %d, ratio_receive = %d WHERE DATE(summary_date) = '%s' AND urlservicekey = '%s' AND country = '%s' AND operator = '%s' AND partner = '%s' AND service = '%s' AND adnet = '%s' AND campaign_id = '%s'"
	UPDATESUMMARYCAMPAIGN      = "UPDATE summary_campaign SET status = %t WHERE DATE(summary_date) = '%s' AND urlservicekey = '%s' AND country = '%s' AND operator = '%s' AND partner = '%s' AND service = '%s' AND adnet = '%s' AND campaign_id = '%s'"
	SUMMARYCAMPAIGN            = "INSERT INTO summary_campaign AS sc (id, status, summary_date, campaign_id, campaign_name, country, partner, operator, urlservicekey, aggregator, service, adnet, short_code, traffic, landing, mo_received, cr_mo, cr_postback, postback, total_fp, success_fp, billrate, po, sbaf, saaf, cpa, revenue, url_after, url_before, mo_limit, ratio_send, ratio_receive, client_type, cost_per_conversion, agency_fee, total_waki_agency_fee, target_daily_budget, budget_usage) VALUES (DEFAULT, %t, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %d, %d, %d, '%s'::double precision, '%s'::double precision, %d, %d, '%s', '%s'::double precision, '%s'::double precision, '%s'::double precision, '%s'::double precision, '%s'::double precision, '%s'::double precision, '%s', '%s', %d, %d, %d, '%s', '%s'::double precision, '%s'::double precision, '%s'::double precision, '%s'::double precision, '%s'::double precision) ON CONFLICT (summary_date, campaign_id, country, partner, operator, urlservicekey, service, adnet) DO UPDATE SET traffic = %d, landing = %d, mo_received = %d, cr_mo = '%s'::double precision, cr_postback = '%s'::double precision, postback = %d, total_fp = %d, success_fp = '%s', billrate = '%s'::double precision, po = '%s'::double precision, sbaf = '%s'::double precision, saaf = '%s'::double precision, cpa = '%s'::double precision, revenue = '%s'::double precision, url_after = '%s', url_before = '%s', mo_limit = %d, ratio_send = %d, ratio_receive = %d, client_type = '%s', cost_per_conversion = '%s'::double precision, agency_fee = '%s'::double precision, total_waki_agency_fee = '%s'::double precision, target_daily_budget = '%s'::double precision, budget_usage = '%s'::double precision, campaign_name = '%s';"
	//SUMMARYCAMPAIGN                             = "INSERT INTO summary_campaign AS sc (id, status, summary_date, campaign_id, campaign_name, country, partner, operator, urlservicekey, aggregator, service, adnet, short_code, traffic, landing, mo_received, cr_mo, cr_postback, postback, total_fp, success_fp, billrate, po, sbaf, saaf, cpa, revenue, url_after, url_before, mo_limit, ratio_send, ratio_receive, client_type, cost_per_conversion, agency_fee, total_waki_agency_fee, target_daily_budget, budget_usage) VALUES (DEFAULT, %t, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %d, %d, %d, '%s'::double precision, '%s'::double precision, %d, %d, '%s', '%s'::double precision, '%s'::double precision, '%s'::double precision, '%s'::double precision, '%s'::double precision, '%s'::double precision, '%s', '%s', %d, %d, %d, '%s', '%s'::double precision, '%s'::double precision, '%s'::double precision, '%s'::double precision, '%s'::double precision) ON CONFLICT (summary_date, campaign_id, country, partner, operator, urlservicekey, service, adnet) DO UPDATE SET traffic = sc.traffic + %d, landing = sc.landing + %d, mo_received = sc.mo_received + %d, cr_mo = '%s'::double precision, cr_postback = '%s'::double precision, postback = sc.postback + %d, total_fp = sc.total_fp + %d, success_fp = '%s', billrate = '%s'::double precision, po = '%s'::double precision, sbaf = '%s'::double precision, saaf = '%s'::double precision, cpa = '%s'::double precision, revenue = '%s'::double precision, url_after = '%s', url_before = '%s', mo_limit = %d, ratio_send = %d, ratio_receive = %d, client_type = '%s', cost_per_conversion = '%s'::double precision, agency_fee = '%s'::double precision, total_waki_agency_fee = '%s'::double precision, target_daily_budget = '%s'::double precision, budget_usage = '%s'::double precision;"
	INSERTTRAFFIC                               = "INSERT INTO data_traffic (id, traffic_time, traffic_added_time, http_status, urlservicekey, campaign_id, country, partner, operator, aggregator, service, short_code, adnet, keyword, subkeyword, is_billable, plan) VALUES (DEFAULT, '%s', %d, %d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, '%s');"
	GETTOTALDATATRAFFIC                         = "SELECT COUNT(1) FROM data_traffic WHERE DATE(traffic_time) = '%s' AND urlservicekey = '%s' AND campaign_id = '%s' AND country = '%s' AND partner = '%s' AND operator = '%s' AND service = '%s' AND short_code = '%s' AND adnet = '%s'"
	INSERTLANDING                               = "INSERT INTO data_landing (id, landing_time, landed_time, http_status, urlservicekey, campaign_id, country, partner, operator, aggregator, service, short_code, adnet, keyword, subkeyword, is_billable, plan) VALUES (DEFAULT, '%s', %d, %d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, '%s');"
	GETTOTALDATALANDING                         = "SELECT COUNT(1) FROM data_landing WHERE DATE(landing_time) = '%s' AND urlservicekey = '%s' AND campaign_id = '%s' AND country = '%s' AND partner = '%s' AND operator = '%s' AND service = '%s' AND short_code = '%s' AND adnet = '%s'"
	INSERTCLICKED                               = "INSERT INTO data_clicked (id, clicked_time, clicked_button_time, http_status, urlservicekey, campaign_id, country, partner, operator, aggregator, service, short_code, adnet, keyword, subkeyword, is_billable, plan) VALUES (DEFAULT, '%s', %d, %d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, '%s');"
	GETTOTALDATACLICKED                         = "SELECT COUNT(1) FROM data_clicked WHERE DATE(clicked_time) = '%s' AND urlservicekey = '%s' AND campaign_id = '%s' AND country = '%s' AND partner = '%s' AND operator = '%s' AND service = '%s' AND short_code = '%s' AND adnet = '%s'"
	INSERTREDIRECT                              = "INSERT INTO data_redirect (id, redirect_time, redirect_added_time, http_status, urlservicekey, campaign_id, country, partner, operator, aggregator, service, short_code, adnet, keyword, subkeyword, is_billable, plan) VALUES (DEFAULT, '%s', %d, %d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, '%s');"
	GETTOTALDATAREDIRECT                        = "SELECT COUNT(1) FROM data_redirect WHERE DATE(redirect_time) = '%s' AND urlservicekey = '%s' AND campaign_id = '%s' AND country = '%s' AND partner = '%s' AND operator = '%s' AND service = '%s' AND short_code = '%s' AND adnet = '%s'"
	UPDATECPAREPORTSUMMARYCAMPAIGN              = "UPDATE summary_campaign SET cost_per_conversion = '%s', agency_fee = '%s' WHERE DATE(summary_date) = '%s' AND urlservicekey = '%s' AND country = '%s' AND operator = '%s' AND partner = '%s' AND service = '%s' AND adnet = '%s' AND campaign_id = '%s'"
	UPDATEREPORTSUMMARYCAMPAIGNMONITORINGBUDGET = "UPDATE summary_campaign SET target_daily_budget = '%s' WHERE DATE(summary_date) = '%s AND country = '%s' AND operator = '%s'"
)

func (r *BaseModel) DelSummaryCampaign(summary_date string, o entity.DataConfig) error {

	SQL := fmt.Sprintf(DELSUMMARYCAMPAIGN, summary_date, o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("(%s) Error %s when preparing SQL statement", SQL, err))

		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when update to table", SQL, err))

		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when finding rows affected", SQL, err))

		return err
	}

	r.Logs.Debug(fmt.Sprintf("SQL : %s, row affected : %d", SQL, rows))
	return nil
}

func (r *BaseModel) EditSettingSummaryCampaign(summary_date string, o entity.DataConfig) error {

	SQL := fmt.Sprintf(EDITSETTINGSUMMARYCAMPAIGN, o.PO, o.MOCapping, o.RatioSend, o.RatioReceive, summary_date, o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("(%s) Error %s when preparing SQL statement", SQL, err))

		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when update to table", SQL, err))

		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when finding rows affected", SQL, err))

		return err
	}

	r.Logs.Debug(fmt.Sprintf("SQL : %s, row affected : %d", SQL, rows))
	return nil
}

func (r *BaseModel) UpdateSummaryCampaign(summary_date string, o entity.DataConfig) error {

	SQL := fmt.Sprintf(UPDATESUMMARYCAMPAIGN, o.IsActive, summary_date, o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("(%s) Error %s when preparing SQL statement", SQL, err))

		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when update to table", SQL, err))

		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when finding rows affected", SQL, err))

		return err
	}

	r.Logs.Debug(fmt.Sprintf("SQL : %s, row affected : %d", SQL, rows))
	return nil
}

func (r *BaseModel) SummaryCampaign(d entity.Summary) int {

	SQL := fmt.Sprintf(SUMMARYCAMPAIGN, d.IsActive, d.SummaryDate, d.CampaignId, d.CampaignName, d.Country, d.Partner, d.Operator, d.URLServiceKey, d.Aggregator, d.Service, d.Adnet, d.ShortCode, d.TotalTraffic, d.TotalLanding, d.TotalMOReceived, d.CRMO, d.CRPostback, d.TotalPostback, d.TotalFP, d.SuccessFP, d.BillRate, d.PO, d.SBAF, d.SAAF, d.CPA, d.Revenue, d.URLWarpLanding, d.URLLanding, d.MOCapping, d.RatioSend, d.RatioReceive, d.ClientType, d.CPCR, d.AgencyFee, d.TotalWakiAgencyFee, d.TDB, d.BudgetUsage, d.TotalTraffic, d.TotalLanding, d.TotalMOReceived, d.CRMO, d.CRPostback, d.TotalPostback, d.TotalFP, d.SuccessFP, d.BillRate, d.PO, d.SBAF, d.SAAF, d.CPA, d.Revenue, d.URLWarpLanding, d.URLLanding, d.MOCapping, d.RatioSend, d.RatioReceive, d.ClientType, d.CPCR, d.AgencyFee, d.TotalWakiAgencyFee, d.TDB, d.BudgetUsage, d.CampaignName)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("(%s) Error %s when preparing SQL statement", SQL, err))

		return 0
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when update to table", SQL, err))

		return 0
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when finding rows affected", SQL, err))

		return 0
	}

	r.Logs.Debug(fmt.Sprintf("SQL : %s, row affected : %d", SQL, rows))
	return int(rows)
}

func (r *BaseModel) DataTraffic(data map[string]string, o entity.DataCounter) int {

	traffic_added_time, _ := strconv.Atoi(data["traffic_added_time"])
	httpstatus, _ := strconv.Atoi(data["http_status"])

	SQL := fmt.Sprintf(INSERTTRAFFIC, data["traffic_time"], traffic_added_time, httpstatus, o.URLServiceKey, o.CampaignId, o.Country, o.Partner, o.Operator, o.Aggregator, o.Service, o.ShortCode, o.Adnet, o.Keyword, o.SubKeyword, o.IsBillable, o.Plan)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("(%s) Error %s when preparing SQL statement", SQL, err))

		return 0
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when update to table", SQL, err))

		return 0
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when finding rows affected", SQL, err))

		return 0
	}

	r.Logs.Debug(fmt.Sprintf("SQL : %s, row affected : %d", SQL, rows))
	return int(rows)
}

func (r *BaseModel) DataLanding(data map[string]string, o entity.DataCounter) int {

	landed_time, _ := strconv.Atoi(data["landed_time"])
	httpstatus, _ := strconv.Atoi(data["http_status"])

	SQL := fmt.Sprintf(INSERTLANDING, data["landing_time"], landed_time, httpstatus, o.URLServiceKey, o.CampaignId, o.Country, o.Partner, o.Operator, o.Aggregator, o.Service, o.ShortCode, o.Adnet, o.Keyword, o.SubKeyword, o.IsBillable, o.Plan)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("(%s) Error %s when preparing SQL statement", SQL, err))

		return 0
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when update to table", SQL, err))

		return 0
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when finding rows affected", SQL, err))

		return 0
	}

	r.Logs.Debug(fmt.Sprintf("SQL : %s, row affected : %d", SQL, rows))
	return int(rows)
}

func (r *BaseModel) DataClicked(data map[string]string, o entity.DataCounter) int {

	clicked_button_time, _ := strconv.Atoi(data["clicked_button_time"])
	httpstatus, _ := strconv.Atoi(data["http_status"])

	SQL := fmt.Sprintf(INSERTCLICKED, data["clicked_time"], clicked_button_time, httpstatus, o.URLServiceKey, o.CampaignId, o.Country, o.Partner, o.Operator, o.Aggregator, o.Service, o.ShortCode, o.Adnet, o.Keyword, o.SubKeyword, o.IsBillable, o.Plan)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("(%s) Error %s when preparing SQL statement", SQL, err))

		return 0
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when update to table", SQL, err))

		return 0
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when finding rows affected", SQL, err))

		return 0
	}

	r.Logs.Debug(fmt.Sprintf("SQL : %s, row affected : %d", SQL, rows))
	return int(rows)
}

func (r *BaseModel) DataRedirect(data map[string]string, o entity.DataCounter) int {

	redirect_added_time, _ := strconv.Atoi(data["redirect_added_time"])
	httpstatus, _ := strconv.Atoi(data["http_status"])

	SQL := fmt.Sprintf(INSERTREDIRECT, data["redirect_time"], redirect_added_time, httpstatus, o.URLServiceKey, o.CampaignId, o.Country, o.Partner, o.Operator, o.Aggregator, o.Service, o.ShortCode, o.Adnet, o.Keyword, o.SubKeyword, o.IsBillable, o.Plan)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("(%s) Error %s when preparing SQL statement", SQL, err))

		return 0
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when update to table", SQL, err))

		return 0
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when finding rows affected", SQL, err))

		return 0
	}

	r.Logs.Debug(fmt.Sprintf("SQL : %s, row affected : %d", SQL, rows))
	return int(rows)
}

func (r *BaseModel) UpdateCPAReportSummaryCampaign(summary_date string, o entity.DataConfig) error {

	SQL := fmt.Sprintf(UPDATECPAREPORTSUMMARYCAMPAIGN, summary_date, o.CPCR, o.AgencyFee, o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("(%s) Error %s when preparing SQL statement", SQL, err))

		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when update to table", SQL, err))

		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when finding rows affected", SQL, err))

		return err
	}

	r.Logs.Debug(fmt.Sprintf("SQL : %s, row affected : %d", SQL, rows))
	return nil
}

func (r *BaseModel) UpdateReportSummaryCampaignMonitoringBudget(summary_date string, o entity.DataConfig) error {

	SQL := fmt.Sprintf(UPDATEREPORTSUMMARYCAMPAIGNMONITORINGBUDGET, o.TargetDailyBudget, summary_date, o.Country, o.Operator)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("(%s) Error %s when preparing SQL statement", SQL, err))

		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when update to table", SQL, err))

		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when finding rows affected", SQL, err))

		return err
	}

	r.Logs.Debug(fmt.Sprintf("SQL : %s, row affected : %d", SQL, rows))
	return nil
} */
