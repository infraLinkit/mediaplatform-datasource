package entity

type (
	GlobalResponse struct {
		Code    int    `json:"code" xml:"code"`
		Message string `json:"message" xml:"message"`
	}

	GlobalResponseWithData struct {
		Code    int         `json:"code" xml:"code"`
		Message string      `json:"message" xml:"message"`
		Data    interface{} `json:"data" xml:"data"`
	}

	GlobalResponseWithDataTable struct {
		Draw            int         `json:"draw" xml:"draw"`
		Code            int         `json:"code" xml:"code"`
		Message         string      `json:"message" xml:"message"`
		Data            interface{} `json:"data" xml:"data"`
		RecordsTotal    int         `json:"recordsTotal" xml:"recordsTotal"`
		RecordsFiltered int         `json:"recordsFiltered" xml:"recordsFiltered"`
	}

	PixelStorageRsp struct {
		Adnet         string `json:"adnet"`
		IsBillable    bool   `json:"is_billable"`
		Pixel         string `json:"pixel"`
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

	GlobalResponseWithTable struct {
		Code            int         `json:"code" xml:"code"`
		Message         string      `json:"message" xml:"message"`
		Data            interface{} `json:"data" xml:"data"`
		Draw            int         `json:"draw" xml:"draw"`
		RecordsTotal    int         `json:"recordsTotal" xml:"recordsTotal"`
		RecordsFiltered int         `json:"recordsFiltered" xml:"recordsFiltered"`
	}
)
