package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware is used to protect routes by validating the JWT token
func (env *env) AuthMiddleware(c *gin.Context) {
	// Retrieve the token from the "Auth" cookie
	tokenStr, err := c.Cookie("Auth")
	if err != nil {
		// If no token is found, redirect to the login page
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	// Parse and validate the JWT token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key used for signing the JWT
		return []byte(env.secretKey), nil
	})
	if err != nil {
		// If the token is invalid, redirect to the login page
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	// Extract claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		// If claims are not valid, redirect to the login page
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	// Check if the token has expired
	if claims["ttl"].(float64) < float64(time.Now().Unix()) {
		// If the token has expired, redirect to the login page
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	// Extract the userID from the token claims
	userID := claims["userID"].(float64)
	if userID == 0 {
		// If the user ID is not valid, redirect to the login page
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	// Store the user ID in the context for future use
	c.Set("userID", userID)

	// Allow the request to proceed to the next step
	c.Next()
}
