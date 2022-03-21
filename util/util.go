package util

func ToUniqStrings(slice []string) []string {
	uniq := make([]string, 0)
	for _, v := range slice {
		if !Exists(uniq, v) {
			uniq = append(uniq, v)
		}
	}
	return uniq
}

func Exists(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
