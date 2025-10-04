package handler

import (
	"net/http"
	"os"
)

const (
	// AuthCookieName is the name of the authentication cookie
	AuthCookieName = "auth_token"

	// AuthCookieMaxAge is the max age of the auth cookie in seconds (24 hours)
	AuthCookieMaxAge = 86400
)

// isProduction checks if the application is running in production mode
func isProduction() bool {
	env := os.Getenv("ENV")
	return env == "production" || env == "prod"
}

// createAuthCookie creates a secure authentication cookie
func createAuthCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     AuthCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   isProduction(), // Only send over HTTPS in production
		SameSite: http.SameSiteLaxMode,
		MaxAge:   AuthCookieMaxAge,
	}
}

// deleteAuthCookie creates a cookie that deletes the auth cookie
func deleteAuthCookie() *http.Cookie {
	return &http.Cookie{
		Name:     AuthCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   isProduction(),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1, // Delete cookie
	}
}
