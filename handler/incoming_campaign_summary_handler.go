package handler

import (
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
			generatedSummary := generateSummaryValue(data, params, startDate, endDate)
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
				generatedSummary := generateSummaryValue(campaignPerAdnet, params, startDate, endDate)

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
			generatedSummary := generateSummaryValue(data, params, startDate, endDate)
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
				generatedCountrySummary := generateSummaryValue(campaignPerCountry, params, startDate, endDate)

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

		generatedSummary := generateSummaryValue(campaignPerOperator, params, startDate, endDate)

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
		generatedSummary := generateSummaryValue(campaignPerPatner, params, startDate, endDate)
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

		generatedSummary := generateSummaryValue(campaignPerService, params, startDate, endDate)
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
		generatedSummary := generateSummaryValue(campaignPerAdnet, params, startDate, endDate)
		placeHolder := map[string]any{
			"level":   "adnet",
			"country": campaignPerAdnet[0].Adnet,
		}
		completeSummary := mergeMapsRevenue(generatedSummary, placeHolder)
		formattedData = append(formattedData, completeSummary)
	}
	return formattedData
}

func generateSummaryValue(data []entity.CampaignSummaryMonitoring, params entity.ParamsCampaignSummary, startDate time.Time, endDate time.Time) map[string]interface{} {
	days := map[string]map[string]map[string]interface{}{}
	totals := make(map[string]float64)

	for _, campaign := range data {
		date := campaign.SummaryDate.Format("2006-01-02")
		prevDate := campaign.SummaryDate.AddDate(0, 0, -1).Format("2006-01-02")

		if params.DataType == "monthly_report" {
			date = campaign.SummaryDate.Format("2006-01")
			prevDate = campaign.SummaryDate.AddDate(0, -1, 0).Format("2006-01")
		}

		// Initialize the day's indicator map if it doesn't exist
		if days[date] == nil {
			days[date] = make(map[string]map[string]interface{})
		}

		for _, indicator := range params.DataIndicators {
			indicatorValue := getIndicatorValueRevenue(campaign, indicator)
			if days[date][indicator] == nil {
				days[date][indicator] = map[string]interface{}{
					"value":      0.0, // Initialize "value" to 0
					"percentage": 0.0, // Initialize "percentage" to 0
				}
			}

			prevValue := getPreviousValueRevenue(days[prevDate], indicator)

			currentValue, ok := days[date][indicator]["value"].(float64)
			if !ok {
				currentValue = 0.0
			}
			newValue := currentValue + indicatorValue

			days[date][indicator]["value"] = newValue
			days[date][indicator]["percentage"] = countPercentageRevenue(newValue, prevValue)

			totals[indicator] += indicatorValue
		}
	}

	// Final calculations
	tmoEnd := countTmoEndRevenue(totals, startDate, endDate)

	// Prepare summary data
	summaryData := map[string]interface{}{
		"data_indicators": params.DataIndicators,
		"total":           totals,
		"avg":             countAverageRevenue(totals, startDate, endDate),
		"t_mo_end":        tmoEnd,
	}

	// Merge with daily breakdowns
	completeSummary := mergeDaysRevenue(summaryData, days)

	return completeSummary
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

func getPreviousValueRevenue(data map[string]map[string]interface{}, key string) float64 {
	value := 0.0
	if val, exists := data[key]; exists {
		value = val["value"].(float64)
	}
	return value
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
