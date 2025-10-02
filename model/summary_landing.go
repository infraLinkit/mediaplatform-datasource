package model

import (
	"fmt"

	"github.com/infraLinkit/mediaplatform-datasource/entity"
)

func (r *BaseModel) GetURLServiceFromSummaryLanding(event_date string, with_limit int) ([]entity.SummaryLanding, error) {

	var (
		ss []entity.SummaryLanding
	)

	query := r.DB.Model(&entity.SummaryLanding{}).Select("DISTINCT url_service_key, country, operator, partner, adnet, service, url_service")

	switch event_date {
	case "1HOURAGO":
		query.Where("DATE(summary_date_hour) = CURRENT_DATE and date_part('hour', summary_date_hour) = date_part('hour', NOW() - INTERVAL '1 hour')")
	case "1DAYAGO":
		query.Where("DATE(summary_date_hour) = CURRENT_DATE - INTERVAL '2 day'")
	default:
		query.Where("DATE(summary_date_hour) = '" + event_date + "'")
	}

	if with_limit > 0 {
		query.Limit(with_limit)
	}

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

func (r *BaseModel) UpdateResponseTimeURLService(event_date string, o entity.SummaryLanding) error {

	var query string

	switch event_date {
	case "1HOURAGO":
		query = "DATE(summary_date_hour) = CURRENT_DATE and date_part('hour', summary_date_hour) = date_part('hour', NOW() - INTERVAL '1 hour')"
	case "1DAYAGO":
		query = "DATE(summary_date_hour) = CURRENT_DATE - INTERVAL '2 day'"
	default:
		query = "DATE(summary_date_hour) = '" + event_date + "'"
	}

	result := r.DB.Exec(fmt.Sprintf("UPDATE summary_landings SET response_url_service_time = '%.2f' WHERE url_service_key = '%s' AND "+query, o.ResponseUrlServiceTime, o.URLServiceKey))

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}
