package model

import (
	"context"
	"fmt"

	"github.com/infraLinkit/mediaplatform-datasource/entity"
)

const (
	NEWMO = "INSERT INTO mo (id, campaign_detail_id, pxdate, urlservicekey, campaign_id, country, partner, operator, aggregator, service, short_code, adnet, keyword, subkeyword, is_billable, plan, url, url_type, pixel, trx_id, token, msisdn, is_used, browser, os, ip, isp, referral_url, pubid, user_agent, traffic_source, traffic_source_data, user_rejected, user_duplicated, handset, handset_code, handset_type, url_landing, url_warp_landing, url_service, url_tfc_or_smartlink, po, cost, is_active, counter_mo_capping, mo_capping, status_capping, counter_mo_ratio, ratio_send, ratio_receive, status_ratio, api_url, pubid, status_postback, is_unique, url_postback, status_url_postback, reason_url_postback) VALUES (DEFAULT, %d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, '%s', %t, %t, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, %d, %d, %t, %d, %d, %d, %t, '%s', '%s', %t, %t, '%s', '%s', '%s')"
)

func (r *BaseModel) NewMO(o entity.PixelStorage, cd entity.DataConfig) error {

	SQL := fmt.Sprintf(NEWMO, o.CampaignDetailId, o.PxDate, o.URLServiceKey, o.CampaignId, o.Country, o.Partner, o.Operator, o.Aggregator, o.Service, o.ShortCode, o.Adnet, o.Keyword, o.Subkeyword, o.IsBillable, o.Plan, o.URL.String, o.URLType, o.Pixel, o.TrxId, o.Token, o.Msisdn, o.IsUsed, o.Browser, o.OS, o.IP, o.ISP, o.ReferralURL.String, o.PubId, o.UserAgent, o.TrafficSource, o.TrafficSourceData, o.UserRejected, o.UserDuplicated, o.Handset, o.HandsetCode, o.HandsetType, o.URLLanding.String, o.URLWarpLanding.String, o.URLService.String, o.URLTFCSmartlink.String, o.PO, o.Cost, cd.IsActive, cd.CounterMOCapping, cd.MOCapping, cd.StatusCapping, cd.CounterMORatio, cd.RatioSend, cd.RatioReceive, cd.StatusRatio, cd.APIURL.String, cd.PubId, o.StatusPostback, o.IsUnique, o.URLPostback.String, o.StatusURLPostback, o.ReasonURLPostback)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("(%s) Error %s when preparing SQL statement", SQL, err))

		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("(%s) Error %s when insert to table", SQL, err))

		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("(%s) Error %s when finding rows affected", SQL, err))

		return err
	}

	r.Logs.Debug(fmt.Sprintf("SQL : %s, row affected : %d", SQL, rows))
	return nil
}
