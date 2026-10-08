package parser

import (
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Tokenize découpe une chaîne en jetons simples.
// Chaque mot, espace, apostrophe et marqueur de transformation est isolé.
func Tokenize(input string) []string {
	if input == "" {
		return nil
	}

	var tokens []string
	var word strings.Builder

	// flushWord ajoute le mot courant au tableau de jetons.
	flushWord := func() {
		if word.Len() == 0 {
			return
		}
		tokens = append(tokens, word.String())
		word.Reset()
	}

	// On lit la chaîne caractère par caractère en tenant compte des caractères UTF-8.
	for i := 0; i < len(input); {
		r, size := utf8.DecodeRuneInString(input[i:])
		switch {
		case unicode.IsSpace(r):
			flushWord()
			tokens = append(tokens, string(r))
			i += size
		case r == '(':
			flushWord()
			end := findMarkerEnd(input[i:])
			if end >= 0 {
				candidate := input[i : i+end+1]
				if isMarker(candidate) {
					tokens = append(tokens, candidate)
					i += end + 1
					continue
				}
			}
			tokens = append(tokens, "(")
			i += size
		case r == ')':
			flushWord()
			tokens = append(tokens, ")")
			i += size
		case r == '\'' || r == '-':
			if word.Len() > 0 && isWordRuneAt(input, i+size) {
				word.WriteRune(r)
				i += size
				continue
			}
			flushWord()
			tokens = append(tokens, string(r))
			i += size
		case unicode.IsLetter(r) || unicode.IsNumber(r):
			word.WriteRune(r)
			i += size
		default:
			flushWord()
			tokens = append(tokens, string(r))
			i += size
		}
	}

	flushWord()
	return tokens
}

// isWordRuneAt indique si le caractère à la position donnée peut continuer un mot.
func isWordRuneAt(input string, index int) bool {
	if index >= len(input) {
		return false
	}
	r, _ := utf8.DecodeRuneInString(input[index:])
	return unicode.IsLetter(r) || unicode.IsNumber(r)
}

// findMarkerEnd cherche la position de la parenthèse fermante d'un marqueur.
func findMarkerEnd(s string) int {
	for i := 1; i < len(s); i++ {
		if s[i] == ')' {
			return i
		}
	}
	return -1
}

// isMarker vérifie si une sous-chaîne correspond à un marqueur valide.
// Exemples : (up), (low, 3), (hex), (bin)
func isMarker(marker string) bool {
	if len(marker) < 3 || marker[0] != '(' || marker[len(marker)-1] != ')' {
		return false
	}

	body := strings.TrimSpace(marker[1 : len(marker)-1])
	if body == "" {
		return false
	}

	name := body
	count := 1
	if strings.Contains(body, ",") {
		parts := strings.SplitN(body, ",", 2)
		if len(parts) != 2 {
			return false
		}
		name = strings.TrimSpace(parts[0])
		count, _ = strconv.Atoi(strings.TrimSpace(parts[1]))
	}

	switch strings.ToLower(name) {
	case "up", "low", "cap":
		return count >= 1
	case "hex", "bin":
		return !strings.Contains(body, ",")
	default:
		return false
	}
}
