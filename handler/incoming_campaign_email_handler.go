package handler

import (
	"net/mail"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource/config"
	"github.com/infraLinkit/mediaplatform-datasource/entity"
)

// CampaignEmailTemplates gives the compose dialog what CampaignManagementController::show passes to the view.
func (h *IncomingHandler) CampaignEmailTemplates(c *fiber.Ctx) error {
	templates := []entity.Email{}
	if err := h.DB.Find(&templates).Error; err != nil {
		h.Logs.Errorf("CampaignEmailTemplates: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch email templates"})
	}
	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithData{Code: fiber.StatusOK, Message: config.OK_DESC, Data: templates})
}

// validMailAddresses keeps what filter_var(FILTER_VALIDATE_EMAIL) would: bare addresses without a display name.
func validMailAddresses(candidates []string) []string {
	valid := []string{}
	for _, candidate := range candidates {
		address, err := mail.ParseAddress(candidate)
		if err == nil && address.Address == candidate {
			valid = append(valid, candidate)
		}
	}
	return valid
}

// SendCampaignEmail is CampaignManagementController::sendEmail: sent as the logged-in user, with their CC list.
func (h *IncomingHandler) SendCampaignEmail(c *fiber.Ctx) error {
	var req struct {
		Subject string   `json:"subject"`
		Body    string   `json:"body"`
		To      []string `json:"to"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}
	var errs []string
	if strings.TrimSpace(req.Subject) == "" {
		errs = append(errs, "The subject field is required.")
	}
	if strings.TrimSpace(req.Body) == "" {
		errs = append(errs, "The body field is required.")
	}
	if len(req.To) == 0 {
		errs = append(errs, "The to field is required.")
	}
	if len(errs) > 0 {
		return validationFailed(c, errs)
	}
	to := validMailAddresses(req.To)
	if len(to) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "No valid email addresses provided."})
	}

	message := mailMessage{To: to, Subject: req.Subject, HTML: strings.ReplaceAll(req.Body, "<b>Body Email :</b><br>", "")}
	userID, identified := c.Locals("user_id").(int)
	if identified {
		var user entity.User
		if err := h.DB.First(&user, userID).Error; err != nil {
			h.Logs.Errorf("SendCampaignEmail: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to send email."})
		}
		name := user.Name
		if name == "" {
			name = user.Username
		}
		message.FromAddress, message.FromName = user.Email, name
		message.ReplyTo = &mail.Address{Name: name, Address: user.Email}
		h.DB.Table("cc_emails").Where("user_id = ?", userID).Pluck("email", &message.Cc)
	}
	if err := sendMailMessage(message); err != nil {
		h.Logs.Errorf("SendCampaignEmail: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to send email."})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "message": "Email sent successfully."})
}
