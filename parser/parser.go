package parser

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Tokenize découpe une chaîne en jetons simples.
// Chaque mot, espace, apostrophe et marqueur de transformation est isolé.
func Tokenize(input string) ([]string, error) {
	// La suite de transformations attend des chaînes UTF-8 valides.
	if !utf8.ValidString(input) {
		return nil, fmt.Errorf("parser : le texte d'entrée n'est pas un UTF-8 valide")
	}
	// nil représente l'absence de jetons, comme pour une entrée vide.
	if input == "" {
		return nil, nil
	}

	var tokens []string
	var word strings.Builder

	// flushWord ajoute le mot courant au tableau de jetons.
	flushWord := func() {
		if word.Len() == 0 {
			return
		}
		// Les séparateurs sont ajoutés séparément : les mots restent des jetons distincts.
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
				if isMarkerName(candidate) {
					// Garder le marqueur entier permet aux transformations de l'identifier.
					tokens = append(tokens, candidate)
					i += end + 1
					continue
				}
			}
			// Une parenthèse qui n'ouvre pas un marqueur reconnu reste un jeton normal.
			tokens = append(tokens, "(")
			i += size
		case r == ')':
			flushWord()
			tokens = append(tokens, ")")
			i += size
		case r == '\'' || r == '-':
			if word.Len() > 0 && isWordRuneAt(input, i+size) {
				// Apostrophe et tiret ne font partie d'un mot que s'ils sont internes.
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
	return tokens, nil
}

// isWordRuneAt indique si le caractère à la position donnée peut continuer un mot.
func isWordRuneAt(input string, index int) bool {
	if index >= len(input) {
		return false
	}
	// Décoder depuis l'index plutôt que lire un octet conserve le support UTF-8.
	r, _ := utf8.DecodeRuneInString(input[index:])
	return unicode.IsLetter(r) || unicode.IsNumber(r)
}

// findMarkerEnd cherche la position de la parenthèse fermante d'un marqueur.
func findMarkerEnd(s string) int {
	// Un marqueur se termine à la première parenthèse fermante rencontrée.
	for i := 1; i < len(s); i++ {
		if s[i] == ')' {
			return i
		}
	}
	return -1
}

// isMarkerName vérifie que le jeton correspond à un nom de marqueur pris en charge.
func isMarkerName(marker string) bool {
	if len(marker) < 3 || marker[0] != '(' || marker[len(marker)-1] != ')' {
		return false
	}

	body := strings.TrimSpace(marker[1 : len(marker)-1])
	if parts := strings.SplitN(body, ",", 2); len(parts) == 2 {
		// Le paramètre éventuel (nombre de mots) ne change pas le nom à valider.
		body = parts[0]
	}
	fields := strings.Fields(body)
	if len(fields) == 0 {
		return false
	}

	// Les noms sont insensibles à la casse, comme lors de l'application des marqueurs.
	switch strings.ToLower(fields[0]) {
	case "up", "low", "cap", "hex", "bin":
		return true
	default:
		return false
	}
}
