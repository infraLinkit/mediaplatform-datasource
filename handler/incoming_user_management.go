package handler

import (
	"bytes"
	"errors"
	"html/template"
	"net/mail"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource/config"
	"github.com/infraLinkit/mediaplatform-datasource/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// userAssignments carries the four switch lists of the user forms; each posts every switch on or off.
type userAssignments struct {
	Services  []entity.DetailUser  `json:"services"`
	Adnets    []entity.UserAdnet   `json:"adnets"`
	Companies []entity.UserCompany `json:"companies"`
	Agencies  []entity.UserAgency  `json:"agencies"`
}

type userRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleID   int    `json:"role_id"`
	userAssignments
}

func (a userAssignments) save(tx *gorm.DB, userID int) error {
	services := make([]assignmentRow, 0, len(a.Services))
	for _, row := range a.Services {
		services = append(services, assignmentRow{Key: []int{row.CountryID, row.OperatorID, row.ServiceID}, Status: row.Status})
	}
	if err := syncAssignments(tx, "detail_users", "user_id", userID, []string{"country_id", "operator_id", "service_id"}, services); err != nil {
		return err
	}
	adnets := make([]assignmentRow, 0, len(a.Adnets))
	for _, row := range a.Adnets {
		adnets = append(adnets, assignmentRow{Key: []int{row.AdnetID}, Status: row.Status})
	}
	if err := syncAssignments(tx, "user_adnets", "user_id", userID, []string{"adnet_id"}, adnets); err != nil {
		return err
	}
	companies := make([]assignmentRow, 0, len(a.Companies))
	for _, row := range a.Companies {
		companies = append(companies, assignmentRow{Key: []int{row.CompanyID}, Status: row.Status})
	}
	if err := syncAssignments(tx, "user_companies", "user_id", userID, []string{"company_id"}, companies); err != nil {
		return err
	}
	agencies := make([]assignmentRow, 0, len(a.Agencies))
	for _, row := range a.Agencies {
		agencies = append(agencies, assignmentRow{Key: []int{row.AgencyID}, Status: row.Status})
	}
	return syncAssignments(tx, "user_agencies", "user_id", userID, []string{"agency_id"}, agencies)
}

// userErrors applies UserManagementController's validation rules; only create requires a password.
func (h *IncomingHandler) userErrors(req *userRequest, exceptID int, passwordRequired bool) ([]string, error) {
	req.Name, req.Username, req.Email = strings.TrimSpace(req.Name), strings.TrimSpace(req.Username), strings.TrimSpace(req.Email)
	var errs []string
	if req.Name == "" {
		errs = append(errs, "The name field is required.")
	}
	if req.Username == "" {
		errs = append(errs, "The username field is required.")
	} else {
		taken, err := h.valueTaken("users", "username", req.Username, exceptID)
		if err != nil {
			return nil, err
		}
		if taken {
			errs = append(errs, "The username has already been taken.")
		}
	}
	if req.Email == "" {
		errs = append(errs, "The email field is required.")
	} else if address, err := mail.ParseAddress(req.Email); err != nil || address.Address != req.Email {
		errs = append(errs, "The email field must be a valid email address.")
	} else {
		taken, err := h.valueTaken("users", "email", req.Email, exceptID)
		if err != nil {
			return nil, err
		}
		if taken {
			errs = append(errs, "The email has already been taken.")
		}
	}
	if passwordRequired && req.Password == "" {
		errs = append(errs, "The password field is required.")
	}
	if req.RoleID == 0 {
		errs = append(errs, "The role id field is required.")
	}
	return errs, nil
}

func (h *IncomingHandler) CreateUser(c *fiber.Ctx) error {
	var req userRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}
	errs, err := h.userErrors(&req, 0, true)
	if err != nil {
		h.Logs.Errorf("CreateUser: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to create user"})
	}
	if len(errs) > 0 {
		return validationFailed(c, errs)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to hash password"})
	}

	err = h.DB.Transaction(func(tx *gorm.DB) error {
		user := entity.User{
			Name:     req.Name,
			Username: req.Username,
			Email:    req.Email,
			Password: string(hashedPassword),
			RoleID:   req.RoleID,
			TypeUser: "internal",
			IsVerify: true,
			Status:   true,
		}
		// Only the columns UserRepository::insert writes; the rest stay NULL.
		if err := tx.Select("name", "username", "email", "password", "role_id", "type_user", "is_verify", "status", "created_at", "updated_at").Create(&user).Error; err != nil {
			return err
		}
		return req.save(tx, user.ID)
	})
	if err != nil {
		h.Logs.Errorf("CreateUser: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to create user"})
	}
	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

// GetUserTable always reads the table: Blade reloads the list right after an action, so a cached copy would hide it.
func (h *IncomingHandler) GetUserTable(c *fiber.Ctx) error {
	users, err := h.DS.GetAllUserWithRelation()
	if err != nil {
		h.Logs.Errorf("GetUserTable: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch user"})
	}
	if users == nil {
		users = []entity.UserManagementData{}
	}
	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithData{Code: fiber.StatusOK, Message: config.OK_DESC, Data: users})
}

func (h *IncomingHandler) GetUserCounts(c *fiber.Ctx) error {
	counts, err := h.DS.GetUserCounts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get user counts",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(counts)
}

func (h *IncomingHandler) UpdateUser(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user ID"})
	}
	var req userRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}
	if _, err := h.DS.GetUserByID(userID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
	}
	errs, err := h.userErrors(&req, userID, false)
	if err != nil {
		h.Logs.Errorf("UpdateUser: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update user"})
	}
	if len(errs) > 0 {
		return validationFailed(c, errs)
	}

	values := map[string]interface{}{
		"name":       req.Name,
		"username":   req.Username,
		"email":      req.Email,
		"role_id":    req.RoleID,
		"updated_at": time.Now(),
	}
	// The password only changes when one was typed.
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to hash password"})
		}
		values["password"] = string(hashedPassword)
	}

	err = h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entity.User{}).Where("id = ?", userID).Updates(values).Error; err != nil {
			return err
		}
		return req.save(tx, userID)
	})
	if err != nil {
		h.Logs.Errorf("UpdateUser: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update user"})
	}
	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) AssignService(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user ID"})
	}
	var req userAssignments
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}
	if _, err := h.DS.GetUserByID(userID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
	}
	err = h.DB.Transaction(func(tx *gorm.DB) error {
		return req.save(tx, userID)
	})
	if err != nil {
		h.Logs.Errorf("AssignService: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to assign service"})
	}
	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) GetUserOptions(c *fiber.Ctx) error {
	options, err := h.DS.GetUserOptions()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch user options",
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithData{
		Code:    fiber.StatusOK,
		Message: config.OK_DESC,
		Data:    options,
	})
}

func (h *IncomingHandler) GetUserDetail(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid user ID",
		})
	}

	detail, err := h.DS.GetUserDetail(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithData{
		Code:    fiber.StatusOK,
		Message: config.OK_DESC,
		Data:    detail,
	})
}

func (h *IncomingHandler) UpdateUserStatus(c *fiber.Ctx) error {
	idParam := c.Params("id")

	userID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid user ID",
		})
	}

	var req struct {
		Status bool `json:"status" binding:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation error",
			"errors":  err.Error(),
		})
	}

	_, err = h.DS.UpdateUserStatus(userID, req.Status)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	statusMessage := "User is now inactive"
	if req.Status {
		statusMessage = "User is now active"
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": statusMessage,
	})
}

func (h *IncomingHandler) DeleteUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid user ID",
		})
	}

	user, err := h.DS.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	if err := h.DS.DeleteUser(user.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

// GetUserApplovalRequestTable reads the table for the same reason as GetUserTable.
func (h *IncomingHandler) GetUserApplovalRequestTable(c *fiber.Ctx) error {
	users, err := h.DS.GetAllUserApprovalRequest()
	if err != nil {
		h.Logs.Errorf("GetUserApplovalRequestTable: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch user"})
	}
	if users == nil {
		users = []entity.UserApprovalRequestData{}
	}
	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithData{Code: fiber.StatusOK, Message: config.OK_DESC, Data: users})
}

var userApprovedMail = template.Must(template.New("user-approved").Parse(`<html>
<body>
    <p>Dear {{.Name}},</p>
    <p>Your account has been approved! You can now log in and start using our services.</p>
    <p>
        <a href="{{.Login}}">Click here to log in to your account</a>
    </p>
    <p>Thank you!</p>
</body>
</html>`))

func sendUserApprovedMail(email, name, frontendURL string) error {
	login := frontendURL + "/#/login"
	var html bytes.Buffer
	if err := userApprovedMail.Execute(&html, struct{ Name, Login string }{name, login}); err != nil {
		return err
	}
	text := "Dear " + name + ",\n\nYour account has been approved! You can now log in and start using our services.\n\n" + login + "\n\nThank you!\n"
	return sendMail(email, "Your Account Has Been Approved", text, html.String())
}

// ApproveUser follows UserRepository::approveUser and sends UserApproved inside the transaction,
// so a mail failure leaves the request pending, as it does in Laravel.
func (h *IncomingHandler) ApproveUser(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user ID"})
	}
	var req struct {
		RoleID int `json:"role_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}
	if req.RoleID == 0 {
		return validationFailed(c, []string{"The role id field is required."})
	}
	user, err := h.DS.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
	}
	// From config only: a link built from the Host header could point the user at an attacker's site.
	frontendURL := strings.TrimRight(os.Getenv("CMS_FRONTEND_URL"), "/")
	if frontendURL == "" {
		h.Logs.Error("ApproveUser: CMS_FRONTEND_URL is not set")
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"message": "Approval email is not configured"})
	}

	verifiedBy := ""
	if approver, ok := c.Locals("user_id").(int); ok {
		verifiedBy = strconv.Itoa(approver)
	}
	err = h.DB.Transaction(func(tx *gorm.DB) error {
		values := map[string]interface{}{
			"role_id":    req.RoleID,
			"verify_at":  time.Now(),
			"is_verify":  true,
			"status":     true,
			"updated_at": time.Now(),
		}
		if verifiedBy != "" {
			values["verify_by"] = verifiedBy
		}
		result := tx.Model(&entity.User{}).Where("id = ?", userID).Updates(values)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("user not found")
		}
		return sendUserApprovedMail(user.Email, user.Name, frontendURL)
	})
	if err != nil {
		h.Logs.Errorf("ApproveUser: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to approve user"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User approved successfully."})
}
