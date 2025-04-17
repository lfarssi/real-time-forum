package middleware

import (
	"database/sql"
	"errors"
	"net/http"
	"time"
)

// VerifyCookie checks the validity of the session cookie, ensuring it exists and is not expired.
// If valid, it returns the cookie, otherwise, it returns an error.
func VerifyCookie(r *http.Request, db *sql.DB) (*http.Cookie, error) {
	cookie, err := r.Cookie("Token")
	if err != nil {
		return nil, err
	}
	exist := ""
	var expired time.Time
	db.QueryRow("SELECT Session, Expared_At FROM Users WHERE Session=?", cookie.Value).Scan(&exist, &expired)
	if exist == "" {
		return nil, errors.New("invalid token")
	}

	if time.Now().UTC().After(expired.UTC()) {
		db.Exec("UPDATE users set Session=? WHERE Session=?", "", exist)
		return nil, errors.New("expiration session")
	}

	return cookie, nil
}
