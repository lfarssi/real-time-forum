package utils

import "regexp"

// IsValidEmail checks if the given email string matches a valid email format using regex.
func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
	return re.MatchString(email)
}