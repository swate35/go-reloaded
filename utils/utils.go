package utils

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Reconstruct reconstitue la chaîne finale à partir des jetons transformés.
func Reconstruct(tokens []string) (string, error) {
	// Valider avant la concaténation évite de renvoyer un résultat avec des octets invalides.
	for i, token := range tokens {
		if !utf8.ValidString(token) {
			return "", fmt.Errorf("utils : le jeton %d n'est pas un UTF-8 valide", i)
		}
	}
	// Les jetons contiennent déjà les espaces et séparateurs à leur place.
	return strings.Join(tokens, ""), nil
}
