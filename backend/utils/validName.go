package utils

func IsValidName(UserName string) bool {
	if len(UserName) == 0 {
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
