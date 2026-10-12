package transform

import (
	"strconv"
	"strings"
	"unicode"
)

// applyMarkers transforme les mots précédant chaque marqueur reconnu.
func applyMarkers(tokens []string) []string {
	result := make([]string, 0, len(tokens))
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		if !isMarkerToken(token) {
			result = append(result, token)
			continue
		}

		markerName, count, _ := parseMarker(token)
		wordCount := 1
		if markerName == "up" || markerName == "low" || markerName == "cap" {
			wordCount = count
		}
		wordIndexes := precedingWordIndexes(result, wordCount)
		if len(wordIndexes) == 0 {
			result = append(result, token)
			continue
		}

		converted := true
		if markerName == "hex" || markerName == "bin" {
			index := wordIndexes[0]
			result[index], converted = convertNumber(result[index], markerName)
		} else {
			for _, index := range wordIndexes {
				result[index] = applyMarker(result[index], markerName)
			}
		}
		if !converted {
			result = append(result, token)
			continue
		}

		hadLeftSpace := false
		for len(result) > 0 && isSpaceToken(result[len(result)-1]) {
			hadLeftSpace = true
			result = result[:len(result)-1]
		}

		rightStart := i + 1
		for rightStart < len(tokens) && isSpaceToken(tokens[rightStart]) {
			rightStart++
		}
		if rightStart == len(tokens) {
			i = rightStart - 1
			continue
		}

		hadRightSpace := rightStart > i+1
		nextIsMarker := isMarkerToken(tokens[rightStart])
		previousIsWord := len(result) > 0 && isWordToken(result[len(result)-1])
		nextIsWord := isWordToken(tokens[rightStart])
		if !nextIsMarker && previousIsWord && nextIsWord {
			result = append(result, " ")
		} else if !nextIsMarker && (hadLeftSpace || hadRightSpace) {
			result = append(result, " ")
		}

		if rightStart > i+1 {
			i = rightStart - 1
		}
	}
	return result
}

// precedingWordIndexes récupère les indices des derniers mots rencontrés.
func precedingWordIndexes(tokens []string, count int) []int {
	capacity := count
	if capacity > len(tokens) {
		capacity = len(tokens)
	}
	indexes := make([]int, 0, capacity)
	for i := len(tokens) - 1; i >= 0 && len(indexes) < count; i-- {
		if isWordToken(tokens[i]) {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

// isMarkerToken vérifie si un jeton est un marqueur de transformation valide.
func isMarkerToken(token string) bool {
	_, _, ok := parseMarker(token)
	return ok
}

// parseMarker valide un marqueur et extrait son nom et son nombre de mots.
func parseMarker(token string) (string, int, bool) {
	if len(token) < 3 || token[0] != '(' || token[len(token)-1] != ')' {
		return "", 0, false
	}
	body := strings.TrimSpace(token[1 : len(token)-1])
	if body == "" {
		return "", 0, false
	}

	name := body
	count := 1
	if strings.Contains(body, ",") {
		parts := strings.SplitN(body, ",", 2)
		if len(parts) != 2 || strings.Contains(parts[1], ",") {
			return "", 0, false
		}
		name = strings.TrimSpace(parts[0])
		var err error
		count, err = strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return "", 0, false
		}
	}

	switch strings.ToLower(name) {
	case "up", "low", "cap":
		return strings.ToLower(name), count, count >= 1
	case "hex", "bin":
		return strings.ToLower(name), 1, !strings.Contains(body, ",")
	default:
		return "", 0, false
	}
}

// applyMarker applique la transformation de casse demandée à un mot.
func applyMarker(word string, markerName string) string {
	switch markerName {
	case "up":
		return strings.ToUpper(word)
	case "low":
		return strings.ToLower(word)
	case "cap":
		return capitalizeWord(word)
	default:
		return word
	}
}

// capitalizeWord met en majuscule la première lettre du mot.
func capitalizeWord(word string) string {
	runes := []rune(word)
	for i, r := range runes {
		if unicode.IsLetter(r) {
			runes[i] = unicode.ToUpper(r)
			break
		}
	}
	return string(runes)
}

// convertNumber convertit un nombre hexadécimal ou binaire en base décimale.
func convertNumber(word string, markerName string) (string, bool) {
	base := 16
	if markerName == "bin" {
		base = 2
	}
	value, err := strconv.ParseUint(word, base, 64)
	if err != nil {
		return word, false
	}
	return strconv.FormatUint(value, 10), true
}
