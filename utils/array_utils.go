package foundation

// StringArrayContains checks if an array contains a specific value
func StringArrayContains(array []string, search string) bool {
	for _, v := range array {
		if v == search {
			return true
		}
	}
	return false
}

// IntArrayContains checks if an array contains a specific value
func IntArrayContains(array []int, search int) bool {
	for _, v := range array {
		if v == search {
			return true
		}
	}
	return false
}
