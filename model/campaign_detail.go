package model

import (
	"context"
	"fmt"
	"strconv"

	"github.com/infraLinkit/mediaplatform-datasource/entity"

	_ "github.com/lib/pq"
)

const (
	GETLASTIDCAMPDETAIL            = "SELECT COALESCE(MAX(id), 0) + 0 FROM %s"
	NEWCAMPAIGN                    = "INSERT INTO campaign (id, campaign_id, name, campaign_objective, country, advertiser) VALUES (DEFAULT, '%s', '%s', '%s', '%s', '%s')"
	GETCAMPAIGNBYCAMPAIGNID        = "SELECT * FROM campaign WHERE campaign_id = '%s'"
	NEWCAMPAIGNDETAIL              = "INSERT INTO campaign_detail (id, urlservicekey, campaign_id, country, operator, partner, aggregator, adnet, service, keyword, subkeyword, is_billable, plan, po, cost, pubid, short_code, device_type, os, url_type, click_type, click_delay, client_type, traffic_source, unique_click, url_banner, url_landing, url_warp_landing, url_service, url_tfc_or_smartlink, glob_post, url_globpost, custom_integration, ip_address, is_active, mo_capping, counter_mo_capping, status_capping, kpi_upper_limit_capping, is_machine_learning_capping, ratio_send, ratio_receive, counter_mo_ratio, status_ratio, kpi_upper_limit_ratio_send, kpi_upper_limit_ratio_receive, is_machine_learning_ratio, api_url, last_update, last_update_capping, cost_per_conversion, agency_fee, target_daily_budget, url_postback) VALUES (%d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %d, %d, '%s', %t, %t, '%s', '%s', '%s', '%s', '%s', %t, '%s', '%s', '%s', %t, %d, %d, %t, %d, %t, %d, %d, %d, %t, %d, %d, %t, '%s', '%s', '%s', '%s'::double precision, '%s'::double precision, '%s'::double precision, '%s')"
	RESETCAPPINGCAMPAIGN           = "UPDATE campaign_detail SET counter_mo_ratio = 0, status_capping = false WHERE is_active = %t"
	GETCAMPAIGNBYCAMPAIGNDETAILID  = "SELECT id, urlservicekey, campaign_id, country, operator, partner, aggregator, adnet, service, keyword, subkeyword, is_billable, plan, po, cost, pubid, short_code, device_type, os, url_type, click_type, click_delay, client_type, traffic_source, unique_click, url_banner, url_landing, url_warp_landing, url_service, url_tfc_or_smartlink, glob_post, url_globpost, custom_integration, ip_address, is_active, mo_capping, counter_mo_capping, status_capping, kpi_upper_limit_capping, is_machine_learning_capping, ratio_send, ratio_receive, counter_mo_ratio, status_ratio, kpi_upper_limit_ratio_send, kpi_upper_limit_ratio_receive, is_machine_learning_ratio, api_url, last_update, last_update_capping, cost_per_conversion, agency_fee, target_daily_budget, url_postback FROM campaign_detail WHERE urlservicekey = '%s' AND country = '%s' AND operator = '%s' AND partner = '%s' AND service = '%s' AND adnet = '%s'"
	UPDATECAMPAIGN                 = "UPDATE campaign SET name = '%s', campaign_objective = '%s', country = '%s', advertiser = '%s' WHERE campaign_id = '%s'"
	UPDATECAMPAIGNDETAIL           = "UPDATE campaign_detail SET campaign_id = '%s', country = '%s', operator = '%s', partner = '%s', aggregator = '%s', adnet = '%s', service = '%s', keyword = '%s', subkeyword = '%s', is_billable = %t, plan = '%s', po = '%s', cost = '%s', pubid = '%s', short_code = '%s', device_type = '%s', os = '%s', url_type = '%s', click_type = %d, click_delay = %d, client_type = '%s', traffic_source = %t, unique_click = %t, url_banner = '%s', url_landing = '%s', url_warp_landing = '%s', url_service = '%s', url_tfc_or_smartlink = '%s', glob_post = %t, url_globpost = '%s', custom_integration = '%s', ip_address = '%s', is_active = %t, mo_capping = %d, counter_mo_capping = %d, status_capping = %t, kpi_upper_limit_capping = %d, is_machine_learning_capping = %t, ratio_send = %d, ratio_receive = %d, counter_mo_ratio = %d, status_ratio = %t, kpi_upper_limit_ratio_send = %d, kpi_upper_limit_ratio_receive = %d, is_machine_learning_ratio = %t, api_url = '%s', last_update = '%s', last_update_capping = '%s', cost_per_conversion = '%s'::double precision, agency_fee = '%s'::double precision, target_daily_budget = '%s'::double precision, url_postback = '%s' WHERE id = %d"
	DELCAMPAIGN                    = "DELETE FROM campaign WHERE campaign_id = '%s'"
	DELCAMPAIGNDETAIL              = "DELETE FROM campaign_detail WHERE urlservicekey = '%s' AND country = '%s' AND operator = '%s' AND partner = '%s' AND service = '%s' AND adnet = '%s' AND campaign_id = '%s'"
	EDITSETTINGCAMPAIGNDETAIL      = "UPDATE campaign_detail SET po = '%s', mo_capping = %d, ratio_send = %d, ratio_receive = %d, last_update = '%s' WHERE urlservicekey = '%s' AND country = '%s' AND operator = '%s' AND partner = '%s' AND service = '%s' AND adnet = '%s' AND campaign_id = '%s'"
	UPDATESTATUSCAMPAIGNDETAIL     = "UPDATE campaign_detail SET is_active = %t WHERE urlservicekey = '%s' AND country = '%s' AND operator = '%s' AND partner = '%s' AND service = '%s' AND adnet = '%s' AND campaign_id = '%s'"
	GETCAMPAIGNDETAIL              = "SELECT id, urlservicekey, is_active, counter_mo_capping, mo_capping, status_capping, counter_mo_ratio, ratio_send, ratio_receive, status_ratio, api_url, pubid, cost, po FROM campaign_detail WHERE id = %d;"
	COUNTERCAPPING                 = "UPDATE campaign_detail SET counter_mo_capping = counter_mo_capping+1, last_update_capping = CASE WHEN counter_mo_capping >= mo_capping THEN '%s'::timestamp(0) END WHERE id = %d;"
	COUNTERRATIO                   = "UPDATE campaign_detail SET counter_mo_ratio = counter_mo_ratio+1 WHERE id = %d;"
	UPDATESTATUSCOUNTER            = "UPDATE campaign_detail SET counter_mo_capping = %d, status_capping = %t, counter_mo_ratio = %d, status_ratio = %t, last_update = '%s'::timestamp(0), last_update_capping = CASE WHEN counter_mo_capping+1 >= mo_capping THEN '%s'::timestamp(0) END WHERE id = %d"
	GETCAMPAIGNDETAILBYSTATUS      = "SELECT c.name, c.campaign_objective, c.advertiser, cd.id, cd.urlservicekey, cd.campaign_id, cd.country, cd.operator, cd.partner, cd.aggregator, cd.adnet, cd.service, cd.keyword, cd.subkeyword, cd.is_billable, cd.plan, cd.po, cd.cost, cd.pubid, cd.short_code, cd.device_type, cd.os, cd.url_type, cd.click_type, cd.click_delay, cd.client_type, cd.traffic_source, cd.unique_click, cd.url_banner, cd.url_landing, cd.url_warp_landing, cd.url_service, cd.url_tfc_or_smartlink, cd.glob_post, cd.url_globpost, cd.custom_integration, cd.ip_address, cd.is_active, cd.mo_capping, cd.counter_mo_capping, cd.status_capping, cd.kpi_upper_limit_capping, cd.is_machine_learning_capping, cd.ratio_send, cd.ratio_receive, cd.counter_mo_ratio, cd.status_ratio, cd.kpi_upper_limit_ratio_send, cd.kpi_upper_limit_ratio_receive, cd.is_machine_learning_ratio, cd.api_url, cd.last_update, cd.cost_per_conversion, cd.agency_fee, cd.target_daily_budget, cd.url_postback FROM campaign_detail cd INNER JOIN campaign c ON cd.campaign_id  = c.campaign_id WHERE is_active = %t;"
	GETCAMPAIGNDETAILALL           = "SELECT * FROM campaign_detail;"
	UPDATECPAREPORT                = "UPDATE campaign_detail SET cost_per_conversion = '%s', agency_fee = '%s' WHERE urlservicekey = '%s' AND country = '%s' AND operator = '%s' AND partner = '%s' AND service = '%s' AND adnet = '%s' AND campaign_id = '%s'"
	UPDATECAMPAIGNMONITORINGBUDGET = "UPDATE campaign_detail SET target_daily_budget = '%s' WHERE country = '%s' AND operator = '%s'"
)

func (r *BaseModel) GetLastCampaignId(tbl string) int {

	SQL := fmt.Sprintf(GETLASTIDCAMPDETAIL, tbl)

	var id int
	err := r.DBPostgre.QueryRow(SQL).Scan(&id)
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("(%s) Error %s when preparing SQL statement", SQL, err))
		return 0
	}

	r.Logs.Debug(fmt.Sprintf("(%s) found %d", SQL, id))
	return id

}

func (r *BaseModel) NewCampaign(o entity.DataCampaignAction) int {

	SQL := fmt.Sprintf(NEWCAMPAIGN, o.CampaignId, o.CampaignName, o.Objective, o.Country, o.Advertiser)

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

func (r *BaseModel) ResetCappingCampaign(o entity.DataConfig) int {

	SQL := fmt.Sprintf(RESETCAPPINGCAMPAIGN, o.IsActive)

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

func (r *BaseModel) GetCampaignByCampaignId(campId string) entity.DataCampaignAction {

	SQL := fmt.Sprintf(GETCAMPAIGNBYCAMPAIGNID, campId)

	var o entity.DataCampaignAction
	err := r.DBPostgre.QueryRow(SQL).Scan(&o.Id, &o.CampaignId, &o.CampaignName, &o.Objective, &o.Country, &o.Advertiser)
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("(%s) Error %s when preparing SQL statement", SQL, err))
		return entity.DataCampaignAction{}
	}

	r.Logs.Debug(fmt.Sprintf("(%s) found %#v", SQL, o))
	return o

}

func (r *BaseModel) NewCampaignDetail(o entity.DataConfig) int {

	ips := "{}"
	if len(o.IPAddress) > 0 {

		ips = "{"
		for _, v := range o.IPAddress {
			ips = ips + strconv.Itoa(int(v))
		}
		ips = ips + "}"
	}

	SQL := fmt.Sprintf(NEWCAMPAIGNDETAIL, o.Id, o.URLServiceKey, o.CampaignId, o.Country, o.Operator, o.Partner, o.Aggregator, o.Adnet, o.Service, o.Keyword, o.SubKeyword, o.IsBillable, o.Plan, o.PO, o.Cost, o.PubId, o.ShortCode, o.DeviceType, o.OS, o.URLType, o.ClickType, o.ClickDelay, o.ClientType, o.TrafficSource, o.UniqueClick, o.URLBanner, o.URLLanding, o.URLWarpLanding, o.URLService, o.URLTFCSmartlink, o.GlobPost, o.URLGlobPost, o.CustomIntegration, ips, o.IsActive, o.MOCapping, o.CounterMOCapping, o.StatusCapping, o.KPIUpperLimitCapping, o.IsMachineLearningCapping, o.RatioSend, o.RatioReceive, o.CounterMORatio, o.StatusRatio, o.KPIUpperLimitRatioSend, o.KPIUpperLimitRatioReceive, o.IsMachineLearningRatio, o.APIURL, o.LastUpdate, o.LastUpdateCapping, o.CPCR, o.AgencyFee, o.TargetDailyBudget, o.URLPostback)

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

func (r *BaseModel) GetCampaignByCampaignDetailId(o entity.DataConfig) entity.DataConfig {

	SQL := fmt.Sprintf(GETCAMPAIGNBYCAMPAIGNDETAILID, o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Adnet, o.Service)

	err := r.DBPostgre.QueryRow(SQL).Scan(&o.Id, &o.URLServiceKey, &o.CampaignId, &o.Country, &o.Operator, &o.Partner, &o.Aggregator, &o.Adnet, &o.Service, &o.Keyword, &o.SubKeyword, &o.IsBillable, &o.Plan, &o.PO, &o.Cost, &o.PubId, &o.ShortCode, &o.DeviceType, &o.OS, &o.URLType, &o.ClickType, &o.ClickDelay, &o.ClientType, &o.TrafficSource, &o.UniqueClick, &o.URLBanner, &o.URLLanding, &o.URLWarpLanding, &o.URLService, &o.URLTFCSmartlink, &o.GlobPost, &o.URLGlobPost, &o.CustomIntegration, &o.IPAddress, &o.IsActive, &o.MOCapping, &o.CounterMOCapping, &o.StatusCapping, &o.KPIUpperLimitCapping, &o.IsMachineLearningCapping, &o.RatioSend, &o.RatioReceive, &o.CounterMORatio, &o.StatusRatio, &o.KPIUpperLimitRatioSend, &o.KPIUpperLimitRatioReceive, &o.IsMachineLearningRatio, &o.APIURL, &o.LastUpdate, &o.LastUpdateCapping, &o.CPCR, &o.AgencyFee, &o.URLPostback)
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("(%s) Error %s when preparing SQL statement", SQL, err))
		return o
	}

	r.Logs.Debug(fmt.Sprintf("(%s) found %#v", SQL, o))
	return o

}

func (r *BaseModel) UpdateCampaign(o entity.DataCampaignAction) error {

	SQL := fmt.Sprintf(UPDATECAMPAIGN, o.CampaignName, o.Objective, o.Country, o.Advertiser, o.CampaignId)

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

func (r *BaseModel) UpdateCampaignDetail(o entity.DataConfig) error {

	ips := "{}"
	if len(o.IPAddress) > 0 {

		ips = "{"
		for _, v := range o.IPAddress {
			ips = ips + strconv.Itoa(int(v))
		}
		ips = ips + "}"
	}

	SQL := fmt.Sprintf(UPDATECAMPAIGNDETAIL, o.CampaignId, o.Country, o.Operator, o.Partner, o.Aggregator, o.Adnet, o.Service, o.Keyword, o.SubKeyword, o.IsBillable, o.Plan, o.PO, o.Cost, o.PubId, o.ShortCode, o.DeviceType, o.OS, o.URLType, o.ClickType, o.ClickDelay, o.ClientType, o.TrafficSource, o.UniqueClick, o.URLBanner, o.URLLanding, o.URLWarpLanding, o.URLService, o.URLTFCSmartlink, o.GlobPost, o.URLGlobPost, o.CustomIntegration, ips, o.IsActive, o.MOCapping, o.CounterMOCapping, o.StatusCapping, o.KPIUpperLimitCapping, o.IsMachineLearningCapping, o.RatioSend, o.RatioReceive, o.CounterMORatio, o.StatusRatio, o.KPIUpperLimitRatioSend, o.KPIUpperLimitRatioReceive, o.IsMachineLearningRatio, o.APIURL, o.LastUpdate, o.LastUpdateCapping, o.CPCR, o.AgencyFee, o.TargetDailyBudget, o.URLPostback, o.Id)

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

func (r *BaseModel) DelCampaign(o entity.DataCampaignAction) error {

	SQL := fmt.Sprintf(DELCAMPAIGN, o.CampaignId)

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

func (r *BaseModel) DelCampaignDetail(o entity.DataConfig) error {

	SQL := fmt.Sprintf(DELCAMPAIGNDETAIL, o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId)

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

func (r *BaseModel) EditSettingCampaignDetail(o entity.DataConfig) error {

	SQL := fmt.Sprintf(EDITSETTINGCAMPAIGNDETAIL, o.PO, o.MOCapping, o.RatioSend, o.RatioReceive, o.LastUpdate, o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId)

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

func (r *BaseModel) UpdateStatusCampaignDetail(o entity.DataConfig) error {

	SQL := fmt.Sprintf(UPDATESTATUSCAMPAIGNDETAIL, o.IsActive, o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId)

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

func (r *BaseModel) GetCampaignDetail(o entity.DataConfig) (entity.DataConfig, error) {

	SQL := fmt.Sprintf(GETCAMPAIGNDETAIL, o.Id)
	rows, err := r.DBPostgre.Query(SQL)
	if err != nil {
		r.Logs.Error(fmt.Sprintf("SQL : %s, error querying occured : %#v", SQL, err))

		return entity.DataConfig{}, err
	}
	defer rows.Close()

	for rows.Next() {

		err = rows.Scan(&o.Id, &o.URLServiceKey, &o.IsActive, &o.CounterMOCapping, &o.MOCapping, &o.StatusCapping, &o.CounterMORatio, &o.RatioSend, &o.RatioReceive, &o.StatusRatio, &o.APIURL, &o.PubId, &o.Cost, &o.PO)

		if err != nil {

			r.Logs.Error(fmt.Sprintf("SQL : %s, error scan occured : %#v", SQL, err))

		}
	}

	r.Logs.Info(fmt.Sprintf("SQL : %s, row selected occured : %#v", SQL, o))
	return o, nil
}

func (r *BaseModel) CounterCappingById(o entity.DataConfig) error {

	SQL := fmt.Sprintf(COUNTERCAPPING, o.LastUpdateCapping, o.Id)

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

func (r *BaseModel) CounterRatioById(id int) error {

	SQL := fmt.Sprintf(COUNTERRATIO, id)

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

func (r *BaseModel) UpdateStatusCounterById(o entity.DataConfig) error {

	SQL := fmt.Sprintf(UPDATESTATUSCOUNTER, o.CounterMOCapping, o.StatusCapping, o.CounterMORatio, o.StatusRatio, o.LastUpdate, o.LastUpdateCapping, o.Id)

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

func (r *BaseModel) GetCampaignDetailByStatus(obj entity.DataConfig, useStatus bool) ([]entity.DataConfig, error) {

	var SQL string
	if useStatus {
		SQL = fmt.Sprintf(GETCAMPAIGNDETAILBYSTATUS, obj.IsActive)
	} else {
		SQL = GETCAMPAIGNDETAILALL
	}

	rows, err := r.DBPostgre.Query(SQL)
	if err != nil {
		r.Logs.Error(fmt.Sprintf("SQL : %s, error querying occured : %#v", SQL, err))

		return []entity.DataConfig{}, err
	}
	defer rows.Close()

	var oo []entity.DataConfig

	for rows.Next() {

		var o entity.DataConfig

		if err = rows.Scan(&o.CampaignName, &o.Objective, &o.Advertiser, &o.Id, &o.URLServiceKey, &o.CampaignId, &o.Country, &o.Operator, &o.Partner, &o.Aggregator, &o.Adnet, &o.Service, &o.Keyword, &o.SubKeyword, &o.IsBillable, &o.Plan, &o.PO, &o.Cost, &o.PubId, &o.ShortCode, &o.DeviceType, &o.OS, &o.URLType, &o.ClickType, &o.ClickDelay, &o.ClientType, &o.TrafficSource, &o.UniqueClick, &o.URLBanner, &o.URLLanding, &o.URLWarpLanding, &o.URLService, &o.URLTFCSmartlink, &o.GlobPost, &o.URLGlobPost, &o.CustomIntegration, &o.IPAddress, &o.IsActive, &o.MOCapping, &o.CounterMOCapping, &o.StatusCapping, &o.KPIUpperLimitCapping, &o.IsMachineLearningCapping, &o.RatioSend, &o.RatioReceive, &o.CounterMORatio, &o.StatusRatio, &o.KPIUpperLimitRatioSend, &o.KPIUpperLimitRatioReceive, &o.IsMachineLearningRatio, &o.APIURL, &o.LastUpdate, &o.CPCR, &o.AgencyFee, &o.TargetDailyBudget, &o.URLPostback); err != nil {

			r.Logs.Error(fmt.Sprintf("SQL : %s, error scan occured : %#v", SQL, err))

			return []entity.DataConfig{}, err
		}

		oo = append(oo, o)
	}

	r.Logs.Info(fmt.Sprintf("SQL : %s, row selected occured : %#v", SQL, len(oo)))

	return oo, nil
}

func (r *BaseModel) UpdateCPAReport(o entity.DataConfig) error {

	SQL := fmt.Sprintf(UPDATECPAREPORT, o.CPCR, o.AgencyFee, o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId)

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

func (r *BaseModel) UpdateCampaignMonitoringBudget(o entity.DataConfig) error {

	SQL := fmt.Sprintf(UPDATECAMPAIGNMONITORINGBUDGET, o.TargetDailyBudget, o.Country, o.Operator)

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
