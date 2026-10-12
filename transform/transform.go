package transform

import (
	"fmt"
	"unicode/utf8"
)

// ApplyTransformations applique les marqueurs, puis normalise le texte.
func ApplyTransformations(tokens []string) ([]string, error) {
	// Vérifier tous les jetons avant de les modifier garantit une erreur cohérente.
	for i, token := range tokens {
		if !utf8.ValidString(token) {
			return nil, fmt.Errorf("transform : le jeton %d n'est pas un UTF-8 valide", i)
		}
	}
	if len(tokens) == 0 {
		return nil, nil
	}

	// L'ordre est important : les transformations précèdent la normalisation typographique.
	tokens = applyMarkers(tokens)
	tokens = formatQuotes(tokens)
	tokens = formatPunctuation(tokens)
	return formatArticles(tokens), nil
}
