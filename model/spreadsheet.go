package model

import (
	"fmt"
	"regexp"

	"github.com/infraLinkit/mediaplatform-datasource/entity"
	"google.golang.org/api/sheets/v4"
)

func (r *BaseModel) UpdateGoogleSheetPixel(ps entity.PixelStorage) {
	sheetId, err := GetSpreadsheetID(ps.GoogleSheet)
	if err == nil {
		values := &sheets.ValueRange{
			Values: [][]interface{}{{
				ps.Pixel,
				ps.CampaignName,
				ps.PixelUsedDate,
				1,
				ps.Currency,
			}},
		}
		_, err := r.GS.Spreadsheets.Values.Append(sheetId, "Sheet1!A:E", values).ValueInputOption("USER_ENTERED").Do()

		if err != nil {
			r.Logs.Error(fmt.Sprintf("Google sheet input failed error:  %#v\n", err))
		}
	} else {
		r.Logs.Info(fmt.Sprintf("Google sheet link not valid for campaign ID:  %#v\n", ps.CampaignId))
		r.Logs.Info(fmt.Sprintf("Google sheet link :  %#v ", ps.GoogleSheet))
	}

}

func GetSpreadsheetID(url string) (string, error) {
	re := regexp.MustCompile(`https://docs\.google\.com/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	matches := re.FindStringSubmatch(url)

	if len(matches) < 2 {
		return "", fmt.Errorf("spreadsheet ID not found in URL")
	}

	return matches[1], nil
}
