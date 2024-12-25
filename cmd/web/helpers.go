package main

import (
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nyaruka/phonenumbers"
)

// Creates and signs a JWT token with a 1-hour TTL for the given user ID
func (env *env) createAndSignJWT(userID int) (string, error) {
	// Create a new JWT token with userID and TTL (time-to-live)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"ttl":    time.Now().Add(time.Hour * 1).Unix(),
	})

	return token.SignedString([]byte(env.secretKey))
}

// Sets the JWT token as a secure, HTTP-only cookie with a 1-hour expiration.
func SetCookie(c *gin.Context, token string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Auth", token, 3600*1, "", "", false, true)
}

// ValidateFullName validates a user's full name
func ValidateFullName(name string) error {
	name = strings.TrimSpace(name)

	// If name is empty, it's invalid
	if name == "" {
		return ErrInvalidFullName
	}

	// Split name into parts and check if there are at least two words
	parts := strings.Fields(name)
	if len(parts) < 2 {
		return ErrInvalidFullName
	}

	// Check if each part has at least two characters
	for _, part := range parts {
		if len(part) < 2 {
			return ErrInvalidFullName
		}
	}

	// Check if name contains valid characters
	if !isValidName(name) {
		return ErrInvalidFullName
	}

	// Ensure name doesn't contain any numeric characters
	for _, char := range name {
		if unicode.IsDigit(char) {
			return ErrInvalidFullName
		}
	}

	return nil
}

// isValidName checks if the name contains only valid characters
func isValidName(name string) bool {
	// Regex pattern that allows letters, spaces, apostrophes, and hyphens
	validNameRegex := `^[A-Za-zÀ-ÖØ-öø-ÿ]+(?:[-'\s][A-Za-zÀ-ÖØ-öø-ÿ]+)*$`
	re := regexp.MustCompile(validNameRegex)
	return re.MatchString(name)
}

// Sanitizes and validates a phone number, returning it in a standardized format
func sanitizePhoneNumber(phone string) (string, error) {
	num, err := phonenumbers.Parse(phone, "UZ")
	if err != nil {
		return "", err
	}

	if !phonenumbers.IsValidNumber(num) {
		return "", ErrInvalidPhoneNumber
	}

	return phonenumbers.Format(num, phonenumbers.E164), nil
}

// CheckPasswordStrength checks the strength of a given password
func CheckPasswordStrength(password string) error {
	if len(password) < 8 {
		return ErrWeakPassword
	}

	if !hasUppercase(password) || !hasLowercase(password) || !hasDigit(password) || !hasSpecialCharacter(password) {
		return ErrWeakPassword
	}

	return nil
}

// hasUppercase checks if the password contains an uppercase letter
func hasUppercase(s string) bool {
	return regexp.MustCompile(`[A-Z]`).MatchString(s)
}

// hasLowercase checks if the password contains a lowercase letter
func hasLowercase(s string) bool {
	return regexp.MustCompile(`[a-z]`).MatchString(s)
}

// hasDigit checks if the password contains a digit
func hasDigit(s string) bool {
	return regexp.MustCompile(`[0-9]`).MatchString(s)
}

// hasSpecialCharacter checks if the password contains a special character
func hasSpecialCharacter(s string) bool {
	return regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]\{\};:\'",<>\./?\\|]`).MatchString(s)
}
