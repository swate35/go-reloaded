package utils

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Reconstruct reconstitue la chaîne finale à partir des jetons transformés.
func Reconstruct(tokens []string) (string, error) {
	for i, token := range tokens {
		if !utf8.ValidString(token) {
			return "", fmt.Errorf("utils : le jeton %d n'est pas un UTF-8 valide", i)
		}
	}
	return strings.Join(tokens, ""), nil
}
