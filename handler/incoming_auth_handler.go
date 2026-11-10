package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/infraLinkit/mediaplatform-datasource/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthGate enforces AuthMiddleware when ENFORCE_API_AUTH=true. Laravel still calls most routes
// without a bearer token, so this stays off until the React SPA replaces it.
func (h *IncomingHandler) AuthGate(c *fiber.Ctx) error {
	if strings.TrimSpace(os.Getenv("ENFORCE_API_AUTH")) != "true" {
		return c.Next()
	}
	return h.AuthMiddleware(c)
}

func (h *IncomingHandler) AuthMiddleware(c *fiber.Ctx) error {
	// Already authenticated by AuthGate on the parent group.
	if c.Locals("user_id") != nil {
		return c.Next()
	}

	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		h.Logs.Warn("AuthMiddleware: Missing or invalid Authorization header")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid token"})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	secret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if secret == "" {
		h.Logs.Error("AuthMiddleware: JWT_SECRET environment variable not set")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "JWT secret not configured"})
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		h.Logs.WithError(err).Warn("AuthMiddleware: Token parse error")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}
	if !token.Valid {
		h.Logs.Warn("AuthMiddleware: Token is invalid")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		h.Logs.Warn("AuthMiddleware: Invalid token claims")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid claims"})
	}

	if exp, ok := claims["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			h.Logs.Warn("AuthMiddleware: Token expired")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token expired"})
		}
	}

	if tokenType, ok := claims["type"].(string); !ok || tokenType != "access" {
		h.Logs.Warn("AuthMiddleware: Invalid token type")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token type"})
	}

	if jti, ok := claims["jti"].(string); !ok || jti == "" {
		h.Logs.Warn("AuthMiddleware: Missing or invalid jti in token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token (jti)"})
	}

	var userID int
	switch v := claims["sub"].(type) {
	case float64:
		userID = int(v)
	case string:
		_, err := fmt.Sscanf(v, "%d", &userID)
		if err != nil {
			h.Logs.WithError(err).Warn("AuthMiddleware: Invalid user ID format in token sub claim")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID in token"})
		}
	default:
		h.Logs.Warn("AuthMiddleware: Unknown type for user ID in token sub claim")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID in token"})
	}

	var companies []entity.Company
	err = h.DB.Raw(`
		SELECT c.id, c.name
		FROM user_companies uc
		JOIN companies c ON uc.company_id = c.id
		WHERE uc.user_id = ? AND uc.status = true
	`, userID).Scan(&companies).Error
	if err != nil {
		h.Logs.WithError(err).Errorf("AuthMiddleware: Database error fetching companies for user %d", userID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	var companyNames []string
	for _, cpy := range companies {
		companyNames = append(companyNames, cpy.Name)
	}

	c.Locals("companies", companyNames)

	var adnets []entity.AdnetList
	err = h.DB.Raw(`
		SELECT a.id, a.code
		FROM user_adnets ua
		JOIN adnet_lists a ON ua.adnet_id = a.id
		WHERE ua.user_id = ? AND ua.status = true
	`, userID).Scan(&adnets).Error
	if err != nil {
		h.Logs.WithError(err).Errorf("AuthMiddleware: Database error fetching adnets for user %d", userID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	var adnetCodes []string
	for _, adn := range adnets {
		adnetCodes = append(adnetCodes, adn.Code)
	}

	c.Locals("adnets", adnetCodes)
	c.Locals("user_id", userID)

	var agencies []entity.Agency
	err = h.DB.Raw(`
		SELECT a.id, upper(a.code) as code, upper(a.name) as name
		FROM user_agencies ua
		JOIN agencies a ON ua.agency_id = a.id
		WHERE ua.user_id = ? AND ua.status = true
	`, userID).Scan(&agencies).Error
	if err != nil {
		h.Logs.WithError(err).Errorf("AuthMiddleware: Database error fetching agencies for user %d", userID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	var agencyNames []string
	for _, agency := range agencies {
		agencyNames = append(agencyNames, agency.Name)
		if agency.Code != "" {
			agencyNames = append(agencyNames, agency.Code)
		}
	}

	c.Locals("agencies", agencyNames)

	var countries []entity.Country
	err = h.DB.Raw(`
		SELECT DISTINCT c.code
		FROM detail_users du
		JOIN countries c ON du.country_id = c.id
		WHERE du.user_id = ? AND du.status = true
	`, userID).Scan(&countries).Error
	if err != nil {
		h.Logs.WithError(err).Errorf("AuthMiddleware: Database error fetching countries for user %d", userID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	var countryCodes []string
	for _, ctr := range countries {
		countryCodes = append(countryCodes, ctr.Code)
	}

	c.Locals("countries", countryCodes)

	return c.Next()
}

type LoginRequest struct {
	Email          string `json:"email"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	RecaptchaToken string `json:"recaptcha_token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type ProfileUpdateRequest struct {
	Username    string `json:"username"`
	CcEmail     string `json:"ccemail"`
	Name        string `json:"name"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func getJWTSecretKey() []byte {
	return []byte(strings.TrimSpace(os.Getenv("JWT_SECRET")))
}

func generateJTI() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func (h *IncomingHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if len(getJWTSecretKey()) == 0 {
		h.Logs.Error("Login: JWT_SECRET environment variable not set")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "JWT secret not configured"})
	}

	if req.Email == "" {
		req.Email = req.Username
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email/Username and password are required"})
	}

	if err := h.verifyRecaptcha(req.RecaptchaToken, c.IP()); err != nil {
		return validationError(c, "recaptcha", err.Error())
	}

	const invalidCredentials = "These credentials do not match our records."

	var user entity.User
	err := h.DB.Preload("Role").Where("email = ? OR username = ?", req.Email, req.Email).First(&user).Error
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": invalidCredentials})
	}

	// Password first: the account state must not be visible to someone who only knows a username.
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": invalidCredentials})
	}

	if !user.IsVerify {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Your account needs verification before you can log in. Please contact the administrator for assistance."})
	}

	if !user.Status {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Your account needs activation before you can log in. Please contact the administrator for assistance."})
	}

	now := time.Now()
	user.IPAddress = c.IP()
	user.Handset = string(c.Request().Header.UserAgent())
	user.LastLogin = &now
	h.DB.Model(&user).Updates(map[string]interface{}{
		"ip_address": user.IPAddress,
		"handset":    user.Handset,
		"last_login": user.LastLogin,
	})

	secret := getJWTSecretKey()

	accessClaims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(15 * time.Minute).Unix(),
		"type":  "access",
		"jti":   generateJTI(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessString, err := accessToken.SignedString(secret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create access token"})
	}

	refreshClaims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(30 * 24 * time.Hour).Unix(),
		"type":  "refresh",
		"jti":   generateJTI(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err := refreshToken.SignedString(secret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create refresh token"})
	}

	permissions := h.rolePermissionCodes(user.RoleID)

	userPayload := fiber.Map{
		"id":          user.ID,
		"name":        user.Name,
		"username":    user.Username,
		"email":       user.Email,
		"role_id":     user.RoleID,
		"role_name":   user.Role.Name,
		"role_code":   user.Role.Code,
		"profile_pic": user.ProfilePic,
		"permissions": permissions,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":          fiber.StatusOK,
		"status":        "success",
		"token":         accessString,
		"access_token":  accessString,
		"refresh_token": refreshString,
		"user":          userPayload,
		"data": fiber.Map{
			"access_token":  accessString,
			"refresh_token": refreshString,
			"user":          userPayload,
		},
	})
}

func (h *IncomingHandler) RefreshToken(c *fiber.Ctx) error {
	var req RefreshRequest
	if err := c.BodyParser(&req); err != nil || req.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Refresh token is required"})
	}

	secret := getJWTSecretKey()
	if len(secret) == 0 {
		h.Logs.Error("RefreshToken: JWT_SECRET environment variable not set")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "JWT secret not configured"})
	}
	token, err := jwt.Parse(req.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return secret, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid refresh token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "refresh" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token type"})
	}

	var userID int
	switch v := claims["sub"].(type) {
	case float64:
		userID = int(v)
	case string:
		_, _ = fmt.Sscanf(v, "%d", &userID)
	}

	var user entity.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
	}

	if !user.Status || !user.IsVerify {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User account is inactive"})
	}

	accessClaims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(15 * time.Minute).Unix(),
		"type":  "access",
		"jti":   generateJTI(),
	}
	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessString, err := newAccessToken.SignedString(secret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": accessString,
	})
}

// rolePermissionCodes returns the menu codes granted to a role, the same check as Laravel hasPermission().
func (h *IncomingHandler) rolePermissionCodes(roleID int) []string {
	codes := []string{}
	h.DB.Table("permissions").
		Joins("JOIN menus ON permissions.menu_id = menus.id").
		Where("permissions.role_id = ? AND permissions.status = true", roleID).
		Pluck("menus.code", &codes)
	return codes
}

func (h *IncomingHandler) GetMe(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int)
	if !ok || userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var user entity.User
	if err := h.DB.Preload("Role").First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	permissions := h.rolePermissionCodes(user.RoleID)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":          user.ID,
		"name":        user.Name,
		"username":    user.Username,
		"email":       user.Email,
		"role_id":     user.RoleID,
		"role_name":   user.Role.Name,
		"role_code":   user.Role.Code,
		"profile_pic": user.ProfilePic,
		"permissions": permissions,
		"companies":   c.Locals("companies"),
		"adnets":      c.Locals("adnets"),
		"agencies":    c.Locals("agencies"),
		"countries":   c.Locals("countries"),
	})
}

// GetProfile returns what the profile form needs beyond /v1/auth/me: the user's CC addresses.
func (h *IncomingHandler) GetProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int)
	if !ok || userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	ccEmails := []string{}
	if err := h.DB.Table("cc_emails").Where("user_id = ?", userID).Order("id ASC").Pluck("email", &ccEmails).Error; err != nil {
		h.Logs.Errorf("GetProfile: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch profile"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"ccemails": ccEmails}})
}

// splitCcEmails reads the comma-separated CC box; the second result is the first invalid address.
func splitCcEmails(value string) ([]string, string) {
	emails := []string{}
	for _, part := range strings.Split(value, ",") {
		email := strings.TrimSpace(part)
		if email == "" {
			continue
		}
		if address, err := mail.ParseAddress(email); err != nil || address.Address != email {
			return nil, email
		}
		emails = append(emails, email)
	}
	return emails, ""
}

// UpdateProfile follows UserProfileController::updateCurrentUser: name, username and the CC list.
// Email is read-only, and an empty CC box keeps the saved list.
func (h *IncomingHandler) UpdateProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int)
	if !ok || userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req ProfileUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var user entity.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	req.Name, req.Username = strings.TrimSpace(req.Name), strings.TrimSpace(req.Username)
	fieldErrors := fiber.Map{}
	if req.Name == "" {
		fieldErrors["name"] = "The name field is required."
	} else if len([]rune(req.Name)) > 255 {
		fieldErrors["name"] = "The name field must not be greater than 255 characters."
	}
	if req.Username == "" {
		fieldErrors["username"] = "The username field is required."
	} else if len([]rune(req.Username)) > 255 {
		fieldErrors["username"] = "The username field must not be greater than 255 characters."
	} else {
		taken, err := h.valueTaken("users", "username", req.Username, userID)
		if err != nil {
			h.Logs.Errorf("UpdateProfile: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update profile"})
		}
		if taken {
			fieldErrors["username"] = "The username has already been taken."
		}
	}
	ccEmails, invalidEmail := splitCcEmails(req.CcEmail)
	if invalidEmail != "" {
		fieldErrors["ccemail"] = "Email tidak valid: " + invalidEmail
	}
	if len(fieldErrors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Validation error", "errors": fieldErrors})
	}

	updates := map[string]interface{}{"name": req.Name, "username": req.Username, "updated_at": time.Now()}
	if req.NewPassword != "" {
		if req.OldPassword == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Current password is required to change password"})
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Current password does not match"})
		}
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to encrypt password"})
		}
		updates["password"] = string(hashed)
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if len(ccEmails) > 0 {
			if err := tx.Where("user_id = ?", userID).Delete(&entity.CcEmail{}).Error; err != nil {
				return err
			}
			for _, email := range ccEmails {
				if err := tx.Create(&entity.CcEmail{UserID: userID, Email: email}).Error; err != nil {
					return err
				}
			}
		}
		return tx.Model(&entity.User{}).Where("id = ?", userID).Updates(updates).Error
	})
	if err != nil {
		h.Logs.Errorf("UpdateProfile: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update profile"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Profile updated successfully!",
	})
}
