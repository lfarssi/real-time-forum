package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"real_time_forum/backend/utils"
)

// Authorization checks if the user is logged in by verifying the session token and expiration time
// from the database. If valid, it attaches user details to the request context and proceeds to the next handler.
func Authorization(next http.Handler, db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Token")
		fmt.Println(cookie, "||", err)
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
		err = db.QueryRow("SELECT ID, UserName, Expared_At FROM Users WHERE Session=?", cookie.Value).Scan(&userId, &userName, &expired)
		if err != nil || userId == 0 {
			utils.ResponseJSON(w, http.StatusUnauthorized, map[string]any{
				"message": "Invalid session",
				"status":  http.StatusUnauthorized,
			})
			return
		}

		if time.Now().UTC().After(expired) {
			db.Exec("UPDATE users SET Session=? WHERE ID=?", "", userId)
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
