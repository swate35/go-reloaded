package main

import (
	"fmt"
	"go-reloaded/parser"
	"go-reloaded/transform"
	"go-reloaded/utils"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	// Vérification du nombre d'arguments.
	if len(os.Args) != 3 {
		return fmt.Errorf("erreur : nombre d'arguments invalide\nUtilisation : go run main.go <chemin-fichier-entree> <fichier-sortie>")
	}
	// Récupération des chemins du fichier d'entrée et de sortie.
	inputPath := os.Args[1]
	outputPath := os.Args[2]
	// Lecture du contenu du fichier d'entrée.
	text, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("erreur : impossible de lire le fichier d'entrée %q : %w", inputPath, err)
	}
	// Découpage du texte en jetons.
	tokens, err := parser.Tokenize(string(text))
	if err != nil {
		return fmt.Errorf("erreur lors de l'analyse du fichier d'entrée : %w", err)
	}
	// Application des transformations demandées par les marqueurs.
	tokens, err = transform.ApplyTransformations(tokens)
	if err != nil {
		return fmt.Errorf("erreur lors de la transformation du texte : %w", err)
	}
	// Reconstruction du texte final à partir des jetons.
	result, err := utils.Reconstruct(tokens)
	if err != nil {
		return fmt.Errorf("erreur lors de la reconstruction du texte : %w", err)
	}
	// Écriture du résultat dans le fichier de sortie.
	err = os.WriteFile(outputPath, []byte(result), 0644)
	if err != nil {
		return fmt.Errorf("erreur : impossible d'écrire le fichier de sortie %q : %w", outputPath, err)
	}
	return nil
}
