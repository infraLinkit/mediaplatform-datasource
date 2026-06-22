package model

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/infraLinkit/mediaplatform-datasource/entity"
	"gorm.io/gorm"
)

func (r *BaseModel) NewPixel(o entity.PixelStorage) int {

	result := r.DB.Create(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return int(o.ID)
}

func (r *BaseModel) GetPx(o entity.PixelStorage) (entity.PixelStorage, bool) {

	result := r.DB.Model(&o).
		Where("url_service_key = ? AND date(pxdate) = CURRENT_DATE AND pixel = ? AND is_used = false AND is_unique = false", o.URLServiceKey, o.Pixel).
		First(&o)

	b := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if b {
		return o, false
	} else {
		r.Logs.Warn(fmt.Sprintf("pixel not found %#v", o))
		return o, true
	}
}

func (r *BaseModel) GetToken(o entity.PixelStorage) (entity.PixelStorage, bool) {

	result := r.DB.Model(&o).
		Where("url_service_key = ? AND DATE(pxdate) = CURRENT_DATE AND token = ? AND is_used = false AND is_unique = false", o.URLServiceKey, o.Pixel).
		First(&o)

	b := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if b {
		return o, false
	} else {
		r.Logs.Warn(fmt.Sprintf("pixel not found %#v", o))
		return o, true
	}
}

func (r *BaseModel) GetByAdnetCode(o entity.PixelStorage) (entity.PixelStorage, bool) {

	result := r.DB.Model(&o).
		Where("url_service_key = ? AND date(pxdate) = CURRENT_DATE AND pixel LIKE ? AND is_used = false AND is_unique = false", o.URLServiceKey, "%"+o.Pixel+"%").
		First(&o)

	b := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if b {
		return o, false
	} else {
		r.Logs.Warn(fmt.Sprintf("pixel found %#v", o))
		return o, true
	}
}

func (r *BaseModel) GetPxByMsisdn(o entity.PixelStorage) (entity.PixelStorage, bool) {

	result := r.DB.Model(&o).
		Where("url_service_key = ? AND DATE(pxdate) = CURRENT_DATE AND msisdn = ? AND is_used = false AND is_unique = false", o.URLServiceKey, o.Pixel).
		First(&o)

	b := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if b {
		return o, false
	} else {
		r.Logs.Warn(fmt.Sprintf("pixel not found %#v", o))
		return o, true
	}
}

func (r *BaseModel) UpdatePixelById(o entity.PixelStorage) error {

	result := r.DB.Model(&o).Select("is_used", "pixel_used_date", "status_postback", "status_postback", "is_unique", "url_postback", "status_url_postback", "reason_url_postback", "reason_url_postback").Updates(o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) SpecialGetPx(o entity.PixelStorage, si int, ei int) (entity.PixelStorage, bool) {

	start := o.Pxdate
	var (
		end          time.Time
		result_query []string
		rpl          *strings.Replacer
		tbl          string
	)

	for x := si; x >= ei; x-- {

		SQL := `SELECT * FROM {TBL} WHERE url_service_key = '{CAMPAIGNID}' AND date(pxdate) = '{DATE}' AND pixel = '{PIXEL}' AND is_used = false AND is_unique = false`

		end = start.AddDate(0, 0, x)

		if x == 0 {
			tbl = "pixel_storages"
		} else {
			tbl = "pixel_storages" + "_" + end.Format("20060102")
		}

		rpl = strings.NewReplacer(
			"{TBL}", tbl,
			"{DATE}", end.Format("2006-01-02"),
			"{CAMPAIGNID}", o.URLServiceKey,
			"{PIXEL}", o.Pixel,
		)
		SQL = rpl.Replace(SQL)

		result_query = append(result_query, SQL)
	}

	result := r.DB.Raw(strings.Join(result_query, ` union `) + " LIMIT 1").Scan(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	b := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if b {
		return o, false
	} else {
		r.Logs.Warn(fmt.Sprintf("pixel not found %#v", o))
		return o, true
	}
}

func pxdateToSQL(pxdate string) (tbl, dateSQL string) {
	if pxdate != "" {
		t, err := time.Parse("20060102", pxdate)
		if err == nil && t.Format("20060102") != time.Now().Format("20060102") {
			return "pixel_storages_" + t.Format("20060102"), "'" + t.Format("2006-01-02") + "'"
		}
	}
	return "pixel_storages", "CURRENT_DATE"
}

func (r *BaseModel) GetPxByDate(o entity.PixelStorage, pxdate string) (entity.PixelStorage, bool) {
	tbl, dateSQL := pxdateToSQL(pxdate)
	result := r.DB.Raw(
		fmt.Sprintf("SELECT * FROM %s WHERE url_service_key = ? AND date(pxdate) = %s AND pixel = ? AND is_used = false AND is_unique = false LIMIT 1", tbl, dateSQL),
		o.URLServiceKey, o.Pixel,
	).Scan(&o)
	if result.RowsAffected == 0 {
		return o, false
	}
	r.Logs.Warn(fmt.Sprintf("pixel found %#v", o))
	return o, true
}

func (r *BaseModel) GetTokenByDate(o entity.PixelStorage, pxdate string) (entity.PixelStorage, bool) {
	tbl, dateSQL := pxdateToSQL(pxdate)
	result := r.DB.Raw(
		fmt.Sprintf("SELECT * FROM %s WHERE url_service_key = ? AND date(pxdate) = %s AND token = ? AND is_used = false AND is_unique = false LIMIT 1", tbl, dateSQL),
		o.URLServiceKey, o.Pixel,
	).Scan(&o)
	if result.RowsAffected == 0 {
		return o, false
	}
	r.Logs.Warn(fmt.Sprintf("pixel found %#v", o))
	return o, true
}

func (r *BaseModel) GetByAdnetCodeByDate(o entity.PixelStorage, pxdate string) (entity.PixelStorage, bool) {
	tbl, dateSQL := pxdateToSQL(pxdate)
	result := r.DB.Raw(
		fmt.Sprintf("SELECT * FROM %s WHERE url_service_key = ? AND date(pxdate) = %s AND pixel LIKE ? AND is_used = false AND is_unique = false LIMIT 1", tbl, dateSQL),
		o.URLServiceKey, "%"+o.Pixel+"%",
	).Scan(&o)
	if result.RowsAffected == 0 {
		return o, false
	}
	r.Logs.Warn(fmt.Sprintf("pixel found %#v", o))
	return o, true
}

func (r *BaseModel) GetPxByMsisdnByDate(o entity.PixelStorage, pxdate string) (entity.PixelStorage, bool) {
	tbl, dateSQL := pxdateToSQL(pxdate)
	result := r.DB.Raw(
		fmt.Sprintf("SELECT * FROM %s WHERE url_service_key = ? AND date(pxdate) = %s AND msisdn = ? AND is_used = false AND is_unique = false LIMIT 1", tbl, dateSQL),
		o.URLServiceKey, o.Pixel,
	).Scan(&o)
	if result.RowsAffected == 0 {
		return o, false
	}
	r.Logs.Warn(fmt.Sprintf("pixel found %#v", o))
	return o, true
}

func (r *BaseModel) UpdatePixelBilled(o entity.PixelStorage, pxdate string) error {

	tbl := "pixel_storages"
	dateClause := "CURRENT_DATE"

	if pxdate != "" {
		t, err := time.Parse("20060102", pxdate)
		if err == nil && t.Format("20060102") != time.Now().Format("20060102") {
			tbl = "pixel_storages_" + t.Format("20060102")
			dateClause = "'" + t.Format("2006-01-02") + "'"
		}
	}

	result := r.DB.Exec(
		fmt.Sprintf(`UPDATE %s
		 SET updated_at = NOW(),
		     m_status_time_charge = NOW(),
		     m_status_charge = ?
		 WHERE DATE(pxdate) = %s
		   AND url_service_key = ?
		   AND pixel = ?
		   AND is_unique = false `, tbl, dateClause),
		o.MStatusCharge,
		o.URLServiceKey,
		o.Pixel,
	)

	if result.Error != nil {
		r.Logs.Error(fmt.Sprintf("DB error update pixel billed: %#v", result.Error))
		return result.Error
	}

	r.Logs.Debug(fmt.Sprintf("update pixel billed affected rows: %d", result.RowsAffected))

	return nil
}