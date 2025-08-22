package handler

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/infraLinkit/mediaplatform-datasource/entity"
)

func (h *IncomingHandler) AuthMiddleware(c *fiber.Ctx) error {
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

	return c.Next()
}
