package model

import (
	"errors"
	"fmt"

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
		Where("url_service_key = ? AND DATE(pxdate) = CURRENT_DATE AND is_used = false AND is_unique = false", o.URLServiceKey).
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
