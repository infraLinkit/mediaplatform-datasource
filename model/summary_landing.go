package model

import (
	"fmt"

	"github.com/infraLinkit/mediaplatform-datasource/entity"
)

func (r *BaseModel) GetURLServiceFromSummaryLanding(event_date string, with_limit int) ([]entity.SummaryLanding, error) {

	var (
		ss []entity.SummaryLanding
	)

	query := r.DB.Model(&entity.SummaryLanding{}).Select("url_service_key", "country", "operator", "partner", "adnet", "service", "url_service", "response_url_service_time")

	switch event_date {
	case "1 HOUR AGO":
		query.Where("summary_date_hour <= NOW() - INTERVAL '1 hour'")
	case "1 DAY AGO":
		query.Where("summary_date_hour <= NOW() - INTERVAL '1 days'")
	default:
		query.Where("DATE(summary_date_hour) = '" + event_date + "'")
	}

	if with_limit > 0 {
		query.Limit(with_limit)
	}

	query.Order("landing DESC")

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
