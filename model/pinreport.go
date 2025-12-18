package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/infraLinkit/mediaplatform-datasource/entity"
	"gorm.io/gorm/clause"
)

func (r *BaseModel) getLastPayout(country, operator, adnet, service, campaignID string, dateSend time.Time) (float64, float64) {

	var prev entity.ApiPinReport

	err := r.DB.
		Select("payout_adn, payout_af").
		Where(
			`country = ?
			 AND operator = ?
			 AND adnet = ?
			 AND service = ?
			 AND campaign_id = ?
			 AND date_send < ?`,
			country, operator, adnet, service, campaignID, dateSend,
		).
		Order("date_send DESC").
		Limit(1).
		First(&prev).Error

	if err != nil {
		return 0, 0
	}

	return prev.PayoutAdn, prev.PayoutAF
}

func (r *BaseModel) PinReport(o entity.ApiPinReport) int {

	if o.PayoutAdn == 0 && o.PayoutAF == 0 {
		prevAdn, prevAF := r.getLastPayout(
			o.Country,
			o.Operator,
			o.Adnet,
			o.Service,
			o.CampaignId,
			o.DateSend,
		)
		o.PayoutAdn = prevAdn
		o.PayoutAF = prevAF
	}

	entity.BuildPinReportCalculation(&o)

	result := r.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "date_send"},
			{Name: "country"},
			{Name: "adnet"},
			{Name: "operator"},
			{Name: "service"},
			{Name: "campaign_id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"payout_adn",
			"payout_af",
			"total_mo",
			"total_postback",
			"sbaf",
			"saaf",
			"price_per_mo",
			"waki_revenue",
			"updated_at",
		}),
	}).Create(&o)

	r.Logs.Debugf(
		"[PIN_REPORT] affected=%d error=%v | %s %s %s %s %s",
		result.RowsAffected,
		result.Error,
		o.DateSend.Format("2006-01-02"),
		o.Country,
		o.Operator,
		o.Adnet,
		o.Service,
	)

	return int(o.ID)
}

func (r *BaseModel) GetDisplayPinReport(o entity.DisplayPinReport) ([]entity.ApiPinReport, int64, error) {
	var totalRows int64
	var ss []entity.ApiPinReport

	query := r.DB.Model(&entity.ApiPinReport{}).Select(`
		api_pin_reports.*,
		(payout_af * total_postback) AS saaf,
		(payout_adn * total_postback) AS sbaf,
		(CASE WHEN total_mo > 0 THEN (payout_af * total_postback) / total_mo ELSE 0 END) AS price_per_mo,
		((payout_af * total_postback) - (payout_adn * total_postback)) AS waki_revenue
	`)

	if o.Action == "Search" {
		if o.CampaignId != "" {
			query = query.Where("campaign_id = ?", o.CampaignId)
		}
		if o.Country != "" {
			query = query.Where("country = ?", o.Country)
		}
		if o.Company != "" {
			query = query.Where("company = ?", o.Company)
		}
		if o.Operator != "" {
			query = query.Where("operator = ?", o.Operator)
		}
		if len(o.Adnets) > 0 {
			query = query.Where("adnet IN ?", o.Adnets)
		}
		if o.Service != "" {
			query = query.Where("service = ?", o.Service)
		}

		if o.DateRange != "" {
			switch strings.ToUpper(o.DateRange) {
			case "TODAY":
				query = query.Where("date_send = CURRENT_DATE")
			case "YESTERDAY":
				query = query.Where("date_send = CURRENT_DATE - INTERVAL '1 DAY'")
			case "LAST7DAY":
				query = query.Where("date_send BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
			case "LAST30DAY":
				query = query.Where("date_send BETWEEN CURRENT_DATE - INTERVAL '30 DAY' AND CURRENT_DATE")
			case "THISMONTH":
				query = query.Where("date_send >= DATE_TRUNC('month', CURRENT_DATE)")
			case "LASTMONTH":
				query = query.Where("date_send BETWEEN DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 DAY'")
			case "CUSTOMRANGE":
				query = query.Where("date_send BETWEEN ? AND ?", o.DateBefore, o.DateAfter)
			default:
				query = query.Where("date_send = ?", o.DateRange)
			}
		} else {
			query = query.Where("date_send = CURRENT_DATE")
		}
	} else {
		query = query.Where("date_send = CURRENT_DATE")
	}

	if err := query.Count(&totalRows).Error; err != nil {
		return []entity.ApiPinReport{}, 0, err
	}

	if o.OrderColumn != "" {
		dir := "ASC"
		if strings.ToUpper(o.OrderDir) == "DESC" {
			dir = "DESC"
		}
	
		switch o.OrderColumn {
		case "saaf":
			query = query.Order(fmt.Sprintf("(payout_af * total_postback) %s", dir))
		case "sbaf":
			query = query.Order(fmt.Sprintf("(payout_adn * total_postback) %s", dir))
		case "price_per_mo":
			query = query.Order(fmt.Sprintf("(CASE WHEN total_mo > 0 THEN (payout_af * total_postback) / total_mo ELSE 0 END) %s", dir))
		case "waki_revenue":
			query = query.Order(fmt.Sprintf("((payout_af * total_postback) - (payout_adn * total_postback)) %s", dir))
		default:
			query = query.Order(fmt.Sprintf("%s %s", o.OrderColumn, dir))
		}
	} else {
		query = query.Order("date_send DESC").Order("id DESC")
	}	

	if err := query.
		Limit(o.PageSize).
		Offset((o.Page - 1) * o.PageSize).
		Find(&ss).Error; err != nil {
		return []entity.ApiPinReport{}, 0, err
	}

	return ss, totalRows, nil
}

func (r *BaseModel) EditPayoutAPIReport(o entity.ApiPinReport) error {

	dateOnly := o.DateSend.Format("2006-01-02")

	var existing entity.ApiPinReport

	err := r.DB.
		Where(`
			date_send = ? AND
			country = ? AND
			operator = ? AND
			service = ? AND
			adnet = ?
		`,
			dateOnly,
			o.Country,
			o.Operator,
			o.Service,
			o.Adnet,
		).
		First(&existing).Error

	if err != nil {
		return err
	}

	if o.PayoutAF == 0 {
		o.PayoutAF = existing.PayoutAF
	}
	if o.PayoutAdn == 0 {
		o.PayoutAdn = existing.PayoutAdn
	}

	o.TotalPostback = existing.TotalPostback
	o.TotalMO = existing.TotalMO

	entity.BuildPinReportCalculation(&o)

	result := r.DB.
		Model(&entity.ApiPinReport{}).
		Where(`
			date_send = ? AND
			country = ? AND
			operator = ? AND
			service = ? AND
			adnet = ?
		`,
			dateOnly,
			o.Country,
			o.Operator,
			o.Service,
			o.Adnet,
		).
		Updates(map[string]interface{}{
			"payout_af":     o.PayoutAF,
			"payout_adn":    o.PayoutAdn,
			"sbaf":         o.SBAF,
			"saaf":         o.SAAF,
			"price_per_mo":  o.PricePerMO,
			"waki_revenue":  o.WakiRevenue,
		})

	r.Logs.Debug(fmt.Sprintf(
		"affected: %d, error: %#v",
		result.RowsAffected,
		result.Error,
	))

	return result.Error
}

func (r *BaseModel) getPrevCPA(country, operator, adnet, service string, dateSend time.Time) (float64, float64) {

	var prev entity.ApiPinPerformance

	err := r.DB.
		Select("cpa, cpa_waki").
		Where(
			"country = ? AND operator = ? AND adnet = ? AND service = ? AND date_send < ?",
			country, operator, adnet, service, dateSend,
		).
		Order("date_send DESC").
		Limit(1).
		First(&prev).Error

	if err != nil {
		return 0, 0
	}

	return prev.CPA, prev.CPAWaki
}

func (r *BaseModel) UpsertPinPerformance(o *entity.ApiPinPerformance) error {

	if o.CPA == 0 && o.CPAWaki == 0 {
		prevCPA, prevCPAWaki := r.getPrevCPA(
			o.Country,
			o.Operator,
			o.Adnet,
			o.Service,
			o.DateSend,
		)
		o.CPA = prevCPA
		o.CPAWaki = prevCPAWaki
	}

	result := r.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "date_send"},
			{Name: "country"},
			{Name: "operator"},
			{Name: "adnet"},
			{Name: "service"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"pin_request",
			"unique_pin_request",
			"pin_sent",
			"pin_failed",
			"verify_request",
			"verify_request_unique",
			"pin_ok",
			"pin_not_ok",
			"pin_ok_send_adnet",
			"charged_mo",
			"subs_cr",
			"adnet_cr",
			"cac",
			"paid_cac",
			"total_spending",
			"saaf",
			"sbaf",
			"estimated_arpu",
			"updated_at",
		}),
	}).Create(o)

	if result.Error != nil {
		return result.Error
	}

	r.Logs.Debugf(
		"[PIN_PERFORMANCE] affected=%d %s %s %s %s %s",
		result.RowsAffected,
		o.DateSend.Format("2006-01-02"),
		o.Country,
		o.Operator,
		o.Adnet,
		o.Service,
	)

	return nil
}

func (r *BaseModel) PinPerformanceReport(o entity.ApiPinPerformance) int {

	result := r.DB.Create(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return int(o.ID)
}

func (r *BaseModel) GetApiPinPerformanceReport(o entity.DisplayPinPerformanceReport) ([]entity.ApiPinPerformance, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.ApiPinPerformance{})
	if o.Action == "Search" {
		if o.Country != "" {
			query = query.Where("country = ?", o.Country)
		}
		if o.Operator != "" {
			query = query.Where("operator = ?", o.Operator)
		}
		if o.Service != "" {
			query = query.Where("service = ?", o.Service)
		}
		if o.DateRange != "" {
			switch strings.ToUpper(o.DateRange) {
			case "TODAY":
				query = query.Where("date_send = CURRENT_DATE")
			case "YESTERDAY":
				query = query.Where("date_send BETWEEN CURRENT_DATE - INTERVAL '1 DAY' AND CURRENT_DATE")
			case "LAST7DAY":
				query = query.Where("date_send BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
			case "LAST30DAY":
				query = query.Where("date_send BETWEEN CURRENT_DATE - INTERVAL '30 DAY' AND CURRENT_DATE")
			case "THISMONTH":
				query = query.Where("date_send >= DATE_TRUNC('month', CURRENT_DATE)")
			case "LASTMONTH":
				query = query.Where("date_send BETWEEN DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 DAY'")
			case "CUSTOMRANGE":
				query = query.Where("date_send BETWEEN ? AND ?", o.DateBefore, o.DateAfter)
			case "ALLDATERANGE":
			default:
				query = query.Where("date_send = ?", o.DateRange)
			}
		} else {
			query = query.Where("date_send = CURRENT_DATE")
		}
	}

	// Get the total count after applying filters
	query.Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Order("date_send").Rows()
	defer rows.Close()

	var ss []entity.ApiPinPerformance
	for rows.Next() {
		var s entity.ApiPinPerformance
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	r.Logs.Debug(fmt.Sprintf("Total data : %d ...\n", len(ss)))

	return ss, total_rows, rows.Err()
}

func (r *BaseModel) GetConversionLogReport(o entity.DisplayConversionLogReport) ([]entity.PixelStorage, int64, error) {
	var (
		rows       *sql.Rows
		total_rows int64
	)

	tableName := "pixel_storages"

	if strings.ToUpper(o.DateRange) == "YESTERDAY" {
		yesterday := time.Now().AddDate(0, 0, -1).Format("20060102")
		tableName = fmt.Sprintf("pixel_storages_%s", yesterday)
	} else if strings.ToUpper(o.DateRange) == "2DAYAGO" {
		twoDaysAgo := time.Now().AddDate(0, 0, -2).Format("20060102")
		tableName = fmt.Sprintf("pixel_storages_%s", twoDaysAgo)
	}

	query := r.DB.Table(tableName)
	query = query.Where("is_used = ?", "true")

	if o.CampaignType == "mainstream" {
		query = query.Where("campaign_objective LIKE ?", "%MAINSTREAM%")
	} else {
		query = query.Where("campaign_objective IN ?", []string{"CPA", "CPC", "CPI", "CPM"})
	}

	if o.Action == "Search" {
		if o.Country != "" {
			query = query.Where("country = ?", o.Country)
		}
		if o.Operator != "" {
			query = query.Where("operator = ?", o.Operator)
		}
		if o.Pixel != "" {
			query = query.Where("pixel = ?", o.Pixel)
		}
		if o.CampaignId != "" {
			query = query.Where("campaign_id = ?", o.CampaignId)
		}
		if o.CampaignName != "" {
			query = query.Where("campaign_name = ?", o.CampaignName)
		}
		if o.CampaignType == "mainstream" {
			if o.StatusPostback != "" {
				query = query.Where("status_postback = ?", o.StatusPostback)
			}
			if o.Agency != "" {
				query = query.Where("adnet = ?", o.Agency)
			}
		}
		if o.CampaignType == "s2s" {
			if o.StatusPostback != "" {
				query = query.Where("status_postback = ?", o.StatusPostback)
			}
			if o.Adnet != "" {
				query = query.Where("adnet = ?", o.Adnet)
			}
		}

		if o.DateRange != "" && strings.ToUpper(o.DateRange) != "YESTERDAY" && strings.ToUpper(o.DateRange) != "2DAYAGO" {
			switch strings.ToUpper(o.DateRange) {
			case "TODAY":
				query = query.Where("pxdate BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '1 day' - INTERVAL '1 second'")
			case "LAST30DAY":
				query = query.Where("pxdate BETWEEN CURRENT_DATE - INTERVAL '30 DAY' AND CURRENT_DATE")
			case "THISMONTH":
				query = query.Where("pxdate >= DATE_TRUNC('month', CURRENT_DATE)")
			case "LASTMONTH":
				query = query.Where("pxdate BETWEEN DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 DAY'")
			case "CUSTOMRANGE":
				dateEnd, _ := time.Parse("2006-01-02", o.DateEnd)
				query = query.Where("pxdate BETWEEN ? AND ?", o.DateStart, dateEnd.AddDate(0, 0, 1))
			case "LAST7DAY":
				query = query.Where("pxdate BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
			case "ALLDATERANGE":
			default:
				query = query.Where("pxdate = ?", o.DateRange)
			}
		}
	}

	query.Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	var startIndex int
	if o.Order == "asc" {
		query_limit = query_limit.Order("pxdate asc")
		startIndex = int(total_rows) - ((o.Page - 1) * o.PageSize)
	} else {
		query_limit = query_limit.Order("pxdate desc")
		startIndex = (o.Page - 1) * o.PageSize
	}

	rows, _ = query_limit.Rows()
	defer rows.Close()

	var ss []entity.PixelStorage
	for rows.Next() {
		var s entity.PixelStorage
		r.DB.ScanRows(rows, &s)

		if o.Order == "asc" {
			s.ID = startIndex
			startIndex--
		} else {
			s.ID = startIndex + 1
			startIndex++
		}
		ss = append(ss, s)
	}
	return ss, total_rows, rows.Err()
}

func (r *BaseModel) GetPerformanceReport(o entity.PerformaceReportParams) ([]entity.PerformanceReport, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.SummaryCampaign{})
	query = query.Where("mo_received > 0")
	query.Select(`country, company, client_type, campaign_name, partner, operator, service, adnet, SUM(mo_received) AS pixel_received, SUM(postback) as pixel_send, SUM(cr_postback) as cr_postback,
	SUM(cr_mo) as cr_mo, SUM(landing) as landing, MAX(ratio_send) as ratio_send, MAX(ratio_receive) as ratio_receive,SUM(po) as price_per_postback,SUM(cost_per_conversion) as cost_per_conversion,
	SUM(agency_fee) as agency_fee, SUM(postback*po) as spending_to_adnets, SUM(total_waki_agency_fee), SUM(total_waki_agency_fee + po*postback) as total_spending,sum(cpa) as e_cpa,
	SUM(total_fp) as total_fp,SUM(success_fp) as success_fp`)

	if o.Action == "Search" {
		if o.Country != "" {
			query = query.Where("country = ?", o.Country)
		}
		if o.Operator != "" {
			query = query.Where("operator = ?", o.Operator)
		}
		if o.Partner != "" {
			query = query.Where("partner = ?", o.Partner)
		}
		if o.Company != "" {
			query = query.Where("company= ?", o.Company)
		}
		if o.CampaignType != "" {
			query = query.Where("campaign_objective = ?", o.CampaignType)
		}
		if o.CampaignId != "" {
			query = query.Where("campaign_id = ?", o.CampaignId)
		}
		if o.CampaignName != "" {
			query = query.Where("campaign_name = ?", o.CampaignName)
		}
		if o.ClientType != "" {
			query = query.Where("client_type = ?", o.ClientType)
		}
		if o.Publisher != "" {
			query = query.Where("adnet = ?", o.Publisher)
		}
		if o.Service != "" {
			query = query.Where("service = ?", o.Service)
		}
	}
	now := time.Now()
	dateStart, errStart := time.Parse("2006-01-02", o.DateStart)
	dateEnd, errEnd := time.Parse("2006-01-02", o.DateEnd)
	if errStart != nil {
		dateStart = now
	}
	if errEnd != nil {
		dateEnd = now
	}

	query = query.Where("summary_date BETWEEN ? AND ?", dateStart, dateEnd)

	query.Group("country, company, client_type, campaign_name, partner, operator, service, adnet")

	// Get the total count after applying filters
	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Order("country").Rows()
	defer rows.Close()

	var ss []entity.PerformanceReport
	for rows.Next() {
		var s entity.PerformanceReport
		r.DB.ScanRows(rows, &s)
		r.GetARPUReport(&s, dateStart, dateEnd)
		ss = append(ss, s)
	}

	return ss, total_rows, rows.Err()
}

func (r *BaseModel) GetDistinctPerformanceReport(o entity.SummaryCampaign) ([]entity.SummaryCampaign, error) {

	var (
		rows *sql.Rows
	)

	query := r.DB.Model(&entity.SummaryCampaign{})

	rows, _ = query.Unscoped().Distinct("summary_date", "country", "operator", "service").Where("summary_date = CURRENT_DATE").Order("summary_date").Rows()

	defer rows.Close()

	var (
		ss []entity.SummaryCampaign
	)

	for rows.Next() {

		var s entity.SummaryCampaign

		// ScanRows scans a row into a struct
		r.DB.ScanRows(rows, &s)

		ss = append(ss, s)
	}

	return ss, rows.Err()
}

func (r *BaseModel) GetARPUReport(s *entity.PerformanceReport, dateStart time.Time, dateEnd time.Time) {

	var apiResponse entity.ARPUResponse
	var isempty bool
	key := fmt.Sprintf("%s_%s_%s_%s_%s", s.Country, s.Operator, s.Service, dateStart, dateEnd)

	if apiResponse, isempty = r.RGetArpuReport(key, "$"); isempty {
		fmt.Println("IS Empty: ", isempty)
		apiBase := os.Getenv("APIARPU")
		if apiBase == "" {
			fmt.Println("Missing APIARPU environment variable")
			return
		}

		// Build base URL
		baseURL, err := url.Parse(apiBase + "/api/v4/arpu/arpu90")
		if err != nil {
			fmt.Println("Failed to parse base URL:", err)
			return
		}

		// Manually encode all query params
		query := fmt.Sprintf(
			"from=%s&to=%s&country=%s&operator=%s&service=%s",
			url.QueryEscape(dateStart.Format("2006-01-02")),
			url.QueryEscape(dateEnd.Format("2006-01-02")),
			url.QueryEscape(s.Country),
			url.QueryEscape(s.Operator),
			url.QueryEscape(s.Service), // encodes spaces as %20
		)

		baseURL.RawQuery = query

		// Make the request
		req, err := http.NewRequest("GET", baseURL.String(), nil)
		if err != nil {
			fmt.Println("Failed to create request:", err)
			return
		}

		req.Header.Add("Authorization", "Basic bWlkZGxld2FyZTpsMW5rMXQzNjA=")
		req.Header.Add("Accept", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Failed to make request:", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Failed to read response:", err)
			return
		}

		err = json.Unmarshal(body, &apiResponse)
		if err != nil {
			fmt.Println("Failed to parse JSON:", err)
			return
		}

		s, _ := json.Marshal(apiResponse)

		r.SetData(key, "$", string(s))
		r.SetExpireData(key, 60)
	}

	if apiResponse.Data == nil || len(apiResponse.Data.Data) == 0 {
		return
	}
	for _, item := range apiResponse.Data.Data {

		if item.Adnet == s.Adnet {
			s.ARPUROI = s.ECPA / (item.Arpu90Net / 3)
			s.ARPU90 = item.Arpu90Net
			if s.TotalFP > 0 {
				s.BillrateFP = s.SuccessFP / s.TotalFP * 100
			}
		}
	}
}
