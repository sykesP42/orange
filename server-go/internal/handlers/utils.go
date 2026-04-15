package handlers

func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}