package entity

type (
	//"[{\"id\":1,\"urlservicekey\":\"idtelgempastelmbv\",\"campaign_id\":\"ID01\",\"name\":\"ID Pass Telesat Gazy\",\"objective\":\"cpa\",\"country\":\"ID\",\"advertiser\":\"PT\",\"operator\":\"telkomsel\",\"partner\":\"pass\",\"aggregator\":\"telesat\",\"adnet\":\"mbv\",\"service\":\"gazy\",\"keyword\":\"gazy\",\"subkeyword\":\"\",\"is_billable\":false,\"plan\":\"\",\"short_code\":\"1234\",\"device_type\":\"all\",\"os\":\"linux\",\"url_type\":\"click2sms\",\"click_type\":2,\"click_delay\":0,\"client_type\":\"internal\",\"traffic_source\":false,\"unique_click\":true,\"url_banner\":\"http://\",\"url_landing\":\"http://\",\"url_warp_landing\":\"http://\",\"url_service\":\"http://\",\"url_tfc_or_smartlink\":\"http://\",\"custom_integration\":\"\",\"ip_address\":\"\",\"is_active\":true,\"mo_capping\":0,\"counter_mo_capping\":0,\"status_capping\":false,\"kpi_upper_limit_capping\":500,\"is_machine_learning_capping\":false,\"ratio_send\":1,\"ratio_receive\":1,\"status_ratio\":true,\"kpi_upper_limit_ratio_send\":1,\"kpi_upper_limit_ratio_receive\":2,\"is_machine_learning_ratio\":false,\"api_url\":\"http://\",\"last_update\":\"\",\"last_update_capping\":\"\"}]"

	DataCampaignAction struct {
		Action       string       `json:"action"`
		Id           int          `json:"id"`
		CampaignId   string       `json:"campaign_id"`
		CampaignName string       `json:"name"`
		Objective    string       `json:"campaign_objective"`
		Country      string       `json:"country"`
		Advertiser   string       `json:"advertiser"`
		IsDelCamp    bool         `json:"is_del_camp"`
		CPCR         string       `json:"cost_per_conversion"`
		AgencyFee    string       `json:"agency_fee"`
		DataConfig   []DataConfig `json:"data"`
	}

	DataConfig struct {
		Id                        int     `json:"id"`            // <-
		URLServiceKey             string  `json:"urlservicekey"` // lp-idxlgazpastelmbv
		CampaignId                string  `json:"campaign_id"`   // ID1
		CampaignName              string  `json:"name"`
		Objective                 string  `json:"objective"`
		Country                   string  `json:"country"`
		Advertiser                string  `json:"advertiser"`
		Operator                  string  `json:"operator"`
		Partner                   string  `json:"partner"`
		Aggregator                string  `json:"aggregator"`
		Adnet                     string  `json:"adnet"`
		Service                   string  `json:"service"`
		Keyword                   string  `json:"keyword"`
		SubKeyword                string  `json:"subkeyword"`
		IsBillable                bool    `json:"is_billable"`
		Plan                      string  `json:"plan"`
		PO                        string  `json:"po"`
		Cost                      string  `json:"cost"`
		PubId                     string  `json:"pubid"`
		ShortCode                 string  `json:"short_code"`
		DeviceType                string  `json:"device_type"`
		OS                        string  `json:"os"`
		URLType                   string  `json:"url_type"`
		ClickType                 int     `json:"click_type"`
		ClickDelay                int     `json:"click_delay"`
		ClientType                string  `json:"client_type"`
		TrafficSource             bool    `json:"traffic_source"`
		UniqueClick               bool    `json:"unique_click"`
		URLBanner                 string  `json:"url_banner"`
		URLLanding                string  `json:"url_landing"`
		URLWarpLanding            string  `json:"url_warp_landing"`
		URLService                string  `json:"url_service"`
		URLTFCSmartlink           string  `json:"url_tfc_or_smartlink"`
		GlobPost                  bool    `json:"glob_post"`
		URLGlobPost               string  `json:"url_glob_post"`
		CustomIntegration         string  `json:"custom_integration"`
		IPAddress                 []uint8 `json:"ip_address"`
		ISP                       string  `json:"isp"`
		IsActive                  bool    `json:"is_active"`
		MOCapping                 int     `json:"mo_capping"`
		CounterMOCapping          int     `json:"counter_mo_capping"`
		StatusCapping             bool    `json:"status_capping"`
		KPIUpperLimitCapping      int     `json:"kpi_upper_limit_capping"`
		IsMachineLearningCapping  bool    `json:"is_machine_learning_capping"`
		RatioSend                 int     `json:"ratio_send"`
		RatioReceive              int     `json:"ratio_receive"`
		CounterMORatio            int     `json:"counter_mo_ratio"`
		StatusRatio               bool    `json:"status_ratio"`
		KPIUpperLimitRatioSend    int     `json:"kpi_upper_limit_ratio_send"`
		KPIUpperLimitRatioReceive int     `json:"kpi_upper_limit_ratio_receive"`
		IsMachineLearningRatio    bool    `json:"is_machine_learning_ratio"`
		APIURL                    string  `json:"api_url"`
		LastUpdate                string  `json:"last_update"`
		LastUpdateCapping         string  `json:"last_update_capping"`
		CPCR                      string  `json:"cost_per_conversion"`
		AgencyFee                 string  `json:"agency_fee"`
	}

	//'{"id":1,"urlservicekey":"idtelgempastelmbv","campaign_id":"ID01","country":"ID","partner":"pass","operator":"telkomsel","aggregator":"telesat","service":"gazy","short_code":"1234","adnet":"mbv","keyword":"gazy","subkeyword":"","is_billable":false,"plan":"","traffic":0,"landing":0,"click":0,"redirect":0,"traffic_data":[],"landing_data":[],"click_data":[],"redirect_data":[]}'

	DataCounter struct {
		CampaignDetailId int                         `json:"campaign_detail_id"`
		URLServiceKey    string                      `json:"urlservicekey"`
		CampaignId       string                      `json:"campaign_id"`
		Country          string                      `json:"country"`
		Partner          string                      `json:"partner"`
		Operator         string                      `json:"operator"`
		Aggregator       string                      `json:"aggregator"`
		Service          string                      `json:"service"`
		ShortCode        string                      `json:"short_code"`
		Adnet            string                      `json:"adnet"`
		Keyword          string                      `json:"keyword"`
		SubKeyword       string                      `json:"subkeyword"`
		IsBillable       bool                        `json:"is_billable"`
		Plan             string                      `json:"plan"`
		Traffic          int                         `json:"traffic"`
		Landing          int                         `json:"landing"`
		Click            int                         `json:"click"`
		Redirect         int                         `json:"redirect"`
		MOReceived       int                         `json:"moreceived"`
		Postback         int                         `json:"postback"`
		TotalFP          int                         `json:"totalfp"`
		TrafficData      []DataCounterDetail         `json:"traffic_data"`
		LandingData      []DataCounterDetail         `json:"landing_data"`
		ClickData        []DataCounterDetail         `json:"click_data"`
		RedirectData     []DataCounterDetail         `json:"redirect_data"`
		MOData           []DataCounterDetailInternal `json:"mo_data"`
		PostbackData     []DataCounterDetailInternal `json:"postback_data"`
		FPData           []DataCounterDetailInternal `json:"fp_data"`
	}

	DataCounterDetail struct {
		Date       string `json:"date"`
		Time       string `json:"time"`
		HTTPStatus string `json:"http_status"`
	}

	DataCounterDetailInternal struct {
		Date   string `json:"date"`
		Msisdn string `json:"msisdn"`
		TrxId  string `json:"trxid"`
		Pixel  string `json:"pixel"`
		Code   string `json:"code"`
		Status bool   `json:"status"`
	}

	PixelStorage struct {
		Id                int    `json:"id"`
		CampaignDetailId  int    `json:"campaign_detail_id"`
		PxDate            string `json:"pxdate"`
		URLServiceKey     string `json:"urlservicekey"`
		CampaignId        string `json:"campaign_id"`
		Country           string `json:"country"`
		Partner           string `json:"partner"`
		Operator          string `json:"operator"`
		Aggregator        string `json:"aggregator"`
		Service           string `json:"service"`
		ShortCode         string `json:"short_code"`
		Adnet             string `json:"adnet"`
		Keyword           string `json:"keyword"`
		Subkeyword        string `json:"subkeyword"`
		IsBillable        bool   `json:"is_billable"`
		Plan              string `json:"plan"`
		URL               string `json:"url"`
		URLType           string `json:"url_type"`
		Pixel             string `json:"pixel"`
		TrxId             string `json:"trx_id"`
		Token             string `json:"token"`
		Msisdn            string `json:"msisdn"`
		IsUsed            bool   `json:"is_used"`
		Browser           string `json:"browser"`
		OS                string `json:"os"`
		IP                string `json:"ip"`
		ISP               string `json:"isp"`
		ReferralURL       string `json:"referral_url"`
		PubId             string `json:"pubid"`
		UserAgent         string `json:"user_agent"`
		TrafficSource     bool   `json:"traffic_source"`
		TrafficSourceData string `json:"traffic_source_data"`
		UserRejected      bool   `json:"user_rejected"`
		UniqueClick       bool   `json:"unique_click"`
		UserDuplicated    bool   `json:"user_duplicated"`
		Handset           string `json:"handset"`
		HandsetCode       string `json:"handset_code"`
		HandsetType       string `json:"handset_type"`
		URLLanding        string `json:"url_landing"`
		URLWarpLanding    string `json:"url_warp_landing"`
		URLService        string `json:"url_service"`
		URLTFCSmartlink   string `json:"url_tfc_or_smartlink"`
		PixelUsedDate     string `json:"pixel_used_date"`
		StatusPostback    bool   `json:"status_postback"`
		IsUnique          bool   `json:"is_unique"`
		URLPostback       string `json:"url_postback"`
		StatusURLPostback string `json:"status_url_postback"`
		ReasonURLPostback string `json:"reason_url_postback"`
		PO                string `json:"po"`
		Cost              string `json:"cost"`
	}
)
