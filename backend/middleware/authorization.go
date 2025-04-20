package middleware

import (
	"context"
	"net/http"
	"time"

	"real_time_forum/backend/database"
	"real_time_forum/backend/utils"
)

// Authorization checks if the user is logged in by verifying the session token and expiration time
// from the database. If valid, it attaches user details to the request context and proceeds to the next handler.
func Authorization(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Token")
		if err != nil {
			utils.ResponseJSON(w, http.StatusForbidden, map[string]any{
				"message": "Missing authentication token",
				"status":  http.StatusForbidden,
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
