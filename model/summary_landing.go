package model

import (
	"fmt"

	"github.com/infraLinkit/mediaplatform-datasource/entity"
)

func (r *BaseModel) GetURLServiceFromSummaryLanding() ([]entity.SummaryLanding, error) {

	var (
		ss []entity.SummaryLanding
	)

	query := r.DB.Model(&entity.SummaryLanding{}).Select("url_service_key", "country", "operator", "partner", "adnet", "service", "url_service", "response_url_service_time").Where("summary_date_hour <= NOW() - INTERVAL '1 hour'")
	if rows, err := query.Rows(); err == nil {

		defer rows.Close()

		for rows.Next() {
			var s entity.SummaryLanding
			// ScanRows scans a row into a struct
			r.DB.ScanRows(rows, &s)
			ss = append(ss, s)
		}

		return ss, nil

	} else {

		return ss, err
	}

}

func (r *BaseModel) UpdateResponseTimeURLService(o entity.SummaryLanding) error {

	result := r.DB.Model(&o).
		Where("summary_date_hour <= NOW() - INTERVAL '1 hour' AND url_service_key = ?", o.URLServiceKey).
		Updates(entity.SummaryLanding{ResponseUrlServiceTime: o.ResponseUrlServiceTime})

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}
