package worker

import (
	"strings"
)

func removeNonPrintableChars(s string) string {
	var result []rune
	for _, r := range s {
		// Allow printable characters, space, and allow typical control characters like \n, \r, and \t
		if (r >= 32 && r <= 126) || r == '\n' || r == '\r' || r == '\t' {
			result = append(result, r)
		}

	}
	return strings.TrimSpace(string(result))
}
