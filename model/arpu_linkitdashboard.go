package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"encoding/base64"

	"github.com/infraLinkit/mediaplatform-datasource/entity"
)

// GetDataArpu method untuk BaseModel
func (r *BaseModel) GetDataArpu(fe entity.ArpuParams) error {
	// Menggunakan environment variable APIARPU dari config
	baseURL := r.Config.APIARPU
	if baseURL == "" {
		return errors.New("APIARPU environment variable is empty")
	}

	// Membuat URL untuk request API
	apiURL, err := url.Parse(baseURL + "/api/v4/arpu/arpu90")
	if err != nil {
		return fmt.Errorf("failed to parse base URL: %v", err)
	}

	// Menambahkan query parameters
	query := url.Values{}
	if fe.From != "" {
		query.Set("from", fe.From)
	}
	if fe.To != "" {
		query.Set("to", fe.To)
	}
	if fe.Country != "" {
		query.Set("country", fe.Country)
	}
	if fe.Operator != "" {
		query.Set("operator", fe.Operator)
	}
	if fe.Service != "" {
		query.Set("service", fe.Service)
	}

	apiURL.RawQuery = query.Encode()

	// Membuat HTTP request
	req, err := http.NewRequest("GET", apiURL.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Menambahkan headers yang diperlukan
	encUsername := r.Config.ARPUUsername
	encPassword := r.Config.ARPUPassword

	username, err := decryptEnv(encUsername)
	if err != nil {
		return fmt.Errorf("failed to decrypt ARPU_USERNAME: %v", err)
	}
	password, err := decryptEnv(encPassword)
	if err != nil {
		return fmt.Errorf("failed to decrypt ARPU_PASSWORD: %v", err)
	}

	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Accept", "application/json")

	// Melakukan request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Membaca response
	var apiResponse entity.ARPUResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	// Validasi response
	if apiResponse.Status != 200 {
		return fmt.Errorf("API returned status %d: %s", apiResponse.Status, apiResponse.Message)
	}

	r.Logs.Info(fmt.Sprintf("Successfully retrieved ARPU data for %s/%s/%s", fe.Country, fe.Operator, fe.Service))
	return nil
}

func decryptEnv(enc string) (string, error) {
	// contoh: jika hanya base64
	data, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (r *BaseModel) SendWakiCallback() error {
	var summaries []entity.SummaryCampaign

	baseURL := r.Config.APILINKITDashboard
	if baseURL == "" {
		return errors.New("APILINKITDashboard environment variable is empty")
	}

	// Ambil semua summary_campaigns untuk hari ini yg mo_received > 0
	if err := r.DB.
		Where("summary_date = DATE(NOW()) AND mo_received > 0 AND deleted_at IS NULL").
		Find(&summaries).Error; err != nil {
		return err
	}

	for _, sc := range summaries {
		// Bangun query URL
		q := url.Values{
			"date":           {time.Now().Format("2006-01-02")},
			"publisher":      {sc.Adnet},
			"adnet":          {sc.Adnet},
			"operator":       {sc.Partner},
			"adn":            {sc.ShortCode},
			"client":         {sc.Partner},
			"aggregator":     {sc.Aggregator},
			"country":        {sc.Country},
			"service":        {sc.Service},
			"mo_received":    {strconv.Itoa(sc.MoReceived)},
			"mo_postback":    {strconv.Itoa(sc.Postback)},
			"total_mo":       {strconv.Itoa(sc.MoReceived)},
			"total_postback": {strconv.Itoa(sc.Postback)},
			"landing":        {strconv.Itoa(sc.Traffic)},
			"cr_mo_received": {strconv.FormatFloat(sc.CrMO, 'f', 2, 64)},
			"cr_mo_postback": {strconv.FormatFloat(sc.CrPostback, 'f', 2, 64)},
			"url_campaign":   {sc.URLAfter},
			"url_service":    {sc.URLBefore},
			"sbaf":           {strconv.FormatFloat(sc.SBAF, 'f', 2, 64)},
			"saaf":           {strconv.FormatFloat(sc.SAAF, 'f', 2, 64)},
			"spending":       {strconv.FormatFloat(sc.SAAF, 'f', 2, 64)},
			"campaign":       {sc.CampaignObjective},
			"payout":         {strconv.FormatFloat(sc.PO, 'f', 2, 64)},
			"price_per_mo":   {strconv.FormatFloat(sc.PricePerMO, 'f', 2, 64)},
		}

		// Gabungkan URL dan query param
		fullURL := fmt.Sprintf("%s?%s", baseURL, q.Encode())

		// Kirim request
		resp, err := http.Get(fullURL)
		if err != nil {
			return fmt.Errorf("failed to send request for campaign %s: %v", sc.CampaignId, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("API returned status %d for campaign %s", resp.StatusCode, sc.CampaignId)
		}

		log.Printf("✅ Sent to LinkIT: %s", fullURL)
	}

	return nil
}
