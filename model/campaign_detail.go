package model

import (
	"context"
	"fmt"

	"github.com/infraLinkit/mediaplatform-datasource/entity"

	_ "github.com/lib/pq"
)

const (
	GETCAMPAIGNDETAIL   = "SELECT * FROM campaign_detail WHERE id = %d;"
	COUNTERCAPPING      = "UPDATE campaign_detail SET counter_mo_capping = counter_mo_capping+1, last_update_capping CASE WHEN counter_mo_capping >= mo_capping THEN '%s' END WHERE id = %d;"
	COUNTERRATIO        = "UPDATE campaign_detail SET counter_mo_ratio = counter_mo_ratio+1 WHERE id = %d;"
	UPDATESTATUSCOUNTER = "UPDATE campaign_detail SET status_capping = %t, status_ratio = %t, last_update = '%s' WHERE id = %d"
)

func (r *BaseModel) GetCampaignDetail(o entity.DataConfig) (entity.DataConfig, error) {

	SQL := fmt.Sprintf(GETCAMPAIGNDETAIL, o.Id)
	rows, err := r.DBPostgre.Query(SQL)
	if err != nil {
		r.Logs.Error(fmt.Sprintf("GetCampaignDetail, SQL : %s, error querying occured : %#v", SQL, err))

		return entity.DataConfig{}, err
	}
	defer rows.Close()

	for rows.Next() {

		err = rows.Scan(&o.Id, &o.URLServiceKey, &o.CampaignId, &o.Country, &o.Operator, &o.Partner, &o.Aggregator, &o.Adnet, &o.Service, &o.Keyword, &o.SubKeyword, &o.IsBillable, &o.Plan, &o.PO, &o.Cost, &o.PubId, &o.ShortCode, &o.DeviceType, &o.OS, &o.URLType, &o.ClickType, &o.ClickDelay, &o.ClientType, &o.TrafficSource, &o.UniqueClick, &o.URLBanner, &o.URLLanding, &o.URLWarpLanding, &o.URLService, &o.URLTFCSmartlink, &o.GlobPost, &o.URLGlobPost, &o.CustomIntegration, &o.IPAddress, &o.IsActive, &o.MOCapping, &o.CounterMOCapping, &o.StatusCapping, &o.KPIUpperLimitCapping, &o.IsMachineLearningCapping, &o.RatioSend, &o.RatioReceive, &o.CounterMORatio, &o.StatusRatio, &o.KPIUpperLimitRatioSend, &o.KPIUpperLimitRatioReceive, &o.IsMachineLearningRatio, &o.APIURL, &o.LastUpdate, &o.LastUpdateCapping)

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

	SQL := fmt.Sprintf(UPDATESTATUSCOUNTER, o.StatusCapping, o.StatusRatio, o.LastUpdate, o.Id)

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
