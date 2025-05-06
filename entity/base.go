package entity

import "encoding/json"

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

	GlobalRequestFromDataTable struct {
		Page     int    `json:"page" xml:"page"`
		Action   string `json:"action" xml:"action"`
		Draw     int    `json:"draw" xml:"draw"`
		PageSize int    `json:"pageSize" xml:"pageSize"`
		Search   string `json:"search" xml:"search"`
	}

	Ack struct {
		Keyword   string `json:"keyword" xml:"keyword"`
		ShortCode string `json:"shortcode" xml:"shortcode"`
		ClickId   string `json:"clickid" xml:"clickid"`
	}
)

func EncodeJsonAck(obj Ack) []byte {
	ack, _ := json.Marshal(obj)

	return ack
}
