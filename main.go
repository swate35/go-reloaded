package main

import (
	"fmt"
	"go-reloaded/parser"
	"go-reloaded/transform"
	"go-reloaded/utils"
	"os"
)

func main() {
	// 1. Vérification du nombre d'arguments.
	if len(os.Args) != 3 {
		fmt.Println("Utilisation : go run main.go <chemin-fichier-entree> <fichier-sortie>")
		return
	}
	// 2. Récupération des chemins du fichier d'entrée et de sortie.
	inputPath := os.Args[1]
	outputPath := os.Args[2]
	// 3. Lecture du contenu du fichier d'entrée.
	text, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier :", err)
		return
	}
	// 4. Découpage du texte en jetons.
	tokens := parser.Tokenize(string(text))
	// 5. Application des transformations demandées par les marqueurs.
	tokens = transform.ApplyTransformations(tokens)
	// 6. Reconstruction du texte final à partir des jetons.
	result := utils.Reconstruct(tokens)
	// 7. Écriture du résultat dans le fichier de sortie.
	err = os.WriteFile(outputPath, []byte(result), 0644)
	if err != nil {
		fmt.Println("Erreur lors de l'écriture du fichier :", err)
		return
	}
}
