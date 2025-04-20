package models

import (
	"time"

	"real_time_forum/backend/database"

	"github.com/gofrs/uuid"
)

// GenerateToken generates a new session token for a user, updates the user's session in the database,
// and returns the generated token along with the expiration time.
func GenerateToken(id int) (string, error) {
	u2, err := uuid.NewV6()
	if err != nil {
		return "", err
	}

	token := u2.String()
	expirationTime := time.Now().UTC().Add(time.Hour)

	_, err = database.DB.Exec("UPDATE users set session=? , expiredAt=?  WHERE id=?", token, expirationTime, id)
	if err != nil {
		return "", err
	}

	return token, nil
}
