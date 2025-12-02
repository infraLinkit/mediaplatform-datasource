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

	ContinentReport struct {
		Continent                 string             `json:"continent"`
		Month                     string             `json:"month"`

		ActualCostWeek1Continent  int                `json:"actual_cost_week1_continent"`
		KPIWeek1Continent         int                `json:"kpi_week1_continent"`
		ActualCostWeek2Continent  int                `json:"actual_cost_week2_continent"`
		KPIWeek2Continent         int                `json:"kpi_week2_continent"`
		ActualCostWeek3Continent  int                `json:"actual_cost_week3_continent"`
		KPIWeek3Continent         int                `json:"kpi_week3_continent"`
		ActualCostWeek4Continent  int                `json:"actual_cost_week4_continent"`
		KPIWeek4Continent         int                `json:"kpi_week4_continent"`

		TotalActualCostContinent  int                `json:"total_actual_cost_continent"`
		TotalIOTargetContinent    int                `json:"total_io_target_continent"`
		BudgetUsageContinent      int                `json:"budget_usage_continent"`

		Countries                 []CountryReport    `json:"countries"`
	}

	CountryReport struct {
		Country                   string `json:"country"`

		ActualCostWeek1Country    int    `json:"actual_cost_week1_country"`
		KPIWeek1Country           int    `json:"kpi_week1_country"`
		ActualCostWeek2Country    int    `json:"actual_cost_week2_country"`
		KPIWeek2Country           int    `json:"kpi_week2_country"`
		ActualCostWeek3Country    int    `json:"actual_cost_week3_country"`
		KPIWeek3Country           int    `json:"kpi_week3_country"`
		ActualCostWeek4Country    int    `json:"actual_cost_week4_country"`
		KPIWeek4Country           int    `json:"kpi_week4_country"`

		TotalActualCostCountry    int    `json:"total_actual_cost_country"`
		TotalIOTargetCountry      int    `json:"total_io_target_country"`
		BudgetUsageCountry        int    `json:"budget_usage_country"`
	}

	DisplaySummaryBudgetIO struct {
		ID                 int       `json:"id"`
		Month              string    `json:"month"`
		Continent          string    `json:"continent"`
		Country            string    `json:"country"`
		Company            string    `json:"company"`
		Partner            string    `json:"partner"`

		TotalMonthlySpendTarget float64 `json:"total_monthly_spend_target"`
		ActualWeek1             float64 `json:"actual_week_1"`
		ActualWeek2             float64 `json:"actual_week_2"`
		ActualWeek3             float64 `json:"actual_week_3"`
		ActualWeek4             float64 `json:"actual_week_4"`

		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`

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
)
