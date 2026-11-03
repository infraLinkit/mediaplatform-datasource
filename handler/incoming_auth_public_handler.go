package handler

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"mime"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// mailRootCAs stays nil in production, which means the system roots. Tests point it at a local CA.
var mailRootCAs *x509.CertPool

const (
	passwordMinLength = 8
	passwordMaxLength = 72 // bcrypt limit
)

type RegisterRequest struct {
	Name                 string `json:"name"`
	Username             string `json:"username"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
	RecaptchaToken       string `json:"recaptcha_token"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token                string `json:"token"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

func validationError(c *fiber.Ctx, field, message string) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
		"error":  message,
		"errors": fiber.Map{field: message},
	})
}

func passwordProblem(password, confirmation string) string {
	switch {
	case len(password) < passwordMinLength:
		return fmt.Sprintf("The password must be at least %d characters.", passwordMinLength)
	case len(password) > passwordMaxLength:
		return fmt.Sprintf("The password must not be greater than %d characters.", passwordMaxLength)
	case password != confirmation:
		return "The password confirmation does not match."
	}
	return ""
}

func validEmail(email string) bool {
	addr, err := mail.ParseAddress(email)
	return err == nil && addr.Address == email
}

// recaptchaEnabled mirrors services.recaptcha in Laravel: on when both keys exist unless disabled.
func recaptchaEnabled() bool {
	return !strings.EqualFold(os.Getenv("RECAPTCHA_ENABLED"), "false") &&
		os.Getenv("RECAPTCHA_SITE_KEY") != "" && os.Getenv("RECAPTCHA_SECRET_KEY") != ""
}

// verifyRecaptcha reads Google's success flag. Laravel only checked the HTTP status, which is 200
// for rejected tokens too.
func (h *IncomingHandler) verifyRecaptcha(token, ip string) error {
	if !recaptchaEnabled() {
		return nil
	}
	if token == "" {
		return errors.New("Please complete the reCAPTCHA verification to continue.")
	}
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.PostForm("https://www.google.com/recaptcha/api/siteverify", url.Values{
		"secret":   {os.Getenv("RECAPTCHA_SECRET_KEY")},
		"response": {token},
		"remoteip": {ip},
	})
	if err != nil {
		h.Logs.Errorf("reCAPTCHA request failed: %v", err)
		return errors.New("Recaptcha Failed to Proceed")
	}
	defer resp.Body.Close()
	var result struct {
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || !result.Success {
		return errors.New("Recaptcha Failed to Proceed")
	}
	return nil
}

// GetAuthConfig tells the public auth pages which widgets to render.
func (h *IncomingHandler) GetAuthConfig(c *fiber.Ctx) error {
	enabled := recaptchaEnabled()
	siteKey := ""
	if enabled {
		siteKey = os.Getenv("RECAPTCHA_SITE_KEY")
	}
	return c.JSON(fiber.Map{"recaptcha_enabled": enabled, "recaptcha_site_key": siteKey})
}

// Register mirrors RegisteredUserController: new users stay unverified until an admin approves them.
func (h *IncomingHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)

	switch {
	case req.Name == "" || len(req.Name) > 255:
		return validationError(c, "name", "The name field is required and must not exceed 255 characters.")
	case req.Username == "" || len(req.Username) > 255:
		return validationError(c, "username", "The username field is required and must not exceed 255 characters.")
	case len(req.Email) > 255 || !validEmail(req.Email):
		return validationError(c, "email", "The email field must be a valid email address.")
	}
	if problem := passwordProblem(req.Password, req.PasswordConfirmation); problem != "" {
		return validationError(c, "password", problem)
	}
	if err := h.verifyRecaptcha(req.RecaptchaToken, c.IP()); err != nil {
		return validationError(c, "recaptcha", err.Error())
	}

	// Unscoped: the unique index also covers soft-deleted users, as Laravel's unique rule does.
	var count int64
	h.DB.Unscoped().Model(&entity.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		return validationError(c, "username", "The username has already been taken.")
	}
	h.DB.Unscoped().Model(&entity.User{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		return validationError(c, "email", "The email has already been taken.")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.Logs.Errorf("Register: hash password: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create account"})
	}

	// A map insert leaves role_id NULL. The struct would send 0, which breaks the roles foreign key.
	now := time.Now()
	if err := h.DB.Model(&entity.User{}).Create(map[string]interface{}{
		"name":       req.Name,
		"username":   req.Username,
		"email":      req.Email,
		"password":   string(hashed),
		"type_user":  "internal",
		"is_verify":  false,
		"status":     true,
		"created_at": now,
		"updated_at": now,
	}).Error; err != nil {
		h.Logs.Errorf("Register: create user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create account"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"code": fiber.StatusCreated, "message": "Register was Successful"})
}

const forgotPasswordResponse = "If that email address is registered, we have emailed a password reset link."

// ForgotPassword answers the same for known and unknown emails, and mails outside the request,
// so neither the message nor the timing reveals which addresses have accounts.
func (h *IncomingHandler) ForgotPassword(c *fiber.Ctx) error {
	var req ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	req.Email = strings.TrimSpace(req.Email)
	if !validEmail(req.Email) {
		return validationError(c, "email", "The email field must be a valid email address.")
	}

	// From config only: a link built from the Host header could send the token to an attacker's site.
	frontendURL := strings.TrimRight(os.Getenv("CMS_FRONTEND_URL"), "/")
	if frontendURL == "" || os.Getenv("MAIL_HOST") == "" {
		h.Logs.Error("ForgotPassword: CMS_FRONTEND_URL or MAIL_HOST is not set")
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Password reset email is not available. Please contact the administrator."})
	}

	email := req.Email
	go func() {
		if err := h.sendPasswordResetLink(email, frontendURL); err != nil {
			h.Logs.Errorf("ForgotPassword: %v", err)
		}
	}()

	return c.JSON(fiber.Map{"code": fiber.StatusOK, "message": forgotPasswordResponse})
}

func (h *IncomingHandler) sendPasswordResetLink(email, frontendURL string) error {
	var user entity.User
	if err := h.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return fmt.Errorf("find user: %w", err)
	}

	// auth.passwords.users.throttle = 60 seconds
	var recent int64
	h.DB.Raw("SELECT COUNT(*) FROM password_reset_tokens WHERE email = ? AND created_at > NOW() - INTERVAL '60 seconds'", email).Scan(&recent)
	if recent > 0 {
		return nil
	}

	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return fmt.Errorf("generate token: %w", err)
	}
	token := hex.EncodeToString(raw)
	hashed, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash token: %w", err)
	}
	if err := h.DB.Exec(`INSERT INTO password_reset_tokens (email, token, created_at) VALUES (?, ?, NOW())
		ON CONFLICT (email) DO UPDATE SET token = EXCLUDED.token, created_at = EXCLUDED.created_at`, email, string(hashed)).Error; err != nil {
		return fmt.Errorf("store token: %w", err)
	}

	link := frontendURL + "/#/reset-password/" + token + "?email=" + url.QueryEscape(email)
	return sendPasswordResetMail(user.Name, email, link)
}

// ResetPassword mirrors NewPasswordController. Tokens last 60 minutes and work once.
func (h *IncomingHandler) ResetPassword(c *fiber.Ctx) error {
	var req ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	req.Email = strings.TrimSpace(req.Email)
	if !validEmail(req.Email) {
		return validationError(c, "email", "The email field must be a valid email address.")
	}
	if problem := passwordProblem(req.Password, req.PasswordConfirmation); problem != "" {
		return validationError(c, "password", problem)
	}

	const invalidToken = "This password reset token is invalid."

	var stored struct{ Token string }
	h.DB.Raw("SELECT token FROM password_reset_tokens WHERE email = ? AND created_at > NOW() - INTERVAL '60 minutes'", req.Email).Scan(&stored)
	if stored.Token == "" || bcrypt.CompareHashAndPassword([]byte(stored.Token), []byte(req.Token)) != nil {
		return validationError(c, "email", invalidToken)
	}

	var user entity.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return validationError(c, "email", invalidToken)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.Logs.Errorf("ResetPassword: hash password: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to reset password"})
	}
	remember := make([]byte, 30)
	if _, err := rand.Read(remember); err != nil {
		h.Logs.Errorf("ResetPassword: remember token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to reset password"})
	}

	err = h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE users SET password = ?, remember_token = ?, updated_at = ? WHERE id = ?",
			string(hashed), hex.EncodeToString(remember), time.Now(), user.ID).Error; err != nil {
			return err
		}
		return tx.Exec("DELETE FROM password_reset_tokens WHERE email = ?", req.Email).Error
	})
	if err != nil {
		h.Logs.Errorf("ResetPassword: save: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to reset password"})
	}

	return c.JSON(fiber.Map{"code": fiber.StatusOK, "message": "Your password has been reset."})
}

var passwordResetTemplate = template.Must(template.New("password-reset").Parse(`<!doctype html>
<html><body style="margin:0;padding:32px 16px;background:#f4f6f8;font-family:Arial,Helvetica,sans-serif;color:#2d3748">
<table role="presentation" width="100%" cellspacing="0" cellpadding="0"><tr><td align="center">
<table role="presentation" width="100%" cellspacing="0" cellpadding="0" style="max-width:520px;background:#ffffff;border-radius:8px">
<tr><td style="padding:32px">
<h1 style="margin:0 0 24px;font-size:20px;color:#00838f">Wakiads</h1>
<p style="margin:0 0 16px;font-size:15px;line-height:22px">Hello {{.Name}},</p>
<p style="margin:0 0 24px;font-size:15px;line-height:22px">You are receiving this email because we received a password reset request for your account.</p>
<p style="margin:0 0 24px"><a href="{{.Link}}" style="display:inline-block;padding:12px 24px;background:#00838f;color:#ffffff;border-radius:6px;font-size:15px;text-decoration:none">Reset Password</a></p>
<p style="margin:0 0 16px;font-size:14px;line-height:21px">This password reset link will expire in 60 minutes.</p>
<p style="margin:0;font-size:14px;line-height:21px">If you did not request a password reset, no further action is required.</p>
</td></tr></table>
</td></tr></table>
</body></html>`))

func sendPasswordResetMail(name, email, link string) error {
	var htmlBody bytes.Buffer
	if err := passwordResetTemplate.Execute(&htmlBody, struct {
		Name string
		Link template.URL
	}{name, template.URL(link)}); err != nil {
		return fmt.Errorf("render mail: %w", err)
	}
	text := fmt.Sprintf("Hello %s,\r\n\r\nYou are receiving this email because we received a password reset request for your account.\r\n\r\nReset your password: %s\r\n\r\nThis password reset link will expire in 60 minutes.\r\nIf you did not request a password reset, no further action is required.\r\n", name, link)
	return sendMail(email, "Reset Password Notification", text, htmlBody.String())
}

// sendMail reads the same MAIL_* variables as Laravel. Port 465 uses implicit TLS, others STARTTLS.
func sendMail(to, subject, textBody, htmlBody string) error {
	return sendMailMessage(mailMessage{To: []string{to}, Subject: subject, Text: textBody, HTML: htmlBody})
}

// mailMessage is one outgoing email. An empty From uses MAIL_FROM_*; an empty Text sends HTML only.
type mailMessage struct {
	FromAddress string
	FromName    string
	ReplyTo     *mail.Address
	To          []string
	Cc          []string
	Subject     string
	Text        string
	HTML        string
}

func joinAddresses(addresses []string) string {
	formatted := make([]string, len(addresses))
	for i, address := range addresses {
		formatted[i] = (&mail.Address{Address: address}).String()
	}
	return strings.Join(formatted, ", ")
}

func buildMailBody(m mailMessage, fromAddr, fromName string) (*bytes.Buffer, error) {
	var msg bytes.Buffer
	fmt.Fprintf(&msg, "From: %s\r\n", (&mail.Address{Name: fromName, Address: fromAddr}).String())
	fmt.Fprintf(&msg, "To: %s\r\n", joinAddresses(m.To))
	if len(m.Cc) > 0 {
		fmt.Fprintf(&msg, "Cc: %s\r\n", joinAddresses(m.Cc))
	}
	if m.ReplyTo != nil {
		fmt.Fprintf(&msg, "Reply-To: %s\r\n", m.ReplyTo.String())
	}
	fmt.Fprintf(&msg, "Subject: %s\r\n", mime.QEncoding.Encode("utf-8", m.Subject))
	fmt.Fprintf(&msg, "Date: %s\r\n", time.Now().Format(time.RFC1123Z))
	if m.Text == "" {
		fmt.Fprintf(&msg, "MIME-Version: 1.0\r\nContent-Type: text/html; charset=utf-8\r\n\r\n%s\r\n", m.HTML)
		return &msg, nil
	}
	raw := make([]byte, 12)
	if _, err := rand.Read(raw); err != nil {
		return nil, err
	}
	boundary := "wakiads-" + hex.EncodeToString(raw)
	fmt.Fprintf(&msg, "MIME-Version: 1.0\r\nContent-Type: multipart/alternative; boundary=%q\r\n\r\n", boundary)
	fmt.Fprintf(&msg, "--%s\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n%s\r\n", boundary, m.Text)
	fmt.Fprintf(&msg, "--%s\r\nContent-Type: text/html; charset=utf-8\r\n\r\n%s\r\n", boundary, m.HTML)
	fmt.Fprintf(&msg, "--%s--\r\n", boundary)
	return &msg, nil
}

func sendMailMessage(m mailMessage) error {
	host := os.Getenv("MAIL_HOST")
	port := os.Getenv("MAIL_PORT")
	if port == "" {
		port = "587"
	}
	fromAddr := m.FromAddress
	if fromAddr == "" {
		fromAddr = os.Getenv("MAIL_FROM_ADDRESS")
	}
	if fromAddr == "" {
		fromAddr = os.Getenv("MAIL_USERNAME")
	}
	fromName := m.FromName
	if fromName == "" {
		fromName = os.Getenv("MAIL_FROM_NAME")
	}
	if fromName == "" {
		fromName = "Wakiads"
	}

	msg, err := buildMailBody(m, fromAddr, fromName)
	if err != nil {
		return err
	}

	addr := net.JoinHostPort(host, port)
	dialer := &net.Dialer{Timeout: 15 * time.Second}
	tlsConfig := &tls.Config{ServerName: host, MinVersion: tls.VersionTLS12, RootCAs: mailRootCAs}

	portNumber, _ := strconv.Atoi(port)
	implicitTLS := portNumber == 465
	var conn net.Conn
	if implicitTLS {
		conn, err = tls.DialWithDialer(dialer, "tcp", addr, tlsConfig)
	} else {
		conn, err = dialer.Dial("tcp", addr)
	}
	if err != nil {
		return fmt.Errorf("connect smtp: %w", err)
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(30 * time.Second))

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("smtp client: %w", err)
	}
	defer client.Close()

	if !implicitTLS {
		// Required: credentials and reset links must never travel in clear text.
		if err := client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("smtp starttls: %w", err)
		}
	}
	if user := os.Getenv("MAIL_USERNAME"); user != "" {
		if err := client.Auth(smtp.PlainAuth("", user, os.Getenv("MAIL_PASSWORD"), host)); err != nil {
			return fmt.Errorf("smtp auth: %w", err)
		}
	}
	if err := client.Mail(fromAddr); err != nil {
		return fmt.Errorf("smtp from: %w", err)
	}
	for _, recipient := range append(append([]string{}, m.To...), m.Cc...) {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("smtp rcpt: %w", err)
		}
	}
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data: %w", err)
	}
	if _, err := w.Write(msg.Bytes()); err != nil {
		return fmt.Errorf("smtp write: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("smtp close: %w", err)
	}
	return client.Quit()
}
