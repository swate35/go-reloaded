package transform

import (
	"unicode"
	"unicode/utf8"
)

// isSpaceToken indique si le jeton ne contient qu'un caractère d'espacement.
func isSpaceToken(token string) bool {
	r, size := utf8.DecodeRuneInString(token)
	return size == len(token) && unicode.IsSpace(r)
}

// isWordToken vérifie qu'un jeton correspond à un mot.
func isWordToken(token string) bool {
	if token == "" || isMarkerToken(token) || isSpaceToken(token) || isPunctuationToken(token) {
		return false
	}
	for _, r := range token {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '\'' || r == '-' {
			continue
		}
		return false
	}
	return true
}
