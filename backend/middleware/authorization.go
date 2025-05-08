package middleware

import (
	"context"
	"net"
	"net/http"
	"time"

	"real_time_forum/backend/database"
	"real_time_forum/backend/utils"
)

// Store to track request counts per IP (in-memory map)
var (
	requestCounts = make(map[string]int)
	lastSeenTimes = make(map[string]time.Time)
)

// Middleware to limit requests based on IP
func RateLimit(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the client IP address
		ip := getClientIP(r)
		if ip == "" {
			http.Error(w, "Unable to determine client IP", http.StatusInternalServerError)
			return
		}

		// Initialize request count for new IPs
		if _, exists := requestCounts[ip]; !exists {
			requestCounts[ip] = 0
			lastSeenTimes[ip] = time.Now()
		}

		// Check if the time window has expired (1 minute in this example)
		if time.Since(lastSeenTimes[ip]) > time.Minute {
			// Reset the count and update the last seen time
			requestCounts[ip] = 0
			lastSeenTimes[ip] = time.Now()
		}

		// Increment the request count
		requestCounts[ip]++

		// Set the rate limit (e.g., 10 requests per minute)
		if requestCounts[ip] > 10 {
			// Deny the request if the limit is exceeded
			http.Error(w, "Too many requests. Slow down!", http.StatusTooManyRequests)
			return
		}

		// Update the last seen time for the IP
		lastSeenTimes[ip] = time.Now()

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	}
}

// Helper function to get the client IP address
func getClientIP(r *http.Request) string {
	// Use the X-Forwarded-For header if it's set
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	// Fall back to the remote address
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	return ip
}

// Authorization checks if the user is logged in by verifying the session token and expiration time
// from the database. If valid, it attaches user details to the request context and proceeds to the next handler.
func Authorization(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Token")
		if err != nil {
			utils.ResponseJSON(w, http.StatusUnauthorized, map[string]any{
				"message": "Missing authentication token",
				"status":  http.StatusUnauthorized,
			})
			return
		}

		var userId int
		var userName string
		var expired time.Time
		err = database.DB.QueryRow("SELECT id, username, expiredAt FROM users WHERE session=?", cookie.Value).Scan(&userId, &userName, &expired)
		if err != nil || userId == 0 {
			utils.ResponseJSON(w, http.StatusUnauthorized, map[string]any{
				"message": "You are not authorized to do this",
				"status":  http.StatusUnauthorized,
			})
			return
		}

		if time.Now().UTC().After(expired) {
			database.DB.Exec("UPDATE users SET session=? WHERE id=?", "", userId)
			utils.ResponseJSON(w, http.StatusUnauthorized, map[string]any{
				"message": "Session expired",
				"status":  http.StatusUnauthorized,
			})
			return
		}

		ctx := context.WithValue(r.Context(), "userId", userId)
		ctx = context.WithValue(ctx, "userName", userName)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func IsLogged(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Token")
	if err != nil {
		utils.ResponseJSON(w, http.StatusUnauthorized, map[string]any{
			"message": "Missing authentication token",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	var userId int
	var userName string
	var firstName string
	var lastName string
	var expired time.Time
	err = database.DB.QueryRow("SELECT id, username, firstName, lastName, expiredAt FROM users WHERE session=?", cookie.Value).Scan(&userId, &userName, &firstName, &lastName, &expired)
	if err != nil || userId == 0 {
		utils.ResponseJSON(w, http.StatusUnauthorized, map[string]any{
			"message": "You are not authorized to do this",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	if time.Now().UTC().After(expired) {
		database.DB.Exec("UPDATE users SET session=? WHERE id=?", "", userId)
		utils.ResponseJSON(w, http.StatusUnauthorized, map[string]any{
			"message": "Session expired",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	utils.ResponseJSON(w, http.StatusOK, map[string]any{
		"message":   "Valid token",
		"status":    http.StatusOK,
		"username":  userName,
		"firstName": firstName,
		"lastName":  lastName,
		"id":        userId,
	})
}
