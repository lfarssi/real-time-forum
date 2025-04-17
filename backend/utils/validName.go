package utils

func ValidName(UserName string) bool {
	for _, ele := range UserName {
		if ele == '_' {
			continue
		}
		if !(ele >= 'a' && ele <= 'z') && !(ele >= 'A' && ele <= 'Z') && !(ele >= '0' && ele <= '9') {
			return false
		}

	}
	return true
}
