package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/infraLinkit/mediaplatform-datasource/entity"
	"gorm.io/gorm/clause"
)

// GetDisplayBudgetIO — BudgetIO is now country+month only; allowedCompanies param kept for signature compat.
func (r *BaseModel) GetDisplayBudgetIO(o entity.DisplayBudgetIO, allowedCompanies []string) ([]entity.BudgetIO, int64, error) {
	var rows *sql.Rows
	var err error
	var total_rows int64

	query := r.DB.Model(&entity.BudgetIO{})

	if o.Action == "Search" {
		if o.Country != "" {
			query = query.Where("country = ?", o.Country)
		}
		if o.DateRange != "" {
			query = query.Where("month = ?", o.DateRange)
		} else {
			query = query.Where("month = TO_CHAR(CURRENT_DATE, 'YYYY-MM')")
		}
	}

	if o.OrderColumn != "" {
		dir := "ASC"
		if strings.ToUpper(o.OrderDir) == "DESC" {
			dir = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", o.OrderColumn, dir))
	} else {
		query = query.Order("month DESC").Order("id DESC")
	}

	query.Unscoped().Count(&total_rows)
	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}
	rows, err = query_limit.Rows()
	if err != nil {
		return []entity.BudgetIO{}, 0, err
	}
	defer rows.Close()

	var ss []entity.BudgetIO
	for rows.Next() {
		var s entity.BudgetIO
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}
	r.Logs.Debug(fmt.Sprintf("Total data : %d ... \n", len(ss)))
	return ss, total_rows, rows.Err()
}

func (r *BaseModel) GetDisplayBudgetIOAll(o entity.DisplayBudgetIO) ([]entity.BudgetIO, int64, error) {
	var rows *sql.Rows
	var err error
	var total_rows int64

	query := r.DB.Model(&entity.BudgetIO{})

	if o.Action == "Search" {
		if o.Country != "" {
			query = query.Where("country = ?", o.Country)
		}
		if o.DateRange != "" {
			query = query.Where("month = ?", o.DateRange)
		} else {
			query = query.Where("month = TO_CHAR(CURRENT_DATE, 'YYYY-MM')")
		}
	}

	if o.OrderColumn != "" {
		dir := "ASC"
		if strings.ToUpper(o.OrderDir) == "DESC" {
			dir = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", o.OrderColumn, dir))
	} else {
		query = query.Order("month DESC").Order("id DESC")
	}

	query.Unscoped().Count(&total_rows)
	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}
	rows, err = query_limit.Rows()
	if err != nil {
		return []entity.BudgetIO{}, 0, err
	}
	defer rows.Close()

	var ss []entity.BudgetIO
	for rows.Next() {
		var s entity.BudgetIO
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}
	r.Logs.Debug(fmt.Sprintf("Total data : %d ... \n", len(ss)))
	return ss, total_rows, rows.Err()
}

func (r *BaseModel) GetDisplayBudgetIOApproved(o entity.DisplayBudgetIO, allowedCompanies []string) ([]entity.BudgetIO, int64, error) {
	return r.GetDisplayBudgetIO(o, allowedCompanies)
}

func (r *BaseModel) GetDisplayBudgetIOApprovedAll(o entity.DisplayBudgetIO) ([]entity.BudgetIO, int64, error) {
	return r.GetDisplayBudgetIOAll(o)
}

func (r *BaseModel) GetDisplaySummaryBudgetIO(o entity.DisplaySummaryBudgetIO) ([]entity.SummaryBudgetIOAgg, int64, error) {
	var whereClause []string
	var args []interface{}

	if o.Action == "Search" {
		if o.Country != "" {
			cs := []string{}
			for _, c := range strings.Split(o.Country, ",") {
				if c = strings.TrimSpace(c); c != "" {
					cs = append(cs, c)
				}
			}
			if len(cs) == 1 {
				whereClause = append(whereClause, "s.country = ?")
				args = append(args, cs[0])
			} else if len(cs) > 1 {
				ph := make([]string, len(cs))
				for i := range ph { ph[i] = "?" }
				whereClause = append(whereClause, fmt.Sprintf("s.country IN (%s)", strings.Join(ph, ",")))
				for _, c := range cs { args = append(args, c) }
			}
		}
		if o.Partner != "" {
			whereClause = append(whereClause, "s.partner = ?")
			args = append(args, o.Partner)
		}
		if o.Continent != "" {
			whereClause = append(whereClause, "s.continent = ?")
			args = append(args, o.Continent)
		}
		if o.CampaignType != "" {
			whereClause = append(whereClause, "s.campaign_type LIKE ?")
			args = append(args, "%"+strings.ToUpper(o.CampaignType)+"%")
		}
		if o.Channel != "" {
			whereClause = append(whereClause, "s.channel = ?")
			args = append(args, o.Channel)
		}
		if o.Company != "" {
			whereClause = append(whereClause, "s.company = ?")
			args = append(args, o.Company)
		}
		if o.Operator != "" {
			whereClause = append(whereClause, "s.operator = ?")
			args = append(args, o.Operator)
		}
		if o.Service != "" {
			whereClause = append(whereClause, "s.service = ?")
			args = append(args, o.Service)
		}
		if o.ClientType != "" {
			whereClause = append(whereClause, "s.client_type = ?")
			args = append(args, o.ClientType)
		}
		if o.DateRange != "" {
			whereClause = append(whereClause, "s.month = ?")
			args = append(args, o.DateRange)
		} else if o.DateBefore != "" && o.DateAfter != "" {
			whereClause = append(whereClause, "s.summary_date BETWEEN ?::date AND ?::date")
			args = append(args, o.DateBefore, o.DateAfter)
		} else {
			whereClause = append(whereClause, "s.month = TO_CHAR(CURRENT_DATE, 'YYYY-MM')")
		}
	} else {
		whereClause = append(whereClause, "s.month = TO_CHAR(CURRENT_DATE, 'YYYY-MM')")
	}

	where := ""
	if len(whereClause) > 0 {
		where = "WHERE " + strings.Join(whereClause, " AND ")
	}

	orderBy := "ORDER BY s.month DESC, s.continent ASC, s.country ASC, s.operator ASC, s.channel ASC"
	if o.OrderColumn != "" {
		dir := "ASC"
		if strings.ToUpper(o.OrderDir) == "DESC" {
			dir = "DESC"
		}
		orderBy = fmt.Sprintf("ORDER BY %s %s", o.OrderColumn, dir)
	}

	rawSQL := fmt.Sprintf(`
		SELECT
			s.country, s.continent, s.company, s.partner, s.operator, s.channel,
			s.campaign_type, s.month,
			MAX(s.summary_date)::text AS last_date,
			SUM(CASE WHEN EXTRACT(DAY FROM s.summary_date) BETWEEN 1  AND 7  THEN s.actual_cost ELSE 0 END) AS actual_week_1,
			SUM(CASE WHEN EXTRACT(DAY FROM s.summary_date) BETWEEN 8  AND 14 THEN s.actual_cost ELSE 0 END) AS actual_week_2,
			SUM(CASE WHEN EXTRACT(DAY FROM s.summary_date) BETWEEN 15 AND 21 THEN s.actual_cost ELSE 0 END) AS actual_week_3,
			SUM(CASE WHEN EXTRACT(DAY FROM s.summary_date) BETWEEN 22 AND 28 THEN s.actual_cost ELSE 0 END) AS actual_week_4,
			SUM(CASE WHEN EXTRACT(DAY FROM s.summary_date) >= 29             THEN s.actual_cost ELSE 0 END) AS actual_week_5,
			SUM(CASE WHEN EXTRACT(DAY FROM s.summary_date) BETWEEN 1  AND 7  THEN s.mo_count ELSE 0 END) AS mo_week1,
			SUM(CASE WHEN EXTRACT(DAY FROM s.summary_date) BETWEEN 8  AND 14 THEN s.mo_count ELSE 0 END) AS mo_week2,
			SUM(CASE WHEN EXTRACT(DAY FROM s.summary_date) BETWEEN 15 AND 21 THEN s.mo_count ELSE 0 END) AS mo_week3,
			SUM(CASE WHEN EXTRACT(DAY FROM s.summary_date) BETWEEN 22 AND 28 THEN s.mo_count ELSE 0 END) AS mo_week4,
			SUM(CASE WHEN EXTRACT(DAY FROM s.summary_date) >= 29             THEN s.mo_count ELSE 0 END) AS mo_week5,
			COALESCE(b.id, 0)         AS budget_io_id,
			COALESCE(b.io_target, 0)  AS io_target,
			COALESCE(b.mo_target, 0)  AS mo_target,
			COALESCE(b.target_cac, 0) AS target_cac,
			COALESCE(b.ltv, 0)        AS ltv,
			COALESCE(b.roas, 0)       AS roas,
			COALESCE(b.roi, 0)        AS roi
		FROM summary_budget_ios s
		LEFT JOIN budget_ios b ON
			b.country = s.country AND b.month = s.month
			AND b.deleted_at IS NULL
		%s
		GROUP BY
			s.country, s.continent, s.company, s.partner, s.operator, s.channel,
			s.campaign_type, s.month,
			b.id, b.io_target, b.mo_target, b.target_cac, b.ltv, b.roas, b.roi
		%s
	`, where, orderBy)

	var ss []entity.SummaryBudgetIOAgg
	if err := r.DB.Raw(rawSQL, args...).Scan(&ss).Error; err != nil {
		return nil, 0, err
	}
	return ss, int64(len(ss)), nil
}

func (r *BaseModel) UpdateSummaryBudgetIO(req entity.UpdateSummaryBudgetIORequest) error {
	updates := map[string]interface{}{"updated_at": time.Now()}
	if req.MOTarget != nil {
		updates["mo_target"] = *req.MOTarget
	}
	if req.IOTarget != nil {
		updates["io_target"] = *req.IOTarget
	}
	if req.TargetCAC != nil {
		updates["target_cac"] = *req.TargetCAC
	}
	if req.LTV != nil {
		updates["ltv"] = *req.LTV
	}
	if req.ROAS != nil {
		updates["roas"] = *req.ROAS
	}
	if req.ROI != nil {
		updates["roi"] = *req.ROI
	}

	if req.ID > 0 {
		result := r.DB.Model(&entity.BudgetIO{}).Where("id = ?", req.ID).Updates(updates)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected > 0 {
			return nil
		}
	}

	if req.Country == "" || req.Month == "" {
		return fmt.Errorf("UpdateSummaryBudgetIO: id or (country, month) required")
	}
	return r.DB.Model(&entity.BudgetIO{}).
		Where("country = ? AND month = ?", req.Country, req.Month).
		Updates(updates).Error
}

type BudgetIOKey struct {
	Country string
	Month   string
}

func (r *BaseModel) UpsertBudgetIOPlaceholders(keys []BudgetIOKey) error {
	for _, k := range keys {
		row := entity.BudgetIO{
			Country: k.Country,
			Month:   k.Month,
		}
		err := r.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "country"}, {Name: "month"}},
			DoNothing: true,
		}).Create(&row).Error
		if err != nil {
			return fmt.Errorf("UpsertBudgetIOPlaceholders: %w", err)
		}
	}
	return nil
}
