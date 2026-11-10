package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource/entity"
	"github.com/infraLinkit/mediaplatform-datasource/helper"
)

// campaignRow is a campaigns or campaign_details row under its column names, as Eloquent serialized it.
type campaignRow map[string]interface{}

// str is PHP's (string) cast on a column; NULL is an empty string.
func (r campaignRow) str(column string) string {
	switch v := r[column].(type) {
	case nil:
		return ""
	case string:
		return v
	case []byte:
		return string(v)
	case bool:
		if v {
			return "1"
		}
		return ""
	case time.Time:
		return v.Format("2006-01-02 15:04:05")
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return fmt.Sprint(v)
	}
}

func (r campaignRow) strOr(column, fallback string) string {
	if r[column] == nil {
		return fallback
	}
	return r.str(column)
}

func (r campaignRow) num(column string) int { return phpInt(r.str(column)) }

func (r campaignRow) numOr(column string, fallback int) int {
	if r[column] == nil {
		return fallback
	}
	return r.num(column)
}

func (r campaignRow) flag(column string) bool {
	if v, ok := r[column].(bool); ok {
		return v
	}
	return phpTruthy(r.str(column))
}

func (r campaignRow) flagOr(column string, fallback bool) bool {
	if r[column] == nil {
		return fallback
	}
	return r.flag(column)
}

// orValue is the row's column, or the fallback when it is NULL; these reach @json untouched.
func (r campaignRow) orValue(column string, fallback interface{}) interface{} {
	if r[column] == nil {
		return fallback
	}
	return r[column]
}

// fillStoredColumns keeps what the form does not manage. applyCampaignUpdate writes every column of the
// entry, so a column left out of the payload would otherwise be blanked by saving the form. Laravel lost
// them this way; the stored value is carried along here instead.
func fillStoredColumns(entry map[string]interface{}, src campaignRow) {
	fields := reflect.TypeOf(entity.DataConfig{})
	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i)
		key := strings.Split(field.Tag.Get("json"), ",")[0]
		if _, sent := entry[key]; sent {
			continue
		}
		column := key
		switch key {
		case "id", "ip_address", "last_update", "last_update_capping":
			continue
		case "urlservicekey":
			column = "url_service_key"
		case "url_tfc_or_smartlink":
			column = "url_tfcor_smartlink"
		}
		if src[column] == nil {
			continue
		}
		switch field.Type.Kind() {
		case reflect.String:
			entry[key] = src.str(column)
		case reflect.Bool:
			entry[key] = src.flag(column)
		case reflect.Int:
			entry[key] = src.num(column)
		case reflect.Float64:
			if number, err := strconv.ParseFloat(src.str(column), 64); err == nil {
				entry[key] = number
			}
		}
	}
}

var htmlEntity = regexp.MustCompile(`(?i)&(#\d+|#x[0-9a-f]+|[a-z][a-z0-9]*);`)

// escapeKeepingEntities is htmlspecialchars with double_encode off. The edit forms send back text that is
// already escaped, so re-encoding it would add another layer of &amp; on every save.
func escapeKeepingEntities(value string) string {
	var out strings.Builder
	last := 0
	for _, match := range htmlEntity.FindAllStringIndex(value, -1) {
		out.WriteString(phpHTMLSpecialChars(value[last:match[0]]))
		out.WriteString(value[match[0]:match[1]])
		last = match[1]
	}
	out.WriteString(phpHTMLSpecialChars(value[last:]))
	return out.String()
}

// strOr reads an input with a ?? fallback.
func (f campaignForm) strOr(fallback string, keys ...string) string {
	if f.value(keys...) == nil {
		return fallback
	}
	return f.str(keys...)
}

// flagOr is filter_var(input ?? column, FILTER_VALIDATE_BOOLEAN).
func (f campaignForm) flagOr(row campaignRow, column string, keys ...string) bool {
	if f.value(keys...) == nil {
		return row.flag(column)
	}
	return phpFilterBool(f.str(keys...))
}

func baseURL(raw string) string {
	scheme, host := "http", ""
	if parsed, err := url.Parse(raw); err == nil {
		if parsed.Scheme != "" {
			scheme = parsed.Scheme
		}
		host = parsed.Host
	}
	return scheme + "://" + host
}

func phpTruthyRow(row campaignRow, column string) bool { return phpTruthy(row.str(column)) }

func nestedSet(root map[string]map[string]interface{}, operator, service string, value interface{}) {
	if root[operator] == nil {
		root[operator] = map[string]interface{}{}
	}
	root[operator][service] = value
}

// buildCampaignEditData computes the variables CampaignManagementController::edit hands to the view.
func buildCampaignEditData(campaign campaignRow, details []campaignRow) fiber.Map {
	selectedOperators := []string{}
	selectedAdnets := []string{}
	selectedServices := map[string][]string{}
	seenOperator, seenAdnet := map[string]bool{}, map[string]bool{}
	seenService := map[string]bool{}
	serviceDomains := map[string]map[string]interface{}{}
	shortCodes := map[string]map[string]interface{}{}
	serviceUrls := map[string]map[string]interface{}{}
	keywords := map[string]map[string]map[string]interface{}{}
	oldUrls := map[string]map[string]map[string]interface{}{}
	oldTnc := map[string]interface{}{}
	oldPage2 := map[string]interface{}{}
	selectedDomain := interface{}(nil)

	for _, detail := range details {
		operator, service, adnet := detail.str("operator"), detail.str("service"), detail.str("adnet")
		if !seenOperator[operator] {
			seenOperator[operator] = true
			selectedOperators = append(selectedOperators, operator)
		}
		if !seenAdnet[adnet] {
			seenAdnet[adnet] = true
			selectedAdnets = append(selectedAdnets, adnet)
		}
		if key := operator + "\x00" + service; !seenService[key] {
			seenService[key] = true
			selectedServices[operator] = append(selectedServices[operator], service)
		}
		if serviceDomains[operator] == nil {
			serviceDomains[operator] = map[string]interface{}{}
		}
		if phpTruthyRow(detail, "url_warp_landing") {
			serviceDomains[operator][service] = baseURL(detail.str("url_warp_landing"))
			nestedSet(shortCodes, operator, service, detail["short_code"])
			serviceURL := detail.str("url_service")
			if adnet != "" && strings.Contains(serviceURL, "="+adnet) {
				serviceURL = strings.ReplaceAll(serviceURL, "="+adnet, "=[adnet]")
			}
			nestedSet(serviceUrls, operator, service, serviceURL)
		}
		if keywords[operator] == nil {
			keywords[operator] = map[string]map[string]interface{}{}
			oldUrls[operator] = map[string]map[string]interface{}{}
		}
		if keywords[operator][service] == nil {
			keywords[operator][service] = map[string]interface{}{}
			oldUrls[operator][service] = map[string]interface{}{}
		}
		keywords[operator][service][adnet] = detail["keyword"]
		oldUrls[operator][service][adnet] = map[string]interface{}{
			"client_type":                  detail["client_type"],
			"unused_traffic_redirect_type": detail["unused_traffic_redirect_type"],
			"url":                          detail["url_tfcor_smartlink"],
		}
		oldTnc[operator] = map[string]interface{}{
			"tnc":           detail["tnc"],
			"tnc_color":     detail["tnc_color"],
			"tnc_style":     detail["tnc_style"],
			"tnc_page_type": detail["tnc_page_type"],
			"tnc_font_size": detail["tnc_font_size"],
		}
		oldPage2[operator] = map[string]interface{}{
			"campaign_detail_name":     detail["campaign_detail_name"],
			"partner":                  detail["partner"],
			"channel":                  detail["channel"],
			"adnet":                    detail["adnet"],
			"service":                  detail["service"],
			"short_code":               detail["short_code"],
			"url_service":              detail["url_service"],
			"currency":                 detail["currency"],
			"is_evina":                 detail["is_evina"],
			"evina_redirect_fraud_url": detail["evina_redirect_fraud_url"],
			"prefix":                   detail["prefix"],
			"keyword":                  detail["keyword"],
			"daily_limit":              detail["mo_capping"],
			"payout":                   detail["po"],
			"conversion_name":          detail["conversion_name"],
			"google_sheet":             detail["google_sheet"],
			"google_sheet_billable":    detail["google_sheet_billable"],
			"url_type":                 detail.orValue("url_type", "wap"),
			"glob_post":                detail.orValue("glob_post", false),
			"click_type":               detail.orValue("click_type", 1),
			"click_delay":              detail.orValue("click_delay", 0),
			"traffic_source":           detail.orValue("traffic_source", false),
			"unique_click":             detail.orValue("unique_click", false),
			"clickable_anywhere":       detail.orValue("clickable_anywhere", false),
			"custom_integration":       detail.orValue("custom_integration", "default"),
			"redirect_type":            detail.strOr("unused_traffic_redirect_type", ""),
			"url_tfc_or_smartlink":     detail.strOr("url_tfcor_smartlink", ""),
		}
	}

	var partner, company, selectedAgency interface{}
	if len(details) > 0 {
		selectedAgency = details[0]["adnet"]
		last := details[len(details)-1]
		if phpTruthyRow(last, "url_warp_landing") {
			selectedDomain = baseURL(last.str("url_warp_landing"))
		}
		partner, company = details[0]["partner"], details[0]["advertiser"]
	}

	detailList := make([]map[string]interface{}, 0, len(details))
	for _, detail := range details {
		detailList = append(detailList, detail)
	}
	withDetails := map[string]interface{}{}
	for key, value := range campaign {
		withDetails[key] = value
	}
	withDetails["details"] = detailList

	return fiber.Map{
		"campaign":             withDetails,
		"partner":              partner,
		"company":              company,
		"selectedOperators":    selectedOperators,
		"selectedServices":     selectedServices,
		"selectedAdnets":       selectedAdnets,
		"serviceDomains":       serviceDomains,
		"shortCodes":           shortCodes,
		"keywords":             keywords,
		"serviceUrls":          serviceUrls,
		"oldUrlTfcOrSmartlink": oldUrls,
		"oldTnc":               oldTnc,
		"oldPage2":             oldPage2,
		"selectedDomain":       selectedDomain,
		"selectedAgency":       selectedAgency,
	}
}

// loadCampaign reads a campaign and its details; selectedIDs narrows the details when given.
func (h *IncomingHandler) loadCampaign(campaignID string, selectedIDs []int) (campaignRow, []campaignRow, error) {
	var campaigns []map[string]interface{}
	if err := h.DB.Table("campaigns").Where("campaign_id = ?", campaignID).Limit(1).Find(&campaigns).Error; err != nil {
		return nil, nil, err
	}
	if len(campaigns) == 0 {
		return nil, nil, nil
	}
	query := h.DB.Table("campaign_details").Where("campaign_id = ?", campaignID)
	if selectedIDs != nil {
		query = query.Where("id IN ?", selectedIDs)
	}
	var rows []map[string]interface{}
	if err := query.Order("id ASC").Find(&rows).Error; err != nil {
		return nil, nil, err
	}
	details := make([]campaignRow, 0, len(rows))
	for _, row := range rows {
		details = append(details, row)
	}
	return campaigns[0], details, nil
}

// CampaignEditData serves the data of CampaignManagementController::edit.
func (h *IncomingHandler) CampaignEditData(c *fiber.Ctx) error {
	var selectedIDs []int
	queries := c.Context().QueryArgs()
	raw := queries.PeekMulti("selected_ids")
	raw = append(raw, queries.PeekMulti("selected_ids[]")...)
	if len(raw) > 0 {
		selectedIDs = []int{}
		for _, item := range raw {
			if id, err := strconv.Atoi(string(item)); err == nil {
				selectedIDs = append(selectedIDs, id)
			}
		}
	}
	campaign, details, err := h.loadCampaign(c.Params("campaign_id"), selectedIDs)
	if err != nil {
		h.Logs.Errorf("campaign edit data: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Something went wrong"})
	}
	if campaign == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Campaign not found."})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"code": fiber.StatusOK, "data": buildCampaignEditData(campaign, details)})
}

func sortedKeys(level map[string]interface{}) []string {
	keys := make([]string, 0, len(level))
	for key := range level {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// updatedDetail is one entry of $campaignData['data'] in updateCampaignData. src is the detail the
// unchanged columns come from: the matching row, or the reference row for an adnet added on the form.
func updatedDetail(f campaignForm, campaign, src, reference campaignRow, isNew bool, mcc, now, operator, service, adnet, keyword, urlService, urlLanding, urlWarp, apiURL, urlRedirect, redirectType string) (map[string]interface{}, error) {
	var globPost bool
	if f.value("glob_post") != nil {
		globPost = phpFilterBool(f.str("glob_post"))
	} else if reference != nil {
		globPost = reference.flag("glob_post")
	} else {
		return nil, errors.New("glob_post has no input and no reference detail")
	}
	channel := f.str("channel")
	if f.value("channel") == nil {
		channel = src.str("channel")
	}
	detail := map[string]interface{}{
		"channel":                       strings.ToLower(channel),
		"campaign_id":                   f.str("campaign_id"),
		"urlservicekey":                 src.str("url_service_key"),
		"name":                          campaign.str("name"),
		"objective":                     campaign.str("campaign_objective"),
		"country":                       src.str("country"),
		"mcc":                           mcc,
		"advertiser":                    f.str("advertiser"),
		"operator":                      src.str("operator"),
		"partner":                       src.str("partner"),
		"url_postback":                  f.str("url_postback"),
		"postback_method":               f.str("postback_method"),
		"aggregator":                    src.str("aggregator"),
		"adnet":                         adnet,
		"service":                       src.str("service"),
		"keyword":                       keyword,
		"subkeyword":                    src.str("subkeyword"),
		"is_billable":                   src.flag("is_billable"),
		"plan":                          src.str("plan"),
		"po":                            src.str("po"),
		"currency":                      src.str("currency"),
		"cost":                          src.str("cost"),
		"pubid":                         src.str("pubid"),
		"short_code":                    src.str("short_code"),
		"device_type":                   src.str("device_type"),
		"os":                            src.str("os"),
		"url_type":                      f.str("url_type"),
		"click_type":                    phpInt(f.str("click_type")),
		"click_delay":                   phpInt(f.str("click_delay")),
		"client_type":                   strings.ToLower(f.str("client_type")),
		"unused_traffic_redirect_type":  redirectType,
		"traffic_source":                f.flagOr(src, "traffic_source", "traffic_source"),
		"unique_click":                  f.flagOr(src, "unique_click", "unique_click"),
		"clickable_anywhere":            f.flagOr(src, "clickable_anywhere", "clickable_anywhere"),
		"url_banner":                    src.str("url_banner"),
		"url_landing":                   urlLanding,
		"url_warp_landing":              urlWarp,
		"url_service":                   urlService,
		"url_tfc_or_smartlink":          urlRedirect,
		"glob_post":                     globPost,
		"url_glob_post":                 "NA",
		"custom_integration":            f.str("custom_integration"),
		"ip_address":                    []interface{}{},
		"isp":                           "NA",
		"is_active":                     src.flag("is_active"),
		"mo_capping":                    src.num("mo_capping"),
		"counter_mo_capping":            src.num("counter_mo_capping"),
		"mo_capping_service":            src.num("mo_capping_service"),
		"counter_mo_capping_service":    src.num("counter_mo_capping_service"),
		"status_capping":                src.flag("status_capping"),
		"kpi_upper_limit_capping":       src.num("kpi_upper_limit_capping"),
		"is_machine_learning_capping":   src.flag("is_machine_learning_capping"),
		"ratio_send":                    src.num("ratio_send"),
		"ratio_receive":                 src.num("ratio_receive"),
		"counter_mo_ratio":              src.num("counter_mo_ratio"),
		"status_ratio":                  src.flag("status_ratio"),
		"kpi_upper_limit_ratio_send":    src.num("kpi_upper_limit_ratio_send"),
		"kpi_upper_limit_ratio_receive": src.num("kpi_upper_limit_ratio_receive"),
		"is_machine_learning_ratio":     src.flag("is_machine_learning_ratio"),
		"api_url":                       apiURL,
		"is_evina":                      f.flagOr(src, "is_evina", "is_evina", operator, service),
		"evina_redirect_fraud_url":      f.strOr(src.str("evina_redirect_fraud_url"), "evina_redirect_fraud_url", operator, service),
		"last_update":                   src.str("last_update"),
		"last_update_capping":           src.str("last_update_capping"),
	}
	if isNew {
		detail["urlservicekey"] = "NEW_KEY"
		detail["operator"] = operator
		detail["service"] = service
		detail["is_active"] = true
		detail["last_update"] = now
		detail["last_update_capping"] = now
	}
	fillStoredColumns(detail, src)
	return detail, nil
}

// buildCampaignUpdatePayload is updateCampaignData: one entry for each adnet and each combination of
// the url inputs, so a form that left one of them out produces no entry for that adnet at all.
func buildCampaignUpdatePayload(f campaignForm, campaign campaignRow, details []campaignRow, mcc, now string) (map[string]interface{}, map[string]interface{}, error) {
	operators, ok := f.value("adnet").(map[string]interface{})
	if !ok {
		return nil, nil, errors.New("adnet is not an array")
	}
	data := []map[string]interface{}{}
	logAfter := map[string]interface{}{}
	for _, operator := range sortedKeys(operators) {
		services, ok := operators[operator].(map[string]interface{})
		if !ok {
			return nil, nil, errors.New("adnet services is not an array")
		}
		for _, service := range sortedKeys(services) {
			urlServices := f.list("url_service", operator, service)
			seen := map[string]bool{}
			for _, adnet := range f.list("adnet", operator, service) {
				if seen[adnet] {
					continue
				}
				seen[adnet] = true
				urlLandings := f.list("url_landing", operator, service, adnet)
				urlWarps := f.list("url_warp_landing", operator, service, adnet)
				apiURLs := f.list("api_url", operator, service, adnet)
				urlRedirects := f.list("url_tfc_or_smartlink", operator, service, adnet)
				redirectTypes := f.list("redirect_type", operator, service, adnet)
				for _, urlService := range urlServices {
					for _, urlLanding := range urlLandings {
						for _, urlWarp := range urlWarps {
							for _, apiURL := range apiURLs {
								for _, urlRedirect := range urlRedirects {
									for _, redirectType := range redirectTypes {
										var reference campaignRow
										isNew := true
										for _, row := range details {
											if row.str("operator") == operator && row.str("service") == service {
												if row.str("adnet") == adnet {
													isNew = false
													break
												}
												reference = row
											}
										}
										if isNew && reference == nil {
											return nil, nil, fmt.Errorf("no reference detail for %s %s %s", operator, service, adnet)
										}
										words := f.list("keyword", operator, service, adnet)
										if len(words) == 0 {
											return nil, nil, fmt.Errorf("no keyword for %s %s %s", operator, service, adnet)
										}
										sources := []campaignRow{reference}
										if !isNew {
											sources = sources[:0]
											for _, row := range details {
												if row.str("operator") == operator && row.str("service") == service && row.str("adnet") == adnet {
													sources = append(sources, row)
												}
											}
										}
										for _, src := range sources {
											entry, err := updatedDetail(f, campaign, src, reference, isNew, mcc, now, operator, service, adnet, words[0], urlService, urlLanding, urlWarp, apiURL, urlRedirect, redirectType)
											if err != nil {
												return nil, nil, err
											}
											data = append(data, entry)
											logAfter[operator+" - "+service+" - "+adnet] = entry
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return map[string]interface{}{
		"action":             "ADD_UPDATE",
		"campaign_id":        f.str("campaign_id"),
		"name":               campaign.str("name"),
		"campaign_objective": campaign.str("campaign_objective"),
		"country":            campaign.str("country"),
		"advertiser":         campaign.str("advertiser"),
		"data":               data,
	}, logAfter, nil
}

// buildUploadSmsUpdatePayload is updateUploadSmsCampaignData; only the channel comes from the form.
func buildUploadSmsUpdatePayload(f campaignForm, campaign, detail campaignRow, mcc string) map[string]interface{} {
	channel := f.str("channel")
	if f.value("channel") == nil {
		channel = detail.str("channel")
	}
	entry := map[string]interface{}{
		"channel":                       strings.ToLower(channel),
		"campaign_id":                   f.str("campaign_id"),
		"urlservicekey":                 detail.str("url_service_key"),
		"name":                          campaign.str("name"),
		"objective":                     campaign.str("campaign_objective"),
		"country":                       campaign.str("country"),
		"mcc":                           mcc,
		"advertiser":                    campaign.str("advertiser"),
		"operator":                      detail.str("operator"),
		"partner":                       detail.str("partner"),
		"url_postback":                  detail.str("url_postback"),
		"postback_method":               detail.str("postback_method"),
		"aggregator":                    detail.str("aggregator"),
		"adnet":                         detail.str("adnet"),
		"service":                       detail.str("service"),
		"keyword":                       detail.str("keyword"),
		"subkeyword":                    detail.str("subkeyword"),
		"is_billable":                   detail.flag("is_billable"),
		"plan":                          detail.strOr("plan", "daily"),
		"po":                            detail.str("po"),
		"currency":                      detail.strOr("currency", "USD"),
		"cost":                          detail.strOr("cost", "0"),
		"pubid":                         detail.strOr("pubid", "0"),
		"short_code":                    detail.str("short_code"),
		"device_type":                   detail.strOr("device_type", "NA"),
		"os":                            detail.strOr("os", "NA"),
		"url_type":                      detail.strOr("url_type", "wap"),
		"click_type":                    detail.numOr("click_type", 1),
		"click_delay":                   detail.numOr("click_delay", 0),
		"client_type":                   strings.ToLower(detail.str("client_type")),
		"traffic_source":                detail.flag("traffic_source"),
		"unique_click":                  detail.flag("unique_click"),
		"url_banner":                    detail.str("url_banner"),
		"url_landing":                   detail.str("url_landing"),
		"url_warp_landing":              detail.str("url_warp_landing"),
		"url_service":                   detail.str("url_service"),
		"url_tfc_or_smartlink":          detail.str("url_tfcor_smartlink"),
		"unused_traffic_redirect_type":  detail.str("unused_traffic_redirect_type"),
		"glob_post":                     detail.flag("glob_post"),
		"url_glob_post":                 "NA",
		"custom_integration":            detail.str("custom_integration"),
		"ip_address":                    []interface{}{},
		"isp":                           "NA",
		"is_active":                     detail.flag("is_active"),
		"mo_capping":                    detail.numOr("mo_capping", 500),
		"counter_mo_capping":            detail.numOr("counter_mo_capping", 0),
		"mo_capping_service":            detail.numOr("mo_capping_service", 0),
		"counter_mo_capping_service":    detail.numOr("counter_mo_capping_service", 0),
		"status_capping":                detail.flagOr("status_capping", false),
		"kpi_upper_limit_capping":       detail.numOr("kpi_upper_limit_capping", 500),
		"is_machine_learning_capping":   detail.flagOr("is_machine_learning_capping", false),
		"ratio_send":                    detail.numOr("ratio_send", 1),
		"ratio_receive":                 detail.numOr("ratio_receive", 4),
		"counter_mo_ratio":              detail.numOr("counter_mo_ratio", 0),
		"status_ratio":                  detail.flagOr("status_ratio", true),
		"kpi_upper_limit_ratio_send":    detail.numOr("kpi_upper_limit_ratio_send", 1),
		"kpi_upper_limit_ratio_receive": detail.numOr("kpi_upper_limit_ratio_receive", 2),
		"is_machine_learning_ratio":     detail.flagOr("is_machine_learning_ratio", false),
		"api_url":                       detail.str("api_url"),
		"last_update":                   detail.str("last_update"),
		"last_update_capping":           detail.str("last_update_capping"),
	}
	fillStoredColumns(entry, detail)
	return map[string]interface{}{
		"action":             "ADD_UPDATE",
		"campaign_id":        f.str("campaign_id"),
		"name":               campaign.str("name"),
		"campaign_objective": campaign.str("campaign_objective"),
		"country":            campaign.str("country"),
		"advertiser":         campaign.str("advertiser"),
		"data":               []map[string]interface{}{entry},
	}
}

// buildSingleUrlUpdatePayload is updateCampaignSingleUrlS2SData after its uploads were handled.
func buildSingleUrlUpdatePayload(f campaignForm, campaign campaignRow, details []campaignRow, urlBanner, urlBackground string) (map[string]interface{}, error) {
	name, country, advertiser := f.str("campaign_name"), f.str("country"), f.str("advertiser")
	bypass := campaign.flag("is_bypass_ip_range")
	if f.value("is_bypass_ip_range") != nil {
		bypass = phpTruthy(f.str("is_bypass_ip_range"))
	}
	data := []map[string]interface{}{}
	for _, operator := range f.list("operator") {
		redirectType := f.str("redirect_type", operator)
		urlRedirect := ""
		switch redirectType {
		case "tfc":
			urlRedirect = f.str("redirect_internal", operator)
		case "smartlink":
			urlRedirect = f.str("redirect_smartlink", operator)
		}
		clickDelay := 0
		if f.strOr("no", "click_delay_radio", operator) == "yes" {
			clickDelay = phpInt(f.str("click_delay", operator))
		}
		var matching campaignRow
		for _, row := range details {
			if row.str("operator") == operator {
				matching = row
				break
			}
		}
		channel := f.str("channel")
		if f.value("channel") == nil {
			if matching == nil {
				return nil, errors.New("channel has no input and no matching detail")
			}
			channel = matching.str("channel")
		}
		isActive, partner, service, serviceKey := true, f.str("partner", operator), f.str("service", operator), ""
		if matching != nil {
			isActive = matching.flagOr("is_active", true)
			partner = matching.strOr("partner", partner)
			service = matching.strOr("service", service)
			serviceKey = matching.str("url_service_key")
		}
		capping := 500
		if f.value("daily_limit", operator) != nil {
			capping = phpInt(f.str("daily_limit", operator))
		}
		backgroundColor := "NA"
		if phpTruthy(f.str("enable_color")) {
			backgroundColor = f.strOr("NA", "background_color")
		}
		title := escapeKeepingEntities(f.str("title"))
		button := escapeKeepingEntities(f.str("button_subscribe"))
		tnc := escapeKeepingEntities(f.str("tnc", operator))
		titleStyle, titlePageType := f.str("title_style"), f.str("title_page_type")
		if matching != nil {
			titleStyle = f.strOr(matching.str("title_style"), "title_style")
			titlePageType = f.strOr(matching.str("title_page_type"), "title_page_type")
		}
		entry := map[string]interface{}{
			"campaign_id":                  f.str("campaign_id"),
			"urlservicekey":                serviceKey,
			"is_active":                    isActive,
			"channel":                      strings.ToLower(channel),
			"name":                         name,
			"campaign_detail_name":         name,
			"country":                      country,
			"mcc":                          f.str("mcc"),
			"is_bypass_ip_range":           bypass,
			"country_dialing_code":         f.str("country_dialing_code"),
			"advertiser":                   advertiser,
			"operator":                     operator,
			"partner":                      partner,
			"url_postback":                 f.str("url_postback", operator),
			"postback_method":              f.str("postback_method", operator),
			"client_type":                  strings.ToLower(f.str("client_type", operator)),
			"adnet":                        f.str("adnet"),
			"api_url":                      f.str("api_url"),
			"url_warp_landing":             f.str("domain"),
			"service":                      service,
			"url_service":                  f.str("url_service", operator),
			"short_code":                   f.str("shortcode", operator),
			"prefix":                       f.str("prefix", operator),
			"keyword":                      f.str("keyword", operator),
			"mo_capping":                   capping,
			"po":                           f.strOr("0", "payout", operator),
			"url_type":                     f.strOr("wap", "url_type", operator),
			"glob_post":                    phpFilterBool(f.str("glob_post", operator)),
			"click_type":                   phpInt(f.strOr("1", "click_type", operator)),
			"click_delay":                  clickDelay,
			"traffic_source":               phpFilterBool(f.str("traffic_source", operator)),
			"unique_click":                 phpFilterBool(f.str("unique_click", operator)),
			"clickable_anywhere":           phpFilterBool(f.str("clickable_anywhere", operator)),
			"is_evina":                     phpFilterBool(f.str("is_evina", operator)),
			"evina_redirect_fraud_url":     f.str("evina_redirect_fraud_url", operator),
			"custom_integration":           f.strOr("default", "custom_integration", operator),
			"unused_traffic_redirect_type": redirectType,
			"url_tfc_or_smartlink":         urlRedirect,
			"mainstream_lp_type":           strings.ToLower(f.str("mainstream_lp_type")),
			"title":                        title,
			"title_original":               title,
			"title_color":                  f.str("title_color"),
			"title_style":                  titleStyle,
			"title_page_type":              titlePageType,
			"title_font_size":              f.str("title_font_size"),
			"background_url":               urlBackground,
			"background_color":             backgroundColor,
			"url_banner":                   urlBanner,
			"url_banner_original":          urlBanner,
			"button_subscribe":             button,
			"button_subscribe_original":    button,
			"button_subscribe_color":       f.str("button_subscribe_color"),
			"tnc":                          tnc,
			"tnc_original":                 tnc,
			"tnc_color":                    f.str("tnc_color", operator),
			"tnc_style":                    f.str("tnc_style", operator),
			"tnc_page_type":                f.str("tnc_page_type", operator),
			"tnc_font_size":                f.str("tnc_font_size", operator),
			"enable_ip_ranges":             true,
		}
		if matching != nil {
			fillStoredColumns(entry, matching)
		}
		data = append(data, entry)
	}
	return map[string]interface{}{
		"action":               "ADD_UPDATE",
		"campaign_id":          f.str("campaign_id"),
		"name":                 name,
		"campaign_objective":   "SINGLE URL S2S",
		"country":              country,
		"advertiser":           advertiser,
		"mcc":                  f.str("mcc"),
		"is_bypass_ip_range":   bypass,
		"country_dialing_code": f.str("country_dialing_code"),
		"data":                 data,
	}, nil
}

func (h *IncomingHandler) countryMCC(code string) string {
	var mobileCode []*string
	h.DB.Table("countries").Where("code = ?", code).Limit(1).Pluck("mobile_country_code", &mobileCode)
	if len(mobileCode) == 0 || mobileCode[0] == nil {
		return "000"
	}
	return *mobileCode[0]
}

func (h *IncomingHandler) applyCampaignPayload(payload map[string]interface{}) error {
	encoded, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	var action entity.DataCampaignAction
	if err := json.Unmarshal(encoded, &action); err != nil {
		return err
	}
	return h.applyCampaignUpdate(action)
}

// UpdateCampaignForm is CampaignManagementController::update for the CPA objective.
func (h *IncomingHandler) UpdateCampaignForm(c *fiber.Ctx) error {
	generic := func() error {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "An error occurred while preparing the campaigns for update. Please try again."})
	}
	form, err := parseCampaignForm(c)
	if err != nil {
		return generic()
	}
	campaign, details, err := h.loadCampaign(form.str("campaign_id"), nil)
	if err != nil || campaign == nil {
		return generic()
	}
	now := helper.GetFormatTime(h.Config.TZ, "2006-01-02 15:04:05")
	payload, logAfter, err := buildCampaignUpdatePayload(form, campaign, details, h.countryMCC(campaign.str("country")), now)
	if err != nil {
		h.Logs.Errorf("campaign update: %v", err)
		return generic()
	}
	if len(payload["data"].([]map[string]interface{})) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Campaign data is empty"})
	}
	applyErr := h.applyCampaignPayload(payload)
	h.writeCampaignAudit(c, "updated", campaign.str("campaign_id"), []interface{}{}, logAfter, "("+campaign.str("campaign_id")+") "+campaign.str("name"))
	if applyErr != nil {
		h.Logs.Errorf("campaign update: %v", applyErr)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Something went wrong"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "message": "Campaign has been updated"})
}

// UpdateUploadSmsCampaign is CampaignManagementController::updateUploadSms.
func (h *IncomingHandler) UpdateUploadSmsCampaign(c *fiber.Ctx) error {
	generic := "An error occurred while updating the campaign. Please try again."
	form, err := parseCampaignForm(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": generic})
	}
	campaign, details, err := h.loadCampaign(form.str("campaign_id"), nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": generic})
	}
	if campaign == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Campaign not found."})
	}
	if len(details) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Campaign detail not found."})
	}
	payload := buildUploadSmsUpdatePayload(form, campaign, details[0], h.countryMCC(campaign.str("country")))
	applyErr := h.applyCampaignPayload(payload)
	h.writeCampaignAudit(c, "updated", campaign.str("campaign_id"), fiber.Map{"channel": details[0]["channel"]}, fiber.Map{"channel": form.value("channel")}, "("+campaign.str("campaign_id")+") "+campaign.str("name"))
	if applyErr != nil {
		h.Logs.Errorf("campaign update upload sms: %v", applyErr)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Something went wrong"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "message": "Campaign has been updated"})
}

// UpdateSingleUrlS2S is CampaignManagementController::updateSingleUrlS2S. Laravel ignored the update's
// response and reported success regardless; as on create, the failure reaches the user here.
func (h *IncomingHandler) UpdateSingleUrlS2S(c *fiber.Ctx) error {
	generic := "An error occurred while updating the campaign. Please try again."
	fail := func(err error) error {
		h.Logs.Errorf("campaign update single url: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": generic})
	}
	form, err := parseCampaignForm(c)
	if err != nil {
		return fail(err)
	}
	campaign, details, err := h.loadCampaign(form.str("campaign_id"), nil)
	if err != nil {
		return fail(err)
	}
	if campaign == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Campaign not found."})
	}
	urlBanner, urlBackground := "", "NA"
	if len(details) > 0 {
		urlBanner = details[0].str("url_banner")
		urlBackground = details[0].strOr("background_url", "NA")
	}
	prefix := strings.ToUpper(firstBytes(form.str("campaign_name"), 3) + "-" + firstBytes(form.str("country"), 3) + "-SUS")
	if file, fileErr := c.FormFile("banner"); fileErr == nil {
		path, saveErr := saveCampaignImage(c, file, "campaigns/s2s/banner/original", prefix)
		if saveErr != nil {
			return fail(saveErr)
		}
		urlBanner = h.publicStorageURL(path)
	}
	if file, fileErr := c.FormFile("background"); fileErr == nil {
		path, saveErr := saveCampaignImage(c, file, "campaigns/s2s/background", prefix)
		if saveErr != nil {
			return fail(saveErr)
		}
		urlBackground = h.publicStorageURL(path)
	}
	if len(form.list("operator")) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "No operator selected."})
	}
	payload, err := buildSingleUrlUpdatePayload(form, campaign, details, urlBanner, urlBackground)
	if err != nil {
		return fail(err)
	}
	applyErr := h.applyCampaignPayload(payload)
	h.writeCampaignAudit(c, "updated", "-1", []interface{}{}, payload, form.str("campaign_name"))
	if applyErr != nil {
		return fail(applyErr)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "message": "Campaign has been updated"})
}
