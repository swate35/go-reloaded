package transform

import (
	"unicode"
	"unicode/utf8"
)

// isSpaceToken indique si le jeton ne contient qu'un caractère d'espacement.
func isSpaceToken(token string) bool {
	r, size := utf8.DecodeRuneInString(token)
	// La taille vérifie que le jeton contient exactement une rune, pas plusieurs espaces.
	return size == len(token) && unicode.IsSpace(r)
}

// isWordToken vérifie qu'un jeton correspond à un mot.
func isWordToken(token string) bool {
	if token == "" || isMarkerToken(token) || isSpaceToken(token) || isPunctuationToken(token) {
		return false
	}
	// Les contractions et mots composés sont admis, mais tout autre séparateur exclut le jeton.
	for _, r := range token {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '\'' || r == '-' {
			continue
		}
		return false
	}
	return true
}
