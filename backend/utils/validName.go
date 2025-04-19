package utils

import "regexp"

func IsValidName(UserName string) bool {
	if len(UserName) < 3 {
		return false
	}

	for _, ele := range UserName {
		if ele == '_' || ele == ' ' {
			continue
		}
		if !(ele >= 'a' && ele <= 'z') && !(ele >= 'A' && ele <= 'Z') && !(ele >= '0' && ele <= '9') {
			return false
		}

	}
	
	return true
}

func IsValidUserName(username string) bool  {
	regex := `^[a-zA-Z0-9_-]{3,13}$`

	regulierExpr := regexp.MustCompile(regex)

	return regulierExpr.MatchString(username)  
}
