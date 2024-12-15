package model

import (
	"context"
	"fmt"
	"strconv"

	"github.com/infraLinkit/mediaplatform-datasource/entity"

	_ "github.com/lib/pq"
)

const (
	DELSUMMARYCAMPAIGN         = "DELETE FROM summary_campaign WHERE DATE(summary_date) = '%s' AND urlservicekey = '%s' AND country = '%s' AND operator = '%s' AND partner = '%s' AND adnet = '%s' AND service = '%s' AND campaign_id = '%s'"
	EDITSETTINGSUMMARYCAMPAIGN = "UPDATE summary_campaign SET po = '%s', mo_capping = %d, ratio_send = %d, ratio_receive = %d, last_update = '%s' WHERE DATE(summary_date) = '%s' AND urlservicekey = '%s' AND country = '%s' AND operator = '%s' AND partner = '%s' AND adnet = '%s' AND service = '%s' AND campaign_id = '%s'"
	UPDATESUMMARYCAMPAIGN      = "UPDATE summary_campaign SET status = %t WHERE DATE(summary_date) = '%s' AND urlservicekey = '%s' AND country = '%s' AND operator = '%s' AND partner = '%s' AND adnet = '%s' AND service = '%s' AND campaign_id = '%s'"
	SUMMARYCAMPAIGN            = "INSERT INTO summary_campaign (id, status, summary_date, campaign_id, campaign_name, country, partner, operator, urlservicekey, aggregator, service, adnet, short_code, traffic, landing, mo_received, cr, postback, total_fp, success_fp, billrate, po, sbaf, saaf, cpa, revenue, url_after, url_before, mo_limit, ratio_send, ratio_receive, client_type, cost_per_conversion, agency_fee, total_waki_agency_fee) VALUES (DEFAULT, %t, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %d, %d, %d, '%s'::double precision, %d, %d, '%s', '%s'::double precision, '%s'::double precision, '%s', '%s', '%s', '%s', '%s', '%s', %d, %d, %d, '%s', '%s', '%s', '%s') ON CONFLICT (summary_date, campaign_id, country, partner, operator, urlservicekey, service, adnet) DO UPDATE SET traffic = traffic + %d, landing = landing + %d, mo_received = mo_received + %d, cr = '%s'::double precision, postback = postback + %d, total_fp = total_fp + %d, success_fp = '%s', billrate = '%s'::double precision, po = '%s'::double precision, sbaf = '%s', saaf = '%s', cpa = '%s', revenue = '%s', url_after = '%s', url_before = '%s', mo_limit = %d, ratio_send = %d, ratio_receive = %d, client_type = '%s', cost_per_conversion = '%s', agency_fee = '%s', total_waki_agency_fee = '%s';"
	INSERTTRAFFIC              = "INSERT INTO data_traffic (id, traffic_time, traffic_added_time, http_status, urlservicekey, campaign_id, country, partner, operator, aggregator, service, short_code, adnet, keyword, subkeyword, is_billable, plan) VALUES (DEFAULT, '%s', %d, %d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, '%s');"
	GETTOTALDATATRAFFIC        = "SELECT COUNT(1) FROM data_traffic WHERE DATE(traffic_time) = '%s' AND urlservicekey = '%s' AND campaign_id = '%s' AND country = '%s' AND partner = '%s' AND operator = '%s' AND service = '%s' AND short_code = '%s' AND adnet = '%s'"
	INSERTLANDING              = "INSERT INTO data_landing (id, landing_time, landed_time, http_status, urlservicekey, campaign_id, country, partner, operator, aggregator, service, short_code, adnet, keyword, subkeyword, is_billable, plan) VALUES (DEFAULT, '%s', %d, %d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, '%s');"
	GETTOTALDATALANDING        = "SELECT COUNT(1) FROM data_landing WHERE DATE(landing_time) = '%s' AND urlservicekey = '%s' AND campaign_id = '%s' AND country = '%s' AND partner = '%s' AND operator = '%s' AND service = '%s' AND short_code = '%s' AND adnet = '%s'"
	INSERTCLICKED              = "INSERT INTO data_clicked (id, clicked_time, clicked_button_time, http_status, urlservicekey, campaign_id, country, partner, operator, aggregator, service, short_code, adnet, keyword, subkeyword, is_billable, plan) VALUES (DEFAULT, '%s', %d, %d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, '%s');"
	GETTOTALDATACLICKED        = "SELECT COUNT(1) FROM data_clicked WHERE DATE(clicked_time) = '%s' AND urlservicekey = '%s' AND campaign_id = '%s' AND country = '%s' AND partner = '%s' AND operator = '%s' AND service = '%s' AND short_code = '%s' AND adnet = '%s'"
	INSERTREDIRECT             = "INSERT INTO data_redirect (id, redirect_time, redirect_added_time, http_status, urlservicekey, campaign_id, country, partner, operator, aggregator, service, short_code, adnet, keyword, subkeyword, is_billable, plan) VALUES (DEFAULT, '%s', %d, %d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, '%s');"
	GETTOTALDATAREDIRECT       = "SELECT COUNT(1) FROM data_redirect WHERE DATE(redirect_time) = '%s' AND urlservicekey = '%s' AND campaign_id = '%s' AND country = '%s' AND partner = '%s' AND operator = '%s' AND service = '%s' AND short_code = '%s' AND adnet = '%s'"
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

	SQL := fmt.Sprintf(EDITSETTINGSUMMARYCAMPAIGN, o.PO, o.MOCapping, o.RatioSend, o.RatioReceive, o.LastUpdate, summary_date, o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId)

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

func (r *BaseModel) SummaryCampaign(data map[string]string, o entity.DataConfig, o2 entity.DataCounter) int {

	SQL := fmt.Sprintf(SUMMARYCAMPAIGN, o.IsActive, data["summary_date"], o.CampaignId, o.CampaignName, o.Country, o.Partner, o.Operator, o.URLServiceKey, o.Aggregator, o.Service, o.Adnet, o.ShortCode, o2.Traffic, o2.Landing, o2.MOReceived, data["cr"], o2.Postback, o2.TotalFP, data["success_fp"], data["billrate"], o.PO, data["sbaf"], data["saaf"], data["cpa"], data["revenue"], o.URLWarpLanding, o.URLLanding, o.MOCapping, o.RatioSend, o.RatioReceive, o.ClientType, data["cost_per_conversion"], data["agency_fee"], data["total_waki_agency_fee"], o2.Traffic, o2.Landing, o2.MOReceived, data["cr"], o2.Postback, o2.TotalFP, data["success_fp"], data["billrate"], o.PO, data["sbaf"], data["saaf"], data["cpa"], data["revenue"], o.URLWarpLanding, o.URLLanding, o.MOCapping, o.RatioSend, o.RatioReceive, o.ClientType, data["cost_per_conversion"], data["agency_fee"], data["total_waki_agency_fee"])

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