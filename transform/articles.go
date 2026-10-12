package transform

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// formatArticles transforme « a » en « an » devant une voyelle ou un h.
func formatArticles(tokens []string) []string {
	result := append([]string(nil), tokens...)
	for i, token := range result {
		if !strings.EqualFold(token, "a") {
			continue
		}
		next := i + 1
		for next < len(result) && isSpaceToken(result[next]) {
			next++
		}
		if next >= len(result) || !isWordToken(result[next]) {
			continue
		}
		// Ignorer les espaces afin de tester la première lettre du mot suivant.
		first, _ := utf8.DecodeRuneInString(result[next])
		switch unicode.ToLower(first) {
		case 'a', 'e', 'i', 'o', 'u', 'h':
			if token == "A" {
				result[i] = "An"
			} else if token == "a" {
				result[i] = "an"
			}
		}
	}
	return result
}
