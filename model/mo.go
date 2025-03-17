package model

import (
	"fmt"

	"github.com/infraLinkit/mediaplatform-datasource/entity"
)

func (r *BaseModel) NewMO(o entity.MO) int {

	result := r.DB.Create(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return int(o.ID)
}

func (r *BaseModel) UpdateMO(o entity.MO) error {

	result := r.DB.Model(&o).
		Where("DATE(pxdate) = ? AND url_service_key = ? AND country = ? AND operator = ? AND service = ? AND adnet = ? AND pixel = ?", o.Pxdate, o.URLServiceKey, o.Country, o.Operator, o.Service, o.Adnet, o.Pixel).
		Updates(entity.MO{Msisdn: o.Msisdn, TrxId: o.TrxId, IsUsed: o.IsUsed, PixelUsedDate: o.PixelUsedDate, StatusPostback: o.StatusPostback, IsUnique: o.IsUnique, URLPostback: o.URLPostback, StatusURLPostback: o.StatusURLPostback, ReasonURLPostback: o.ReasonURLPostback})

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

/*
package model

import (
	"context"
	"fmt"

	"github.com/infraLinkit/mediaplatform-datasource/entity"
)

const (
	NEWMO = "INSERT INTO mo (id, campaign_detail_id, pxdate, urlservicekey, campaign_id, country, partner, operator, aggregator, service, short_code, adnet, keyword, subkeyword, is_billable, plan, url, url_type, pixel, trx_id, token, msisdn, is_used, browser, os, ip, isp, referral_url, pubid, user_agent, traffic_source, traffic_source_data, user_rejected, user_duplicated, handset, handset_code, handset_type, url_landing, url_warp_landing, url_service, url_tfc_or_smartlink, po, cost, is_active, counter_mo_capping, mo_capping, status_capping, counter_mo_ratio, ratio_send, ratio_receive, status_ratio, api_url, status_postback, is_unique, url_postback, status_url_postback, reason_url_postback) VALUES (DEFAULT, %d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, '%s', %t, %t, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, %d, %d, %t, %d, %d, %d, %t, '%s', %t, %t, '%s', '%s', '%s')"
)

func (r *BaseModel) NewMO(o entity.PixelStorage, cd entity.DataConfig) error {

	SQL := fmt.Sprintf(NEWMO, o.CampaignDetailId, o.PxDate, o.URLServiceKey, o.CampaignId, o.Country, o.Partner, o.Operator, o.Aggregator, o.Service, o.ShortCode, o.Adnet, o.Keyword, o.Subkeyword, o.IsBillable, o.Plan, o.URL, o.URLType, o.Pixel, o.TrxId, o.Token, o.Msisdn, o.IsUsed, o.Browser, o.OS, o.IP, o.ISP, o.ReferralURL, o.PubId, o.UserAgent, o.TrafficSource, o.TrafficSourceData, o.UserRejected, o.UserDuplicated, o.Handset, o.HandsetCode, o.HandsetType, o.URLLanding, o.URLWarpLanding, o.URLService, o.URLTFCSmartlink, o.PO, o.Cost, cd.IsActive, cd.CounterMOCapping, cd.MOCapping, cd.StatusCapping, cd.CounterMORatio, cd.RatioSend, cd.RatioReceive, cd.StatusRatio, cd.APIURL, o.StatusPostback, o.IsUnique, o.URLPostback, o.StatusURLPostback, o.ReasonURLPostback)

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
} */
