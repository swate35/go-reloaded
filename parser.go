package parser

import (
	"strings"
	"unicode"
)

func Tokenize(input string) []string {
	// 1. Créer une liste vide de tokens
	var tokens []string
	// 2. Parcourir le texte caractère par caractère
	for _, char := range input {
		// Vérifier si le caractère est un espace ou une ponctuation
		if unicode.IsSpace(char) || unicode.IsPunct(char) {
			
	// 3. Construire des mots
	// 4. Séparer la ponctuation
	// 5. Détecter les marqueurs (hex), (bin), (up), (low, 3), etc.
	// 6. Gérer les apostrophes '
	// 7. Ajouter chaque token dans la liste
	// 8. Retourner la liste de tokens
