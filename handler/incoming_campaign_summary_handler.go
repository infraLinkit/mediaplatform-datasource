package handler

import (
	"fmt"
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

func (h *IncomingHandler) DisplayCampaignSummary(c *fiber.Ctx) error {
	dataIndicators := extractQueryArrayRevenue(c, "data-indicators[]")
	if len(dataIndicators) == 0 {
		switch dataType := c.Query("data-type"); dataType {
		case "spending":
			dataIndicators = append(dataIndicators, "spending")
		default:
			dataIndicators = append(dataIndicators, "traffic")
		}

	}

	params := entity.ParamsCampaignSummary{
		DataType:             c.Query("data-type"),
		ReportType:           c.Query("report-type"),
		Country:              c.Query("country"),
		Operator:             c.Query("operator"),
		PartnerName:          c.Query("partner-name"),
		Adnet:                c.Query("adnet"),
		Service:              c.Query("service"),
		DataIndicators:       dataIndicators,
		DataBasedOn:          c.Query("data-based-on"),
		DataBasedOnIndicator: c.Query("data-based-on-indicator"),
		DateRange:            c.Query("date-range"),
		DateStart:            c.Query("date-start"),
		DateEnd:              c.Query("date-end"),
		DateCustomRange:      c.Query("date-custom-range"),
		All:                  c.Query("all"),
	}

	r := h.GenerateCampaignSummary(c, params)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayCampaignSummaryChart(c *fiber.Ctx) error {
	dataIndicators := extractQueryArrayRevenue(c, "data-indicators[]")
	if len(dataIndicators) == 0 {
		switch dataType := c.Query("data-type"); dataType {
		case "spending":
			dataIndicators = append(dataIndicators, "spending")
		default:
			dataIndicators = append(dataIndicators, "traffic")
		}

	}

	params := entity.ParamsCampaignSummary{
		DataType:             c.Query("data-type"),
		ChartType:            c.Query("chart-type"),
		ReportType:           c.Query("report-type"),
		Country:              c.Query("country"),
		Operator:             c.Query("operator"),
		PartnerName:          c.Query("partner-name"),
		CampaignName:         c.Query("campaign-name"),
		Adnet:                c.Query("adnet"),
		Service:              c.Query("service"),
		DataIndicators:       dataIndicators,
		DataBasedOn:          c.Query("data-based-on"),
		DataBasedOnIndicator: c.Query("data-based-on-indicator"),
		DateRange:            c.Query("date-range"),
		DateStart:            c.Query("date-start"),
		DateEnd:              c.Query("date-end"),
		DateCustomRange:      c.Query("date-custom-range"),
		All:                  c.Query("all"),
	}

	summaryChart, _, _, err := h.DS.GetSummaryCampaignChart(params)

	var response entity.ReturnResponse
	if err == nil {
		response = entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithData{
				Code:    fiber.StatusOK,
				Message: config.OK_DESC,
				Data:    summaryChart,
			},
		}
	} else {
		response = entity.ReturnResponse{
			HttpStatus: fiber.StatusNotFound,
			Rsp: entity.GlobalResponse{
				Code:    fiber.StatusNotFound,
				Message: "empty",
			},
		}
	}

	r := response
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) GenerateCampaignSummary(c *fiber.Ctx, params entity.ParamsCampaignSummary) entity.ReturnResponse {

	summaryCampaign, startDate, endDate, err := h.DS.GetSummaryCampaignMonitoring(params)
	summary := formatSummaryDataValue(summaryCampaign, params, startDate, endDate)
	sortedSummary := sortDataRevenue(summary, params.DataBasedOn, params.DataBasedOnIndicator)

	if err == nil {
		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithData{
				Code:    fiber.StatusOK,
				Message: config.OK_DESC,
				Data:    sortedSummary,
			},
		}
	} else {
		return entity.ReturnResponse{
			HttpStatus: fiber.StatusNotFound,
			Rsp: entity.GlobalResponse{
				Code:    fiber.StatusNotFound,
				Message: "empty",
			},
		}
	}
}

func formatSummaryDataValue(data []entity.CampaignSummaryMonitoring, params entity.ParamsCampaignSummary, startDate time.Time, endDate time.Time) []map[string]interface{} {
	var formattedData []map[string]interface{}

	if params.DataType == "cr" || params.DataType == "spending" {
		if params.All == "true" {
			generatedSummary := generateSummaryValue(data, params, startDate, endDate, "all")
			placeHolder := map[string]any{
				"all": "All Campaign",
			}
			completeSummary := mergeMapsRevenue(generatedSummary, placeHolder)
			formattedData = append(formattedData, completeSummary)
		} else {
			groupedAdnet := goterators.Group(data, func(campaign entity.CampaignSummaryMonitoring) string {
				return campaign.Country + "|" + campaign.Operator + "|" + campaign.Service + "|" + campaign.Adnet
			})
			for _, campaignPerAdnet := range groupedAdnet {
				generatedSummary := generateSummaryValue(campaignPerAdnet, params, startDate, endDate, "adnet")

				placeHolder := map[string]interface{}{
					"level":         "country",
					"campaign_id":   campaignPerAdnet[0].CampaignId,
					"campaign_name": campaignPerAdnet[0].CampaignName,
					"country":       campaignPerAdnet[0].Country,
					"operator":      campaignPerAdnet[0].Operator,
					"service":       campaignPerAdnet[0].Service,
					"adnet":         campaignPerAdnet[0].Adnet,
					"date":          campaignPerAdnet[0].SummaryDate,
				}
				completeSummary := mergeMapsRevenue(generatedSummary, placeHolder)
				formattedData = append(formattedData, completeSummary)
			}
		}
	} else {
		if params.All == "true" {
			generatedSummary := generateSummaryValue(data, params, startDate, endDate, "all")
			placeHolder := map[string]interface{}{
				"level":   "country",
				"country": "All",
			}
			completeSummary := mergeMapsRevenue(generatedSummary, placeHolder)
			formattedData = append(formattedData, completeSummary)
		} else {

			groupedCountry := goterators.Group(data, func(campaign entity.CampaignSummaryMonitoring) string {
				return campaign.Country
			})

			for _, campaignPerCountry := range groupedCountry {
				generatedCountrySummary := generateSummaryValue(campaignPerCountry, params, startDate, endDate, "country")

				var children []any

				switch params.ReportType {
				case "campaign_summary":
					children = groupPartnerValue(campaignPerCountry, params, startDate, endDate)
				case "url_service_summary":
					children = groupServiceValue(campaignPerCountry, params, startDate, endDate)
				case "adnet_summary":
					children = groupAdnetValue(campaignPerCountry, params, startDate, endDate)
				default:
					children = groupOperatorValue(campaignPerCountry, params, startDate, endDate)
				}

				placeHolder := map[string]interface{}{
					"level":     "country",
					"country":   campaignPerCountry[0].Country,
					"_children": children,
				}
				completeSummary := mergeMapsRevenue(generatedCountrySummary, placeHolder)
				formattedData = append(formattedData, completeSummary)
			}

		}
	}
	return formattedData
}

func groupOperatorValue(campaings []entity.CampaignSummaryMonitoring, params entity.ParamsCampaignSummary, startDate time.Time, endDate time.Time) []interface{} {
	var formattedData []any
	groupedOperator := goterators.Group(campaings, func(campaign entity.CampaignSummaryMonitoring) string {
		return campaign.Operator
	})

	for _, campaignPerOperator := range groupedOperator {
		var children []any

		generatedSummary := generateSummaryValue(campaignPerOperator, params, startDate, endDate, "operator")

		children = groupPartnerValue(campaignPerOperator, params, startDate, endDate)

		placeHolder := map[string]any{
			"level":     "operator",
			"country":   campaignPerOperator[0].Operator,
			"_children": children,
		}
		completeSummary := mergeMapsRevenue(generatedSummary, placeHolder)
		formattedData = append(formattedData, completeSummary)
	}
	return formattedData
}

func groupPartnerValue(campaings []entity.CampaignSummaryMonitoring, params entity.ParamsCampaignSummary, startDate time.Time, endDate time.Time) []interface{} {
	var formattedData []any

	groupedPartner := goterators.Group(campaings, func(campaign entity.CampaignSummaryMonitoring) string {
		return campaign.Partner
	})

	for _, campaignPerPatner := range groupedPartner {
		var children []any
		generatedSummary := generateSummaryValue(campaignPerPatner, params, startDate, endDate, "parnter")
		children = groupServiceValue(campaignPerPatner, params, startDate, endDate)

		placeHolder := map[string]any{
			"level":     "partner",
			"country":   campaignPerPatner[0].Partner,
			"_children": children,
		}
		completeSummary := mergeMapsRevenue(generatedSummary, placeHolder)
		formattedData = append(formattedData, completeSummary)
	}
	return formattedData
}

func groupServiceValue(campaings []entity.CampaignSummaryMonitoring, params entity.ParamsCampaignSummary, startDate time.Time, endDate time.Time) []interface{} {
	var formattedData []any
	groupedService := goterators.Group(campaings, func(campaign entity.CampaignSummaryMonitoring) string {
		return campaign.Service
	})

	for _, campaignPerService := range groupedService {
		var children []any

		generatedSummary := generateSummaryValue(campaignPerService, params, startDate, endDate, "service")
		children = groupAdnetValue(campaignPerService, params, startDate, endDate)

		placeHolder := map[string]any{
			"level":     "service",
			"country":   campaignPerService[0].Service,
			"_children": children,
		}
		completeSummary := mergeMapsRevenue(generatedSummary, placeHolder)
		formattedData = append(formattedData, completeSummary)
	}
	return formattedData
}

func groupAdnetValue(campaings []entity.CampaignSummaryMonitoring, params entity.ParamsCampaignSummary, startDate time.Time, endDate time.Time) []interface{} {
	var formattedData []any
	groupedAdnet := goterators.Group(campaings, func(campaign entity.CampaignSummaryMonitoring) string {
		return campaign.Adnet
	})

	for _, campaignPerAdnet := range groupedAdnet {
		generatedSummary := generateSummaryValue(campaignPerAdnet, params, startDate, endDate, "adnet")
		placeHolder := map[string]any{
			"level":   "adnet",
			"country": campaignPerAdnet[0].Adnet,
		}
		completeSummary := mergeMapsRevenue(generatedSummary, placeHolder)
		formattedData = append(formattedData, completeSummary)
	}
	return formattedData
}

func generateSummaryValue(data []entity.CampaignSummaryMonitoring, params entity.ParamsCampaignSummary, startDate time.Time, endDate time.Time, groupType string) map[string]interface{} {
	days := map[string]map[string]map[string]interface{}{}
	totals := make(map[string]float64)
	monthlyBudgets := make(map[string]map[string]float64)
	countryDailyBudgets := make(map[string]map[string]float64)
	operatorDailyBudgets := make(map[string]map[string]float64)

	// Initialize dates
	currentDate := startDate
	for !currentDate.After(endDate) {
		date := formatDate(currentDate, params.DataType)

		days[date] = make(map[string]map[string]interface{})
		for _, indicator := range params.DataIndicators {
			days[date][indicator] = map[string]interface{}{
				"value":      0.0,
				"percentage": 0.0,
			}
		}
		currentDate = incrementDate(currentDate, params.DataType)
	}

	// Collect monthly budgets
	for _, campaign := range data {
		if containsString(params.DataIndicators, "target_daily_budget") {
			month := campaign.SummaryDate.Format("2006-01")
			key := fmt.Sprintf("%s|%s|%s", campaign.Country, campaign.Operator, month)
			budgetValue := getIndicatorValueRevenue(campaign, "target_daily_budget")
			if budgetValue > 0 {
				if monthlyBudgets[key] == nil {
					monthlyBudgets[key] = make(map[string]float64)
				}
				monthlyBudgets[key][month] = budgetValue
			}
		}
	}

	// Group by country and operator
	groupedData := make(map[string]map[string][]entity.CampaignSummaryMonitoring)
	for _, campaign := range data {
		if groupedData[campaign.Country] == nil {
			groupedData[campaign.Country] = make(map[string][]entity.CampaignSummaryMonitoring)
		}
		groupedData[campaign.Country][campaign.Operator] = append(groupedData[campaign.Country][campaign.Operator], campaign)
	}

	var results []map[string]interface{}

	// Accumulators for budget_usage correct calculation
	totalBudgetUsage := 0.0
	totalBudgetUsageDays := 0

	for country, operators := range groupedData {
		var countryTotal float64
		var operatorResults []map[string]interface{}

		if countryDailyBudgets[country] == nil {
			countryDailyBudgets[country] = make(map[string]float64)
		}

		for operator, campaigns := range operators {
			operatorKey := fmt.Sprintf("%s|%s", country, operator)
			operatorData := map[string]interface{}{}

			for _, indicator := range params.DataIndicators {
				if indicator != "target_daily_budget" && indicator != "budget_usage" {
					operatorData[indicator] = 0.0
				}
			}

			if operatorDailyBudgets[operatorKey] == nil {
				operatorDailyBudgets[operatorKey] = make(map[string]float64)
			}

			for _, campaign := range campaigns {
				date := formatDate(campaign.SummaryDate, params.DataType)

				for _, indicator := range params.DataIndicators {
					if indicator == "target_daily_budget" || indicator == "budget_usage" {
						continue
					}

					indicatorValue := getIndicatorValueRevenue(campaign, indicator)
					if days[date][indicator] != nil {
						days[date][indicator]["value"] = safeFloat(days[date][indicator], "value") + indicatorValue
					}
					if val, ok := operatorData[indicator].(float64); ok {
						operatorData[indicator] = val + indicatorValue
					}
					totals[indicator] += indicatorValue
				}
			}

			if containsString(params.DataIndicators, "target_daily_budget") {
				operatorTotal := 0.0
				currentDate := startDate
				for !currentDate.After(endDate) {
					month := currentDate.Format("2006-01")
					date := formatDate(currentDate, params.DataType)

					if budgets, exists := monthlyBudgets[operatorKey+"|"+month]; exists {
						if budget, ok := budgets[month]; ok {
							operatorDailyBudgets[operatorKey][date] = budget
							days[date]["target_daily_budget"]["value"] = safeFloat(days[date]["target_daily_budget"], "value") + budget
							countryDailyBudgets[country][date] += budget
							operatorTotal += budget
						}
					}
					currentDate = incrementDate(currentDate, params.DataType)
				}
				operatorData["target_daily_budget"] = operatorTotal
				totals["target_daily_budget"] += operatorTotal
				countryTotal += operatorTotal
			}

			operatorData["operator"] = operator
			operatorResults = append(operatorResults, operatorData)
		}

		countryData := map[string]interface{}{
			"country":   country,
			"_children": operatorResults,
		}

		if containsString(params.DataIndicators, "target_daily_budget") {
			countryData["target_daily_budget"] = map[string]interface{}{
				"total": countryTotal,
				"avg":   countryTotal / float64(len(countryDailyBudgets[country])),
			}
		}

		results = append(results, countryData)
	}

	// Handle budget_usage separately after grouping
	if containsString(params.DataIndicators, "budget_usage") {
		currentDate = startDate
		for !currentDate.After(endDate) {
			date := formatDate(currentDate, params.DataType)
			dailySpending := safeFloat(days[date]["spending_to_adnets"], "value")
			dailyBudget := safeFloat(days[date]["target_daily_budget"], "value")
			dailyUsage := 0.0
			if dailyBudget > 0 {
				dailyUsage = (dailySpending / dailyBudget) * 100
			}
			days[date]["budget_usage"]["value"] = dailyUsage

			// accumulate to get average
			totalBudgetUsage += dailyUsage
			totalBudgetUsageDays++

			currentDate = incrementDate(currentDate, params.DataType)
		}

		if totalBudgetUsageDays > 0 {
			totals["budget_usage"] = totalBudgetUsage
		}
	}

	// Calculate percentages
	currentDate = startDate
	for !currentDate.After(endDate) {
		date := formatDate(currentDate, params.DataType)
		prevDate := formatPreviousDate(currentDate, params.DataType)

		for _, indicator := range params.DataIndicators {
			if days[date][indicator] != nil {
				currentValue := safeFloat(days[date][indicator], "value")
				prevValue := safeFloat(days[prevDate][indicator], "value")
				days[date][indicator]["percentage"] = safePercentage(currentValue, prevValue)
			}
		}
		currentDate = incrementDate(currentDate, params.DataType)
	}

	// Final calculations
	tmoEnd := countTmoEndRevenue(totals, startDate, endDate)
	avg := countAverageRevenue(totals, startDate, endDate)

	// Special fix for budget_usage avg (because it's percentage, not amount)
	if containsString(params.DataIndicators, "budget_usage") && totalBudgetUsageDays > 0 {
		avg["budget_usage"] = totalBudgetUsage / float64(totalBudgetUsageDays)
	}

	return mergeDaysRevenue(map[string]interface{}{
		"data_indicators": params.DataIndicators,
		"total":           totals,
		"avg":             avg,
		"t_mo_end":        tmoEnd,
		"results":         results,
	}, days)
}

func getIndicatorValueRevenue(item entity.CampaignSummaryMonitoring, key string) float64 {

	key = SnakeToCamelValue(key)
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
func getPreviousValueRevenue(day map[string]map[string]interface{}, indicator string) float64 {
	if day == nil {
		return 0
	}
	if v, ok := day[indicator]; ok {
		if value, ok := v["value"].(float64); ok {
			return value
		}
	}
	return 0
}
func countPercentageRevenue(now, prev float64) float64 {
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

func countTmoEndRevenue(totals map[string]float64, startDate time.Time, endDate time.Time) map[string]float64 {
	tmoEnd := map[string]float64{}
	totalDaysRunning := int(endDate.Sub(startDate).Hours() / 24)
	if totalDaysRunning < 1 {
		totalDaysRunning = 1
	}

	// Calculate total days in the last month
	lastMonthEnd := endDate.AddDate(0, 0, -endDate.Day())
	lastMonthStart := lastMonthEnd.AddDate(0, 0, -lastMonthEnd.Day()+1)
	totalDaysLastMonth := int(lastMonthEnd.Sub(lastMonthStart).Hours()/24) + 1

	for key, value := range totals {
		result := (value / float64(totalDaysRunning)) * float64(totalDaysLastMonth)
		tmoEnd[key] = result
	}

	return tmoEnd
}

func countAverageRevenue(totals map[string]float64, startDate, endDate time.Time) map[string]float64 {
	averages := map[string]float64{}
	totalDaysRunning := int(endDate.Sub(startDate).Hours() / 24)
	if totalDaysRunning < 1 {
		totalDaysRunning = 1
	}

	for key, value := range totals {
		averages[key] = value / float64(totalDaysRunning)
	}

	return averages
}

func countBudgetUsage(totals map[string]float64) float64 {
	budget_usage := 0.0
	if totals["target_daily_budgets"] != 0.0 && totals["spending_to_adnets"] != 0 {
		budget_usage = totals["spending_to_adnets"] / totals["target_daily_budget"] * 100
	}
	return budget_usage
}

func extractQueryArrayRevenue(c *fiber.Ctx, key string) []string {
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

func SnakeToCamelValue(snake string) string {
	words := strings.Split(snake, "_")
	for i, word := range words {
		words[i] = strings.ToLower(word)
		if len(word) > 0 {
			words[i] = strings.ToUpper(words[i][:1]) + words[i][1:]
		}
	}
	return strings.Join(words, "")
}

func mergeDaysRevenue(summaryData map[string]interface{}, days map[string]map[string]map[string]interface{}) map[string]interface{} {
	for key, value := range days {
		summaryData[key] = value
	}
	return summaryData
}

func mergeMapsRevenue(map1, map2 map[string]interface{}) map[string]interface{} {
	mergedMap := make(map[string]interface{})
	for key, value := range map1 {
		mergedMap[key] = value
	}
	for key, value := range map2 {
		mergedMap[key] = value
	}
	return mergedMap
}

func sortDataRevenue(data []map[string]interface{}, dataBasedOn string, dataBasedOnIndicator string) []map[string]interface{} {
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

// Helper function to replace goterators.Contains
func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func formatDate(date time.Time, dataType string) string {
	if dataType == "monthly_report" {
		return date.Format("2006-01")
	}
	return date.Format("2006-01-02")
}

func incrementDate(date time.Time, dataType string) time.Time {
	if dataType == "monthly_report" {
		return date.AddDate(0, 1, 0)
	}
	return date.AddDate(0, 0, 1)
}

func formatPreviousDate(date time.Time, dataType string) string {
	if dataType == "monthly_report" {
		return date.AddDate(0, -1, 0).Format("2006-01")
	}
	return date.AddDate(0, 0, -1).Format("2006-01-02")
}

func safeFloat(m map[string]interface{}, key string) float64 {
	if m == nil {
		return 0
	}
	if val, ok := m[key]; ok && val != nil {
		if f, ok := val.(float64); ok {
			return f
		}
	}
	return 0
}

func safePercentage(current, previous float64) float64 {
	if previous == 0 {
		return 0
	}
	return ((current - previous) / previous) * 100
}
