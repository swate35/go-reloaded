package main

import (
	"fmt"
	"go-reloaded/parser"
	"go-reloaded/transform"
	"go-reloaded/utils"
	"os"
)

func main() {
	// 1. Vérification des arguments
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <path-to-go-file> <output-file>")
		return
	}
	// 2. Récupération des chemins
	inputPath := os.Args[1]
	outputPath := os.Args[2]
	// 3. Lecture du fichier d'entrée
	text, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	// 4. Tokenisation
	tokens := parser.Tokenize(string(text))
	// 5. Transformations
	tokens = transform.ApplyTransformations(tokens)
	// 6. Reconstruction
	result := utils.Reconstruct(tokens)
	// 7. Écriture du fichier de sortie
	err = os.WriteFile(outputPath, []byte(result), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
}
