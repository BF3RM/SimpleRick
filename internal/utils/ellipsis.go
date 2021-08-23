package utils

func Ellipsis(text string, length int) string {
	textLen := len(text)
	if textLen <= length {
		return text
	}
	return text[:textLen-1-3] + "..."
}
