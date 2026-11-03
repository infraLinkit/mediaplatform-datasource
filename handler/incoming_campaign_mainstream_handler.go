package handler

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// laravelValidator collects errors the way Laravel's validator reports them, with its default messages.
type laravelValidator struct {
	messages []string
	errors   map[string][]string
}

func newLaravelValidator() *laravelValidator {
	return &laravelValidator{errors: map[string][]string{}}
}

func (v *laravelValidator) add(field, message string) {
	v.messages = append(v.messages, message)
	v.errors[field] = append(v.errors[field], message)
}

func validatorLabel(field string) string { return strings.ReplaceAll(field, "_", " ") }

func (v *laravelValidator) failed(c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": laravelValidationMessage(v.messages), "errors": v.errors})
}

func phpIsNumeric(text string) bool {
	_, err := strconv.ParseFloat(strings.TrimSpace(text), 64)
	return err == nil && text != ""
}

// campaignUpload is what the image and max rules need to know about an uploaded file.
type campaignUpload struct {
	present   bool
	image     bool
	kilobytes int64
}

func inspectCampaignUpload(c *fiber.Ctx, field string) campaignUpload {
	file, err := c.FormFile(field)
	if err != nil {
		return campaignUpload{}
	}
	return campaignUpload{present: true, image: isLaravelImage(file), kilobytes: file.Size / 1024}
}

func isLaravelImage(file *multipart.FileHeader) bool {
	if strings.EqualFold(filepath.Ext(file.Filename), ".svg") {
		return true
	}
	opened, err := file.Open()
	if err != nil {
		return false
	}
	defer opened.Close()
	head := make([]byte, 512)
	n, _ := opened.Read(head)
	switch http.DetectContentType(head[:n]) {
	case "image/jpeg", "image/png", "image/gif", "image/bmp", "image/webp":
		return true
	}
	return false
}

func (v *laravelValidator) imageRules(field string, upload campaignUpload) {
	if !upload.image {
		v.add(field, "The "+validatorLabel(field)+" field must be an image.")
	}
	if upload.kilobytes > 512000 {
		v.add(field, "The "+validatorLabel(field)+" field must not be greater than 512000 kilobytes.")
	}
}

func (v *laravelValidator) requiredImage(field string, upload campaignUpload) {
	if !upload.present {
		v.add(field, "The "+validatorLabel(field)+" field is required.")
		return
	}
	v.imageRules(field, upload)
}

func (v *laravelValidator) backgroundRules(f campaignForm, background campaignUpload, message string) {
	if f.value("background_color") == nil && !background.present {
		v.add("background_color", message)
	}
	if background.present {
		v.imageRules("background", background)
	} else if f.value("background_color") == nil {
		v.add("background", message)
	}
}

// mainstreamErrors applies the rules of CampaignMainstreamController::store.
func mainstreamErrors(f campaignForm, uploads map[string]campaignUpload) *laravelValidator {
	v := newLaravelValidator()
	required := func(fields ...string) {
		for _, field := range fields {
			if f.value(field) == nil {
				v.add(field, "The "+validatorLabel(field)+" field is required.")
			}
		}
	}
	maxText := func(field string, limit int) {
		if len([]rune(f.str(field))) > limit {
			v.add(field, fmt.Sprintf("The %s field must not be greater than %d characters.", validatorLabel(field), limit))
		}
	}
	numeric := func(field string) {
		if f.value(field) == nil {
			v.add(field, "The "+validatorLabel(field)+" field is required.")
		} else if !phpIsNumeric(f.str(field)) {
			v.add(field, "The "+validatorLabel(field)+" field must be a number.")
		}
	}
	required("campaign_name")
	maxText("campaign_name", 20)
	if f.value("service") == nil {
		v.add("service", "Mainstream group required")
	}
	required("country", "operator", "partner", "url_service")
	numeric("daily_limit")
	numeric("payout")
	required("currency")
	maxText("currency", 3)
	required("google_sheet", "mainstream_lp_type", "title", "title_color", "title_style", "title_page_type", "title_original")
	v.backgroundRules(f, uploads["background"], "Background image or background color is required.")
	v.requiredImage("logo", uploads["logo"])
	v.requiredImage("banner", uploads["banner"])
	v.requiredImage("banner_original", uploads["banner_original"])
	required("tnc", "tnc_color", "tnc_style", "tnc_page_type", "tnc_original", "button_subscribe", "button_subscribe_color", "button_subscribe_original")
	return v
}

// mainstreamSingleUrlErrors applies the rules of CampaignMainstreamController::storeSingleUrl. A wildcard
// rule reports the raw key, as Laravel does for attributes it expanded from a *.
func mainstreamSingleUrlErrors(f campaignForm, uploads map[string]campaignUpload) *laravelValidator {
	v := newLaravelValidator()
	required := func(fields ...string) {
		for _, field := range fields {
			if f.value(field) == nil {
				v.add(field, "The "+validatorLabel(field)+" field is required.")
			}
		}
	}
	perKey := func(field string, check func(key, text string)) {
		level, ok := f.value(field).(map[string]interface{})
		if !ok {
			if list, isList := f.value(field).([]interface{}); isList {
				level = map[string]interface{}{}
				for i, item := range list {
					level[strconv.Itoa(i)] = item
				}
			} else {
				if f.value(field) != nil {
					v.add(field, "The "+validatorLabel(field)+" field must be an array.")
				}
				return
			}
		}
		keys := make([]string, 0, len(level))
		for key := range level {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			text, _ := level[key].(string)
			check(field+"."+key, text)
		}
	}
	requiredList := func(fields ...string) {
		for _, field := range fields {
			if f.value(field) == nil {
				v.add(field, "The "+validatorLabel(field)+" field is required.")
				continue
			}
			perKey(field, func(key, text string) {
				if text == "" {
					v.add(key, "The "+key+" field is required.")
				}
			})
		}
	}
	required("campaign_name")
	if len([]rune(f.str("campaign_name"))) > 20 {
		v.add("campaign_name", "The campaign name field must not be greater than 20 characters.")
	}
	required("country")
	switch f.value("operator").(type) {
	case nil:
		v.add("operator", "At least one operator is required")
	case []interface{}, map[string]interface{}:
		perKey("operator", func(key, text string) {
			if text == "" {
				v.add(key, "The "+key+" field is required.")
			}
		})
	default:
		v.add("operator", "Operators must be an array")
	}
	if raw := f.str("non_target_url"); raw != "" {
		if parsed, err := url.Parse(raw); err != nil || parsed.Scheme == "" || parsed.Host == "" {
			v.add("non_target_url", "The non target url field must be a valid URL.")
		}
	}
	required("mainstream_lp_type", "title", "title_color", "title_style", "title_page_type")
	background := uploads["background"]
	v.backgroundRules(f, background, "Background image or background color is required.")
	v.requiredImage("banner", uploads["banner"])
	required("button_subscribe", "button_subscribe_color")
	if f.value("campaign_detail_name") == nil {
		v.add("campaign_detail_name", "The campaign detail name field is required.")
	} else {
		perKey("campaign_detail_name", func(key, text string) {
			if text == "" {
				v.add(key, "The "+key+" field is required.")
			} else if len([]rune(text)) > 20 {
				v.add(key, "The "+key+" field must not be greater than 20 characters.")
			}
		})
	}
	requiredList("partner", "channel", "short_code", "url_service", "prefix")
	required("advertiser")
	requiredList("conversion_name", "google_sheet", "currency", "keyword")
	for _, field := range []string{"daily_limit", "payout"} {
		if f.value(field) == nil {
			v.add(field, "The "+validatorLabel(field)+" field is required.")
			continue
		}
		perKey(field, func(key, text string) {
			if text == "" {
				v.add(key, "The "+key+" field is required.")
			} else if !phpIsNumeric(text) {
				v.add(key, "The "+key+" field must be a number.")
			}
		})
	}
	requiredList("tnc", "tnc_color", "tnc_style", "tnc_page_type")
	return v
}

// buildMainstreamPayload is CampaignRepository::storeCampaignMainstreamData.
func buildMainstreamPayload(f campaignForm, mcc string, urls map[string]string) map[string]interface{} {
	name, country, advertiser := f.str("campaign_name"), f.str("country"), f.str("advertiser")
	return map[string]interface{}{
		"action": "ADD", "name": name, "campaign_objective": "MAINSTREAM", "country": country, "advertiser": advertiser,
		"data": []map[string]interface{}{{
			"is_active":                 true,
			"name":                      name,
			"country":                   country,
			"mcc":                       mcc,
			"advertiser":                advertiser,
			"channel":                   strings.ToLower(f.str("channel")),
			"adnet":                     f.str("agency"),
			"url_warp_landing":          f.str("domain"),
			"operator":                  f.str("operator"),
			"partner":                   f.str("partner"),
			"url_postback":              f.str("url_postback"),
			"postback_method":           f.str("postback_method"),
			"client_type":               strings.ToLower(f.str("client_type")),
			"currency":                  f.str("currency"),
			"google_sheet":              f.str("google_sheet"),
			"google_sheet_billable":     f.str("google_sheet_billable"),
			"service":                   f.str("service"),
			"url_service":               f.str("url_service"),
			"non_target_url":            f.str("non_target_url"),
			"show_prelander_non_target": phpFilterBool(f.str("show_prelander_non_target")),
			"mo_capping":                phpInt(f.str("daily_limit")),
			"po":                        f.str("payout"),
			"url_type":                  f.str("url_type"),
			"glob_post":                 phpFilterBool(f.str("glob_post")),
			"click_type":                phpInt(f.str("click_type")),
			"click_delay":               phpInt(f.str("click_delay")),
			"traffic_source":            phpFilterBool(f.str("traffic_source")),
			"unique_click":              phpFilterBool(f.str("unique_click")),
			"clickable_anywhere":        phpFilterBool(f.str("clickable_anywhere")),
			"is_evina":                  phpFilterBool(f.str("is_evina")),
			"evina_redirect_fraud_url":  f.str("evina_redirect_fraud_url"),
			"custom_integration":        f.str("custom_integration"),
			"mainstream_lp_type":        strings.ToLower(f.str("mainstream_lp_type")),
			"title":                     phpHTMLSpecialChars(f.str("title")),
			"title_original":            phpHTMLSpecialChars(f.str("title_original")),
			"title_color":               f.str("title_color"),
			"title_style":               f.str("title_style"),
			"title_page_type":           f.str("title_page_type"),
			"title_font_size":           f.str("title_font_size"),
			"sub_title":                 phpHTMLSpecialChars(f.str("sub_title")),
			"sub_title_original":        phpHTMLSpecialChars(f.str("sub_title_original")),
			"sub_title_color":           f.str("sub_title_color"),
			"sub_title_style":           f.str("sub_title_style"),
			"sub_title_page_type":       f.str("sub_title_page_type"),
			"sub_title_font_size":       f.str("sub_title_font_size"),
			"background_url":            urls["background"],
			"background_color":          f.str("background_color"),
			"logo_url":                  urls["logo"],
			"url_banner":                urls["banner"],
			"url_banner_original":       urls["banner_original"],
			"tnc":                       phpHTMLSpecialChars(f.str("tnc")),
			"tnc_original":              phpHTMLSpecialChars(f.str("tnc_original")),
			"tnc_color":                 f.str("tnc_color"),
			"tnc_style":                 f.str("tnc_style"),
			"tnc_page_type":             f.str("tnc_page_type"),
			"tnc_font_size":             f.str("tnc_font_size"),
			"button_subscribe":          phpHTMLSpecialChars(f.str("button_subscribe")),
			"button_subscribe_original": phpHTMLSpecialChars(f.str("button_subscribe_original")),
			"button_subscribe_color":    f.str("button_subscribe_color"),
			"enable_ip_ranges":          phpTruthy(f.str("enable_ip_ranges")),
			"conversion_name":           f.str("conversion_name"),
			"domain_service":            f.str("domain_service"),
			"api_url":                   f.str("api_url"),
			"service_price":             phpInt(f.str("service_price")),
			"portal_url":                f.str("portal_url"),
			"company_legal_name":        f.str("company_legal_name"),
			"company_address":           f.str("company_address"),
			"company_phone":             f.str("company_phone"),
			"company_email":             f.str("company_email"),
		}},
	}
}

// buildUploadSmsMainstreamPayload is CampaignRepository::storeUploadSmsMainstreamCampaignData. Unlike the
// other objectives it keeps channel and client_type in the case the form sent.
func buildUploadSmsMainstreamPayload(f campaignForm, mcc string) map[string]interface{} {
	name, country, advertiser := f.str("campaign_name"), f.str("country"), f.str("advertiser")
	return map[string]interface{}{
		"action": "ADD", "name": name, "campaign_objective": "UPLOAD SMS MAINSTREAM", "country": country, "advertiser": advertiser,
		"data": []map[string]interface{}{{
			"name":                          name,
			"objective":                     "UPLOAD SMS MAINSTREAM",
			"country":                       country,
			"mcc":                           mcc,
			"advertiser":                    advertiser,
			"operator":                      f.str("operator"),
			"channel":                       f.str("channel"),
			"partner":                       f.str("partner"),
			"url_postback":                  f.str("url_postback"),
			"postback_method":               f.str("postback_method"),
			"adnet":                         f.str("agency"),
			"service":                       f.str("service"),
			"keyword":                       f.str("keyword"),
			"is_billable":                   true,
			"plan":                          "daily",
			"po":                            f.str("payout"),
			"currency":                      "USD",
			"cost":                          "0",
			"pubid":                         "0",
			"short_code":                    f.str("short_code"),
			"device_type":                   "NA",
			"os":                            "NA",
			"url_type":                      "wap",
			"click_type":                    1,
			"click_delay":                   0,
			"client_type":                   f.str("client_type"),
			"traffic_source":                false,
			"unique_click":                  false,
			"url_glob_post":                 "NA",
			"ip_address":                    []interface{}{},
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
			"last_update":                   "",
			"last_update_capping":           "",
		}},
	}
}

// phpFallback is CampaignRepository::fallback: PHP's empty() treats "0" as missing too.
func phpFallback(primary, secondary string) string {
	if phpTruthy(primary) {
		return primary
	}
	if phpTruthy(secondary) {
		return secondary
	}
	return ""
}

// buildMainstreamUpdatePayload is CampaignRepository::updateCampaignMainstreamData.
func buildMainstreamUpdatePayload(f campaignForm, campaign, detail campaignRow, mcc string, urls map[string]string) map[string]interface{} {
	name, country, advertiser := campaign.str("name"), campaign.str("country"), f.str("advertiser")
	pick := func(field string) string { return phpFallback(f.str(field), detail.str(field)) }
	// Laravel escaped these on create but not here, so an edit could store raw HTML.
	pickText := func(field string) string { return escapeKeepingEntities(pick(field)) }
	showPrelander := detail.flag("show_prelander_non_target")
	if _, sent := f["show_prelander_non_target"]; sent {
		showPrelander = phpFilterBool(f.str("show_prelander_non_target"))
	}
	lpType := f.str("mainstream_lp_type")
	if f.value("mainstream_lp_type") == nil {
		lpType = detail.str("mainstream_lp_type")
	}
	entry := map[string]interface{}{
		"campaign_id":               f.str("campaign_id"),
		"urlservicekey":             detail.str("url_service_key"),
		"is_active":                 detail.flag("is_active"),
		"name":                      name,
		"country":                   country,
		"mcc":                       mcc,
		"advertiser":                advertiser,
		"channel":                   strings.ToLower(detail.str("channel")),
		"adnet":                     detail.str("adnet"),
		"url_warp_landing":          detail.str("url_warp_landing"),
		"operator":                  detail.str("operator"),
		"partner":                   detail.str("partner"),
		"url_postback":              f.str("url_postback"),
		"postback_method":           f.str("postback_method"),
		"client_type":               strings.ToLower(f.str("client_type")),
		"currency":                  f.str("currency"),
		"google_sheet":              f.str("google_sheet"),
		"google_sheet_billable":     f.str("google_sheet_billable"),
		"service":                   detail.str("service"),
		"url_service":               f.str("url_service"),
		"non_target_url":            f.str("non_target_url"),
		"show_prelander_non_target": showPrelander,
		"mo_capping":                detail.num("mo_capping"),
		"po":                        detail.str("po"),
		"url_type":                  f.str("url_type"),
		"glob_post":                 f.flagOr(detail, "glob_post", "glob_post"),
		"click_type":                phpInt(f.str("click_type")),
		"click_delay":               phpInt(f.str("click_delay")),
		"traffic_source":            f.flagOr(detail, "traffic_source", "traffic_source"),
		"unique_click":              f.flagOr(detail, "unique_click", "unique_click"),
		"clickable_anywhere":        f.flagOr(detail, "clickable_anywhere", "clickable_anywhere"),
		"is_evina":                  f.flagOr(detail, "is_evina", "is_evina"),
		"evina_redirect_fraud_url":  f.str("evina_redirect_fraud_url"),
		"custom_integration":        f.str("custom_integration"),
		"mainstream_lp_type":        strings.ToLower(lpType),
		"title":                     pickText("title"),
		"title_original":            pickText("title_original"),
		"title_color":               pick("title_color"),
		"title_style":               pick("title_style"),
		"title_page_type":           pick("title_page_type"),
		"title_font_size":           pick("title_font_size"),
		"sub_title":                 pickText("sub_title"),
		"sub_title_original":        pickText("sub_title_original"),
		"sub_title_color":           pick("sub_title_color"),
		"sub_title_style":           pick("sub_title_style"),
		"sub_title_page_type":       pick("sub_title_page_type"),
		"sub_title_font_size":       pick("sub_title_font_size"),
		"background_url":            phpFallback(urls["background"], detail.str("background_url")),
		"background_color":          pick("background_color"),
		"logo_url":                  phpFallback(urls["logo"], detail.str("logo_url")),
		"url_banner":                phpFallback(urls["banner"], detail.str("url_banner")),
		"url_banner_original":       phpFallback(urls["banner_original"], detail.str("url_banner_original")),
		"tnc":                       pickText("tnc"),
		"tnc_original":              pickText("tnc_original"),
		"tnc_color":                 pick("tnc_color"),
		"tnc_style":                 pick("tnc_style"),
		"tnc_page_type":             pick("tnc_page_type"),
		"tnc_font_size":             pick("tnc_font_size"),
		"button_subscribe":          pickText("button_subscribe"),
		"button_subscribe_original": pickText("button_subscribe_original"),
		"button_subscribe_color":    pick("button_subscribe_color"),
		"enable_ip_ranges":          phpFilterBool(f.str("enable_ip_ranges")),
		"conversion_name":           pick("conversion_name"),
		"domain_service":            pick("domain_service"),
		"api_url":                   pick("api_url"),
		"service_price":             phpInt(pick("service_price")),
		"portal_url":                pick("portal_url"),
		"company_legal_name":        pick("company_legal_name"),
		"company_address":           pick("company_address"),
		"company_phone":             pick("company_phone"),
		"company_email":             pick("company_email"),
	}
	fillStoredColumns(entry, detail)
	return map[string]interface{}{
		"action": "ADD_UPDATE", "campaign_id": f.str("campaign_id"), "name": name, "campaign_objective": "MAINSTREAM", "country": country, "advertiser": advertiser,
		"data": []map[string]interface{}{entry},
	}
}

func oneOf(value string, allowed ...string) bool {
	for _, item := range allowed {
		if value == item {
			return true
		}
	}
	return false
}

// buildMainstreamSingleUrlPayload is CampaignRepository::storeCampaignMainstreamSingleUrlData.
func buildMainstreamSingleUrlPayload(f campaignForm, urlBanner, urlBackground string) map[string]interface{} {
	name, country, advertiser, mcc, dialing := f.str("campaign_name"), f.str("country"), f.str("advertiser"), f.str("mcc"), f.str("country_dialing_code")
	bypass := phpTruthy(f.str("is_bypass_ip_range"))
	urlType := f.str("url_type")
	if !oneOf(urlType, "wap", "click2sms") {
		urlType = "wap"
	}
	clickType := 1
	if phpIsNumeric(f.str("click_type")) && phpInt(f.str("click_type")) == 2 {
		clickType = 2
	}
	integration := f.str("custom_integration")
	if !oneOf(integration, "click2sms", "telesat") {
		integration = "default"
	}
	data := []map[string]interface{}{}
	for _, operator := range f.list("operator") {
		capping := 500
		if f.value("daily_limit", operator) != nil {
			capping = phpInt(f.str("daily_limit", operator))
		}
		data = append(data, map[string]interface{}{
			"is_active":                 true,
			"name":                      name,
			"campaign_detail_name":      f.strOr(name, "campaign_detail_name", operator),
			"country":                   country,
			"country_dialing_code":      dialing,
			"mcc":                       mcc,
			"is_bypass_ip_range":        bypass,
			"advertiser":                advertiser,
			"company_legal_name":        f.str("company_legal_name"),
			"company_address":           f.str("company_address"),
			"company_phone":             f.str("company_phone"),
			"company_email":             f.str("company_email"),
			"channel":                   strings.ToLower(f.str("channel", operator)),
			"adnet":                     f.str("agency", operator),
			"url_warp_landing":          f.str("unique_domain"),
			"operator":                  operator,
			"partner":                   f.str("partner", operator),
			"url_postback":              f.str("url_postback", operator),
			"postback_method":           f.str("postback_method", operator),
			"client_type":               "",
			"currency":                  f.strOr("USD", "currency", operator),
			"is_evina":                  phpTruthy(f.str("is_evina", operator)),
			"evina_redirect_fraud_url":  f.str("evina_redirect_fraud_url", operator),
			"google_sheet":              f.strOr("NA", "google_sheet", operator),
			"google_sheet_billable":     f.strOr("NA", "google_sheet_billable", operator),
			"service":                   f.str("service", operator),
			"url_service":               f.str("url_service", operator),
			"non_target_url":            f.str("non_target_url"),
			"show_prelander_non_target": phpFilterBool(f.str("show_prelander_non_target")),
			"mo_capping":                capping,
			"po":                        f.strOr("0", "payout", operator),
			"url_type":                  urlType,
			"glob_post":                 phpFilterBool(f.strOr("false", "glob_post")),
			"click_type":                clickType,
			"click_delay":               phpInt(f.str("click_delay")),
			"traffic_source":            phpFilterBool(f.strOr("false", "traffic_source")),
			"unique_click":              phpFilterBool(f.strOr("false", "unique_click")),
			"clickable_anywhere":        phpFilterBool(f.strOr("false", "clickable_anywhere")),
			"custom_integration":        integration,
			"mainstream_lp_type":        strings.ToLower(f.str("mainstream_lp_type")),
			"title":                     phpHTMLSpecialChars(f.str("title")),
			"title_original":            phpHTMLSpecialChars(f.str("title")),
			"title_color":               f.str("title_color"),
			"title_style":               f.str("title_style"),
			"title_page_type":           f.str("title_page_type"),
			"title_font_size":           f.str("title_font_size"),
			"background_url":            urlBackground,
			"background_color":          f.str("background_color"),
			"url_banner":                urlBanner,
			"url_banner_original":       urlBanner,
			"tnc":                       phpHTMLSpecialChars(f.str("tnc", operator)),
			"tnc_original":              phpHTMLSpecialChars(f.str("tnc", operator)),
			"tnc_color":                 f.str("tnc_color", operator),
			"tnc_style":                 f.str("tnc_style", operator),
			"tnc_page_type":             f.str("tnc_page_type", operator),
			"tnc_font_size":             f.str("tnc_font_size", operator),
			"button_subscribe":          phpHTMLSpecialChars(f.str("button_subscribe")),
			"button_subscribe_original": phpHTMLSpecialChars(f.str("button_subscribe")),
			"button_subscribe_color":    f.str("button_subscribe_color"),
			"enable_ip_ranges":          true,
			"conversion_name":           f.str("conversion_name", operator),
			"domain_service":            f.str("domain_service", operator),
			"api_url":                   f.str("api_url", operator),
			"service_price":             phpInt(f.str("service_price", operator)),
			"portal_url":                f.str("portal_url", operator),
			"short_code":                f.str("short_code", operator),
			"prefix":                    f.str("prefix", operator),
			"keyword":                   f.str("keyword", operator),
		})
	}
	return map[string]interface{}{
		"action": "ADD", "name": name, "campaign_objective": "SINGLE URL MAINSTREAM", "country": country, "advertiser": advertiser,
		"mcc": mcc, "is_bypass_ip_range": bypass, "country_dialing_code": dialing, "data": data,
	}
}

// buildMainstreamSingleUrlUpdatePayload is CampaignRepository::updateCampaignMainstreamSingleUrlData.
func buildMainstreamSingleUrlUpdatePayload(f campaignForm, campaign campaignRow, details []campaignRow, urlBanner, urlBackground string) (map[string]interface{}, error) {
	if _, isList := f.value("operator").([]interface{}); !isList {
		return nil, errors.New("operator is not an array")
	}
	if len(details) == 0 {
		return nil, errors.New("campaign has no details")
	}
	name, country := campaign.str("name"), campaign.str("country")
	bypass := campaign.flag("is_bypass_ip_range")
	if f.value("is_bypass_ip_range") != nil {
		bypass = phpTruthy(f.str("is_bypass_ip_range"))
	}
	data := []map[string]interface{}{}
	for _, operator := range f.list("operator") {
		detail := details[0]
		for _, row := range details {
			if row.str("operator") == operator {
				detail = row
				break
			}
		}
		perOperator := func(field, column string) string {
			return f.strOr(detail.str(column), field, operator)
		}
		shared := func(field string) string { return f.strOr(detail.str(field), field) }
		booleanOr := func(field string) bool {
			if text := f.str(field); oneOf(text, "true", "false") {
				return text == "true"
			}
			return detail.flag(field)
		}
		showPrelander := detail.flag("show_prelander_non_target")
		if _, sent := f["show_prelander_non_target"]; sent {
			showPrelander = phpFilterBool(f.str("show_prelander_non_target"))
		}
		urlType := f.str("url_type")
		if !oneOf(urlType, "wap", "click2sms") {
			urlType = detail.str("url_type")
			if !phpTruthy(urlType) {
				urlType = "wap"
			}
		}
		clickType := phpInt(f.str("click_type"))
		if clickType != 1 && clickType != 2 {
			clickType = detail.num("click_type")
			if clickType == 0 {
				clickType = 1
			}
		}
		clickDelay := detail.num("click_delay")
		if _, sent := f["click_delay"]; sent {
			clickDelay = phpInt(f.str("click_delay"))
		}
		integration := f.str("custom_integration")
		if !oneOf(integration, "default", "click2sms", "telesat") {
			integration = detail.str("custom_integration")
			if !phpTruthy(integration) {
				integration = "default"
			}
		}
		isEvina := detail.flag("is_evina")
		if f.value("is_evina", operator) != nil {
			isEvina = phpTruthy(f.str("is_evina", operator))
		}
		capping := detail.num("mo_capping")
		if f.value("daily_limit", operator) != nil {
			capping = phpInt(f.str("daily_limit", operator))
		}
		titleOriginal := detail.str("title_original")
		if f.value("title") != nil {
			titleOriginal = f.str("title")
		}
		tncOriginal := perOperator("tnc", "tnc_original")
		buttonOriginal := detail.str("button_subscribe_original")
		if f.value("button_subscribe") != nil {
			buttonOriginal = f.str("button_subscribe")
		}
		entry := map[string]interface{}{
			"campaign_id":               f.str("campaign_id"),
			"urlservicekey":             detail.str("url_service_key"),
			"is_active":                 detail.flag("is_active"),
			"name":                      name,
			"campaign_detail_name":      perOperator("campaign_detail_name", "campaign_detail_name"),
			"country":                   country,
			"country_dialing_code":      f.str("country_dialing_code"),
			"mcc":                       shared("mcc"),
			"is_bypass_ip_range":        bypass,
			"advertiser":                shared("advertiser"),
			"company_legal_name":        shared("company_legal_name"),
			"company_address":           shared("company_address"),
			"company_phone":             shared("company_phone"),
			"company_email":             shared("company_email"),
			"channel":                   strings.ToLower(perOperator("channel", "channel")),
			"adnet":                     perOperator("agency", "adnet"),
			"url_warp_landing":          detail.str("url_warp_landing"),
			"operator":                  operator,
			"partner":                   perOperator("partner", "partner"),
			"url_postback":              perOperator("url_postback", "url_postback"),
			"postback_method":           perOperator("postback_method", "postback_method"),
			"currency":                  perOperator("currency", "currency"),
			"is_evina":                  isEvina,
			"evina_redirect_fraud_url":  perOperator("evina_redirect_fraud_url", "evina_redirect_fraud_url"),
			"google_sheet":              perOperator("google_sheet", "google_sheet"),
			"google_sheet_billable":     perOperator("google_sheet_billable", "google_sheet_billable"),
			"service":                   perOperator("service", "service"),
			"url_service":               perOperator("url_service", "url_service"),
			"non_target_url":            shared("non_target_url"),
			"show_prelander_non_target": showPrelander,
			"mo_capping":                capping,
			"po":                        perOperator("payout", "po"),
			"url_type":                  urlType,
			"glob_post":                 booleanOr("glob_post"),
			"click_type":                clickType,
			"click_delay":               clickDelay,
			"traffic_source":            booleanOr("traffic_source"),
			"unique_click":              booleanOr("unique_click"),
			"clickable_anywhere":        booleanOr("clickable_anywhere"),
			"custom_integration":        integration,
			"mainstream_lp_type":        strings.ToLower(shared("mainstream_lp_type")),
			"title":                     escapeKeepingEntities(shared("title")),
			"title_original":            escapeKeepingEntities(titleOriginal),
			"title_color":               shared("title_color"),
			"title_style":               shared("title_style"),
			"title_page_type":           shared("title_page_type"),
			"title_font_size":           shared("title_font_size"),
			"background_url":            phpFallback(urlBackground, detail.str("background_url")),
			"background_color":          phpFallback(f.str("background_color"), detail.str("background_color")),
			"url_banner":                phpFallback(urlBanner, detail.str("url_banner")),
			"url_banner_original":       phpFallback(urlBanner, detail.str("url_banner_original")),
			"tnc":                       escapeKeepingEntities(perOperator("tnc", "tnc")),
			"tnc_original":              escapeKeepingEntities(tncOriginal),
			"tnc_color":                 perOperator("tnc_color", "tnc_color"),
			"tnc_style":                 perOperator("tnc_style", "tnc_style"),
			"tnc_page_type":             perOperator("tnc_page_type", "tnc_page_type"),
			"tnc_font_size":             perOperator("tnc_font_size", "tnc_font_size"),
			"button_subscribe":          escapeKeepingEntities(shared("button_subscribe")),
			"button_subscribe_original": escapeKeepingEntities(buttonOriginal),
			"button_subscribe_color":    shared("button_subscribe_color"),
			"conversion_name":           perOperator("conversion_name", "conversion_name"),
			"domain_service":            perOperator("domain_service", "domain_service"),
			"api_url":                   perOperator("api_url", "api_url"),
			"service_price":             phpInt(perOperator("service_price", "service_price")),
			"portal_url":                perOperator("portal_url", "portal_url"),
			"short_code":                perOperator("short_code", "short_code"),
			"prefix":                    perOperator("prefix", "prefix"),
			"keyword":                   perOperator("keyword", "keyword"),
			"enable_ip_ranges":          true,
		}
		fillStoredColumns(entry, detail)
		data = append(data, entry)
	}
	return map[string]interface{}{
		"action": "ADD_UPDATE", "campaign_id": f.str("campaign_id"), "name": name, "campaign_objective": "SINGLE URL MAINSTREAM", "country": country,
		"advertiser": f.str("advertiser"), "mcc": f.str("mcc"), "is_bypass_ip_range": bypass, "country_dialing_code": f.str("country_dialing_code"), "data": data,
	}, nil
}

// saveMainstreamImages stores the files present and returns their public urls by field.
func (h *IncomingHandler) saveMainstreamImages(c *fiber.Ctx, prefix string, folders map[string]string) (map[string]string, error) {
	urls := map[string]string{}
	for field, folder := range folders {
		file, err := c.FormFile(field)
		if err != nil {
			continue
		}
		path, err := saveCampaignImage(c, file, folder, prefix)
		if err != nil {
			return nil, err
		}
		urls[field] = h.publicStorageURL(path)
	}
	return urls, nil
}

var mainstreamFolders = map[string]string{
	"banner_original": "campaigns/mainstream/banner/original",
	"banner":          "campaigns/mainstream/banner/en",
	"logo":            "campaigns/mainstream/logo",
	"background":      "campaigns/mainstream/background",
}

var mainstreamSingleUrlFolders = map[string]string{
	"banner":     "campaigns/mainstream/banner/original",
	"background": "campaigns/mainstream/background",
}

func campaignImagePrefix(name, country, channel string) string {
	return strings.ToUpper(firstBytes(name, 3) + "-" + firstBytes(country, 3) + "-" + firstBytes(channel, 3))
}

func (h *IncomingHandler) nameTakenResponse(c *fiber.Ctx, name string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Campaign with the same name (" + name + ") already exists."})
}

// StoreCampaignMainstream is CampaignMainstreamController::store. Laravel ignored the worker's response and
// reported success regardless; as with the S2S forms, the failure reaches the user here.
func (h *IncomingHandler) StoreCampaignMainstream(c *fiber.Ctx) error {
	form, err := parseCampaignForm(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid request body"})
	}
	name := form.str("campaign_name")
	if form.str("objective") == "UPLOAD SMS MAINSTREAM" {
		generic := "An error occurred while preparing the campaign for create. Please try again."
		if h.campaignNameTaken(name) {
			return h.nameTakenResponse(c, name)
		}
		key := campaignServiceKey(form.str("country"), form.str("operator"), form.str("partner"), form.str("service"), form.str("agency"))
		if h.campaignLookups().keyExists(key) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Campaign with the same URL service key (" + key + ") already exists."})
		}
		if err := h.submitCampaign(c, buildUploadSmsMainstreamPayload(form, h.countryMCC(form.str("country"))), name); err != nil {
			return h.campaignResult(c, err, generic)
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "message": "Campaign has been created"})
	}
	uploads := map[string]campaignUpload{}
	for field := range mainstreamFolders {
		uploads[field] = inspectCampaignUpload(c, field)
	}
	if v := mainstreamErrors(form, uploads); len(v.messages) > 0 {
		return v.failed(c)
	}
	generic := "An error occurred while preparing the campaigns for create. Please try again."
	if h.campaignNameTaken(name) {
		return h.nameTakenResponse(c, name)
	}
	urls, err := h.saveMainstreamImages(c, campaignImagePrefix(name, form.str("country"), form.str("channel")), mainstreamFolders)
	if err != nil {
		return h.campaignResult(c, err, generic)
	}
	if err := h.submitCampaign(c, buildMainstreamPayload(form, h.countryMCC(form.str("country")), urls), name); err != nil {
		return h.campaignResult(c, err, generic)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "message": "Campaign has been created"})
}

// StoreCampaignMainstreamSingleUrl is CampaignMainstreamController::storeSingleUrl.
func (h *IncomingHandler) StoreCampaignMainstreamSingleUrl(c *fiber.Ctx) error {
	form, err := parseCampaignForm(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid request body"})
	}
	uploads := map[string]campaignUpload{"banner": inspectCampaignUpload(c, "banner"), "background": inspectCampaignUpload(c, "background")}
	if v := mainstreamSingleUrlErrors(form, uploads); len(v.messages) > 0 {
		return v.failed(c)
	}
	generic := "An error occurred while preparing the campaigns for create. Please try again."
	name := form.str("campaign_name")
	if h.campaignNameTaken(name) {
		return h.nameTakenResponse(c, name)
	}
	urls, err := h.saveMainstreamImages(c, campaignImagePrefix(name, form.str("country"), "SUL"), mainstreamSingleUrlFolders)
	if err != nil {
		return h.campaignResult(c, err, generic)
	}
	if err := h.submitCampaign(c, buildMainstreamSingleUrlPayload(form, urls["banner"], urls["background"]), name); err != nil {
		return h.campaignResult(c, err, generic)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "message": "Campaign has been created"})
}

const mainstreamUpdateGeneric = "An error occurred while preparing the campaigns for update. Please try again."

func (h *IncomingHandler) finishMainstreamUpdate(c *fiber.Ctx, payload map[string]interface{}, auditID string, after interface{}, name string) error {
	applyErr := h.applyCampaignPayload(payload)
	h.writeCampaignAudit(c, "updated", auditID, []interface{}{}, after, name)
	if applyErr != nil {
		h.Logs.Errorf("mainstream update: %v", applyErr)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Something went wrong"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "message": "Campaign has been updated"})
}

// UpdateCampaignMainstream is CampaignMainstreamController::update.
func (h *IncomingHandler) UpdateCampaignMainstream(c *fiber.Ctx) error {
	generic := func(err error) error {
		h.Logs.Errorf("mainstream update: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": mainstreamUpdateGeneric})
	}
	form, err := parseCampaignForm(c)
	if err != nil {
		return generic(err)
	}
	campaign, details, err := h.loadCampaign(form.str("campaign_id"), nil)
	if err != nil {
		return generic(err)
	}
	if campaign == nil || len(details) == 0 {
		return generic(errors.New("campaign or its details not found"))
	}
	urls, err := h.saveMainstreamImages(c, campaignImagePrefix(campaign.str("name"), campaign.str("country"), form.str("channel")), mainstreamFolders)
	if err != nil {
		return generic(err)
	}
	payload := buildMainstreamUpdatePayload(form, campaign, details[0], h.countryMCC(campaign.str("country")), urls)
	id := campaign.str("campaign_id")
	after := map[string]interface{}{details[0].str("operator") + " - " + details[0].str("service") + " - " + details[0].str("adnet"): payload}
	return h.finishMainstreamUpdate(c, payload, id, after, "("+id+") "+campaign.str("name"))
}

// UpdateCampaignMainstreamSingleUrl is CampaignMainstreamController::updateSingleUrl.
func (h *IncomingHandler) UpdateCampaignMainstreamSingleUrl(c *fiber.Ctx) error {
	generic := func(err error) error {
		h.Logs.Errorf("mainstream single url update: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": mainstreamUpdateGeneric})
	}
	form, err := parseCampaignForm(c)
	if err != nil {
		return generic(err)
	}
	campaign, details, err := h.loadCampaign(form.str("campaign_id"), nil)
	if err != nil {
		return generic(err)
	}
	if campaign == nil {
		return generic(errors.New("campaign not found"))
	}
	urls, err := h.saveMainstreamImages(c, campaignImagePrefix(campaign.str("name"), campaign.str("country"), "SUL"), mainstreamSingleUrlFolders)
	if err != nil {
		return generic(err)
	}
	payload, err := buildMainstreamSingleUrlUpdatePayload(form, campaign, details, urls["banner"], urls["background"])
	if err != nil {
		return generic(err)
	}
	return h.finishMainstreamUpdate(c, payload, "-1", payload, campaign.str("name"))
}
