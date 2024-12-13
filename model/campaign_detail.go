package model

import (
	"context"
	"fmt"

	"github.com/infraLinkit/mediaplatform-datasource/entity"

	_ "github.com/lib/pq"
)

const (
	GETLASTIDCAMPDETAIL           = "SELECT MAX(id) FROM %s"
	NEWCAMPAIGN                   = "INSERT INTO campaign (id, campaign_id, name, campaign_objective, country, advertiser) VALUES (DEFAULT, '%s', '%s', '%s', '%s', '%s')"
	GETCAMPAIGNBYCAMPAIGNID       = "SELECT * FROM campaign WHERE campaign_id = '%s'"
	NEWCAMPAIGNDETAIL             = "INSERT INTO campaign_detail (id, urlservicekey, campaign_id, country, operator, partner, aggregator, adnet, service, keyword, subkeyword, is_billable, plan, po, cost, pubid, short_code, device_type, os, url_type, click_type, click_delay, client_type, traffic_source, unique_click, url_banner, url_landing, url_warp_landing, url_service, url_tfc_or_smartlink, glob_post, url_globpost, custom_integration, ip_address, is_active, mo_capping, counter_mo_capping, status_capping, kpi_upper_limit_capping, is_machine_learning_capping, ratio_send, ratio_receive, counter_mo_ratio, status_ratio, kpi_upper_limit_ratio_send, kpi_upper_limit_ratio_receive, is_machine_learning_ratio, api_url, last_update, last_update_capping) VALUES (DEFAULT, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %d, %d, '%s', %t, %t, '%s', '%s', '%s', '%s', '%s', %t, '%s', '%s', '%s', %t, %d, %d, %t, %d, %t, %d, %d, %d, %t, %d, %d, %t, '%s', '%s', '%s')"
	GETCAMPAIGNBYCAMPAIGNDETAILID = "SELECT * FROM campaign_detail WHERE urlservicekey = '%s' AND country = '%s' AND operator = '%s' AND partner = '%s' AND adnet = '%s' AND service = '%s'"
	UPDATECAMPAIGN                = "UPDATE campaign SET name = '%s', campaign_objective = '%s', country = '%s', advertiser = '%s' WHERE campaign_id = '%s'"
	UPDATECAMPAIGNDETAIL          = "UPDATE campaign_detail SET campaign_id = '%s', country = '%s', operator = '%s', partner = '%s', aggregator = '%s', adnet = '%s', service = '%s', keyword = '%s', subkeyword = '%s', is_billable = %t, plan = '%s', po = '%s', cost = '%s', pubid = '%s', short_code = '%s', device_type = '%s', os = '%s', url_type = '%s', click_type = %d, click_delay = %d, client_type = '%s', traffic_source = %t, unique_click = %t, url_banner = '%s', url_landing = '%s', url_warp_landing = '%s', url_service = '%s', url_tfc_or_smartlink = '%s', glob_post = %t, url_globpost = '%s', custom_integration = '%s', ip_address = '%s', is_active = %t, mo_capping = %d, counter_mo_capping = %d, status_capping = %t, kpi_upper_limit_capping = %d, is_machine_learning_capping = %t, ratio_send = %d, ratio_receive = %d, counter_mo_ratio = %d, status_ratio = %t, kpi_upper_limit_ratio_send = %d, kpi_upper_limit_ratio_receive = %d, is_machine_learning_ratio = %t, api_url = '%s', last_update = '%s', last_update_capping = '%s' WHERE id = %d"
	DELCAMPAIGN                   = "DELETE FROM campaign WHERE campaign_id = '%s'"
	DELCAMPAIGNDETAIL             = "DELETE FROM campaign_detail WHERE urlservicekey = '%s' AND country = '%s' AND operator = '%s' AND partner = '%s' AND adnet = '%s' AND service = '%s' AND campaign_id = '%s'"
	DELSUMMARYCAMPAIGN            = "DELETE FROM summary_campaign WHERE id = %d"
	EDITSETTINGCAMPAIGNDETAIL     = "UPDATE campaign_detail SET po = '%s', mo_capping = %d, ratio_send = %d, ratio_receive = %d, last_update = '%s' WHERE urlservicekey = '%s' AND country = '%s' AND operator = '%s' AND partner = '%s' AND adnet = '%s' AND service = '%s' AND campaign_id = '%s'"
	EDITSETTINGSUMMARYCAMPAIGN    = "UPDATE summary_campaign SET po = '%s', mo_capping = %d, ratio_send = %d, ratio_receive = %d, last_update = '%s' WHERE id = %d"
	UPDATESTATUSCAMPAIGNDETAIL    = "UPDATE campaign_detail SET is_active = %t WHERE urlservicekey = '%s' AND country = '%s' AND operator = '%s' AND partner = '%s' AND adnet = '%s' AND service = '%s' AND campaign_id = '%s'"
	UPDATESUMMARYCAMPAIGN         = "UPDATE summary_campaign SET is_active = %t WHERE id = %d'"
	GETCAMPAIGNDETAIL             = "SELECT id, urlservicekey, is_active, counter_mo_capping, mo_capping, status_capping, counter_mo_ratio, ratio_send, ratio_receive, status_ratio, api_url, pubid, cost, po FROM campaign_detail WHERE id = %d;"
	COUNTERCAPPING                = "UPDATE campaign_detail SET counter_mo_capping = counter_mo_capping+1, last_update_capping = CASE WHEN counter_mo_capping >= mo_capping THEN '%s'::timestamp(0) END WHERE id = %d;"
	COUNTERRATIO                  = "UPDATE campaign_detail SET counter_mo_ratio = counter_mo_ratio+1 WHERE id = %d;"
	UPDATESTATUSCOUNTER           = "UPDATE campaign_detail SET counter_mo_capping = %d, status_capping = %t, counter_mo_ratio = %d, status_ratio = %t, last_update = '%s'::timestamp(0), last_update_capping = CASE WHEN counter_mo_capping+1 >= mo_capping THEN '%s'::timestamp(0) END WHERE id = %d"
)

func (r *BaseModel) GetLastCampaignId(tbl string) int {

	SQL := fmt.Sprintf(GETLASTIDCAMPDETAIL, tbl)

	var id int
	err := r.DBPostgre.QueryRow(SQL).Scan(&id)
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("GetLastCampaignId (%s) Error %s when preparing SQL statement", SQL, err))
		return 0
	}

	r.Logs.Debug(fmt.Sprintf("GetLastCampaignId (%s) found %d", SQL, id))
	return id

}

func (r *BaseModel) NewCampaign(o entity.DataCampaignAction) int {

	SQL := fmt.Sprintf(NEWCAMPAIGN, o.CampaignId, o.CampaignName, o.Objective, o.Country, o.Advertiser)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("NewCampaign (%s) Error %s when preparing SQL statement", SQL, err))

		return 0
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("NewCampaign, SQL : %s, Error %s when update to table", SQL, err))

		return 0
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("NewCampaign, SQL : %s, Error %s when finding rows affected", SQL, err))

		return 0
	}

	r.Logs.Debug(fmt.Sprintf("NewCampaign, SQL : %s, row affected : %d", SQL, rows))
	return int(rows)
}

func (r *BaseModel) GetCampaignByCampaignId(campId string) entity.DataCampaignAction {

	SQL := fmt.Sprintf(GETCAMPAIGNBYCAMPAIGNID, campId)

	var o entity.DataCampaignAction
	err := r.DBPostgre.QueryRow(SQL).Scan(&o.Id, &o.CampaignId, &o.CampaignName, &o.Objective, &o.Country, &o.Advertiser)
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("GetCampaignByCampaignId (%s) Error %s when preparing SQL statement", SQL, err))
		return entity.DataCampaignAction{}
	}

	r.Logs.Debug(fmt.Sprintf("GetCampaignByCampaignId (%s) found %#v", SQL, o))
	return o

}

func (r *BaseModel) NewCampaignDetail(o entity.DataConfig) int {

	ips := "{}"
	if len(o.IPAddress) > 0 {

		ips = "{"
		for _, v := range o.IPAddress {
			ips = ips + v
		}
		ips = ips + "}"
	}

	SQL := fmt.Sprintf(NEWCAMPAIGNDETAIL, o.URLServiceKey, o.CampaignId, o.Country, o.Operator, o.Partner, o.Aggregator, o.Adnet, o.Service, o.Keyword, o.SubKeyword, o.IsBillable, o.Plan, o.PO, o.Cost, o.PubId, o.ShortCode, o.DeviceType, o.OS, o.URLType, o.ClickType, o.ClickDelay, o.ClientType, o.TrafficSource, o.UniqueClick, o.URLBanner, o.URLLanding, o.URLWarpLanding, o.URLService, o.URLTFCSmartlink, o.GlobPost, o.URLGlobPost, o.CustomIntegration, ips, o.IsActive, o.MOCapping, o.CounterMOCapping, o.StatusCapping, o.KPIUpperLimitCapping, o.IsMachineLearningCapping, o.RatioSend, o.RatioReceive, o.CounterMORatio, o.StatusRatio, o.KPIUpperLimitRatioSend, o.KPIUpperLimitRatioReceive, o.IsMachineLearningRatio, o.APIURL, o.LastUpdate, o.LastUpdateCapping)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("NewCampaignDetail (%s) Error %s when preparing SQL statement", SQL, err))

		return 0
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("NewCampaignDetail, SQL : %s, Error %s when update to table", SQL, err))

		return 0
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("NewCampaignDetail, SQL : %s, Error %s when finding rows affected", SQL, err))

		return 0
	}

	r.Logs.Debug(fmt.Sprintf("NewCampaignDetail, SQL : %s, row affected : %d", SQL, rows))
	return int(rows)
}

func (r *BaseModel) GetCampaignByCampaignDetailId(o entity.DataConfig) entity.DataConfig {

	SQL := fmt.Sprintf(GETCAMPAIGNBYCAMPAIGNDETAILID, o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Adnet, o.Service)

	err := r.DBPostgre.QueryRow(SQL).Scan(&o.Id, &o.URLServiceKey, &o.CampaignId, &o.Country, &o.Operator, &o.Partner, &o.Aggregator, &o.Adnet, &o.Service, &o.Keyword, &o.SubKeyword, &o.IsBillable, &o.Plan, &o.PO, &o.Cost, &o.PubId, &o.ShortCode, &o.DeviceType, &o.OS, &o.URLType, &o.ClickType, &o.ClickDelay, &o.ClientType, &o.TrafficSource, &o.UniqueClick, &o.URLBanner, &o.URLLanding, &o.URLWarpLanding, &o.URLService, &o.URLTFCSmartlink, &o.GlobPost, &o.URLGlobPost, &o.CustomIntegration, &o.IPAddress, &o.IsActive, &o.MOCapping, &o.CounterMOCapping, &o.StatusCapping, &o.KPIUpperLimitCapping, &o.IsMachineLearningCapping, &o.RatioSend, &o.RatioReceive, &o.CounterMORatio, &o.StatusRatio, &o.KPIUpperLimitRatioSend, &o.KPIUpperLimitRatioReceive, &o.IsMachineLearningRatio, &o.APIURL, &o.LastUpdate, &o.LastUpdateCapping)
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("GetCampaignByCampaignDetailId (%s) Error %s when preparing SQL statement", SQL, err))
		return o
	}

	r.Logs.Debug(fmt.Sprintf("GetCampaignByCampaignDetailId (%s) found %#v", SQL, o))
	return o

}

func (r *BaseModel) UpdateCampaign(o entity.DataCampaignAction) error {

	SQL := fmt.Sprintf(UPDATECAMPAIGN, o.CampaignName, o.Objective, o.Country, o.Advertiser, o.CampaignId)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("UpdateCampaign (%s) Error %s when preparing SQL statement", SQL, err))

		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("UpdateCampaign, SQL : %s, Error %s when update to table", SQL, err))

		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("UpdateCampaign, SQL : %s, Error %s when finding rows affected", SQL, err))

		return err
	}

	r.Logs.Debug(fmt.Sprintf("UpdateCampaign, SQL : %s, row affected : %d", SQL, rows))
	return nil
}

func (r *BaseModel) UpdateCampaignDetail(o entity.DataConfig) error {

	ips := "{}"
	if len(o.IPAddress) > 0 {

		ips = "{"
		for _, v := range o.IPAddress {
			ips = ips + v
		}
		ips = ips + "}"
	}

	SQL := fmt.Sprintf(UPDATECAMPAIGNDETAIL, o.CampaignId, o.Country, o.Operator, o.Partner, o.Aggregator, o.Adnet, o.Service, o.Keyword, o.SubKeyword, o.IsBillable, o.Plan, o.PO, o.Cost, o.PubId, o.ShortCode, o.DeviceType, o.OS, o.URLType, o.ClickType, o.ClickDelay, o.ClientType, o.TrafficSource, o.UniqueClick, o.URLBanner, o.URLLanding, o.URLWarpLanding, o.URLService, o.URLTFCSmartlink, o.GlobPost, o.URLGlobPost, o.CustomIntegration, ips, o.IsActive, o.MOCapping, o.CounterMOCapping, o.StatusCapping, o.KPIUpperLimitCapping, o.IsMachineLearningCapping, o.RatioSend, o.RatioReceive, o.CounterMORatio, o.StatusRatio, o.KPIUpperLimitRatioSend, o.KPIUpperLimitRatioReceive, o.IsMachineLearningRatio, o.APIURL, o.LastUpdate, o.LastUpdateCapping, o.Id)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("NewCampaignDetail (%s) Error %s when preparing SQL statement", SQL, err))

		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("NewCampaignDetail, SQL : %s, Error %s when update to table", SQL, err))

		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("NewCampaignDetail, SQL : %s, Error %s when finding rows affected", SQL, err))

		return err
	}

	r.Logs.Debug(fmt.Sprintf("NewCampaignDetail, SQL : %s, row affected : %d", SQL, rows))
	return nil
}

func (r *BaseModel) DelCampaign(o entity.DataCampaignAction) error {

	SQL := fmt.Sprintf(DELCAMPAIGN, o.CampaignId)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("DelCampaign (%s) Error %s when preparing SQL statement", SQL, err))

		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("DelCampaign, SQL : %s, Error %s when update to table", SQL, err))

		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("DelCampaign, SQL : %s, Error %s when finding rows affected", SQL, err))

		return err
	}

	r.Logs.Debug(fmt.Sprintf("DelCampaign, SQL : %s, row affected : %d", SQL, rows))
	return nil
}

func (r *BaseModel) DelCampaignDetail(o entity.DataConfig) error {

	SQL := fmt.Sprintf(DELCAMPAIGNDETAIL, o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("DelCampaignDetail (%s) Error %s when preparing SQL statement", SQL, err))

		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("DelCampaignDetail, SQL : %s, Error %s when update to table", SQL, err))

		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("DelCampaignDetail, SQL : %s, Error %s when finding rows affected", SQL, err))

		return err
	}

	r.Logs.Debug(fmt.Sprintf("DelCampaignDetail, SQL : %s, row affected : %d", SQL, rows))
	return nil
}

func (r *BaseModel) DelSummaryCampaign(o entity.DataConfig) error {

	SQL := fmt.Sprintf(DELSUMMARYCAMPAIGN, o.Id)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("DelSummaryCampaign (%s) Error %s when preparing SQL statement", SQL, err))

		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("DelSummaryCampaign, SQL : %s, Error %s when update to table", SQL, err))

		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("DelSummaryCampaign, SQL : %s, Error %s when finding rows affected", SQL, err))

		return err
	}

	r.Logs.Debug(fmt.Sprintf("DelSummaryCampaign, SQL : %s, row affected : %d", SQL, rows))
	return nil
}

func (r *BaseModel) EditSettingCampaignDetail(o entity.DataConfig) error {

	SQL := fmt.Sprintf(EDITSETTINGCAMPAIGNDETAIL, o.PO, o.MOCapping, o.RatioSend, o.RatioReceive, o.LastUpdate, o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("EditSettingCampaignDetail (%s) Error %s when preparing SQL statement", SQL, err))

		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("EditSettingCampaignDetail, SQL : %s, Error %s when update to table", SQL, err))

		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("EditSettingCampaignDetail, SQL : %s, Error %s when finding rows affected", SQL, err))

		return err
	}

	r.Logs.Debug(fmt.Sprintf("EditSettingCampaignDetail, SQL : %s, row affected : %d", SQL, rows))
	return nil
}

func (r *BaseModel) EditSettingSummaryCampaign(o entity.DataConfig) error {

	SQL := fmt.Sprintf(EDITSETTINGSUMMARYCAMPAIGN, o.PO, o.MOCapping, o.RatioSend, o.RatioReceive, o.LastUpdate, o.Id)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("EditSettingSummaryCampaign (%s) Error %s when preparing SQL statement", SQL, err))

		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("EditSettingSummaryCampaign, SQL : %s, Error %s when update to table", SQL, err))

		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("EditSettingSummaryCampaign, SQL : %s, Error %s when finding rows affected", SQL, err))

		return err
	}

	r.Logs.Debug(fmt.Sprintf("EditSettingSummaryCampaign, SQL : %s, row affected : %d", SQL, rows))
	return nil
}

func (r *BaseModel) UpdateStatusCampaignDetail(o entity.DataConfig) error {

	SQL := fmt.Sprintf(UPDATESTATUSCAMPAIGNDETAIL, o.IsActive, o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("UpdateStatusCampaignDetail (%s) Error %s when preparing SQL statement", SQL, err))

		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("UpdateStatusCampaignDetail, SQL : %s, Error %s when update to table", SQL, err))

		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("UpdateStatusCampaignDetail, SQL : %s, Error %s when finding rows affected", SQL, err))

		return err
	}

	r.Logs.Debug(fmt.Sprintf("UpdateStatusCampaignDetail, SQL : %s, row affected : %d", SQL, rows))
	return nil
}

func (r *BaseModel) UpdateSummaryCampaign(o entity.DataConfig) error {

	SQL := fmt.Sprintf(UPDATESUMMARYCAMPAIGN, o.IsActive, o.Id)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("UpdateSummaryCampaign (%s) Error %s when preparing SQL statement", SQL, err))

		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("UpdateSummaryCampaign, SQL : %s, Error %s when update to table", SQL, err))

		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("UpdateSummaryCampaign, SQL : %s, Error %s when finding rows affected", SQL, err))

		return err
	}

	r.Logs.Debug(fmt.Sprintf("UpdateSummaryCampaign, SQL : %s, row affected : %d", SQL, rows))
	return nil
}

func (r *BaseModel) GetCampaignDetail(o entity.DataConfig) (entity.DataConfig, error) {

	SQL := fmt.Sprintf(GETCAMPAIGNDETAIL, o.Id)
	rows, err := r.DBPostgre.Query(SQL)
	if err != nil {
		r.Logs.Error(fmt.Sprintf("GetCampaignDetail, SQL : %s, error querying occured : %#v", SQL, err))

		return entity.DataConfig{}, err
	}
	defer rows.Close()

	for rows.Next() {

		err = rows.Scan(&o.Id, &o.URLServiceKey, &o.IsActive, &o.CounterMOCapping, &o.MOCapping, &o.StatusCapping, &o.CounterMORatio, &o.RatioSend, &o.RatioReceive, &o.StatusRatio, &o.APIURL, &o.PubId, &o.Cost, &o.PO)

		if err != nil {

			r.Logs.Error(fmt.Sprintf("GetCampaignDetail, SQL : %s, error scan occured : %#v", SQL, err))

		}
	}

	r.Logs.Info(fmt.Sprintf("GetCampaignDetail, SQL : %s, row selected occured : %#v", SQL, o))
	return o, nil
}

func (r *BaseModel) CounterCappingById(o entity.DataConfig) error {

	SQL := fmt.Sprintf(COUNTERCAPPING, o.LastUpdateCapping, o.Id)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("CounterCappingById (%s) Error %s when preparing SQL statement", SQL, err))

		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("CounterCappingById, SQL : %s, Error %s when update to table", SQL, err))

		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("CounterCappingById, SQL : %s, Error %s when finding rows affected", SQL, err))

		return err
	}

	r.Logs.Debug(fmt.Sprintf("CounterCappingById, SQL : %s, row affected : %d", SQL, rows))
	return nil
}

func (r *BaseModel) CounterRatioById(id int) error {

	SQL := fmt.Sprintf(COUNTERRATIO, id)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("CounterRatioById (%s) Error %s when preparing SQL statement", SQL, err))

		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("CounterRatioById, SQL : %s, Error %s when update to table", SQL, err))

		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("CounterRatioById, SQL : %s, Error %s when finding rows affected", SQL, err))

		return err
	}

	r.Logs.Debug(fmt.Sprintf("CounterRatioById, SQL : %s, row affected : %d", SQL, rows))
	return nil
}

func (r *BaseModel) UpdateStatusCounterById(o entity.DataConfig) error {

	SQL := fmt.Sprintf(UPDATESTATUSCOUNTER, o.CounterMOCapping, o.StatusCapping, o.CounterMORatio, o.StatusRatio, o.LastUpdate, o.LastUpdateCapping, o.Id)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("UpdateStatusCounterById (%s) Error %s when preparing SQL statement", SQL, err))

		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("UpdateStatusCounterById, SQL : %s, Error %s when update to table", SQL, err))

		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("UpdateStatusCounterById, SQL : %s, Error %s when finding rows affected", SQL, err))

		return err
	}

	r.Logs.Debug(fmt.Sprintf("UpdateStatusCounterById, SQL : %s, row affected : %d", SQL, rows))
	return nil
}
