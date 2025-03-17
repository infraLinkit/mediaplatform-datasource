package model

import (
	"fmt"

	"github.com/infraLinkit/mediaplatform-datasource/entity"
)

func (r *BaseModel) NewPixel(o entity.PixelStorage) int {

	result := r.DB.Create(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return int(o.ID)
}

func (r *BaseModel) GetPx(o entity.PixelStorage) (entity.PixelStorage, error) {

	result := r.DB.Exec("SELECT * FROM (SELECT id, campaign_detail_id, pxdate, url_service_key, campaign_id, country, partner, operator, aggregator, service, short_code, adnet, keyword, subkeyword, is_billable, plan, url, url_type, pixel, trx_id, token, msisdn, is_used, browser, os, ip, isp, referral_url, pub_id, user_agent, traffic_source, traffic_source_data, user_rejected, user_duplicated, handset, handset_code, handset_type, url_landing, url_warp_landing, url_service, url_tfcor_smartlink, po, cost, is_unique, campaign_objective FROM pixel_storages WHERE url_service_key = ? AND pixel = ?) AS px ORDER by px.pxdate DESC", o.URLServiceKey, o.Pixel).Scan(&o)

	return o, result.Error
}

func (r *BaseModel) GetToken(o entity.PixelStorage) (entity.PixelStorage, error) {

	result := r.DB.Exec("SELECT * FROM (SELECT id, campaign_detail_id, pxdate, url_service_key, campaign_id, country, partner, operator, aggregator, service, short_code, adnet, keyword, subkeyword, is_billable, plan, url, url_type, pixel, trx_id, token, msisdn, is_used, browser, os, ip, isp, referral_url, pub_id, user_agent, traffic_source, traffic_source_data, user_rejected, user_duplicated, handset, handset_code, handset_type, url_landing, url_warp_landing, url_service, url_tfcor_smartlink, po, cost, is_unique, campaign_objective FROM pixel_storages WHERE url_service_key = ? AND token = ?) AS px ORDER BY px.pxdate DESC", o.URLServiceKey, o.IsUsed).Scan(&o)

	return o, result.Error
}

func (r *BaseModel) GetByAdnetCode(o entity.PixelStorage) (entity.PixelStorage, error) {

	result := r.DB.Exec("SELECT * FROM (SELECT id, campaign_detail_id, pxdate, url_service_key, campaign_id, country, partner, operator, aggregator, service, short_code, adnet, keyword, subkeyword, is_billable, plan, url, url_type, pixel, trx_id, token, msisdn, is_used, browser, os, ip, isp, referral_url, pub_id, user_agent, traffic_source, traffic_source_data, user_rejected, user_duplicated, handset, handset_code, handset_type, url_landing, url_warp_landing, url_service, url_tfcor_smartlink, po, cost, is_unique, campaign_objective FROM pixel_storages WHERE url_service_key = ? AND is_used = ?) AS px ORDER BY px.pxdate ASC", o.URLServiceKey, o.IsUsed).Scan(&o)

	return o, result.Error
}

func (r *BaseModel) UpdatePixelById(o entity.PixelStorage) error {

	result := r.DB.Model(&o).Select("msisdn = ?, trx_id = ?, is_used = ?, pixel_used_date = ?, status_postback = ?, is_unique = ?, url_postback = ?, status_url_postback = ?, reason_url_postback = ?", o.Msisdn, o.TrxId, o.IsUsed, o.PixelUsedDate, o.StatusPostback, o.IsUnique, o.URLPostback, o.StatusURLPostback, o.ReasonURLPostback)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

/* package model

import (
	"context"
	"fmt"

	"github.com/infraLinkit/mediaplatform-datasource/entity"

	_ "github.com/lib/pq"
)

const (
	NEWPX        = "INSERT INTO pixel_storage (id, campaign_detail_id, pxdate, urlservicekey, campaign_id, country, partner, operator, aggregator, service, short_code, adnet, keyword, subkeyword, is_billable, plan, url, url_type, pixel, trx_id, token, msisdn, is_used, browser, os, ip, isp, referral_url, pubid, user_agent, traffic_source, traffic_source_data, user_rejected, user_duplicated, handset, handset_code, handset_type, url_landing, url_warp_landing, url_service, url_tfc_or_smartlink, po, cost) VALUES (DEFAULT, %d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, '%s', %t, %t, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')"
	GETPX        = "SELECT * FROM (SELECT id, campaign_detail_id, pxdate, urlservicekey, campaign_id, country, partner, operator, aggregator, service, short_code, adnet, keyword, subkeyword, is_billable, plan, url, url_type, pixel, trx_id, token, msisdn, is_used, browser, os, ip, isp, referral_url, pubid, user_agent, traffic_source, traffic_source_data, user_rejected, user_duplicated, handset, handset_code, handset_type, url_landing, url_warp_landing, url_service, url_tfc_or_smartlink, po, cost, is_unique FROM pixel_storage WHERE urlservicekey = '%s' AND pixel = '%s') AS px ORDER by px.pxdate DESC LIMIT 1"
	GETTOKEN     = "SELECT * FROM (SELECT id, campaign_detail_id, pxdate, urlservicekey, campaign_id, country, partner, operator, aggregator, service, short_code, adnet, keyword, subkeyword, is_billable, plan, url, url_type, pixel, trx_id, token, msisdn, is_used, browser, os, ip, isp, referral_url, pubid, user_agent, traffic_source, traffic_source_data, user_rejected, user_duplicated, handset, handset_code, handset_type, url_landing, url_warp_landing, url_service, url_tfc_or_smartlink, po, cost, is_unique FROM pixel_storage WHERE urlservicekey = '%s' AND token = '%s') AS px ORDER BY px.pxdate DESC LIMIT 1"
	GETADNETCODE = "SELECT * FROM (SELECT id, campaign_detail_id, pxdate, urlservicekey, campaign_id, country, partner, operator, aggregator, service, short_code, adnet, keyword, subkeyword, is_billable, plan, url, url_type, pixel, trx_id, token, msisdn, is_used, browser, os, ip, isp, referral_url, pubid, user_agent, traffic_source, traffic_source_data, user_rejected, user_duplicated, handset, handset_code, handset_type, url_landing, url_warp_landing, url_service, url_tfc_or_smartlink, po, cost, is_unique FROM pixel_storage WHERE urlservicekey = '%s' AND is_used = false) AS px ORDER BY px.pxdate ASC LIMIT 1"
	UPDATEPX     = "UPDATE pixel_storage SET msisdn = '%s', trx_id = '%s', is_used = %t, pixel_used_date = '%s', status_postback = %t, is_unique = %t, url_postback = '%s', status_url_postback = '%s', reason_url_postback = '%s' WHERE id = %d"
)

func (r *BaseModel) NewPixel(o entity.PixelStorage) error {

	SQL := fmt.Sprintf(NEWPX, o.CampaignDetailId, o.PxDate, o.URLServiceKey, o.CampaignId, o.Country, o.Partner, o.Operator, o.Aggregator, o.Service, o.ShortCode, o.Adnet, o.Keyword, o.Subkeyword, o.IsBillable, o.Plan, o.URL, o.URLType, o.Pixel, o.TrxId, o.Token, o.Msisdn, o.IsUsed, o.Browser, o.OS, o.IP, o.ISP, o.ReferralURL, o.PubId, o.UserAgent, o.TrafficSource, o.TrafficSourceData, o.UserRejected, o.UserDuplicated, o.Handset, o.HandsetCode, o.HandsetType, o.URLLanding, o.URLWarpLanding, o.URLService, o.URLTFCSmartlink, o.PO, o.Cost)

	//L.Write(L.LogName, "debug", fmt.Sprintf("NewSubs (%s)", SQL))

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("(%s) Error %s when preparing SQL statement", SQL, err))

		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("Error %s when insert to table", err))

		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("Error %s when finding rows affected", err))

		return err
	}

	r.Logs.Debug(fmt.Sprintf("SQL : %s, row affected : %d", SQL, rows))
	return nil
}

func (r *BaseModel) GetPx(o entity.PixelStorage) (entity.PixelStorage, error) {

	SQL := fmt.Sprintf(GETPX, o.URLServiceKey, o.Pixel)
	rows, err := r.DBPostgre.Query(SQL)
	if err != nil {
		r.Logs.Error(fmt.Sprintf("SQL : %s, error querying occured : %#v", SQL, err))
		return entity.PixelStorage{}, err
	}
	defer rows.Close()

	var px entity.PixelStorage

	for rows.Next() {

		err = rows.Scan(&px.Id, &px.CampaignDetailId, &px.PxDate, &px.URLServiceKey, &px.CampaignId, &px.Country, &px.Partner, &px.Operator, &px.Aggregator, &px.Service, &px.ShortCode, &px.Adnet, &px.Keyword, &px.Subkeyword, &px.IsBillable, &px.Plan, &px.URL, &px.URLType, &px.Pixel, &px.TrxId, &px.Token, &px.Msisdn, &px.IsUsed, &px.Browser, &px.OS, &px.IP, &px.ISP, &px.ReferralURL, &px.PubId, &px.UserAgent, &px.TrafficSource, &px.TrafficSourceData, &px.UserRejected, &px.UserDuplicated, &px.Handset, &px.HandsetCode, &px.HandsetType, &px.URLLanding, &px.URLWarpLanding, &px.URLService, &px.URLTFCSmartlink, &px.PO, &px.Cost, &px.IsUnique)

		if err != nil {

			r.Logs.Error(fmt.Sprintf("SQL : %s, error scan occured : %#v", SQL, err))
			return entity.PixelStorage{}, err
		}
	}

	r.Logs.Info(fmt.Sprintf("SQL : %s, row selected occured : %#v", SQL, px))
	return px, nil
}

func (r *BaseModel) GetToken(o entity.PixelStorage) (entity.PixelStorage, error) {

	SQL := fmt.Sprintf(GETTOKEN, o.URLServiceKey, o.Pixel)
	rows, err := r.DBPostgre.Query(SQL)
	if err != nil {
		r.Logs.Error(fmt.Sprintf("SQL : %s, error querying occured : %#v", SQL, err))
		return entity.PixelStorage{}, err
	}
	defer rows.Close()

	var px entity.PixelStorage

	for rows.Next() {

		err = rows.Scan(&px.Id, &px.CampaignDetailId, &px.PxDate, &px.URLServiceKey, &px.CampaignId, &px.Country, &px.Partner, &px.Operator, &px.Aggregator, &px.Service, &px.ShortCode, &px.Adnet, &px.Keyword, &px.Subkeyword, &px.IsBillable, &px.Plan, &px.URL, &px.URLType, &px.Pixel, &px.TrxId, &px.Token, &px.Msisdn, &px.IsUsed, &px.Browser, &px.OS, &px.IP, &px.ISP, &px.ReferralURL, &px.PubId, &px.UserAgent, &px.TrafficSource, &px.TrafficSourceData, &px.UserRejected, &px.UserDuplicated, &px.Handset, &px.HandsetCode, &px.HandsetType, &px.URLLanding, &px.URLWarpLanding, &px.URLService, &px.URLTFCSmartlink, &px.PO, &px.Cost, &px.IsUnique)

		if err != nil {

			r.Logs.Error(fmt.Sprintf("SQL : %s, error scan occured : %#v", SQL, err))
			return entity.PixelStorage{}, err
		}
	}

	r.Logs.Info(fmt.Sprintf("SQL : %s, row selected occured : %#v", SQL, px))
	return px, nil
}

func (r *BaseModel) GetByAdnetCode(o entity.PixelStorage) (entity.PixelStorage, error) {

	SQL := fmt.Sprintf(GETADNETCODE, o.URLServiceKey)
	rows, err := r.DBPostgre.Query(SQL)
	if err != nil {
		r.Logs.Error(fmt.Sprintf("SQL : %s, error querying occured : %#v", SQL, err))
		return entity.PixelStorage{}, err
	}
	defer rows.Close()

	var px entity.PixelStorage

	for rows.Next() {

		err = rows.Scan(&px.Id, &px.CampaignDetailId, &px.PxDate, &px.URLServiceKey, &px.CampaignId, &px.Country, &px.Partner, &px.Operator, &px.Aggregator, &px.Service, &px.ShortCode, &px.Adnet, &px.Keyword, &px.Subkeyword, &px.IsBillable, &px.Plan, &px.URL, &px.URLType, &px.Pixel, &px.TrxId, &px.Token, &px.Msisdn, &px.IsUsed, &px.Browser, &px.OS, &px.IP, &px.ISP, &px.ReferralURL, &px.PubId, &px.UserAgent, &px.TrafficSource, &px.TrafficSourceData, &px.UserRejected, &px.UserDuplicated, &px.Handset, &px.HandsetCode, &px.HandsetType, &px.URLLanding, &px.URLWarpLanding, &px.URLService, &px.URLTFCSmartlink, &px.PO, &px.Cost, &px.IsUnique)

		if err != nil {

			r.Logs.Error(fmt.Sprintf("SQL : %s, error scan occured : %#v", SQL, err))
			return entity.PixelStorage{}, err
		}
	}

	r.Logs.Info(fmt.Sprintf("SQL : %s, row selected occured : %#v", SQL, px))
	return px, nil
}

func (r *BaseModel) UpdatePixelById(o entity.PixelStorage) error {

	SQL := fmt.Sprintf(UPDATEPX, o.Msisdn, o.TrxId, o.IsUsed, o.PixelUsedDate, o.StatusPostback, o.IsUnique, o.URLPostback, o.StatusURLPostback, o.ReasonURLPostback, o.Id)

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
