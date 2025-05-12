package middleware

import (
	"context"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"real_time_forum/backend/database"
	"real_time_forum/backend/utils"
)

var (
	requestCounts = make(map[string]int)
	lastSeenTimes = make(map[string]time.Time)
	mu            sync.Mutex
)

// RateLimit limits requests per IP (10 per minute)
func RateLimit(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)

		if ip == "" {
			utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
				"message": "Unable to determine client IP",
				"status":  http.StatusInternalServerError,
			})
			return
		}

		mu.Lock()

		if _, exists := requestCounts[ip]; !exists {
			requestCounts[ip] = 0
			lastSeenTimes[ip] = time.Now()
		}

		if time.Since(lastSeenTimes[ip]) > time.Minute {
			requestCounts[ip] = 0
			lastSeenTimes[ip] = time.Now()
		}

		requestCounts[ip]++
		if requestCounts[ip] > 10 {
			mu.Unlock()
			utils.ResponseJSON(w, http.StatusTooManyRequests, map[string]any{
				"message": "Too many requests. Slow down!",
				"status":  http.StatusTooManyRequests,
			})

			return
		}

		lastSeenTimes[ip] = time.Now()
		mu.Unlock()

		next.ServeHTTP(w, r)
	}
}

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header, which can contain multiple IPs
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// X-Forwarded-For can be a comma-separated list of IPs; take the first one
		ips := strings.Split(xff, ",")
		ip := strings.TrimSpace(ips[0])
		if net.ParseIP(ip) != nil {
			return ip
		}
	}

	// Fallback to X-Real-IP header
	xri := r.Header.Get("X-Real-IP")
	if xri != "" && net.ParseIP(xri) != nil {
		return xri
	}

	// Fallback to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	if ip == "::1" {
		return "127.0.0.1"
	}
	return ip
}

// Authorization middleware: checks valid session and adds user info to context
func Authorization(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		err = database.DB.QueryRow(`SELECT id, username, expiredAt FROM users WHERE session=?`, cookie.Value).
			Scan(&userId, &userName, &expired)

		if err != nil || userId == 0 {
			log.Println("Authorization error:", err)
			utils.ResponseJSON(w, http.StatusUnauthorized, map[string]any{
				"message": "You are not authorized to do this",
				"status":  http.StatusUnauthorized,
			})
			return
		}

		if expired.IsZero() || time.Now().UTC().After(expired) {
			_, _ = database.DB.Exec(`UPDATE users SET session=? WHERE id=?`, "", userId)
			utils.ResponseJSON(w, http.StatusUnauthorized, map[string]any{
				"message": "Session expired",
				"status":  http.StatusUnauthorized,
			})
			return
		}

		ctx := context.WithValue(r.Context(), "userId", userId)
		ctx = context.WithValue(ctx, "userName", userName)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// IsLogged is an endpoint to verify user session
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
	var userName, firstName, lastName string
	var expired time.Time
	err = database.DB.QueryRow(`
		SELECT id, username, firstName, lastName, expiredAt 
		FROM users 
		WHERE session=?`, cookie.Value).
		Scan(&userId, &userName, &firstName, &lastName, &expired)

	if err != nil || userId == 0 {
		utils.ResponseJSON(w, http.StatusUnauthorized, map[string]any{
			"message": "You are not authorized to do this",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	if expired.IsZero() || time.Now().UTC().After(expired) {
		_, _ = database.DB.Exec(`UPDATE users SET session=? WHERE id=?`, "", userId)
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
