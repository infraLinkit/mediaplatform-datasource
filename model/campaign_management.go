package model

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"github.com/infraLinkit/mediaplatform-datasource/entity"
	"github.com/lib/pq"
)

func (r *BaseModel) GetCampaignManagement(o entity.DisplayCampaignManagement) ([]entity.CampaignManagementData, entity.CampaignCounts, error) {
	var rows *sql.Rows
	query := r.DB.Model(&entity.CampaignDetail{}).
		Select(`
			campaigns.campaign_id AS campaign_id,
			campaigns.name AS campaign_name,
			campaigns.campaign_objective AS campaign_objective, 
			campaign_details.country, 
			campaign_details.partner, 
			COUNT(DISTINCT campaign_details.operator) AS total_operator, 
			COUNT(DISTINCT campaign_details.service) AS service, 
			COUNT(DISTINCT campaign_details.adnet) AS total_adnet, 
			COUNT(DISTINCT campaign_details.short_code) AS short_code,
			campaign_details.is_active
		`).
		Joins("INNER JOIN campaigns ON campaigns.campaign_id = campaign_details.campaign_id").
		Group("campaigns.campaign_id, campaigns.name, campaigns.campaign_objective, campaign_details.country, campaign_details.partner, campaign_details.is_active, campaigns.created_at")

		orderColumn := map[string]string{
			"total_operator": "COUNT(DISTINCT campaign_details.operator)",
			"service":        "COUNT(DISTINCT campaign_details.service)",
			"total_adnet":    "COUNT(DISTINCT campaign_details.adnet)",
		}
		
		if col, ok := orderColumn[o.OrderColumn]; ok {
			dir := "ASC"
			if strings.ToUpper(o.OrderDir) == "DESC" {
				dir = "DESC"
			}
			query = query.Order(fmt.Sprintf("%s %s", col, dir))
		} else {
			query = query.Order("campaigns.created_at DESC")
		}		
		

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
		if o.CampaignType != "" {
			if o.CampaignType == "mainstream" {
				query = query.Where("campaigns.campaign_objective = ?", "MAINSTREAM")
			} else {
				query = query.Where("campaigns.campaign_objective IN ?", []string{"CPA", "CPC", "CPI", "CPM"})
			}
		}
		if o.URLServiceKey != "" {
			query = query.Where("campaign_details.url_service_key ILIKE ?", "%"+o.URLServiceKey+"%")
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
		var urlKeys []string

		r.DB.ScanRows(rows, &campaign)

		// Fetch all campaign IDs associated with this campaign
		err := r.DB.Table("campaign_details").
			Where("campaign_id = ?", campaign.CampaignID).
			Pluck("id", &campaignIDs).Error

		if err != nil {
			return nil, entity.CampaignCounts{}, err
		}

		campaign.ID = campaignIDs // Assign campaign IDs

		// Get all url_service_key for this group
		err = r.DB.Table("campaign_details").
			Where("campaign_id = ? AND country = ? AND partner = ?", campaign.CampaignID, campaign.Country, campaign.Partner).
			Pluck("DISTINCT url_service_key", &urlKeys).Error
		if err != nil {
			return nil, entity.CampaignCounts{}, err
		}
		campaign.URLServiceKey = urlKeys
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
	// First, get the campaign objective to determine if we should include cc_email
	var campaignObjective string
	err := r.DB.Model(&entity.Campaign{}).
		Select("campaign_objective").
		Where("campaign_id = ?", o.CampaignId).
		Scan(&campaignObjective).Error

	if err != nil {
		return nil, err
	}

	// Build the SELECT clause based on campaign objective
	selectClause := `
		campaign_details.id,
		campaign_details.campaign_id,
        campaign_details.operator, 
        campaign_details.service, 
        campaigns.name AS campaign_name,
		campaigns.campaign_objective AS campaign_objective, 
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
		campaign_details.url_service_key,
		campaign_details.channel,
		campaign_details.url_type,
		campaign_details.device_type,
		campaign_details.is_billable`

	// Add cc_email only if campaign objective is not MAINSTREAM
	if campaignObjective != "MAINSTREAM" {
		selectClause += `,
			adnet_lists.cc_email`
	} else {
		selectClause += `,
			NULL AS cc_email`
	}

	query := r.DB.Model(&entity.CampaignDetail{}).
		Select(selectClause).
		Joins("INNER JOIN campaigns ON campaigns.campaign_id = campaign_details.campaign_id")

	// Only join adnet_lists if campaign objective is not MAINSTREAM
	if campaignObjective != "MAINSTREAM" {
		query = query.Joins("INNER JOIN adnet_lists ON adnet_lists.code = campaign_details.adnet")
	}

	query = query.Where("campaigns.campaign_id = ?", o.CampaignId).
		Where("campaign_details.is_active = ?", o.Status).
		Order("CAST(SUBSTRING(campaign_details.url_service_key FROM '[0-9]+') AS INTEGER) ASC, campaign_details.operator, campaign_details.service")

	rows, err := query.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	campaignMap := make(map[string]map[string]*entity.CampaignManagementDataDetail)

	for rows.Next() {
		var detail entity.CampaignManagementDetail
		var ccEmail interface{}

		// Create scan args based on campaign objective
		scanArgs := []interface{}{
			&detail.ID, &detail.CampaignID, &detail.Operator, &detail.Service, &detail.CampaignName, &detail.CampaignObjective, &detail.Country,
			&detail.Partner, &detail.Adnet, &detail.ShortCode, &detail.MOLimit, &detail.Payout,
			&detail.RatioSend, &detail.RatioReceive, &detail.URLPostback, &detail.URLService,
			&detail.URLanding, &detail.URLWarpLanding, &detail.APIURL, &detail.IsActive, &detail.UrlServiceKey, &detail.Channel, &detail.URLType,
			&detail.DeviceType, &detail.IsBillable, &ccEmail,
		}

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}

		// Handle cc_email based on campaign objective
		if campaignObjective != "MAINSTREAM" {
			if ccEmail != nil {
				// Convert string to pq.StringArray if needed
				if ccEmailStr, ok := ccEmail.(string); ok {
					// If it's a string, we need to parse it or handle it appropriately
					// For now, let's create a single-element array
					detail.CCEmail = pq.StringArray{ccEmailStr}
				} else if ccEmailArr, ok := ccEmail.(pq.StringArray); ok {
					detail.CCEmail = ccEmailArr
				} else {
					// Fallback to empty array
					detail.CCEmail = pq.StringArray{}
				}
			} else {
				detail.CCEmail = pq.StringArray{}
			}
		} else {
			// For MAINSTREAM campaigns, set empty array
			detail.CCEmail = pq.StringArray{}
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