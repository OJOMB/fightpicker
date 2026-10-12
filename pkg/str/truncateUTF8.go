package str

import "unicode/utf8"

// TruncateUTF8 shortens s to at most maxBytes without splitting a multi-byte character.
func TruncateUTF8(s string, maxBytes int) string {
	if len(s) <= maxBytes {
		return s
	}
	// back up until we're at the start of a rune
	for maxBytes > 0 && !utf8.RuneStart(s[maxBytes]) {
		maxBytes--
	}
	return s[:maxBytes]
}
