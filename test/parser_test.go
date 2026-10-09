package test

import (
	"go-reloaded/parser"
	"go-reloaded/transform"
	"go-reloaded/utils"
	"testing"
)

// les tests unitaires pour vérifier les transformations de texte, on été copier du pdf.
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
		{
			name:  "marqueurs insensibles à la casse",
			input: "This is a test (CAP, 2).",
			want:  "This is A Test.",
		},
		{
			name:  "marqueurs non documentés rejetés",
			input: "10 (bin, 2) and 5 (hex, 3)",
			want:  "10 (bin, 2) and 5 (hex, 3)",
		},
		{
			name:  "espace conservé autour du marqueur",
			input: "hello   (up) world",
			want:  "HELLO world",
		},
		{
			name:  "espaces multiples autour du marqueur",
			input: "hello   (up)   world",
			want:  "HELLO world",
		},
		{
			name:  "marqueur collé entre deux mots",
			input: "hello(up)world",
			want:  "HELLO world",
		},
		{
			name:  "ponctuation avant marqueur",
			input: "hello ,(up) world",
			want:  "HELLO, world",
		},
		{
			name:  "nombre invalide conserve ses espaces",
			input: "G1   (hex) and 10202 (bin)",
			want:  "G1   (hex) and 10202 (bin)",
		},
		{
			name:  "marqueurs enchaînés",
			input: "hello (up) (low)",
			want:  "hello",
		},
		{
			name:  "espaces finaux après marqueur",
			input: "hello (up)   ",
			want:  "HELLO",
		},
		{
			name:  "count invalide préservé",
			input: "hello (up,abc)",
			want:  "hello (up,abc)",
		},
		{
			name:  "marqueur vide préservé",
			input: "hello (up,)",
			want:  "hello (up,)",
		},
		{
			name:  "marqueur avec count supplémentaire préservé",
			input: "hello (up,3,4)",
			want:  "hello (up,3,4)",
		},
		{
			name:  "marqueur sans virgule invalide préservé",
			input: "hello (up 3)",
			want:  "hello (up 3)",
		},
		{
			name:  "espaces dans le marqueur valide",
			input: "hello (up , 1)",
			want:  "HELLO",
		},
		{
			name:  "count négatif préservé",
			input: "hello (low, -3)",
			want:  "hello (low, -3)",
		},
		{
			name:  "count trop grand préservé",
			input: "hello (cap, 999999999999999999999999)",
			want:  "hello (cap, 999999999999999999999999)",
		},
		{
			name:  "count valide supérieur aux mots disponibles",
			input: "hello (up, 999999999)",
			want:  "HELLO",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens, err := parser.Tokenize(test.input)
			if err != nil {
				t.Fatalf("erreur inattendue lors de l'analyse : %v", err)
			}
			tokens, err = transform.ApplyTransformations(tokens)
			if err != nil {
				t.Fatalf("erreur inattendue lors de la transformation : %v", err)
			}
			got, err := utils.Reconstruct(tokens)
			if err != nil {
				t.Fatalf("erreur inattendue lors de la reconstruction : %v", err)
			}
			if got != test.want {
				t.Fatalf("résultat inattendu :\nobtenu  : %q\nattendu : %q", got, test.want)
			}
		})
	}
}

func TestInvalidUTF8ReturnsErrors(t *testing.T) {
	invalidUTF8 := string([]byte{0xff})

	if _, err := parser.Tokenize(invalidUTF8); err == nil {
		t.Fatal("Tokenize aurait dû signaler l'UTF-8 invalide")
	}
	if _, err := transform.ApplyTransformations([]string{invalidUTF8}); err == nil {
		t.Fatal("ApplyTransformations aurait dû signaler l'UTF-8 invalide")
	}
	if _, err := utils.Reconstruct([]string{invalidUTF8}); err == nil {
		t.Fatal("Reconstruct aurait dû signaler l'UTF-8 invalide")
	}
}
