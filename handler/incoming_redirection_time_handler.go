package handler

import (
	"fmt"
	"math"
	"net/url"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource/config"
	"github.com/infraLinkit/mediaplatform-datasource/entity"
	"github.com/ledongthuc/goterators"
)

func (h *IncomingHandler) DisplayRedirectionTime(c *fiber.Ctx) error {
	dataIndicators := extractQueryArrayRedirection(c, "data-indicators[]")
	if len(dataIndicators) == 0 {
		dataIndicators = append(dataIndicators, "total_load_time", "response_time", "landing", "success_rate", "click_ios", "click_android", "click_operator", "click_non_operator")
	}

	params := entity.RedirectionTimeParams{
		DataType:             c.Query("data-type"),
		ReportType:           c.Query("report-type"),
		Country:              c.Query("country"),
		Operator:             c.Query("operator"),
		PartnerName:          c.Query("partner-name"),
		Service:              c.Query("service"),
		Adnet:                c.Query("adnet"),
		TypeData:             c.Query("type-data"),
		CampaignId:           c.Query("campaign_id"),
		CampaignName:         c.Query("campaign"),
		URLServiceKey:        c.Query("url_service_key"),
		DataIndicators:       dataIndicators,
		DataBasedOn:          c.Query("data-based-on"),
		DataBasedOnIndicator: c.Query("data-based-on-indicator"),
		DateRange:            c.Query("date-range"),
		DateStart:            c.Query("date-start"),
		DateEnd:              c.Query("date-end"),
		DateCustomRange:      c.Query("date-custom-range"),
		All:                  c.Query("all"),
	}

	r := h.GenerateRedirection(c, params)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) GenerateRedirection(c *fiber.Ctx, params entity.RedirectionTimeParams) entity.ReturnResponse {
	var (
		summaryCampaign []entity.SummaryLanding
		startDate       time.Time
		endDate         time.Time
		kpiStats        entity.RedirectionKPIStats
		err             error
	)

	var chartData interface{} // ✅ interface{} agar bisa assign berbagai tipe

	// Panggil GetRedirectionTimeHourly khusus untuk hourly
	if params.DataType == "hourly_report" {
		summaryCampaign, startDate, endDate, err = h.DS.GetRedirectionTimeHourly(params)
		chartData = generateHourlyChartRedirection(summaryCampaign, params, startDate)
	} else {
		summaryCampaign, startDate, endDate, kpiStats, err = h.DS.GetRedirectionTime(params)
	}

	if err != nil {
		return entity.ReturnResponse{
			HttpStatus: fiber.StatusNotFound,
			Rsp: entity.GlobalResponse{
				Code:    fiber.StatusNotFound,
				Message: "empty",
			},
		}
	}

	summary := formatSummaryDataValRedirection(summaryCampaign, params, startDate, endDate)
	sortedSummary := sortDataRedirection(summary, params.DataBasedOn, params.DataBasedOnIndicator)

	return entity.ReturnResponse{
		HttpStatus: fiber.StatusOK,
		Rsp: entity.GlobalResponseWithData{
			Code:    fiber.StatusOK,
			Message: config.OK_DESC,
			Data: map[string]interface{}{
				"summary":    sortedSummary,
				"chart_data": chartData,
				"kpi_stats":  kpiStats,
			},
		},
	}
}


func formatSummaryDataValRedirection(data []entity.SummaryLanding, params entity.RedirectionTimeParams, startDate, endDate time.Time) []map[string]interface{} {
	var formattedData []map[string]interface{}

	if params.DataType == "hourly_report" {
		// Filter data by URLServiceKey and date (detail daily)
		var filtered []entity.SummaryLanding
		for _, item := range data {
			if item.URLServiceKey == params.URLServiceKey && sameDay(item.SummaryDateHour, startDate) {
				filtered = append(filtered, item)
			}
		}

		if len(filtered) > 0 {
			generatedSummary := generateHourlySummaryRedirection(filtered, params, startDate)

			placeHolder := map[string]interface{}{
				"level":         "detail",
				"campaign_id":   filtered[0].URLServiceKey,
				"campaign_name": filtered[0].CampaignName,
				"country":       filtered[0].Country,
				"operator":      filtered[0].Operator,
				"service":       filtered[0].Service,
				"adnet":         filtered[0].Adnet,
				"url_campaign":  filtered[0].URLCampaign,
				"date":          startDate.Format("2006-01-02"),
			}
			completeSummary := mergeMapsRedirection(generatedSummary, placeHolder)
			formattedData = append(formattedData, completeSummary)
		}

	} else if params.DataType == "daily_report" {

		if params.All == "true" {
			generatedSummary := generateSummaryRedirection(data, params, startDate, endDate)
			placeHolder := map[string]interface{}{
				"all": "All URL",
			}
			completeSummary := mergeMapsRedirection(generatedSummary, placeHolder)
			formattedData = append(formattedData, completeSummary)
		} else {
			groupedAdnet := goterators.Group(data, func(campaign entity.SummaryLanding) string {
				return campaign.URLServiceKey
			})
			for _, campaignPerAdnet := range groupedAdnet {
				generatedSummary := generateSummaryRedirection(campaignPerAdnet, params, startDate, endDate)

				placeHolder := map[string]interface{}{
					"level":         "country",
					"campaign_id":   campaignPerAdnet[0].URLServiceKey,
					"campaign_name": campaignPerAdnet[0].CampaignName,
					"country":       campaignPerAdnet[0].Country,
					"operator":      campaignPerAdnet[0].Operator,
					"service":       campaignPerAdnet[0].Service,
					"adnet":         campaignPerAdnet[0].Adnet,
					"url_campaign":  campaignPerAdnet[0].URLCampaign,
					"date":          campaignPerAdnet[0].SummaryDateHour,
				}
				completeSummary := mergeMapsRedirection(generatedSummary, placeHolder)
				formattedData = append(formattedData, completeSummary)
			}
		}
	}

	return formattedData
}


func generateSummaryRedirection(data []entity.SummaryLanding, params entity.RedirectionTimeParams, startDate time.Time, endDate time.Time) map[string]interface{} {
	days := map[string]map[string]map[string]interface{}{}
	totals := make(map[string]float64)
	counts := make(map[string]int)

	for _, campaign := range data {
		date := campaign.SummaryDateHour.Format("2006-01-02")
		prevDate := campaign.SummaryDateHour.AddDate(0, 0, -1).Format("2006-01-02")

		if params.DataType == "monthly_report" {
			date = campaign.SummaryDateHour.Format("2006-01")
			prevDate = campaign.SummaryDateHour.AddDate(0, -1, 0).Format("2006-01")
		}

		if days[date] == nil {
			days[date] = make(map[string]map[string]interface{})
		}

		for _, indicator := range params.DataIndicators {
			indicatorValue := getIndicatorValueRedirection(campaign, indicator)

			if days[date][indicator] == nil {
				days[date][indicator] = map[string]interface{}{
					"value":      0.0,
					"percentage": 0.0,
				}
			}

			prevValue := getPreviousValueRedirection(days[prevDate], indicator)

			// Hitung Total & Count
			totals[indicator] += indicatorValue
			counts[indicator]++

			var newValue float64

			// KHUSUS untuk "All", hitung average untuk indikator tertentu
			if params.All == "true" && (indicator == "success_rate" || indicator == "response_time" || indicator == "total_load_time") {
				if counts[indicator] > 0 {
					newValue = totals[indicator] / float64(counts[indicator])
				} else {
					newValue = 0.0
				}
			} else {
				// Default: jumlahkan value
				currentValue, ok := days[date][indicator]["value"].(float64)
				if !ok {
					currentValue = 0.0
				}
				newValue = currentValue + indicatorValue
			}

			days[date][indicator]["value"] = newValue
			days[date][indicator]["percentage"] = countPercentageRedirection(newValue, prevValue)
		}
	}

	// Final calculations
	tmoEnd := countTmoEndRedirection(totals, startDate, endDate)

	// Prepare summary data
	summaryData := map[string]interface{}{
		"data_indicators": params.DataIndicators,
		"total":           totals,
		"avg":             countAverageRedirection(totals, startDate, endDate),
		"t_mo_end":        tmoEnd,
	}

	// Merge with daily breakdowns
	completeSummary := mergeDaysRedirection(summaryData, days)

	return completeSummary
}


func generateHourlySummaryRedirection(data []entity.SummaryLanding, params entity.RedirectionTimeParams, date time.Time) map[string]interface{} {
	hours := make(map[string]map[string]map[string]interface{})
	totals := make(map[string]float64)

	// Siapkan struktur untuk 24 jam (00–23)
	for i := 0; i < 24; i++ {
		hourStr := fmt.Sprintf("%02d", i)
		hours[hourStr] = make(map[string]map[string]interface{})
		for _, indicator := range params.DataIndicators {
			hours[hourStr][indicator] = map[string]interface{}{
				"value":      0.0,
				"percentage": 0.0,
			}
		}
	}

	for _, campaign := range data {
		// Pastikan hanya data dengan tanggal yang sama
		if !sameDay(campaign.SummaryDateHour, date) {
			continue
		}

		hour := campaign.SummaryDateHour.Format("15")
		prevHour := campaign.SummaryDateHour.Add(-1 * time.Hour).Format("15")

		for _, indicator := range params.DataIndicators {
			indicatorValue := getIndicatorValueRedirection(campaign, indicator)

			currentValue := hours[hour][indicator]["value"].(float64)
			newValue := currentValue + indicatorValue
			hours[hour][indicator]["value"] = newValue

			prevValue := 0.0
			if val, ok := hours[prevHour]; ok {
				if pv, exists := val[indicator]; exists {
					prevValue, _ = pv["value"].(float64)
				}
			}

			hours[hour][indicator]["percentage"] = countPercentageRedirection(newValue, prevValue)
			totals[indicator] += indicatorValue
		}
	}

	// Hapus jam-jam yang belum terjadi jika date adalah hari ini
	if sameDay(date, time.Now()) {
		currentHour := time.Now().Hour()
		for i := currentHour + 1; i < 24; i++ {
			hourStr := fmt.Sprintf("%02d", i)
			delete(hours, hourStr)
		}
	}

	summaryData := map[string]interface{}{
		"data_indicators": params.DataIndicators,
		"total":           totals,
		"avg":             countAverageRedirectionHourly(totals, hours),
	}

	completeSummary := mergeHoursRedirection(summaryData, hours)

	return completeSummary
}

func mergeDaysRedirection(summaryData map[string]interface{}, days map[string]map[string]map[string]interface{}) map[string]interface{} {
	for key, value := range days {
		summaryData[key] = value
	}
	return summaryData
}

func mergeHoursRedirection(summaryData map[string]interface{}, hours map[string]map[string]map[string]interface{}) map[string]interface{} {
	for hour, indicators := range hours {
		summaryData[hour] = indicators
	}
	return summaryData
}

func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}



func groupOperatorRedirection(campaings []entity.SummaryLanding, params entity.RedirectionTimeParams, startDate time.Time, endDate time.Time) []interface{} {
	var formattedData []any
	groupedOperator := goterators.Group(campaings, func(campaign entity.SummaryLanding) string {
		return campaign.Operator
	})

	for _, campaignPerOperator := range groupedOperator {
		var children []any

		generatedSummary := generateSummaryRedirection(campaignPerOperator, params, startDate, endDate)

		children = groupPartnerRedirection(campaignPerOperator, params, startDate, endDate)

		placeHolder := map[string]any{
			"level":     "operator",
			"country":   campaignPerOperator[0].Operator,
			"_children": children,
		}
		completeSummary := mergeMapsRedirection(generatedSummary, placeHolder)
		formattedData = append(formattedData, completeSummary)
	}
	return formattedData
}

func groupPartnerRedirection(campaings []entity.SummaryLanding, params entity.RedirectionTimeParams, startDate time.Time, endDate time.Time) []interface{} {
	var formattedData []any

	groupedPartner := goterators.Group(campaings, func(campaign entity.SummaryLanding) string {
		return campaign.Partner
	})

	for _, campaignPerPatner := range groupedPartner {
		var children []any
		generatedSummary := generateSummaryRedirection(campaignPerPatner, params, startDate, endDate)
		children = groupServiceRedirection(campaignPerPatner, params, startDate, endDate)

		placeHolder := map[string]any{
			"level":     "partner",
			"country":   campaignPerPatner[0].Partner,
			"_children": children,
		}
		completeSummary := mergeMapsRedirection(generatedSummary, placeHolder)
		formattedData = append(formattedData, completeSummary)
	}
	return formattedData
}

func groupServiceRedirection(campaings []entity.SummaryLanding, params entity.RedirectionTimeParams, startDate time.Time, endDate time.Time) []interface{} {
	var formattedData []any
	groupedService := goterators.Group(campaings, func(campaign entity.SummaryLanding) string {
		return campaign.Service
	})

	for _, campaignPerService := range groupedService {
		var children []any

		generatedSummary := generateSummaryRedirection(campaignPerService, params, startDate, endDate)
		children = groupAdnetRedirection(campaignPerService, params, startDate, endDate)

		placeHolder := map[string]any{
			"level":     "service",
			"country":   campaignPerService[0].Service,
			"_children": children,
		}
		completeSummary := mergeMapsRedirection(generatedSummary, placeHolder)
		formattedData = append(formattedData, completeSummary)
	}
	return formattedData
}

func groupAdnetRedirection(campaings []entity.SummaryLanding, params entity.RedirectionTimeParams, startDate time.Time, endDate time.Time) []interface{} {
	var formattedData []any
	groupedAdnet := goterators.Group(campaings, func(campaign entity.SummaryLanding) string {
		return campaign.Adnet
	})

	for _, campaignPerAdnet := range groupedAdnet {
		generatedSummary := generateSummaryRedirection(campaignPerAdnet, params, startDate, endDate)
		placeHolder := map[string]any{
			"level":   "adnet",
			"country": campaignPerAdnet[0].Adnet,
		}
		completeSummary := mergeMapsRedirection(generatedSummary, placeHolder)
		formattedData = append(formattedData, completeSummary)
	}
	return formattedData
}


func getIndicatorValueRedirection(item entity.SummaryLanding, key string) float64 {

	key = SnakeToCamelRedirection(key)
	values := reflect.ValueOf(item)
	keyValues := values.FieldByName(key)

	if !keyValues.IsValid() {
		return 0
	}

	switch keyValues.Kind() {
	case reflect.Float64:
		return keyValues.Float()
	case reflect.Int, reflect.Int64:
		return float64(keyValues.Int())
	default:
		return 0
	}

}

func getPreviousValueRedirection(data map[string]map[string]interface{}, key string) float64 {
	value := 0.0
	if val, exists := data[key]; exists {
		value = val["value"].(float64)
	}
	return value
}

func countPercentageRedirection(now, prev float64) float64 {
	if prev == 0 {
		if now == 0 {
			return 0
		}
		return 100
	}
	if now == 0 {
		return -100
	}
	percentage := ((now - prev) / prev) * 100
	return percentage
}

func countTmoEndRedirection(totals map[string]float64, startDate time.Time, endDate time.Time) map[string]float64 {
	tmoEnd := map[string]float64{}
	totalDaysRunning := int(math.Ceil(endDate.Sub(startDate).Hours() / 24))
	if totalDaysRunning < 1 {
		totalDaysRunning = 1
	}

	// Calculate total days in the last month
	lastMonthEnd := endDate.AddDate(0, 0, -endDate.Day())
	lastMonthStart := lastMonthEnd.AddDate(0, 0, -lastMonthEnd.Day())
	totalDaysLastMonth := int(math.Ceil(lastMonthEnd.Sub(lastMonthStart).Hours()/24)) + 1
	for key, value := range totals {
		result := (value / float64(totalDaysRunning)) * float64(totalDaysLastMonth)
		tmoEnd[key] = result
	}

	return tmoEnd
}

func countAverageRedirection(totals map[string]float64, startDate, endDate time.Time) map[string]float64 {
	averages := map[string]float64{}
	totalDaysRunning := int(endDate.Sub(startDate).Hours() / 24)
	if totalDaysRunning < 1 {
		totalDaysRunning = 1
	}

	for key, value := range totals {
		avg := value / float64(totalDaysRunning)
		if key == "success_rate" && avg > 100 {
			avg = 100
		}
		averages[key] = avg
	}

	return averages
}

func countAverageRedirectionHourly(totals map[string]float64, hours map[string]map[string]map[string]interface{}) map[string]float64 {
	averages := map[string]float64{}

	numHours := len(hours)
	if numHours == 0 {
		numHours = 1
	}

	for indicator, totalValue := range totals {
		avg := totalValue / float64(numHours)
		if indicator == "success_rate" && avg > 100 {
			avg = 100
		}
		averages[indicator] = avg
	}

	return averages
}



func extractQueryArrayRedirection(c *fiber.Ctx, key string) []string {
	rawQuery, _ := url.QueryUnescape(string(c.Request().URI().QueryString()))
	values := []string{}

	// Adjust logic if the input `rawQuery` needs different parsing
	params := strings.Split(rawQuery, "&")

	for _, param := range params {
		parts := strings.SplitN(param, "=", 2)
		if len(parts) == 2 && parts[0] == key {
			values = append(values, parts[1])
		}
	}
	return values
}

func SnakeToCamelRedirection(snake string) string {
	if snake == "click_ios" {
		return "ClickIOS"
	}

	words := strings.Split(snake, "_")
	for i, word := range words {
		words[i] = strings.ToLower(word)
		if len(word) > 0 {
			words[i] = strings.ToUpper(words[i][:1]) + words[i][1:]
		}
	}
	return strings.Join(words, "")
}

func mergeMapsRedirection(map1, map2 map[string]interface{}) map[string]interface{} {
	mergedMap := make(map[string]interface{})
	for key, value := range map1 {
		mergedMap[key] = value
	}
	for key, value := range map2 {
		mergedMap[key] = value
	}
	return mergedMap
}

func sortDataRedirection(data []map[string]interface{}, dataBasedOn string, dataBasedOnIndicator string) []map[string]interface{} {
	// Make a copy of the original slice to avoid modifying it directly
	sortedData := make([]map[string]interface{}, len(data))
	copy(sortedData, data)

	// Define the sorting function based on the parameters
	sort.Slice(sortedData, func(i, j int) bool {
		// Get the "total" map from each item
		totalI, okI := sortedData[i]["total"].(map[string]float64)
		totalJ, okJ := sortedData[j]["total"].(map[string]float64)

		// If either total is missing or doesn't have the indicator, maintain order
		if !okI || !okJ {
			return i < j
		}

		// Get the indicator values
		valI := totalI[dataBasedOnIndicator]
		valJ := totalJ[dataBasedOnIndicator]

		// Sort based on the specified direction
		if dataBasedOn == "highest" {
			return valI > valJ // Sort descending
		} else {
			return valI < valJ // Sort ascending
		}
	})

	return sortedData
}

func generateHourlyChartRedirection(data []entity.SummaryLanding, params entity.RedirectionTimeParams, date time.Time) map[string][]map[string]interface{} {
	chart := make(map[string][]map[string]interface{})

	// Inisialisasi untuk semua indikator
	for _, indicator := range params.DataIndicators {
		chart[indicator] = []map[string]interface{}{}
	}

	// Siapkan nilai awal per jam
	hourlyMap := make(map[string]map[string]float64) // hour => indicator => value
	for i := 0; i < 24; i++ {
		hourStr := fmt.Sprintf("%02d", i)
		hourlyMap[hourStr] = make(map[string]float64)
		for _, indicator := range params.DataIndicators {
			hourlyMap[hourStr][indicator] = 0.0
		}
	}

	// Kelompokkan nilai berdasarkan jam dan indikator
	for _, item := range data {
		if !sameDay(item.SummaryDateHour, date) {
			continue
		}

		hour := item.SummaryDateHour.Format("15")
		for _, indicator := range params.DataIndicators {
			val := getIndicatorValueRedirection(item, indicator)
			hourlyMap[hour][indicator] += val
		}
	}

	// Hitung percentage dan bentuk data chart
	for i := 0; i < 24; i++ {
		hourStr := fmt.Sprintf("%02d", i)
		prevHourStr := fmt.Sprintf("%02d", i-1)

		for _, indicator := range params.DataIndicators {
			current := hourlyMap[hourStr][indicator]
			previous := 0.0
			if i > 0 {
				previous = hourlyMap[prevHourStr][indicator]
			}
			percentage := countPercentageRedirection(current, previous)

			chart[indicator] = append(chart[indicator], map[string]interface{}{
				"time":       hourStr + ":00",
				"value":      current,
				"percentage": percentage,
			})
		}
	}

	// Hapus jam-jam yang belum terjadi jika tanggal adalah hari ini
	if sameDay(date, time.Now()) {
		currentHour := time.Now().Hour()
		for _, indicator := range params.DataIndicators {
			chart[indicator] = chart[indicator][:currentHour+1]
		}
	}

	return chart
}
