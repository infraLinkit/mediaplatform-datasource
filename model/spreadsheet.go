package model

import (
	"fmt"
	"regexp"

	"github.com/infraLinkit/mediaplatform-datasource/entity"
	"google.golang.org/api/sheets/v4"
)

func (r *BaseModel) UpdateGoogleSheetPixel(ps entity.PixelStorage) {
	sheetId, err := GetSpreadsheetID(ps.GoogleSheet)
	if err != nil {
		r.Logs.Info(fmt.Sprintf("Google sheet link not valid for campaign ID:  %#v\n", ps.CampaignId))
		r.Logs.Info(fmt.Sprintf("Google sheet link :  %#v ", ps.GoogleSheet))
		return
	}

	resp, err := r.GS.Spreadsheets.Values.Get(sheetId, "Sheet1!A1:E7").Do()
	if err != nil {
		r.Logs.Error(fmt.Sprintf("Failed to read sheet: %#v\n", err))
		return
	}

	if len(resp.Values) < 7 {
		header := &sheets.ValueRange{
			Range: "Sheet1!A1:E7",
			Values: [][]interface{}{
				{"#### INSTRUCTIONS ####"},
				{"# IMPORTANT: Remember to set the TimeZone value in the \"parameters\" row and/or in your Conversion Time column"},
				{"# For instructions on how to setup your data, visit http://goo.gl/T1C5Ov"},
				{}, // empty row
				{"#### TEMPLATE ####"},
				{"Parameters: TimeZone=+0700"},
				{"Google Click ID", "Conversion Name", "Conversion Time", "Conversion Value", "Conversion Currency"},
			},
		}
		_, err := r.GS.Spreadsheets.Values.Update(sheetId, "Sheet1!A1:E7", header).ValueInputOption("RAW").Do()
		if err != nil {
			r.Logs.Error(fmt.Sprintf("Failed to insert header: %#v\n", err))
			return
		}
	}

	values := &sheets.ValueRange{
		Values: [][]interface{}{{
			ps.Pixel,
			ps.CampaignName,
			ps.PixelUsedDate,
			1,
			ps.Currency,
		}},
	}
	_, err = r.GS.Spreadsheets.Values.Append(sheetId, "Sheet1!A:E", values).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		r.Logs.Error(fmt.Sprintf("Google sheet input failed error:  %#v\n", err))
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
