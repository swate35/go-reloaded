package transform

// formatQuotes colle les apostrophes de citation au texte qu'elles encadrent.
func formatQuotes(tokens []string) []string {
	quoteIndexes := make([]int, 0)
	for i, token := range tokens {
		if token == "'" {
			quoteIndexes = append(quoteIndexes, i)
		}
	}

	remove := make(map[int]bool)
	for i := 0; i+1 < len(quoteIndexes); i += 2 {
		opening, closing := quoteIndexes[i], quoteIndexes[i+1]
		for j := opening + 1; j < closing && isSpaceToken(tokens[j]); j++ {
			remove[j] = true
		}
		for j := closing - 1; j > opening && isSpaceToken(tokens[j]); j-- {
			remove[j] = true
		}
	}

	result := make([]string, 0, len(tokens)-len(remove))
	for i, token := range tokens {
		if !remove[i] {
			result = append(result, token)
		}
	}
	return result
}
