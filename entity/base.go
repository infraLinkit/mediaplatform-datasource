package entity

type (
	GlobalResponse struct {
		Code    int    `json:"code" xml:"code"`
		Message string `json:"message" xml:"message"`
	}

	ReturnResponse struct {
		HttpStatus int
		Rsp        interface{}
	}
)
