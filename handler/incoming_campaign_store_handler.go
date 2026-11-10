package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource/entity"
)

// campaignForm is the nested request the campaign forms post, read with PHP's loose rules so the
// payload sent to the worker matches what CampaignRepository built.
type campaignForm map[string]interface{}

func (f campaignForm) value(keys ...string) interface{} {
	var current interface{} = map[string]interface{}(f)
	for _, key := range keys {
		level, ok := current.(map[string]interface{})
		if !ok {
			return nil
		}
		if current, ok = level[key]; !ok {
			return nil
		}
	}
	return current
}

// str is PHP's (string) cast on an input; a missing input is an empty string.
func (f campaignForm) str(keys ...string) string {
	switch v := f.value(keys...).(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		if v {
			return "1"
		}
	}
	return ""
}

// list reads an input a foreach walks: nothing is an empty list and a single value a one-item list.
func (f campaignForm) list(keys ...string) []string {
	switch v := f.value(keys...).(type) {
	case nil:
		return nil
	case []interface{}:
		items := make([]string, 0, len(v))
		for _, item := range v {
			switch text := item.(type) {
			case string:
				items = append(items, text)
			case nil:
				items = append(items, "")
			}
		}
		return items
	case string:
		return []string{v}
	}
	return nil
}

// phpInt is PHP's (int) cast on text: the leading number, otherwise 0.
func phpInt(value string) int {
	value = strings.TrimSpace(value)
	end := 0
	for end < len(value) && (value[end] >= '0' && value[end] <= '9' || (end == 0 && (value[end] == '-' || value[end] == '+'))) {
		end++
	}
	number, err := strconv.Atoi(value[:end])
	if err != nil {
		return 0
	}
	return number
}

// phpFilterBool is filter_var(..., FILTER_VALIDATE_BOOLEAN).
func phpFilterBool(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "on", "yes":
		return true
	}
	return false
}

// phpTruthy is PHP's (bool) cast on text.
func phpTruthy(value string) bool {
	return value != "" && value != "0"
}

func phpHTMLSpecialChars(value string) string {
	return strings.NewReplacer("&", "&amp;", "\"", "&quot;", "'", "&#039;", "<", "&lt;", ">", "&gt;").Replace(value)
}

func firstBytes(value string, count int) string {
	if len(value) > count {
		return value[:count]
	}
	return value
}

// campaignServiceKey is the key the store methods check for duplicates; the worker generates the stored one.
func campaignServiceKey(country, operator, partner, service, adnet string) string {
	part := func(value string) string {
		lowered := strings.ToLower(firstBytes(value, 3))
		return strings.Repeat("0", 3-len(lowered)) + lowered
	}
	return strings.ToLower(firstBytes(country, 2)) + part(operator) + part(partner) + part(service) + part(adnet)
}

type campaignCapping struct {
	MOCappingService        int
	CounterMOCappingService int
	StatusCapping           bool
}

// campaignLookups is what the payload builders need from the database.
type campaignLookups struct {
	capping   func(country, operator, partner, service string) campaignCapping
	keyExists func(key string) bool
}

type campaignFormError struct{ message string }

func (e campaignFormError) Error() string { return e.message }

// buildCPAPayload is CampaignRepository::storeCampaignData. Every list is walked in the same nested order, so
// a missing list (no redirect chosen, say) produces no rows, and an empty campaign is refused below.
func buildCPAPayload(f campaignForm, mcc, urlBanner string, look campaignLookups) (map[string]interface{}, error) {
	name, objective, country, advertiser := f.str("campaign_name"), f.str("objective"), f.str("country"), f.str("advertiser")
	partner := f.str("partner")
	details := []map[string]interface{}{}
	for _, operator := range f.list("operator") {
		for _, service := range f.list("service_name", operator) {
			adnets := f.list("adnet", operator, service)
			isEvina := f.str("is_evina", operator, service)
			if isEvina == "" {
				isEvina = "false"
			}
			evinaURL := f.str("evina_redirect_fraud_url", operator, service)
			shortCodes := f.list("short_code", operator, service)
			urlServices := f.list("url_service", operator, service)
			capping := look.capping(country, operator, partner, service)
			for _, adnet := range adnets {
				keyword := ""
				if words := f.list("keyword", operator, service, adnet); len(words) > 0 {
					keyword = words[0]
				}
				for _, shortCode := range shortCodes {
					for _, urlService := range urlServices {
						key := campaignServiceKey(country, operator, partner, service, adnet)
						if look.keyExists(key) {
							return nil, campaignFormError{"Campaign with the same URL service key (" + key + ") already exists."}
						}
						for _, urlLanding := range f.list("url_landing", operator, service, adnet) {
							for _, urlWarp := range f.list("url_warp_landing", operator, service, adnet) {
								for _, apiURL := range f.list("api_url", operator, service, adnet) {
									for _, urlRedirect := range f.list("url_tfc_or_smartlink", operator, service, adnet) {
										for _, redirectType := range f.list("redirect_type", operator, service, adnet) {
											details = append(details, map[string]interface{}{
												"channel":                       strings.ToLower(f.str("channel")),
												"name":                          name,
												"objective":                     objective,
												"country":                       country,
												"mcc":                           mcc,
												"advertiser":                    advertiser,
												"operator":                      operator,
												"partner":                       partner,
												"url_postback":                  f.str("url_postback"),
												"postback_method":               f.str("postback_method"),
												"aggregator":                    f.str("aggregator"),
												"adnet":                         adnet,
												"service":                       service,
												"keyword":                       keyword,
												"subkeyword":                    f.str("subkeyword"),
												"is_billable":                   phpFilterBool(f.str("is_billable")),
												"plan":                          f.str("plan"),
												"po":                            f.str("payout"),
												"currency":                      f.str("currency"),
												"cost":                          "0",
												"pubid":                         "0",
												"short_code":                    shortCode,
												"device_type":                   f.str("device_type"),
												"os":                            f.str("os_type"),
												"url_type":                      f.str("url_type"),
												"click_type":                    phpInt(f.str("click_type")),
												"click_delay":                   phpInt(f.str("click_delay")),
												"client_type":                   strings.ToLower(f.str("client_type")),
												"traffic_source":                phpFilterBool(f.str("traffic_source")),
												"unique_click":                  phpFilterBool(f.str("unique_click")),
												"clickable_anywhere":            phpFilterBool(f.str("clickable_anywhere")),
												"url_banner":                    urlBanner,
												"url_landing":                   urlLanding,
												"url_warp_landing":              urlWarp,
												"url_service":                   urlService,
												"url_tfc_or_smartlink":          urlRedirect,
												"unused_traffic_redirect_type":  redirectType,
												"is_evina":                      phpFilterBool(isEvina),
												"evina_redirect_fraud_url":      evinaURL,
												"glob_post":                     phpFilterBool(f.str("glob_post")),
												"url_glob_post":                 "NA",
												"custom_integration":            f.str("custom_integration"),
												"ip_address":                    []string{},
												"isp":                           "NA",
												"is_active":                     true,
												"mo_capping":                    phpInt(f.str("daily_limit")),
												"mo_capping_service":            capping.MOCappingService,
												"counter_mo_capping_service":    capping.CounterMOCappingService,
												"status_capping":                capping.StatusCapping,
												"kpi_upper_limit_capping":       500,
												"is_machine_learning_capping":   false,
												"ratio_send":                    1,
												"ratio_receive":                 4,
												"counter_mo_ratio":              0,
												"status_ratio":                  true,
												"kpi_upper_limit_ratio_send":    1,
												"kpi_upper_limit_ratio_receive": 2,
												"is_machine_learning_ratio":     false,
												"api_url":                       apiURL,
												"last_update":                   "",
												"last_update_capping":           "",
											})
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
	if len(details) == 0 {
		return nil, campaignFormError{"Campaign data is empty"}
	}
	return map[string]interface{}{"action": "ADD", "name": name, "campaign_objective": objective, "country": country, "advertiser": advertiser, "data": details}, nil
}

// buildUploadSmsPayload is storeUploadSmsCampaignData and storeUploadSmsApiCampaignData: one row of single values.
func buildUploadSmsPayload(f campaignForm, mcc, urlBanner string, look campaignLookups) (map[string]interface{}, error) {
	name, objective, country, advertiser := f.str("campaign_name"), f.str("objective"), f.str("country"), f.str("advertiser")
	key := campaignServiceKey(country, f.str("operator"), f.str("partner"), f.str("service"), f.str("adnet"))
	if look.keyExists(key) {
		return nil, campaignFormError{"Campaign with the same URL service key (" + key + ") already exists."}
	}
	detail := map[string]interface{}{
		"name":                          name,
		"objective":                     objective,
		"country":                       country,
		"mcc":                           mcc,
		"advertiser":                    advertiser,
		"operator":                      f.str("operator"),
		"channel":                       f.str("channel"),
		"partner":                       f.str("partner"),
		"url_postback":                  f.str("url_postback"),
		"postback_method":               f.str("postback_method"),
		"aggregator":                    f.str("aggregator"),
		"adnet":                         f.str("adnet"),
		"service":                       f.str("service"),
		"keyword":                       f.str("keyword"),
		"subkeyword":                    f.str("subkeyword"),
		"is_billable":                   f.value("is_billable") == nil || phpFilterBool(f.str("is_billable")),
		"plan":                          f.str("plan"),
		"po":                            f.str("payout"),
		"currency":                      f.str("currency"),
		"cost":                          "0",
		"pubid":                         "0",
		"short_code":                    f.str("short_code"),
		"device_type":                   f.str("device_type"),
		"os":                            f.str("os_type"),
		"url_type":                      f.str("url_type"),
		"click_type":                    phpInt(f.str("click_type")),
		"click_delay":                   phpInt(f.str("click_delay")),
		"client_type":                   f.str("client_type"),
		"traffic_source":                phpFilterBool(f.str("traffic_source")),
		"unique_click":                  phpFilterBool(f.str("unique_click")),
		"url_banner":                    urlBanner,
		"url_landing":                   f.str("url_landing"),
		"url_warp_landing":              f.str("url_warp_landing"),
		"url_service":                   f.str("url_service"),
		"url_tfc_or_smartlink":          f.str("url_tfc_or_smartlink"),
		"glob_post":                     phpFilterBool(f.str("glob_post")),
		"url_glob_post":                 "NA",
		"custom_integration":            f.str("custom_integration"),
		"ip_address":                    []string{},
		"isp":                           "NA",
		"is_active":                     true,
		"mo_capping":                    500,
		"mo_capping_service":            0,
		"counter_mo_capping":            0,
		"status_capping":                false,
		"kpi_upper_limit_capping":       500,
		"is_machine_learning_capping":   false,
		"ratio_send":                    1,
		"ratio_receive":                 4,
		"counter_mo_ratio":              0,
		"status_ratio":                  true,
		"kpi_upper_limit_ratio_send":    1,
		"kpi_upper_limit_ratio_receive": 2,
		"is_machine_learning_ratio":     false,
		"api_url":                       f.str("api_url"),
		"last_update":                   "",
		"last_update_capping":           "",
	}
	return map[string]interface{}{"action": "ADD", "name": name, "campaign_objective": objective, "country": country, "advertiser": advertiser, "data": []map[string]interface{}{detail}}, nil
}

// buildSingleUrlS2SPayload is CampaignRepository::storeCampaignSingleUrlS2SData.
func buildSingleUrlS2SPayload(f campaignForm, urlBanner, urlBackground string, look campaignLookups) (map[string]interface{}, error) {
	name, country, advertiser := f.str("campaign_name"), f.str("country"), f.str("advertiser")
	bypass := phpTruthy(f.str("is_bypass_ip_range"))
	operators := f.list("operator")
	if len(operators) == 0 {
		return nil, campaignFormError{"No operator selected."}
	}
	background := "NA"
	if urlBackground != "" {
		background = urlBackground
	}
	backgroundColor := "NA"
	if phpTruthy(f.str("enable_color")) && f.str("background_color") != "" {
		backgroundColor = f.str("background_color")
	}
	details := make([]map[string]interface{}, 0, len(operators))
	for _, operator := range operators {
		redirectType := f.str("redirect_type", operator)
		redirectURL := ""
		switch redirectType {
		case "tfc":
			redirectURL = f.str("redirect_internal", operator)
		case "smartlink":
			redirectURL = f.str("redirect_smartlink", operator)
		}
		partner, service := f.str("partner", operator), f.str("service", operator)
		clickDelay := 0
		if f.str("click_delay_radio", operator) == "yes" {
			clickDelay = phpInt(f.str("click_delay", operator))
		}
		orDefault := func(field, fallback string) string {
			if f.value(field, operator) == nil {
				return fallback
			}
			return f.str(field, operator)
		}
		clickType := 1
		if f.value("click_type", operator) != nil {
			clickType = phpInt(f.str("click_type", operator))
		}
		moCapping := 500
		if f.value("daily_limit", operator) != nil {
			moCapping = phpInt(f.str("daily_limit", operator))
		}
		capping := look.capping(country, operator, partner, service)
		title := phpHTMLSpecialChars(f.str("title"))
		button := phpHTMLSpecialChars(f.str("button_subscribe"))
		tnc := phpHTMLSpecialChars(f.str("tnc", operator))
		details = append(details, map[string]interface{}{
			"is_active":                    true,
			"channel":                      strings.ToLower(f.str("channel")),
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
			"mo_capping":                   moCapping,
			"po":                           orDefault("payout", "0"),
			"mo_capping_service":           capping.MOCappingService,
			"counter_mo_capping_service":   capping.CounterMOCappingService,
			"status_capping":               capping.StatusCapping,
			"url_type":                     orDefault("url_type", "wap"),
			"glob_post":                    phpFilterBool(f.str("glob_post", operator)),
			"click_type":                   clickType,
			"click_delay":                  clickDelay,
			"traffic_source":               phpFilterBool(f.str("traffic_source", operator)),
			"unique_click":                 phpFilterBool(f.str("unique_click", operator)),
			"clickable_anywhere":           phpFilterBool(f.str("clickable_anywhere", operator)),
			"is_evina":                     phpFilterBool(f.str("is_evina", operator)),
			"evina_redirect_fraud_url":     f.str("evina_redirect_fraud_url", operator),
			"custom_integration":           orDefault("custom_integration", "default"),
			"unused_traffic_redirect_type": redirectType,
			"url_tfc_or_smartlink":         redirectURL,
			"mainstream_lp_type":           strings.ToLower(f.str("mainstream_lp_type")),
			"title":                        title,
			"title_original":               title,
			"title_color":                  f.str("title_color"),
			"title_style":                  f.str("title_style"),
			"title_page_type":              f.str("title_page_type"),
			"title_font_size":              f.str("title_font_size"),
			"background_url":               background,
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
		})
	}
	return map[string]interface{}{
		"action": "ADD", "name": name, "campaign_objective": "SINGLE URL S2S", "country": country, "advertiser": advertiser,
		"mcc": f.str("mcc"), "is_bypass_ip_range": bypass, "country_dialing_code": f.str("country_dialing_code"), "data": details,
	}, nil
}

// singleUrlErrors applies CampaignManagementController::storeSingleUrlS2S's rules, keyed like Laravel's errors bag.
func singleUrlErrors(f campaignForm) map[string][]string {
	errs := map[string][]string{}
	add := func(field, message string) { errs[field] = append(errs[field], message) }
	required := func(field, message string) {
		if strings.TrimSpace(f.str(field)) == "" {
			if message == "" {
				message = "The " + strings.ReplaceAll(field, "_", " ") + " field is required."
			}
			add(field, message)
		}
	}
	required("campaign_name", "")
	if len([]rune(f.str("campaign_name"))) > 20 {
		add("campaign_name", "The campaign name field must not be greater than 20 characters.")
	}
	required("country", "")
	required("advertiser", "")
	if text := f.str("advertiser"); text != "" && len([]rune(text)) < 3 {
		add("advertiser", "Advertiser must be at least 3 characters.")
	}
	for _, field := range []string{"adnet", "domain", "mainstream_lp_type", "title", "title_color", "title_style", "title_page_type", "button_subscribe", "button_subscribe_color"} {
		required(field, "")
	}
	operators := f.list("operator")
	if len(operators) == 0 {
		add("operator", "At least one operator is required.")
	}
	if _, isMap := f.value("service").(map[string]interface{}); !isMap {
		add("service", "Service is required for each operator.")
	}
	perOperator := []struct{ field, message string }{
		{"partner", "The partner field is required."},
		{"service", "Service is required for each operator."},
		{"shortcode", "Short code is required for each operator."},
		{"url_service", "URL service is required for each operator."},
		{"prefix", "Prefix is required for each operator."},
		{"keyword", "Keyword is required for each operator."},
		{"daily_limit", "Daily limit is required for each operator."},
		{"payout", "Payout is required for each operator."},
	}
	for _, operator := range operators {
		for _, rule := range perOperator {
			if strings.TrimSpace(f.str(rule.field, operator)) == "" {
				add(rule.field+"."+operator, rule.message)
			}
		}
		for field, message := range map[string]string{"daily_limit": "Daily limit must be a number.", "payout": "Payout must be a number."} {
			if text := strings.TrimSpace(f.str(field, operator)); text != "" {
				if number, err := strconv.ParseFloat(text, 64); err != nil {
					add(field+"."+operator, message)
				} else if number < 0 {
					add(field+"."+operator, "The "+strings.ReplaceAll(field, "_", " ")+" field must be at least 0.")
				}
			}
		}
	}
	return errs
}

func (h *IncomingHandler) campaignLookups() campaignLookups {
	return campaignLookups{
		capping: func(country, operator, partner, service string) campaignCapping {
			var found struct {
				MOCappingService        int
				CounterMOCappingService int
				StatusCapping           bool
			}
			h.DB.Table("campaign_details AS d").
				Select("d.mo_capping_service, d.counter_mo_capping_service, d.status_capping").
				Joins("JOIN campaigns c ON c.campaign_id = d.campaign_id").
				Where("d.country = ? AND d.operator = ? AND d.partner = ? AND d.service = ? AND d.is_active = true", country, operator, partner, service).
				Where("c.campaign_objective IN ?", []string{"CPA", "CPI", "CPC", "CPM", "SINGLE URL S2S"}).
				Limit(1).Scan(&found)
			return campaignCapping{found.MOCappingService, found.CounterMOCappingService, found.StatusCapping}
		},
		keyExists: func(key string) bool {
			var count int64
			h.DB.Table("campaign_details").Where("url_service_key = ?", key).Count(&count)
			return count > 0
		},
	}
}

func (h *IncomingHandler) campaignNameTaken(name string) bool {
	var count int64
	h.DB.Table("campaigns").Where("name = ?", name).Count(&count)
	return count > 0
}

// publicStorageURL gives a stored file an absolute address, since landing pages load it from any domain.
// It uses APPAPI because this server serves /storage; APPHOST is the landing page host.
func (h *IncomingHandler) publicStorageURL(path string) string {
	return strings.TrimRight(h.Config.AppApi, "/") + path
}

// saveCampaignImage stores an uploaded image under ./storage. The type comes from the bytes and the
// name from the server, because /storage is public.
func saveCampaignImage(c *fiber.Ctx, file *multipart.FileHeader, folder, prefix string) (string, error) {
	opened, err := file.Open()
	if err != nil {
		return "", err
	}
	head := make([]byte, 512)
	n, _ := opened.Read(head)
	opened.Close()
	// Laravel's image rule also admits svg; it is refused here because a stored svg can carry script.
	detected := http.DetectContentType(head[:n])
	extension := lpBannerExtension(detected)
	switch detected {
	case "image/webp":
		extension = ".webp"
	case "image/bmp":
		extension = ".bmp"
	}
	if extension == "" {
		return "", campaignFormError{"Campaign images must be a jpeg, png, gif, bmp or webp file."}
	}
	dir := filepath.Join(".", "storage", folder)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	name := fmt.Sprintf("%s-%d-%s%s", strings.Map(func(r rune) rune {
		if r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-' {
			return r
		}
		return -1
	}, prefix), time.Now().Unix(), randomHex(4), extension)
	if err := c.SaveFile(file, filepath.Join(dir, name)); err != nil {
		return "", err
	}
	return "/storage/" + folder + "/" + name, nil
}

func parseCampaignForm(c *fiber.Ctx) (campaignForm, error) {
	form := campaignForm{}
	if err := json.Unmarshal([]byte(c.FormValue("payload")), &form); err != nil {
		return form, err
	}
	normalizeCampaignForm(form)
	return form, nil
}

// normalizeCampaignForm does what Laravel's TrimStrings and ConvertEmptyStringsToNull middleware
// did before the repository saw the request, which decides whether a ?? fallback applies.
func normalizeCampaignForm(level interface{}) interface{} {
	switch v := level.(type) {
	case string:
		text := strings.TrimFunc(v, func(r rune) bool {
			return unicode.IsSpace(r) || r == '\uFEFF' || r == '\u200B' || r == '\u200E'
		})
		if text == "" {
			return nil
		}
		return text
	case campaignForm:
		for key, item := range v {
			v[key] = normalizeCampaignForm(item)
		}
	case map[string]interface{}:
		for key, item := range v {
			v[key] = normalizeCampaignForm(item)
		}
	case []interface{}:
		for i, item := range v {
			v[i] = normalizeCampaignForm(item)
		}
	}
	return level
}

// submitCampaign posts the built payload to the worker. Laravel ignored this response and reported
// success regardless, so a failed campaign looked created; here the failure reaches the user.
func (h *IncomingHandler) submitCampaign(c *fiber.Ctx, payload map[string]interface{}, name string) error {
	encoded, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	var action entity.DataCampaignAction
	if err := json.Unmarshal(encoded, &action); err != nil {
		return err
	}
	if _, _, err := h.sendCampaignAction(c.UserContext(), &action); err != nil {
		return err
	}
	h.writeCampaignAudit(c, "created", "-1", []interface{}{}, payload, name)
	return nil
}

// writeCampaignAudit is the audit_log() helper; a failed audit row does not fail the request.
func (h *IncomingHandler) writeCampaignAudit(c *fiber.Ctx, event, auditableID string, before, after interface{}, name string) {
	oldValues, _ := json.Marshal(before)
	newValues, _ := json.Marshal(after)
	userID, _ := c.Locals("user_id").(int)
	audit := entity.Audit{
		UserType:      "App\\Models\\User",
		UserID:        userID,
		Event:         event,
		AuditableType: "App\\Models\\Campaign",
		AuditableID:   auditableID,
		OldValues:     string(oldValues),
		NewValues:     string(newValues),
		URL:           c.OriginalURL(),
		IPAddress:     c.IP(),
		UserAgent:     firstBytes(c.Get("User-Agent"), 1023),
		ActionName:    name,
	}
	if err := h.DB.Create(&audit).Error; err != nil {
		h.Logs.Errorf("campaign audit: %v", err)
	}
}

func (h *IncomingHandler) campaignResult(c *fiber.Ctx, err error, generic string) error {
	var refused campaignFormError
	if errors.As(err, &refused) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": refused.message})
	}
	h.Logs.Errorf("campaign store: %v", err)
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": generic})
}

// StoreCampaign is CampaignManagementController::store for the CPA and Upload SMS objectives.
func (h *IncomingHandler) StoreCampaign(c *fiber.Ctx) error {
	form, err := parseCampaignForm(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid request body"})
	}
	generic := "An error occurred while preparing the campaigns for create. Please try again."
	name := form.str("campaign_name")
	if h.campaignNameTaken(name) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Campaign with the same name (" + name + ") already exists."})
	}
	mcc := "000"
	var mobileCode []string
	h.DB.Table("countries").Where("code = ?", form.str("country")).Limit(1).Pluck("mobile_country_code", &mobileCode)
	if len(mobileCode) > 0 && mobileCode[0] != "" {
		mcc = mobileCode[0]
	}
	urlBanner := ""
	if file, fileErr := c.FormFile("image_upload"); fileErr == nil {
		prefix := strings.ToUpper(firstBytes(name, 3) + "-" + firstBytes(form.str("objective"), 3) + "-" + firstBytes(form.str("country"), 3) + "-" + firstBytes(form.str("advertiser"), 3))
		path, saveErr := saveCampaignImage(c, file, "campaigns/s2s/banner", prefix)
		if saveErr != nil {
			return h.campaignResult(c, saveErr, generic)
		}
		urlBanner = h.publicStorageURL(path)
	}
	var payload map[string]interface{}
	switch form.str("objective") {
	case "UPLOAD SMS", "UPLOAD SMS API":
		payload, err = buildUploadSmsPayload(form, mcc, urlBanner, h.campaignLookups())
	default:
		payload, err = buildCPAPayload(form, mcc, urlBanner, h.campaignLookups())
	}
	if err == nil {
		err = h.submitCampaign(c, payload, name)
	}
	if err != nil {
		return h.campaignResult(c, err, generic)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "message": "Campaign has been created"})
}

// StoreSingleUrlS2S is CampaignManagementController::storeSingleUrlS2S.
func (h *IncomingHandler) StoreSingleUrlS2S(c *fiber.Ctx) error {
	form, err := parseCampaignForm(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid request body"})
	}
	errs := singleUrlErrors(form)
	banner, bannerErr := c.FormFile("banner")
	if bannerErr != nil {
		errs["banner"] = []string{"Banner image is required."}
	}
	if len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": "Input Not Valid", "errors": errs})
	}
	generic := "An error occurred while preparing the campaign. Please try again."
	name := form.str("campaign_name")
	if h.campaignNameTaken(name) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Campaign with the same name (" + name + ") already exists."})
	}
	prefix := strings.ToUpper(firstBytes(name, 3) + "-" + firstBytes(form.str("country"), 3) + "-SUS")
	bannerPath, saveErr := saveCampaignImage(c, banner, "campaigns/s2s/banner/original", prefix)
	if saveErr != nil {
		return h.campaignResult(c, saveErr, generic)
	}
	urlBackground := ""
	if file, fileErr := c.FormFile("background"); fileErr == nil {
		path, saveErr := saveCampaignImage(c, file, "campaigns/s2s/background", prefix)
		if saveErr != nil {
			return h.campaignResult(c, saveErr, generic)
		}
		urlBackground = h.publicStorageURL(path)
	}
	payload, err := buildSingleUrlS2SPayload(form, h.publicStorageURL(bannerPath), urlBackground, h.campaignLookups())
	if err == nil {
		err = h.submitCampaign(c, payload, name)
	}
	if err != nil {
		return h.campaignResult(c, err, generic)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "message": "Campaign has been created"})
}
