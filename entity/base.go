package entity

type (
	GlobalResponse struct {
		Code    int    `json:"code" xml:"code"`
		Message string `json:"message" xml:"message"`
	}

	GlobalResponseWithData struct {
		Code    int             `json:"code" xml:"code"`
		Message string          `json:"message" xml:"message"`
		Data    PixelStorageRsp `json:"data" xml:"data"`
	}

	PixelStorageRsp struct {
		Adnet         string `json:"adnet"`
		IsBillable    bool   `json:"is_billable"`
		Pixel         string `json:"pixel"`
		TrxId         string `json:"trx_id"`
		Msisdn        string `json:"msisdn"`
		Browser       string `json:"browser"`
		OS            string `json:"os"`
		PubId         string `json:"pubid"`
		Handset       string `json:"handset"`
		PixelUsedDate string `json:"pixel_used_date"`
	}

	ReturnResponse struct {
		HttpStatus int
		Rsp        interface{}
	}
)
