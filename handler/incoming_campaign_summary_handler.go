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
)

func (h *IncomingHandler) DisplayCampaignSummary(c *fiber.Ctx) error {
	dataIndicators := extractQueryArray(c, "filter[data-indicators][]")

	fe := entity.DisplayCampaignSummary{
		DataType:       c.Query("filter[data-type]"),
		ReportType:     c.Query("filter[report-type]"),
		Country:        c.Query("filter[country]"),
		Operator:       c.Query("filter[operator]"),
		PartnerName:    c.Query("filter[partner_name]"),
		Adnet:          c.Query("filter[adnet]"),
		Service:        c.Query("filter[service]"),
		DataIndicators: dataIndicators,
		DataBasedOn:    c.Query("filter[data-based-on]"),
		DateRange:      c.Query("filter[date-range]"),
		CustomRange:    c.Query("filter[custom-range]"),
	}

	r := h.GenerateCampaignSummary(c, fe)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) GenerateCampaignSummary(c *fiber.Ctx, fe entity.DisplayCampaignSummary) entity.ReturnResponse {
	summaryCampaign, startDate, endDate, err := h.DS.GetSummaryCampaignMonitoring(fe)

	formattedSummaryCampaign := formatData(summaryCampaign, fe.ReportType, fe.DataIndicators, startDate, endDate)

	if err == nil {
		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithData{
				Code:    fiber.StatusOK,
				Message: config.OK_DESC,
				Data:    formattedSummaryCampaign,
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

type ReportSummary struct {
	DataIndicators []string
	Total          map[string]float64
	Avg            map[string]float64
	TmoEnd         map[string]float64
	Days           map[string]map[string]map[string]interface{}
}

func formatData(cmsData []entity.CampaignSummaryMonitoring, reportType string, dataIndicators []string, startDate, endDate time.Time) []map[string]interface{} {
	var formattedData []map[string]interface{}

	if reportType == "All" {
		summary := generateSummary(cmsData, map[string]string{"country": "All"}, "country", dataIndicators, startDate, endDate)
		data := map[string]interface{}{
			"level":   "country",
			"country": "All",
		}
		mergedData := mergeMaps(summary, data)
		formattedData = append(formattedData, mergedData)
	} else {
		uniqueCountries := getUniqueField(cmsData, "country")
		for _, country := range uniqueCountries {
			filters := map[string]string{"country": country}
			var children []interface{}

			switch reportType {
			case "campaign_summary":
				children = generatePartner(cmsData, filters, dataIndicators, startDate, endDate)
			case "url_service_summary":
				children = generateService(cmsData, filters, dataIndicators, startDate, endDate)
			case "adnet_summary":
				children = generateAdnet(cmsData, filters, dataIndicators, startDate, endDate)
			default:
				children = generateOperator(cmsData, filters, dataIndicators, startDate, endDate)
			}

			summary := generateSummary(cmsData, filters, "country", dataIndicators, startDate, endDate)

			data := map[string]interface{}{
				"level":     "country",
				"country":   country,
				"_children": children,
			}
			mergedData := mergeMaps(summary, data)
			formattedData = append(formattedData, mergedData)
		}
	}

	return formattedData
}

func sortData(data []map[string]interface{}, basedOn, indicator string) {
	sort.Slice(data, func(i, j int) bool {
		totalI := data[i]["summary"].(ReportSummary).Total[indicator]
		totalJ := data[j]["summary"].(ReportSummary).Total[indicator]
		if basedOn == "highest" {
			return totalI > totalJ
		}
		return totalI < totalJ
	})
}

func generateOperator(cmsData []entity.CampaignSummaryMonitoring, filters map[string]string, dataIndicators []string, startDate, endDate time.Time) []interface{} {
	filteredData := applyFilters(cmsData, filters)
	uniqueOperators := getUniqueField(filteredData, "operator")
	var formattedOperators []interface{}

	for _, operator := range uniqueOperators {
		filters["operator"] = operator
		summary := generateSummary(cmsData, filters, "operator", dataIndicators, startDate, endDate)
		children := generatePartner(cmsData, filters, dataIndicators, startDate, endDate)

		data := map[string]interface{}{
			"level":     "operator",
			"country":   operator,
			"_children": children,
		}
		mergedData := mergeMaps(summary, data)
		formattedOperators = append(formattedOperators, mergedData)
	}

	return formattedOperators
}

func generatePartner(cmsData []entity.CampaignSummaryMonitoring, filters map[string]string, dataIndicators []string, startDate, endDate time.Time) []interface{} {
	filteredData := applyFilters(cmsData, filters)
	uniquePartners := getUniqueField(filteredData, "partner")
	var formattedPartners []interface{}

	for _, partner := range uniquePartners {
		filters["partner"] = partner
		summary := generateSummary(cmsData, filters, "partner", dataIndicators, startDate, endDate)
		children := generateService(cmsData, filters, dataIndicators, startDate, endDate)

		data := map[string]interface{}{
			"level":     "partner",
			"country":   partner,
			"_children": children,
		}
		mergedData := mergeMaps(summary, data)

		formattedPartners = append(formattedPartners, mergedData)
	}

	return formattedPartners
}

func generateService(cmsData []entity.CampaignSummaryMonitoring, filters map[string]string, dataIndicators []string, startDate, endDate time.Time) []interface{} {
	filteredData := applyFilters(cmsData, filters)
	uniqueServices := getUniqueField(filteredData, "service")
	var formattedServices []interface{}

	for _, service := range uniqueServices {
		filters["service"] = service
		summary := generateSummary(cmsData, filters, "service", dataIndicators, startDate, endDate)
		children := generateAdnet(cmsData, filters, dataIndicators, startDate, endDate)

		data := map[string]interface{}{
			"level":     "service",
			"country":   service,
			"_children": children,
		}
		mergedData := mergeMaps(summary, data)
		formattedServices = append(formattedServices, mergedData)
	}

	return formattedServices
}

func generateAdnet(cmsData []entity.CampaignSummaryMonitoring, filters map[string]string, dataIndicators []string, startDate, endDate time.Time) []interface{} {
	filteredData := applyFilters(cmsData, filters)
	uniqueAdnets := getUniqueField(filteredData, "adnet")
	var formattedAdnets []interface{}

	for _, adnet := range uniqueAdnets {
		filters["adnet"] = adnet
		summary := generateSummary(cmsData, filters, "adnet", dataIndicators, startDate, endDate)

		data := map[string]interface{}{
			"level":   "adnet",
			"country": adnet,
			"summary": summary,
		}
		mergedData := mergeMaps(summary, data)
		formattedAdnets = append(formattedAdnets, mergedData)
	}

	return formattedAdnets
}

func generateSummary(cmsData []entity.CampaignSummaryMonitoring, filters map[string]string, level string, dataIndicators []string, startDate, endDate time.Time) map[string]interface{} {
	var data []entity.CampaignSummaryMonitoring
	if filters["country"] == "All" {
		data = cmsData
	} else {
		data = applyFilters(cmsData, filters)
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].SummaryDate.Before(data[j].SummaryDate)
	})

	days := map[string]map[string]map[string]interface{}{}
	totals := map[string]float64{}
	budgets := map[string]map[string]float64{}
	operator := map[string]map[string][]string{}
	budgetsTotal := map[string]float64{}
	counted := map[string]bool{}

	for _, item := range data {
		date := item.SummaryDate.Format("2006-01-02")
		prevDate := item.SummaryDate.AddDate(0, 0, -1).Format("2006-01-02")
		if level == "monthly_report" {
			date = item.SummaryDate.Format("2006-01")
			prevDate = item.SummaryDate.AddDate(0, -1, 0).Format("2006-01")
		}
		for _, dataIndicator := range dataIndicators {
			var value float64

			// Dynamically get the field value from item
			value = getField(item, dataIndicator)
			switch dataIndicator {
			case "target_daily_budget":
				value = item.TargetDailyBudget
			}

			if dataIndicator == "target_daily_budget" {
				if level == "country" {
					if _, exists := operator[date]; !exists {
						operator[date] = make(map[string][]string)
					}
					if _, exists := budgets[date]; !exists {
						budgets[date] = make(map[string]float64)
					}

					if !contains(operator[date][dataIndicator], item.Operator) {
						operator[date][dataIndicator] = append(operator[date][dataIndicator], item.Operator)
						budgets[date][dataIndicator] += value
					}
					value = budgets[date][dataIndicator]
				}

				if _, exists := budgetsTotal[dataIndicator]; !exists {
					budgetsTotal[dataIndicator] = 0
				}

				if !counted[date+item.Operator] {
					counted[date+item.Operator] = true
					budgetsTotal[dataIndicator] += value
				}

				if days[date] == nil {
					days[date] = make(map[string]map[string]interface{})
				}
				days[date][dataIndicator] = map[string]interface{}{
					"value":      value,
					"percentage": countPercentage(value, getNestedValue(days[prevDate], dataIndicator)),
				}
				totals[dataIndicator] = budgetsTotal[dataIndicator]
			} else {
				var prevValue float64
				if days[date] == nil {
					days[date] = make(map[string]map[string]interface{})
				}
				if days[date][dataIndicator] == nil {
					days[date][dataIndicator] = map[string]interface{}{}
				}
				if currentData, exists := days[date][dataIndicator]; exists {
					if val, ok := currentData["value"]; ok && val != nil {
						prevValue = val.(float64)
					} else {
						prevValue = 0
					}
				} else {
					days[date][dataIndicator] = map[string]interface{}{}
				}
				updatedValue := prevValue + value
				days[date][dataIndicator] = map[string]interface{}{
					"value":      updatedValue,
					"percentage": countPercentage(updatedValue, getNestedValue(days[prevDate], dataIndicator)),
				}
				totals[dataIndicator] += value
			}
		}

	}

	tmoEnd := countTmoEnd(totals, startDate, endDate)

	collectedData := map[string]interface{}{
		"data_indicators": dataIndicators,
		"total":           totals,
		"avg":             countAverage(totals, len(days), startDate, endDate),
		"t_mo_end":        tmoEnd,
	}
	mergedData := mergeDays(collectedData, days)

	return mergedData
}

func getNestedValue(data map[string]map[string]interface{}, key string) float64 {
	if val, exists := data[key]; exists {
		return val["value"].(float64)
	}
	return 0
}

func countAverage(totals map[string]float64, startDate, endDate time.Time) map[string]float64 {
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

func countTmoEnd(totals map[string]float64, startDate, endDate time.Time) map[string]float64 {
	tmoEnd := map[string]float64{}
	totalDaysRunning := int(endDate.Sub(startDate).Hours() / 24)
	if totalDaysRunning < 1 {
		totalDaysRunning = 1
	}

	totalDaysLastMonth := int(endDate.AddDate(0, 1, -1).Sub(startDate).Hours() / 24)

	for key, value := range totals {
		result := (value / float64(totalDaysRunning)) * float64(totalDaysLastMonth)
		tmoEnd[key] = result
	}

	return tmoEnd
}

func countPercentage(now, prev float64) float64 {
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

func applyFilters(cmsData []entity.CampaignSummaryMonitoring, filters map[string]string) []entity.CampaignSummaryMonitoring {
	var filteredData []entity.CampaignSummaryMonitoring
	for _, data := range cmsData {
		match := true
		for key, value := range filters {
			switch key {
			case "country":
				if data.Country != value {
					match = false
				}
			case "operator":
				if data.Operator != value {
					match = false
				}
			case "partner":
				if data.Partner != value {
					match = false
				}
			case "service":
				if data.Service != value {
					match = false
				}
			case "adnet":
				if data.Adnet != value {
					match = false
				}
			}
			if !match {
				break
			}
		}
		if match {
			filteredData = append(filteredData, data)
		}
	}
	return filteredData
}

// Helper functions to get unique values, adjust as needed.
func getUniqueField(cmsData []entity.CampaignSummaryMonitoring, field string) []string {
	switch field {
	case "country":
		return getUniqueValues(cmsData, func(data entity.CampaignSummaryMonitoring) string { return data.Country })
	case "operator":
		return getUniqueValues(cmsData, func(data entity.CampaignSummaryMonitoring) string { return data.Operator })
	case "partner":
		return getUniqueValues(cmsData, func(data entity.CampaignSummaryMonitoring) string { return data.Partner })
	case "service":
		return getUniqueValues(cmsData, func(data entity.CampaignSummaryMonitoring) string { return data.Service })
	case "adnet":
		return getUniqueValues(cmsData, func(data entity.CampaignSummaryMonitoring) string { return data.Adnet })
	default:
		return nil
	}
}
func getUniqueValues(cmsData []entity.CampaignSummaryMonitoring, keyFunc func(entity.CampaignSummaryMonitoring) string) []string {
	uniqueMap := make(map[string]struct{})
	var uniqueList []string

	for _, data := range cmsData {
		key := keyFunc(data)
		if _, exists := uniqueMap[key]; !exists {
			uniqueMap[key] = struct{}{}
			uniqueList = append(uniqueList, key)
		}
	}

	return uniqueList
}

func contains(slice []string, item string) bool {
	for _, elem := range slice {
		if elem == item {
			return true
		}
	}
	return false
}

// Helper function, consider moving to a utility package
func extractQueryArray(c *fiber.Ctx, key string) []string {
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

func mergeMaps(map1, map2 map[string]interface{}) map[string]interface{} {
	mergedMap := make(map[string]interface{})

	for key, value := range map1 {
		mergedMap[key] = value
	}

	for key, value := range map2 {
		mergedMap[key] = value
	}

	return mergedMap
}

func mergeDays(map1 map[string]interface{}, map2 map[string]map[string]map[string]interface{}) map[string]interface{} {
	mergedMap := make(map[string]interface{})

	for key, value := range map1 {
		mergedMap[key] = value
	}

	for key, value := range map2 {
		mergedMap[key] = value
	}

	return mergedMap
}

func getField(item interface{}, field string) float64 {
	val := reflect.ValueOf(item)
	field = SnakeToCamel(field)
	fieldVal := val.FieldByName(field)

	if !fieldVal.IsValid() {
		return 0
	}

	switch fieldVal.Kind() {
	case reflect.Float64:
		return fieldVal.Float()
	case reflect.Int, reflect.Int64:
		return float64(fieldVal.Int())
	default:
		return 0
	}
}

func SnakeToCamel(snake string) string {
	words := strings.Split(snake, "_")
	for i, word := range words {
		words[i] = strings.ToLower(word)
		if len(word) > 0 {
			words[i] = strings.ToUpper(words[i][:1]) + words[i][1:]
		}
	}
	return strings.Join(words, "")
}
