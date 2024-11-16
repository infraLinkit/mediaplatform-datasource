package model

import (
	"context"
	"fmt"

	"github.com/infraLinkit/mediaplatform-datasource/entity"

	_ "github.com/lib/pq"
)

const (
	NEWPX = "INSERT INTO pixel_storage (id, campaign_detail_id, pxdate, urlservicekey, campaign_id, country, partner, operator, aggregator, service, short_code, adnet, keyword, subkeyword, is_billable, plan, url, url_type, pixel, trx_id, token, is_used, browser, os, ip, isp, referral_url, pubid, user_agent, traffic_source, traffic_source_data, user_rejected, user_duplicated, handset, handset_code, handset_type, url_landing, url_warp_landing, url_service, url_tfc_or_smartlink) VALUES (DEFAULT, %d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, '%s', '%s', '%s', '%s', '%s', '%s', %t, '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, '%s', %t, %t, '%s', '%s', '%s', '%s', '%s', '%s', '%s')"
)

func (r *BaseModel) NewPixel(o entity.PixelStorage) error {

	SQL := fmt.Sprintf(NEWPX, o.CampaignDetailId, o.PxDate, o.URLServiceKey, o.CampaignId, o.Country, o.Partner, o.Partner, o.Aggregator, o.Service, o.ShortCode, o.Adnet, o.Keyword, o.Subkeyword, o.IsBillable, o.Plan, o.URL, o.URLType, o.Pixel, o.TrxId, o.Token, o.IsUsed, o.Browser, o.OS, o.IP, o.ISP, o.ReferralURL, o.PubId, o.UserAgent, o.TrafficSource, o.TrafficSourceData, o.UserRejected, o.UserDuplicated, o.Handset, o.HandsetCode, o.HandsetType, o.URLLanding, o.URLWarpLanding, o.URLService, o.URLTFCSmartlink)

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
