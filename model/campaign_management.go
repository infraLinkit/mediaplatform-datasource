package model

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/infraLinkit/mediaplatform-datasource/entity"
)

func (r *BaseModel) GetCampaignManagement(o entity.DisplayCampaignManagement) ([]entity.CampaignManagementData, entity.CampaignCounts, error) {
	var rows *sql.Rows
	query := r.DB.Model(&entity.CampaignDetail{}).
		Select(`
			campaigns.campaign_id AS campaign_id,
			campaigns.name AS campaign_name, 
			campaign_details.country, 
			campaign_details.partner, 
			COUNT(DISTINCT campaign_details.operator) AS total_operator, 
			COUNT(DISTINCT campaign_details.service) AS service, 
			COUNT(DISTINCT campaign_details.adnet) AS total_adnet, 
			campaign_details.short_code, 
			campaign_details.is_active
		`).
		Joins("INNER JOIN campaigns ON campaigns.campaign_id = campaign_details.campaign_id").
		Group("campaigns.campaign_id, campaigns.name, campaign_details.country, campaign_details.partner, campaign_details.short_code, campaign_details.is_active")

	if o.Action == "Search" {
		if o.Country != "" {
			query = query.Where("campaign_details.country = ?", o.Country)
		}
		if o.Operator != "" {
			query = query.Where("campaign_details.operator = ?", o.Operator)
		}
		if o.Service != "" {
			query = query.Where("campaign_details.service = ?", o.Service)
		}
		if o.Adnet != "" {
			query = query.Where("campaign_details.adnet = ?", o.Adnet)
		}
		if o.Partner != "" {
			query = query.Where("campaign_details.partner = ?", o.Partner)
		}
		if o.Status != "" {
			query = query.Where("campaign_details.is_active = ?", o.Status)
		}
		if o.CampaignName != "" {
			query = query.Where("campaigns.name = ?", o.CampaignName)
		}
	}

	rows, err := query.Rows()
	if err != nil {
		return nil, entity.CampaignCounts{}, err
	}
	defer rows.Close()

	var campaigns []entity.CampaignManagementData
	var total, active, nonActive int

	for rows.Next() {
		var campaign entity.CampaignManagementData
		var campaignIDs []int // Slice to store all IDs for a campaign

		r.DB.ScanRows(rows, &campaign)

		// Fetch all campaign IDs associated with this campaign
		err := r.DB.Table("campaign_details").
			Where("campaign_id = ?", campaign.CampaignID).
			Pluck("id", &campaignIDs).Error

		if err != nil {
			return nil, entity.CampaignCounts{}, err
		}

		campaign.ID = campaignIDs // Assign campaign IDs
		campaigns = append(campaigns, campaign)

		// Count total campaigns
		total++
		if campaign.IsActive {
			active++
		} else {
			nonActive++
		}
	}

	r.Logs.Debug(fmt.Sprintf("Total data : %d ...\n", len(campaigns)))

	return campaigns, entity.CampaignCounts{
		TotalCampaigns:          total,
		TotalActiveCampaigns:    active,
		TotalNonActiveCampaigns: nonActive,
	}, nil
}

func (r *BaseModel) GetCampaignManagementDetail(o entity.DisplayCampaignManagement) ([]entity.CampaignManagementDataDetail, error) {
	query := r.DB.Model(&entity.CampaignDetail{}).
		Select(`
			campaign_details.id,
			campaign_details.campaign_id,
            campaign_details.operator, 
            campaign_details.service, 
            campaigns.name AS campaign_name, 
            campaign_details.country, 
            campaign_details.partner,
			campaign_details.adnet,
            campaign_details.short_code, 
            campaign_details.mo_capping AS mo_limit, 
            campaign_details.po, 
            campaign_details.ratio_send, 
            campaign_details.ratio_receive, 
            campaign_details.url_postback, 
            campaign_details.url_service, 
            campaign_details.url_landing, 
            campaign_details.url_warp_landing, 
            campaign_details.api_url, 
            campaign_details.is_active,
			campaign_details.url_service_key
        `).
		Joins("INNER JOIN campaigns ON campaigns.campaign_id = campaign_details.campaign_id").
		Where("campaigns.campaign_id = ?", o.CampaignId).
		Where("campaign_details.is_active = ?", o.Status).
		Order("campaign_details.operator, campaign_details.service")

	rows, err := query.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	campaignMap := make(map[string]map[string]*entity.CampaignManagementDataDetail)

	for rows.Next() {
		var detail entity.CampaignManagementDetail
		if err := rows.Scan(
			&detail.ID, &detail.CampaignID, &detail.Operator, &detail.Service, &detail.CampaignName, &detail.Country,
			&detail.Partner, &detail.Adnet, &detail.ShortCode, &detail.MOLimit, &detail.Payout,
			&detail.RatioSend, &detail.RatioReceive, &detail.URLPostback, &detail.URLService,
			&detail.URLanding, &detail.URLWarpLanding, &detail.APIURL, &detail.IsActive, &detail.UrlServiceKey,
		); err != nil {
			return nil, err
		}

		if _, exists := campaignMap[detail.Operator]; !exists {
			campaignMap[detail.Operator] = make(map[string]*entity.CampaignManagementDataDetail)
		}

		if _, exists := campaignMap[detail.Operator][detail.Service]; !exists {
			campaignMap[detail.Operator][detail.Service] = &entity.CampaignManagementDataDetail{
				Operator: detail.Operator,
				Service:  detail.Service,
				Details:  []entity.CampaignManagementDetail{},
			}
		}

		campaignMap[detail.Operator][detail.Service].Details = append(
			campaignMap[detail.Operator][detail.Service].Details, detail,
		)
	}

	var campaigns []entity.CampaignManagementDataDetail
	for _, serviceMap := range campaignMap {
		for _, campaign := range serviceMap {
			campaigns = append(campaigns, *campaign)
		}
	}

	// Sorting berdasarkan Operator dan Service
	sort.SliceStable(campaigns, func(i, j int) bool {
		if campaigns[i].Operator == campaigns[j].Operator {
			return campaigns[i].Service < campaigns[j].Service
		}
		return campaigns[i].Operator < campaigns[j].Operator
	})

	r.Logs.Debug(fmt.Sprintf("Total data: %d ...\n", len(campaigns)))
	return campaigns, nil
}