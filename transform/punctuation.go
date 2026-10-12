package transform

import "strings"

// isPunctuationToken reconnaît les signes dont la mise en forme est définie.
func isPunctuationToken(token string) bool {
	return token == "." || token == "," || token == "!" || token == "?" || token == ":" || token == ";"
}

// formatPunctuation colle les signes au mot précédent et sépare le suivant.
func formatPunctuation(tokens []string) []string {
	result := make([]string, 0, len(tokens))
	for i := 0; i < len(tokens); {
		if !isPunctuationToken(tokens[i]) {
			result = append(result, tokens[i])
			i++
			continue
		}

		for len(result) > 0 && isSpaceToken(result[len(result)-1]) {
			result = result[:len(result)-1]
		}

		// Rassembler les signes adjacents (même séparés par des espaces) en un seul groupe.
		var group strings.Builder
		for i < len(tokens) {
			if isPunctuationToken(tokens[i]) {
				group.WriteString(tokens[i])
				i++
				continue
			}
			if isSpaceToken(tokens[i]) {
				next := i + 1
				for next < len(tokens) && isSpaceToken(tokens[next]) {
					next++
				}
				if next < len(tokens) && isPunctuationToken(tokens[next]) {
					i = next
					continue
				}
			}
			break
		}
		result = append(result, group.String())

		for i < len(tokens) && isSpaceToken(tokens[i]) {
			i++
		}
		// La ponctuation est collée à gauche ; un seul espace la sépare du jeton suivant.
		if i < len(tokens) {
			result = append(result, " ")
		}
	}
	return result
}
