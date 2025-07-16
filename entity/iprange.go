package entity

type (
	ImplementIPRangeRequest struct {
		IPType      string `json:"ip_type"`
		UploadMonth string `json:"upload_month"`
	}

	ResultIPRange struct {
		IPType   string `json:"ip_type"`
		Month    string `json:"month"`
		Filename string `json:"file_name"`
	}

	DownloadReq struct {
		IPType string `json:"ip_type"`
		Month  string `json:"month"`
	}
)
