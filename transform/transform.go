package transform

import (
	"strconv"
	"strings"
	"unicode"
)

// ApplyTransformations applique les marqueurs de transformation sur le mot précédent.
// Par exemple : "Hello (up)" devient "HELLO".
func ApplyTransformations(tokens []string) []string {
	if len(tokens) == 0 {
		return nil
	}

	result := make([]string, 0, len(tokens))
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		if !isMarkerToken(token) {
			result = append(result, token)
			continue
		}

		markerName, count := parseMarker(token)
		if len(result) > 0 && isSpaceToken(result[len(result)-1]) {
			result = result[:len(result)-1]
		}
		index := len(result) - 1
		for index >= 0 && isSpaceToken(result[index]) {
			index--
		}
		if index >= 0 && isWordToken(result[index]) {
			result[index] = applyMarker(result[index], markerName, count)
		}
	}
	return result
}

// isMarkerToken vérifie si un jeton est un marqueur valide.
func isMarkerToken(token string) bool {
	if len(token) < 3 || token[0] != '(' || token[len(token)-1] != ')' {
		return false
	}
	body := strings.TrimSpace(token[1 : len(token)-1])
	if body == "" {
		return false
	}
	if strings.Contains(body, ",") {
		parts := strings.SplitN(body, ",", 2)
		if len(parts) != 2 {
			return false
		}
		if _, err := strconv.Atoi(strings.TrimSpace(parts[1])); err != nil {
			return false
		}
		body = strings.TrimSpace(parts[0])
	}
	return body == "up" || body == "low" || body == "cap" || body == "hex" || body == "bin"
}

// parseMarker récupère le nom du marqueur et son éventuel compteur.
func parseMarker(token string) (string, int) {
	body := strings.TrimSpace(token[1 : len(token)-1])
	count := 0
	if strings.Contains(body, ",") {
		parts := strings.SplitN(body, ",", 2)
		body = strings.TrimSpace(parts[0])
		if n, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
			count = n
		}
	}
	return strings.ToLower(body), count
}

// isSpaceToken indique si le jeton est un espace unique.
func isSpaceToken(token string) bool {
	return len(token) == 1 && unicode.IsSpace(rune(token[0]))
}

// isWordToken vérifie qu'un jeton est un mot valide.
func isWordToken(token string) bool {
	if token == "" || isMarkerToken(token) || isSpaceToken(token) {
		return false
	}
	for _, r := range token {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '\'' {
			continue
		}
		return false
	}
	return true
}

// applyMarker applique la transformation correspondant au marqueur.
func applyMarker(word string, markerName string, count int) string {
	switch markerName {
	case "up":
		return applyCount(word, strings.ToUpper, count, true)
	case "low":
		return applyCount(word, strings.ToLower, count, true)
	case "cap":
		return applyCount(word, capitalizeWord, count, false)
	case "hex":
		return toHex(word)
	case "bin":
		return toBinary(word)
	default:
		return word
	}
}

// applyCount applique une transformation sur un nombre précis de caractères.
func applyCount(word string, fn func(string) string, count int, wholeWord bool) string {
	if count == 0 || !wholeWord {
		return fn(word)
	}
	if count >= len(word) {
		return fn(word)
	}
	return fn(word[:count]) + word[count:]
}

// capitalizeWord met la première lettre en majuscule.
func capitalizeWord(word string) string {
	if word == "" {
		return word
	}
	runes := []rune(word)
	if len(runes) == 0 {
		return word
	}
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// toHex convertit chaque caractère en hexadécimal.
func toHex(word string) string {
	if word == "" {
		return word
	}
	var b strings.Builder
	for _, r := range word {
		b.WriteString(strconv.FormatInt(int64(r), 16))
	}
	return b.String()
}

// toBinary convertit chaque caractère en binaire.
func toBinary(word string) string {
	if word == "" {
		return word
	}
	var b strings.Builder
	for _, r := range word {
		b.WriteString(strconv.FormatInt(int64(r), 2))
	}
	return b.String()
}
