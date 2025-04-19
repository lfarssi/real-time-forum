package models

import (
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
)

// GenerateToken generates a new session token for a user, updates the user's session in the database,
// and returns the generated token along with the expiration time.
func GenerateToken(id int, db *sql.DB) (string, error) {
	u2, err := uuid.NewV6()
	if err != nil {
		return "", err
	}

	token := u2.String()
	expirationTime := time.Now().UTC().Add(time.Second * 10)

	_, err = db.Exec("UPDATE users set Session=? , Expared_At=?  WHERE ID=?", token, expirationTime, id)
	if err != nil {
		return "", err
	}

	return token, nil
}
