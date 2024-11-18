package model

import (
	"context"
	"fmt"

	"github.com/infraLinkit/mediaplatform-datasource/entity"

	_ "github.com/lib/pq"
)

const (
	GETCAMPAIGNDETAIL   = "SELECT id, urlservicekey, is_active, counter_mo_capping, mo_capping, status_capping, counter_mo_ratio, ratio_send, ratio_receive, status_ratio, api_url, pubid, cost, po FROM campaign_detail WHERE id = %d;"
	COUNTERCAPPING      = "UPDATE campaign_detail SET counter_mo_capping = counter_mo_capping+1, last_update_capping = CASE WHEN counter_mo_capping >= mo_capping THEN '%s'::timestamp(0) END WHERE id = %d;"
	COUNTERRATIO        = "UPDATE campaign_detail SET counter_mo_ratio = counter_mo_ratio+1 WHERE id = %d;"
	UPDATESTATUSCOUNTER = "UPDATE campaign_detail SET counter_mo_capping = %d, status_capping = %t, counter_mo_ratio = %d, status_ratio = %t, last_update = '%s'::timestamp(0), last_update_capping = CASE WHEN counter_mo_capping >= mo_capping THEN '%s'::timestamp(0) WHERE id = %d"
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
