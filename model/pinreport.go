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
)

func (r *BaseModel) PinReport(o entity.ApiPinReport) int {

	result := r.DB.Create(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return int(o.ID)
}

func (r *BaseModel) GetApiPinReport(o entity.DisplayPinReport) ([]entity.ApiPinReport, error) {

	var (
		rows *sql.Rows
	)

	query := r.DB.Model(&entity.ApiPinReport{})

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
			default:
				query = query.Where("date_send = ?", o.DateRange)
			}
		}

		rows, _ = query.Order("date_send").Rows()
	} else {
		rows, _ = query.Rows()
	}

	defer rows.Close()

	var (
		ss []entity.ApiPinReport
	)

	for rows.Next() {

		var s entity.ApiPinReport

		// ScanRows scans a row into a struct
		r.DB.ScanRows(rows, &s)

		ss = append(ss, s)
	}

	r.Logs.Debug(fmt.Sprintf("Total data : %d ...\n", len(ss)))

	return ss, rows.Err()
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

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.PixelStorage{})
	query = query.Where("is_used = ?", "true")
	if o.CampaignType == "mainstream" {
		query = query.Where("campaign_objective = ?", "MAINSTREAM").Where("status_postback = ? ", "true")
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
		if o.CampaignType == "mainstream" {
			if o.Agency != "" {
				//compare to adnet maybe will change in the future
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
		if o.DateRange != "" {
			switch strings.ToUpper(o.DateRange) {
			case "TODAY":
				query = query.Where("pxdate BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '1 day' - INTERVAL '1 second'")
			case "YESTERDAY":
				query = query.Where("pxdate BETWEEN CURRENT_DATE - INTERVAL '1 DAY' AND CURRENT_DATE")
				query = query.Where("pxdate BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
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
		} else {
			query = query.Where("pxdate BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '1 day' - INTERVAL '1 second'")
		}
	}

	// Get the total count after applying filters
	query.Unscoped().Count(&total_rows)

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
