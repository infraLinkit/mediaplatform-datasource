package model

import (
	"context"
	"fmt"

	"github.com/infraLinkit/mediaplatform-datasource/entity"

	_ "github.com/lib/pq"
)

const (
	NEWPX       = "INSERT INTO pixel_storage (id, campaign_detail_id, pxdate, urlservicekey, campaign_id, country, partner, operator, aggregator, service, short_code, adnet, keyword, subkeyword, is_billable, plan, url, url_type, pixel, trx_id, token, msisdn, is_used, browser, os, ip, isp, referral_url, pubid, user_agent, traffic_source, traffic_source_data, user_rejected, user_duplicated, handset, handset_code, handset_type, url_landing, url_warp_landing, url_service, url_tfc_or_smartlink, po, cost) VALUES (DEFAULT, %d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, '%s', %t, %t, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')"
	GETPX       = "SELECT id, campaign_detail_id, pxdate, urlservicekey, campaign_id, country, partner, operator, aggregator, service, short_code, adnet, keyword, subkeyword, is_billable, plan, url, url_type, pixel, trx_id, token, msisdn, is_used, browser, os, ip, isp, referral_url, pubid, user_agent, traffic_source, traffic_source_data, user_rejected, user_duplicated, handset, handset_code, handset_type, url_landing, url_warp_landing, url_service, url_tfc_or_smartlink, po, cost FROM pixel_storage WHERE country = '%s' AND operator = '%s' AND partner = '%s' AND service = '%s' AND keyword = '%s' AND is_billable = %t AND pixel = '%s'"
	UPDATEPX    = "UPDATE pixel_storage SET msisdn = '%s', trx_id = '%s', is_used = %t, pixel_used_date = '%s', status_postback = %t, is_unique = %t, url_postback = '%s', status_url_postback = '%s', reason_url_postback = '%s' WHERE id = %d"
	NEWPOSTBACK = "INSERT INTO postback (id, campaign_detail_id, pxdate, urlservicekey, campaign_id, country, partner, operator, aggregator, service, short_code, adnet, keyword, subkeyword, is_billable, plan, url, url_type, pixel, trx_id, token, msisdn, is_used, browser, os, ip, isp, referral_url, pubid, user_agent, traffic_source, traffic_source_data, user_rejected, user_duplicated, handset, handset_code, handset_type, url_landing, url_warp_landing, url_service, url_tfc_or_smartlink, po, cost, is_active, counter_mo_capping, mo_capping, status_capping, counter_mo_ratio, ratio_send, ratio_receive, status_ratio, api_url, pubid, status_postback, is_unique, url_postback, status_url_postback, reason_url_postback) VALUES (DEFAULT, %d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, '%s', %t, %t, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, %d, %d, %t, %d, %d, %d, %t, '%s', '%s', %t, %t, '%s', '%s', '%s')"
)

func (r *BaseModel) NewPixel(o entity.PixelStorage) error {

	SQL := fmt.Sprintf(NEWPX, o.CampaignDetailId, o.PxDate, o.URLServiceKey, o.CampaignId, o.Country, o.Partner, o.Operator, o.Aggregator, o.Service, o.ShortCode, o.Adnet, o.Keyword, o.Subkeyword, o.IsBillable, o.Plan, o.URL, o.URLType, o.Pixel, o.TrxId, o.Token, o.Msisdn, o.IsUsed, o.Browser, o.OS, o.IP, o.ISP, o.ReferralURL, o.PubId, o.UserAgent, o.TrafficSource, o.TrafficSourceData, o.UserRejected, o.UserDuplicated, o.Handset, o.HandsetCode, o.HandsetType, o.URLLanding, o.URLWarpLanding, o.URLService, o.URLTFCSmartlink, o.PO, o.Cost)

	//L.Write(L.LogName, "debug", fmt.Sprintf("NewSubs (%s)", SQL))

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("NewPixel (%s) Error %s when preparing SQL statement", SQL, err))

		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("NewPixel Error %s when insert to table", err))

		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("NewPixel Error %s when finding rows affected", err))

		return err
	}

	r.Logs.Debug(fmt.Sprintf("NewPixel, SQL : %s, row affected : %d", SQL, rows))
	return nil
}

func (r *BaseModel) GetPx(o entity.PixelStorage) (entity.PixelStorage, error) {

	SQL := fmt.Sprintf(GETPX, o.Country, o.Operator, o.Partner, o.Service, o.Keyword, o.IsBillable, o.Pixel)
	rows, err := r.DBPostgre.Query(SQL)
	if err != nil {
		r.Logs.Error(fmt.Sprintf("GetPx, SQL : %s, error querying occured : %#v", SQL, err))
		return entity.PixelStorage{}, err
	}
	defer rows.Close()

	var px entity.PixelStorage

	for rows.Next() {

		err = rows.Scan(&px.Id, &px.CampaignDetailId, &px.PxDate, &px.URLServiceKey, &px.CampaignId, &px.Country, &px.Partner, &px.Operator, &px.Aggregator, &px.Service, &px.ShortCode, &px.Adnet, &px.Keyword, &px.Subkeyword, &px.IsBillable, &px.Plan, &px.URL, &px.URLType, &px.Pixel, &px.TrxId, &px.Token, &px.Msisdn, &px.IsUsed, &px.Browser, &px.OS, &px.IP, &px.ISP, &px.ReferralURL, &px.PubId, &px.UserAgent, &px.TrafficSource, &px.TrafficSourceData, &px.UserRejected, &px.UserDuplicated, &px.Handset, &px.HandsetCode, &px.HandsetType, &px.URLLanding, &px.URLWarpLanding, &px.URLService, &px.URLTFCSmartlink, &px.PO, &px.Cost)

		if err != nil {

			r.Logs.Error(fmt.Sprintf("GetPx, SQL : %s, error scan occured : %#v", SQL, err))
			return entity.PixelStorage{}, err
		}
	}

	r.Logs.Info(fmt.Sprintf("GetPx, SQL : %s, row selected occured : %#v", SQL, px))
	return px, nil
}

func (r *BaseModel) UpdatePixelById(o entity.PixelStorage) error {

	SQL := fmt.Sprintf(UPDATEPX, o.Msisdn, o.TrxId, o.IsUsed, o.PixelUsedDate, o.StatusPostback, o.IsUnique, o.URLPostback, o.StatusURLPostback, o.ReasonURLPostback, o.Id)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("UpdatePixelById (%s) Error %s when preparing SQL statement", SQL, err))

		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("UpdatePixelById, SQL : %s, Error %s when update to table", SQL, err))

		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("UpdatePixelById, SQL : %s, Error %s when finding rows affected", SQL, err))

		return err
	}

	r.Logs.Debug(fmt.Sprintf("UpdatePixelById, SQL : %s, row affected : %d", SQL, rows))
	return nil
}

func (r *BaseModel) NewPostback(o entity.PixelStorage, cd entity.DataConfig) error {

	SQL := fmt.Sprintf(NEWPOSTBACK, o.CampaignDetailId, o.PxDate, o.URLServiceKey, o.CampaignId, o.Country, o.Partner, o.Operator, o.Aggregator, o.Service, o.ShortCode, o.Adnet, o.Keyword, o.Subkeyword, o.IsBillable, o.Plan, o.URL, o.URLType, o.Pixel, o.TrxId, o.Token, o.Msisdn, o.IsUsed, o.Browser, o.OS, o.IP, o.ISP, o.ReferralURL, o.PubId, o.UserAgent, o.TrafficSource, o.TrafficSourceData, o.UserRejected, o.UserDuplicated, o.Handset, o.HandsetCode, o.HandsetType, o.URLLanding, o.URLWarpLanding, o.URLService, o.URLTFCSmartlink, o.PO, o.Cost, cd.IsActive, cd.CounterMOCapping, cd.MOCapping, cd.StatusCapping, cd.CounterMORatio, cd.RatioSend, cd.RatioReceive, cd.StatusRatio, cd.APIURL, cd.PubId, o.StatusPostback, o.IsUnique, o.URLPostback, o.StatusURLPostback, o.ReasonURLPostback)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("NewPostback (%s) Error %s when preparing SQL statement", SQL, err))

		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("NewPostback Error %s when insert to table", err))

		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("NewPostback Error %s when finding rows affected", err))

		return err
	}

	r.Logs.Debug(fmt.Sprintf("NewPostback, SQL : %s, row affected : %d", SQL, rows))
	return nil
}
