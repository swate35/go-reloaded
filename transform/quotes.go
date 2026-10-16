package transform

// formatQuotes colle les guillemets au texte qu'ils encadrent.
func formatQuotes(tokens []string) []string {
	quoteIndexes := map[string][]int{
		"'": {},
		`"`: {},
	}
	remove := make(map[int]bool)
	for i, token := range tokens {
		if token != "'" {
			continue
		}
		left, leftSpaces := adjacentWord(tokens, i, -1)
		right, rightSpaces := adjacentWord(tokens, i, 1)
		if left && right && (leftSpaces || rightSpaces) && !hasClosingQuoteBeforePunctuation(tokens, i) {
			for j := i - 1; j >= 0 && isSpaceToken(tokens[j]); j-- {
				remove[j] = true
			}
			for j := i + 1; j < len(tokens) && isSpaceToken(tokens[j]); j++ {
				remove[j] = true
			}
			continue
		}
		quoteIndexes[token] = append(quoteIndexes[token], i)
	}

	quoteIndexes[`"`] = make([]int, 0)
	for i, token := range tokens {
		if token == `"` {
			quoteIndexes[token] = append(quoteIndexes[token], i)
		}
	}

	// Les guillemets sont traités par paires ; seules les espaces aux bords sont retirées.
	for _, indexes := range quoteIndexes {
		for i := 0; i+1 < len(indexes); i += 2 {
			opening, closing := indexes[i], indexes[i+1]
			for j := opening + 1; j < closing && isSpaceToken(tokens[j]); j++ {
				remove[j] = true
			}
			for j := closing - 1; j > opening && isSpaceToken(tokens[j]); j-- {
				remove[j] = true
			}
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

func adjacentWord(tokens []string, quoteIndex, direction int) (bool, bool) {
	spaces := false
	for i := quoteIndex + direction; i >= 0 && i < len(tokens); i += direction {
		if isSpaceToken(tokens[i]) {
			spaces = true
			continue
		}
		return isWordToken(tokens[i]), spaces
	}
	return false, spaces
}

func hasClosingQuoteBeforePunctuation(tokens []string, quoteIndex int) bool {
	for i := quoteIndex + 1; i < len(tokens); i++ {
		if tokens[i] == "'" {
			return true
		}
		if isPunctuationToken(tokens[i]) {
			return false
		}
	}
	return false
}
