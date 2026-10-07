package test

import (
	"go-reloaded/parser"
	"go-reloaded/transform"
	"go-reloaded/utils"
	"testing"
)

// les tests unitaires pour vérifier les transformations de texte, on été generés par Ia.
func TestTextTransformations(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "hexadecimal et binaire",
			input: "Simply add 42 (hex) and 10 (bin) and you will see the result is 68.",
			want:  "Simply add 66 and 2 and you will see the result is 68.",
		},
		{
			name:  "casse et nombre de mots",
			input: "Welcome to the Brooklyn bridge (cap). This is so exciting (up, 2).",
			want:  "Welcome to the Brooklyn Bridge. This is SO EXCITING.",
		},
		{
			name:  "ponctuation simple et groupe",
			input: "Punctuation tests are ... kinda boring ,what do you think ?",
			want:  "Punctuation tests are... kinda boring, what do you think?",
		},
		{
			name:  "groupe de ponctuation espace",
			input: "I was sitting over there ,and then BAMM !!",
			want:  "I was sitting over there, and then BAMM!!",
		},
		{
			name:  "apostrophes de citation",
			input: "As Elton John said: ' I am the most well-known homosexual in the world '",
			want:  "As Elton John said: 'I am the most well-known homosexual in the world'",
		},
		{
			name:  "espaces multiples dans une citation",
			input: "'   texte cité   '",
			want:  "'texte cité'",
		},
		{
			name:  "article devant voyelle et h",
			input: "There it was. A amazing rock and a honest answer.",
			want:  "There it was. An amazing rock and an honest answer.",
		},
		{
			name:  "contraction préservée",
			input: "I don't know (up).",
			want:  "I don't KNOW.",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens := parser.Tokenize(test.input)
			tokens = transform.ApplyTransformations(tokens)
			got := utils.Reconstruct(tokens)
			if got != test.want {
				t.Fatalf("résultat inattendu :\nobtenu  : %q\nattendu : %q", got, test.want)
			}
		})
	}
}
