package model

import (
	"fmt"
	"regexp"

	"github.com/infraLinkit/mediaplatform-datasource/entity"
	"google.golang.org/api/sheets/v4"
)

type StatusData struct {
	Status       string
	StatusCode   string
	StatusDetail string
}

func (r *BaseModel) UpdateGoogleSheetPixel(GS *sheets.Service, ps entity.PixelStorage, s StatusData) {
	sheetId, err := GetSpreadsheetID(ps.GoogleSheet)
	if err != nil {
		r.Logs.Info(fmt.Sprintf("Google sheet link not valid for campaign ID:  %#v\n", ps.CampaignId))
		r.Logs.Info(fmt.Sprintf("Google sheet link :  %#v ", ps.GoogleSheet))
	}

	/* prop, err := GS.Spreadsheets.Get(sheetId).Fields("properties.title").Context(context.Background()).Do()
	if err != nil {
		Logs.Error(fmt.Sprintf("Failed to read title: %#v\n", err))
	} */

	resp, err := GS.Spreadsheets.Values.Get(sheetId, "Sheet1!A1:E8").Do()
	if err != nil {
		r.Logs.Error(fmt.Sprintf("Failed to read sheet: %#v\n", err))
	}

	if len(resp.Values) < 7 {
		header := &sheets.ValueRange{
			Range: "Sheet1!A1:E8",
			Values: [][]interface{}{
				{"#### INSTRUCTIONS ####"},
				{"# IMPORTANT: Remember to set the TimeZone value in the \"parameters\" row and/or in your Conversion Time column"},
				{"# For instructions on how to setup your data, visit http://goo.gl/T1C5Ov"},
				{}, // empty row
				{"#### TEMPLATE ####"},
				{"Parameters: TimeZone=+0700"},
				{"Google Click ID", "Time", "MSISDN", "Status", "StatusCode", "StatusDetail"},
			},
		}
		_, err := GS.Spreadsheets.Values.Update(sheetId, "Sheet1!A1:E8", header).ValueInputOption("RAW").Do()
		if err != nil {
			r.Logs.Error(fmt.Sprintf("Failed to insert header: %#v\n", err))
		}
	}

	values := &sheets.ValueRange{
		Values: [][]interface{}{{
			ps.Pixel,
			ps.PixelUsedDate,
			ps.Msisdn,
			s.Status,
			s.StatusCode,
			s.StatusDetail,
		}},
	}
	_, err = GS.Spreadsheets.Values.Append(sheetId, "Sheet1!A:E", values).ValueInputOption("USER_ENTERED").Do()
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
