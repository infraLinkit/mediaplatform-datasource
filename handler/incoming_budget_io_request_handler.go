package handler

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/mail"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource/config"
	"github.com/infraLinkit/mediaplatform-datasource/entity"
	"gorm.io/gorm"
)

const (
	budgetIOSignatureDir = "storage/budget_io/signature"
	budgetIOMaxSignature = 2048 * 1024
)

// budgetIOOrderColumns are the only columns a list may sort on; the value goes into the SQL text.
var budgetIOOrderColumns = map[string]bool{
	"id": true, "io_id": true, "campaign_type": true, "month": true, "country": true, "continent": true,
	"company_group_name": true, "service": true, "submitted_by": true, "target_cac": true, "target_roi": true,
	"monthly_mo_target": true, "monthly_spend_target": true, "created_at": true, "status": true,
}

func budgetIOOrder(column, dir string) string {
	if !budgetIOOrderColumns[column] {
		return ""
	}
	if strings.ToUpper(dir) == "DESC" {
		return column + " DESC"
	}
	return column + " ASC"
}

func (h *IncomingHandler) budgetIORequests(c *fiber.Ctx, approved, allCompanies bool) error {
	query := h.DB.Model(&entity.BudgetIORequest{})
	if !allCompanies {
		allowed, _ := c.Locals("companies").([]string)
		query = query.Where("company IN ?", allowed)
	}
	if approved {
		query = query.Where("status = ?", "approved")
	} else {
		query = query.Where("status <> ?", "approved")
	}
	for _, column := range []string{"campaign_type", "country", "partner", "continent"} {
		if value := c.Query(column); value != "" {
			query = query.Where(column+" = ?", value)
		}
	}

	if approved {
		if keyword := c.Query("keyword"); keyword != "" {
			pattern := "%" + keyword + "%"
			query = query.Where("io_id ILIKE ? OR company_group_name ILIKE ? OR submitted_by ILIKE ?", pattern, pattern, pattern)
		}
		if month := c.Query("date_range"); month != "" {
			query = query.Where("month = ?", month)
		} else {
			query = query.Where("month = TO_CHAR(CURRENT_DATE, 'YYYY-MM')")
		}
	} else {
		switch dateRange := strings.ToUpper(c.Query("date_range")); dateRange {
		case "LASTMONTH":
			query = query.Where("DATE(created_at) >= DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE(created_at) < DATE_TRUNC('month', CURRENT_DATE)")
		case "CUSTOMRANGE":
			query = query.Where("DATE(created_at) BETWEEN ? AND ?", c.Query("date_before"), c.Query("date_after"))
		default:
			query = query.Where("DATE(created_at) >= DATE_TRUNC('month', CURRENT_DATE)")
		}
	}

	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		h.Logs.Errorf("budgetIORequests: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch budget IO"})
	}

	if order := budgetIOOrder(c.Query("order_column"), c.Query("order_dir")); order != "" {
		query = query.Order(order)
	} else {
		query = query.Order("DATE(created_at) DESC").Order("id DESC")
	}
	pageSize, err := strconv.Atoi(c.Query("page_size"))
	if err != nil || pageSize == 0 {
		pageSize = PAGESIZE
	}
	page, _ := strconv.Atoi(c.Query("page"))
	if page > 1 && pageSize > 0 {
		query = query.Offset((page - 1) * pageSize)
	}

	rows := []entity.BudgetIORequest{}
	if err := query.Limit(pageSize).Find(&rows).Error; err != nil {
		h.Logs.Errorf("budgetIORequests: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch budget IO"})
	}
	draw, _ := strconv.Atoi(c.Query("draw"))
	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithDataTable{
		Code:            fiber.StatusOK,
		Message:         config.OK_DESC,
		Data:            rows,
		Draw:            draw,
		RecordsTotal:    int(total),
		RecordsFiltered: int(total),
	})
}

func (h *IncomingHandler) ListBudgetIORequests(c *fiber.Ctx) error {
	return h.budgetIORequests(c, false, false)
}

func (h *IncomingHandler) ListAllBudgetIORequests(c *fiber.Ctx) error {
	return h.budgetIORequests(c, false, true)
}

func (h *IncomingHandler) ListApprovedBudgetIORequests(c *fiber.Ctx) error {
	return h.budgetIORequests(c, true, false)
}

func (h *IncomingHandler) ListAllApprovedBudgetIORequests(c *fiber.Ctx) error {
	return h.budgetIORequests(c, true, true)
}

func (h *IncomingHandler) GetBudgetIORequest(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid ID"})
	}
	var row entity.BudgetIORequest
	if err := h.DB.First(&row, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Budget IO not found"})
	}
	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithData{Code: fiber.StatusOK, Message: config.OK_DESC, Data: row})
}

func (h *IncomingHandler) budgetIOCurrentUserName(c *fiber.Ctx) string {
	userID, _ := c.Locals("user_id").(int)
	var name string
	h.DB.Model(&entity.User{}).Where("id = ?", userID).Pluck("name", &name)
	return name
}

// ApproveBudgetIORequest only changes the request's status: IO Report targets stay a manual edit.
func (h *IncomingHandler) ApproveBudgetIORequest(c *fiber.Ctx) error {
	return h.setBudgetIOStatus(c, map[string]interface{}{
		"status":      "approved",
		"verified_at": time.Now(),
		"approved_by": h.budgetIOCurrentUserName(c),
		"updated_at":  time.Now(),
	}, "Approved")
}

func (h *IncomingHandler) RejectBudgetIORequest(c *fiber.Ctx) error {
	return h.setBudgetIOStatus(c, map[string]interface{}{"status": "rejected", "updated_at": time.Now()}, "Rejected")
}

func (h *IncomingHandler) setBudgetIOStatus(c *fiber.Ctx, values map[string]interface{}, done string) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid ID"})
	}
	result := h.DB.Model(&entity.BudgetIORequest{}).Where("id = ?", id).Updates(values)
	if result.Error != nil {
		h.Logs.Errorf("setBudgetIOStatus: %v", result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update budget IO"})
	}
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Budget IO not found"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": done})
}

// budgetIOForm is the JSON "payload" field of the multipart create and revise requests.
// Numbers stay text until validated, as they do when Laravel reads form fields.
type budgetIOForm struct {
	CPName             string            `json:"cp_name"`
	PICName            string            `json:"pic_name"`
	ContactEmail       string            `json:"contact_email"`
	CPBusinessPICName  string            `json:"cp_business_pic_name"`
	TargetCAC          string            `json:"target_cac"`
	TargetROI          string            `json:"target_roi"`
	MonthlyMOTarget    string            `json:"monthly_mo_target"`
	MonthlySpendTarget string            `json:"monthly_spend_target"`
	Data               []budgetIOFormRow `json:"data"`
}

type budgetIOFormRow struct {
	IOID               string `json:"io_id"`
	CampaignType       string `json:"campaign_type"`
	Month              string `json:"month"`
	Country            string `json:"country"`
	CountryName        string `json:"country_name"`
	Continent          string `json:"continent"`
	CompanyGroupName   string `json:"company_group_name"`
	Company            string `json:"company"`
	Service            string `json:"service"`
	Partner            string `json:"partner"`
	TargetCAC          string `json:"target_cac"`
	TargetROI          string `json:"target_roi"`
	MonthlyMOTarget    string `json:"monthly_mo_target"`
	MonthlySpendTarget string `json:"monthly_spend_target"`
}

func budgetIOHeaderErrors(f budgetIOForm) []string {
	var errs []string
	text := func(name, value string) {
		if strings.TrimSpace(value) == "" {
			errs = append(errs, "The "+name+" field is required.")
		} else if len([]rune(value)) > 255 {
			errs = append(errs, "The "+name+" field must not be greater than 255 characters.")
		}
	}
	text("cp name", f.CPName)
	text("pic name", f.PICName)
	if strings.TrimSpace(f.ContactEmail) == "" {
		errs = append(errs, "The contact email field is required.")
	} else if address, err := mail.ParseAddress(f.ContactEmail); err != nil || address.Address != f.ContactEmail {
		errs = append(errs, "The contact email field must be a valid email address.")
	}
	text("cp business pic name", f.CPBusinessPICName)
	return errs
}

func budgetIONumberErrors(prefix string, fields [][2]string) []string {
	var errs []string
	for _, field := range fields {
		name, value := prefix+field[0], strings.TrimSpace(field[1])
		if value == "" {
			errs = append(errs, "The "+name+" field is required.")
		} else if _, err := strconv.ParseFloat(value, 64); err != nil {
			errs = append(errs, "The "+name+" field must be a number.")
		}
	}
	return errs
}

// budgetIORowErrors applies the store rules to the table rows, numbered from 1 as the form posts them.
func budgetIORowErrors(rows []budgetIOFormRow) []string {
	if len(rows) == 0 {
		return []string{"The data field is required."}
	}
	var errs []string
	for i, row := range rows {
		prefix := fmt.Sprintf("data.%d.", i+1)
		for _, field := range [][2]string{
			{"io_id", row.IOID}, {"campaign_type", row.CampaignType}, {"month", row.Month}, {"country", row.Country},
			{"country_name", row.CountryName}, {"continent", row.Continent}, {"company_group_name", row.CompanyGroupName},
		} {
			if strings.TrimSpace(field[1]) == "" {
				errs = append(errs, "The "+prefix+field[0]+" field is required.")
			}
		}
		errs = append(errs, budgetIONumberErrors(prefix, [][2]string{
			{"target_cac", row.TargetCAC}, {"target_roi", row.TargetROI},
			{"monthly_mo_target", row.MonthlyMOTarget}, {"monthly_spend_target", row.MonthlySpendTarget},
		})...)
	}
	return errs
}

func budgetIONumber(value string) float64 {
	number, _ := strconv.ParseFloat(strings.TrimSpace(value), 64)
	return number
}

// budgetIOSignatureError checks the uploaded bytes, not the file name or the declared type.
func budgetIOSignatureError(size int64, head []byte) string {
	switch http.DetectContentType(head) {
	case "image/png", "image/jpeg":
	default:
		return "The signature field must be a file of type: png, jpg, jpeg."
	}
	if size > budgetIOMaxSignature {
		return "The signature field must not be greater than 2048 kilobytes."
	}
	return ""
}

func checkBudgetIOSignature(file *multipart.FileHeader) (string, error) {
	opened, err := file.Open()
	if err != nil {
		return "", err
	}
	defer opened.Close()
	head := make([]byte, 512)
	n, err := opened.Read(head)
	if err != nil {
		return "", err
	}
	return budgetIOSignatureError(file.Size, head[:n]), nil
}

// saveBudgetIOSignature names the file itself, so a client file name can never reach the path.
func saveBudgetIOSignature(c *fiber.Ctx, file *multipart.FileHeader) (string, error) {
	if err := os.MkdirAll(budgetIOSignatureDir, 0755); err != nil {
		return "", err
	}
	ext := ".jpg"
	if strings.EqualFold(filepath.Ext(file.Filename), ".png") {
		ext = ".png"
	}
	name := fmt.Sprintf("%d-%s%s", time.Now().Unix(), randomHex(8), ext)
	if err := c.SaveFile(file, filepath.Join(budgetIOSignatureDir, name)); err != nil {
		return "", err
	}
	return "/" + budgetIOSignatureDir + "/" + name, nil
}

func removeBudgetIOSignature(url string) {
	if url != "" {
		os.Remove(filepath.Join(budgetIOSignatureDir, filepath.Base(url)))
	}
}

func parseBudgetIOForm(c *fiber.Ctx) (budgetIOForm, error) {
	var form budgetIOForm
	err := json.Unmarshal([]byte(c.FormValue("payload")), &form)
	return form, err
}

func (h *IncomingHandler) CreateBudgetIORequest(c *fiber.Ctx) error {
	form, err := parseBudgetIOForm(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}
	errs := budgetIOHeaderErrors(form)
	file, fileErr := c.FormFile("signature")
	if fileErr != nil {
		errs = append(errs, "The signature field is required.")
	} else {
		message, err := checkBudgetIOSignature(file)
		if err != nil {
			errs = append(errs, "The signature failed to upload.")
		} else if message != "" {
			errs = append(errs, message)
		}
	}
	errs = append(errs, budgetIORowErrors(form.Data)...)
	if len(errs) > 0 {
		return validationFailed(c, errs)
	}

	signature, err := saveBudgetIOSignature(c, file)
	if err != nil {
		h.Logs.Errorf("CreateBudgetIORequest: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to create new budget io"})
	}
	submittedBy := h.budgetIOCurrentUserName(c)
	rows := make([]entity.BudgetIORequest, 0, len(form.Data))
	for _, row := range form.Data {
		rows = append(rows, entity.BudgetIORequest{
			Status:             "submitted",
			SubmittedBy:        submittedBy,
			IOID:               row.IOID,
			CampaignType:       row.CampaignType,
			Month:              row.Month,
			Country:            row.Country,
			CountryName:        row.CountryName,
			Continent:          row.Continent,
			CompanyGroupName:   row.CompanyGroupName,
			Company:            row.Company,
			Partner:            row.Partner,
			Service:            row.Service,
			TargetCAC:          budgetIONumber(row.TargetCAC),
			TargetROI:          int(budgetIONumber(row.TargetROI)),
			MonthlyMOTarget:    budgetIONumber(row.MonthlyMOTarget),
			MonthlySpendTarget: budgetIONumber(row.MonthlySpendTarget),
			CPName:             form.CPName,
			PICName:            form.PICName,
			ContactEmail:       form.ContactEmail,
			CPBusinessPICName:  form.CPBusinessPICName,
			Signature:          signature,
		})
	}
	if err := h.DB.Create(&rows).Error; err != nil {
		removeBudgetIOSignature(signature)
		h.Logs.Errorf("CreateBudgetIORequest: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to create new budget io"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "message": "Budget IO created successfully"})
}

// UpdateBudgetIORequest is the revise save: the request goes back to submitted with a new submit date.
func (h *IncomingHandler) UpdateBudgetIORequest(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid ID"})
	}
	var existing entity.BudgetIORequest
	if err := h.DB.First(&existing, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Budget IO not found"})
	}
	form, err := parseBudgetIOForm(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}
	errs := budgetIOHeaderErrors(form)
	errs = append(errs, budgetIONumberErrors("", [][2]string{
		{"target cac", form.TargetCAC}, {"target roi", form.TargetROI},
		{"monthly mo target", form.MonthlyMOTarget}, {"monthly spend target", form.MonthlySpendTarget},
	})...)
	file, fileErr := c.FormFile("signature")
	if fileErr == nil {
		message, err := checkBudgetIOSignature(file)
		if err != nil {
			errs = append(errs, "The signature failed to upload.")
		} else if message != "" {
			errs = append(errs, message)
		}
	}
	if len(errs) > 0 {
		return validationFailed(c, errs)
	}

	now := time.Now()
	values := map[string]interface{}{
		"cp_name":              form.CPName,
		"pic_name":             form.PICName,
		"contact_email":        form.ContactEmail,
		"cp_business_pic_name": form.CPBusinessPICName,
		"target_cac":           budgetIONumber(form.TargetCAC),
		"target_roi":           int(budgetIONumber(form.TargetROI)),
		"monthly_mo_target":    budgetIONumber(form.MonthlyMOTarget),
		"monthly_spend_target": budgetIONumber(form.MonthlySpendTarget),
		"status":               "submitted",
		"submitted_by":         h.budgetIOCurrentUserName(c),
		"created_at":           now,
		"updated_at":           now,
	}
	newSignature := ""
	if fileErr == nil {
		newSignature, err = saveBudgetIOSignature(c, file)
		if err != nil {
			h.Logs.Errorf("UpdateBudgetIORequest: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update budget IO"})
		}
		values["signature"] = newSignature
	}
	if err := h.DB.Model(&entity.BudgetIORequest{}).Where("id = ?", id).Updates(values).Error; err != nil {
		removeBudgetIOSignature(newSignature)
		h.Logs.Errorf("UpdateBudgetIORequest: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update budget IO"})
	}
	// Rows from one form share a signature file, so it stays while another request still points at it.
	if newSignature != "" && existing.Signature != "" {
		var others int64
		h.DB.Model(&entity.BudgetIORequest{}).Where("signature = ?", existing.Signature).Count(&others)
		if others == 0 {
			removeBudgetIOSignature(existing.Signature)
		}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "message": "Budget IO updated successfully"})
}
