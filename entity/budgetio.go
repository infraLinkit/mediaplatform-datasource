package entity

import "time"


type (
	BudgetIORequest struct {
		CPName            string               `json:"cp_name"`
		PICName           string               `json:"pic_name"`
		ContactEmail      string               `json:"contact_email"`
		CPBusinessPICName string               `json:"cp_business_pic_name"`
		Signature         string               `json:"signature"`
		SubmittedBy 	  string			   `json:"submitted_by"`
		Data              []BudgetIORowRequest `json:"data"`
	}

	BudgetIORowRequest struct {
		IOID               string  `json:"io_id"`
		CampaignType       string  `json:"campaign_type"`
		Month              string  `json:"month"`
		Country            string  `json:"country"`
		CountryName        string  `json:"country_name"`
		Continent          string  `json:"continent"`
		CompanyGroupName   string  `json:"company_group_name"`
		Company            string  `json:"company"`
		Partner            string  `json:"partner"`
		TargetCAC          float64 `json:"target_cac"`
		TargetROI          int     `json:"target_roi"`
		MonthlyMOTarget    float64 `json:"monthly_mo_target"`
		MonthlySpendTarget float64 `json:"monthly_spend_target"`
	}

	DisplayBudgetIO struct {
		ID                 int       `json:"id"`
		Status             string    `json:"status"`
		VerifiedAt         time.Time `json:"verified_at"`

		SubmittedBy        string 	 `json:"submitted_by"`
		VerifiedBy         string 	 `json:"verified_by"`

		IOID               string    `json:"io_id"`
		CampaignType       string    `json:"campaign_type"`
		Month              string    `json:"month"`
		Country            string    `json:"country"`
		CountryName        string    `json:"country_name"`
		Continent          string    `json:"continent"`
		CompanyGroupName   string    `json:"company_group_name"`
		Company            string    `json:"company"`
		Partner            string    `json:"partner"`

		TargetCAC          float64   `json:"target_cac"`
		TargetROI          int       `json:"target_roi"`
		MonthlyMOTarget    float64   `json:"monthly_mo_target"`
		MonthlySpendTarget float64   `json:"monthly_spend_target"`

		CPName             string    `json:"cp_name"`
		PICName            string    `json:"pic_name"`
		ContactEmail       string    `json:"contact_email"`
		CPBusinessPICName  string    `json:"business_pic_name"`
		Signature          string    `json:"signature"`

		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`

		// ===== Extra fields for datatables / filtering =====
		Keyword			  string  `json:"keyword"`
		PageSize          int     `json:"page_size"`
		Page              int     `json:"page"`
		Action            string  `json:"action"`
		DateRange         string  `json:"date_range"`
		DateBefore        string  `json:"date_before"`
		DateAfter         string  `json:"date_after"`
		Draw              int     `json:"draw"`
		Reload            string  `json:"reload"`
		OrderColumn       string  `json:"order_column"`
		OrderDir          string  `json:"order_dir"`
	}


	TotalBudgetIO struct {
		MonthlySpendTarget float64 `json:"monthly_spend_target"`
	}
)
