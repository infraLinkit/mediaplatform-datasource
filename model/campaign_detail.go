package model

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/infraLinkit/mediaplatform-datasource/entity"
	"gorm.io/gorm"
)

func (r *BaseModel) GetLastCampaignId(tbl string) int {

	var result int
	row := r.DB.Table(tbl).
		Select("COALESCE(MAX(id), 0) + 0").Row()
	row.Scan(&result)

	return result
}

func (r *BaseModel) GetCampaignByCampaignId(o entity.Campaign) (entity.Campaign, bool) {

	result := r.DB.Model(&o).
		Where("campaign_id = ?", o.CampaignId).
		First(&o)

	b := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if b {
		return o, false
	} else {
		r.Logs.Warn(fmt.Sprintf("Campaign id not found %#v", o))
		return o, true
	}
}

func (r *BaseModel) NewCampaign(o entity.Campaign) int {

	result := r.DB.Create(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return int(o.ID)
}

func (r *BaseModel) NewCampaignDetail(o entity.CampaignDetail) int {

	result := r.DB.Create(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return int(o.ID)
}

func (r *BaseModel) ResetCappingCampaign(o entity.CampaignDetail) error {

	result := r.DB.Model(&o).
		Where("is_active = ?", o.IsActive).
		Updates(entity.CampaignDetail{CounterMORatio: 0, StatusCapping: false})

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) GetCampaignByCampaignDetailId(o entity.CampaignDetail) (entity.CampaignDetail, bool) {

	result := r.DB.Model(&o).
		Where("url_service_key = ? AND country = ? AND operator = ? AND partner = ? AND service = ? AND adnet = ?", o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet).
		First(&o)

	b := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if b {
		return o, false
	} else {
		r.Logs.Warn(fmt.Sprintf("Campaign existed or data found %#v", o))
		return o, true
	}
}

func (r *BaseModel) UpdateCampaign(o entity.Campaign) error {

	result := r.DB.Model(&o).
		Where("campaign_id = ?", o.CampaignId).
		Updates(entity.Campaign{Name: o.Name, CampaignObjective: o.CampaignObjective, Country: o.Country, Advertiser: o.Advertiser})

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) UpdateCampaignDetail(o entity.CampaignDetail) error {

	result := r.DB.Exec("UPDATE campaign_details SET is_active = ? WHERE id = ?", o.IsActive, o.ID)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) SaveCampaignDetail(o entity.CampaignDetail) error {

	result := r.DB.Save(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) DelCampaign(o entity.Campaign) error {

	result := r.DB.
		Unscoped().
		Where("campaign_id = ?", o.CampaignId).
		Delete(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) DelCampaignDetail(o entity.CampaignDetail) error {

	result := r.DB.
		Unscoped().
		Where("url_service_key = ? AND country = ? AND operator = ? AND partner = ? AND service = ? AND adnet = ? AND campaign_id = ?", o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId).
		Delete(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) EditSettingCampaignDetail(o entity.CampaignDetail) error {

	result := r.DB.Model(&o).
		Where("url_service_key = ? AND country = ? AND operator = ? AND partner = ? AND service = ? AND adnet = ? AND campaign_id = ?", o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId).
		Updates(entity.CampaignDetail{PO: o.PO, MOCapping: o.MOCapping, RatioSend: o.RatioSend, RatioReceive: o.RatioReceive, LastUpdate: o.LastUpdate})

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) UpdateStatusCampaignDetail(o entity.CampaignDetail) error {
	result := r.DB.Model(&o).
		Where("url_service_key = ? AND country = ? AND operator = ? AND partner = ? AND service = ? AND adnet = ? AND campaign_id = ?",
			o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId).
		Updates(entity.CampaignDetail{IsActive: o.IsActive})

	r.Logs.Debug(fmt.Sprintf("Query Affected: %d, Error: %v", result.RowsAffected, result.Error))

	if result.RowsAffected == 0 {
		r.Logs.Debug("No rows updated. Check if the WHERE condition matches any records.")
	}

	return result.Error
}

func (r *BaseModel) GetCampaignDetail(o entity.CampaignDetail) (entity.CampaignDetail, bool) {

	result := r.DB.Model(&o).First(&o)

	b := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if b {
		return o, false
	} else {
		r.Logs.Warn(fmt.Sprintf("Campaign detail existed or data found %#v", o))
		return o, true
	}
}

func (r *BaseModel) CounterRatioById(o entity.CampaignDetail) error {

	result := r.DB.Exec("UPDATE campaign_details SET counter_mo_ratio = counter_mo_ratio+1 WHERE id = ?", o.ID)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) UpdateStatusCounterById(o entity.CampaignDetail) error {

	//result := r.DB.Model(&o).Where("id = ?", o.ID).Updates(entity.CampaignDetail{CounterMOCapping: o.CounterMOCapping, StatusCapping: o.StatusCapping, CounterMORatio: o.CounterMORatio, StatusRatio: o.StatusRatio, LastUpdate: o.LastUpdate})

	result := r.DB.Model(&o).Select("counter_mo_capping", "status_capping", "counter_mo_ratio", "status_ratio", "last_update").Updates(o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) GetCampaignDetailByStatus(o entity.CampaignDetail, useStatus bool) ([]entity.ResultCampaign, error) {

	var rows *sql.Rows

	/* if useStatus {
		rows, _ = r.DB.Model(&entity.CampaignDetail{}).Select("campaigns.name, campaigns.campaign_objective, campaigns.advertiser, campaign_details.id, campaign_details.url_service_key, campaign_details.campaign_id, campaign_details.country, campaign_details.operator, campaign_details.partner, campaign_details.aggregator, campaign_details.adnet, campaign_details.service, campaign_details.keyword, campaign_details.subkeyword, campaign_details.is_billable, campaign_details.plan, campaign_details.po, campaign_details.cost, campaign_details.pub_id, campaign_details.short_code, campaign_details.device_type, campaign_details.os, campaign_details.url_type, campaign_details.click_type, campaign_details.click_delay, campaign_details.client_type, campaign_details.traffic_source, campaign_details.unique_click, campaign_details.url_banner, campaign_details.url_landing, campaign_details.url_warp_landing, campaign_details.url_service, campaign_details.url_tfcor_smartlink, campaign_details.glob_post, campaign_details.url_glob_post, campaign_details.custom_integration, campaign_details.ip_address, campaign_details.is_active, campaign_details.mo_capping, campaign_details.counter_mo_capping, campaign_details.status_capping, campaign_details.kpi_upper_limit_capping, campaign_details.is_machine_learning_capping, campaign_details.ratio_send, campaign_details.ratio_receive, campaign_details.counter_mo_ratio, campaign_details.status_ratio, campaign_details.kpi_upper_limit_ratio_send, campaign_details.kpi_upper_limit_ratio_receive, campaign_details.is_machine_learning_ratio, campaign_details.api_url, campaign_details.last_update, campaign_details.cost_per_conversion, campaign_details.agency_fee, campaign_details.target_daily_budget, campaign_details.url_postback").Joins("JOIN campaigns ON campaigns.campaign_id = campaign_details.campaign_id").Where("campaign_details.is_active = ?", o.IsActive).Rows()
	} else { */
	rows, _ = r.DB.Model(&entity.CampaignDetail{}).Where("is_active = ?", o.IsActive).Rows()
	//}

	defer rows.Close()

	var (
		ss []entity.ResultCampaign
	)

	for rows.Next() {

		var s entity.CampaignDetail

		// ScanRows scans a row into a struct
		r.DB.ScanRows(rows, &s)

		camp, _ := r.GetCampaignByCampaignId(entity.Campaign{
			CampaignId: s.CampaignId,
		})

		ss = append(ss, entity.ResultCampaign{
			ID:                        s.ID,
			Name:                      camp.Name,
			CampaignObjective:         camp.CampaignObjective,
			Advertiser:                camp.Advertiser,
			URLServiceKey:             s.URLServiceKey,
			CampaignId:                s.CampaignId,
			Country:                   s.Country,
			Operator:                  s.Operator,
			Partner:                   s.Partner,
			Aggregator:                s.Aggregator,
			Adnet:                     s.Adnet,
			Service:                   s.Service,
			Keyword:                   s.Keyword,
			Subkeyword:                s.Subkeyword,
			IsBillable:                s.IsBillable,
			Plan:                      s.Plan,
			PO:                        s.PO,
			Cost:                      s.Cost,
			PubId:                     s.PubId,
			ShortCode:                 s.ShortCode,
			DeviceType:                s.DeviceType,
			OS:                        s.OS,
			URLType:                   s.URLType,
			ClickType:                 s.ClickType,
			ClickDelay:                s.ClickDelay,
			ClientType:                s.ClientType,
			TrafficSource:             s.TrafficSource,
			UniqueClick:               s.UniqueClick,
			URLBanner:                 s.URLBanner,
			URLLanding:                s.URLLanding,
			URLWarpLanding:            s.URLWarpLanding,
			URLService:                s.URLService,
			URLTFCORSmartlink:         s.URLTFCORSmartlink,
			GlobPost:                  s.GlobPost,
			URLGlobPost:               s.URLGlobPost,
			CustomIntegration:         s.CustomIntegration,
			IpAddress:                 s.IpAddress,
			IsActive:                  s.IsActive,
			MOCapping:                 s.MOCapping,
			CounterMOCapping:          s.CounterMOCapping,
			StatusCapping:             s.StatusCapping,
			KPIUpperLimitCapping:      s.KPIUpperLimitCapping,
			IsMachineLearningCapping:  s.IsMachineLearningCapping,
			RatioSend:                 s.RatioSend,
			RatioReceive:              s.RatioReceive,
			CounterMORatio:            s.CounterMORatio,
			StatusRatio:               s.StatusRatio,
			KPIUpperLimitRatioSend:    s.KPIUpperLimitRatioSend,
			KPIUpperLimitRatioReceive: s.KPIUpperLimitRatioReceive,
			IsMachineLearningRatio:    s.IsMachineLearningRatio,
			APIURL:                    s.APIURL,
			LastUpdate:                s.LastUpdate,
			LastUpdateCapping:         s.LastUpdateCapping,
			CostPerConversion:         s.CostPerConversion,
			AgencyFee:                 s.AgencyFee,
			TargetDailyBudget:         s.TargetDailyBudget,
			TechnicalFee:              s.TechnicalFee,
			URLPostback:               s.URLPostback,
			Channel:                   s.Channel,
		})
	}

	return ss, rows.Err()
}

func (r *BaseModel) UpdateCPAReport(o entity.CampaignDetail) error {

	result := r.DB.Model(&o).Where("url_service_key = ? AND country = ? AND operator = ? AND partner = ? AND service = ? AND adnet = ? AND campaign_id = ?", o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId).Updates(entity.CampaignDetail{CostPerConversion: o.CostPerConversion, AgencyFee: o.AgencyFee})

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) UpdateCampaignMonitoringBudget(o entity.CampaignDetail) error {

	result := r.DB.Model(&o).Where("country = ? AND operator = ?", o.Country, o.Operator).Updates(entity.CampaignDetail{TargetDailyBudget: o.TargetDailyBudget})

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) UpdateKeyMainstreamCampaignDetail(o entity.CampaignDetail) error {

	result := r.DB.Exec(`
		UPDATE campaign_details 
		SET status_submit_key_mainstream = ?, key_mainstream = ? 
		WHERE url_service_key = ? AND campaign_id = ?`,
		o.StatusSubmitKeyMainstream, o.KeyMainstream, o.URLServiceKey, o.CampaignId,
	)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) UpdateGoogleSheetCampaignDetail(o entity.CampaignDetail) error {

	result := r.DB.Exec(`
		UPDATE campaign_details 
		SET google_sheet = ? 
		WHERE url_service_key = ? AND campaign_id = ?`,
		o.GoogleSheet, o.URLServiceKey, o.CampaignId,
	)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}
