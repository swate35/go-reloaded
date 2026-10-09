# Application CLI & Parsing en Go : TPs

**Public :** Bachelor 1 · **Format :** FFP 12h = 4 séances de 3h · **Fil rouge :** un outil CLI `go-reloaded` qui lit un fichier texte, applique des règles de transformation et écrit le résultat.

## 1. Vue d'ensemble

| Séance | Durée | Thème | Objectifs UP couverts | TP |
|---|---|---|---|---|
| S1 | 3h | Architecture d'une CLI : modules, `main` vs `utils`, flux I/O | Architecture, fichiers | TP1 : squelette + lecture/écriture |
| S2 | 3h | Gestion d'erreurs « à la Go » | Erreurs, sécurité des fichiers | TP2 : conversions `hex` / `bin` |
| S3 | 3h | Parsing de texte : tokenisation, modificateurs, ponctuation | Parser robuste | TP3 : `up/low/cap(n)`, ponctuation, `a`→`an` |
| S4 | 3h | Qualité : tests manuels, README, revue de code | Tests, documentation | TP4 : jeu de tests + README + revue croisée |

**Déroulé type d'une séance (3h) :** 45 min de cours (PPT) · 30 min de démo live · 90 min de TP en autonomie · 15 min de correction et bilan.

## 2. Structure cible du projet

```
go-reloaded/
├── go.mod
├── main.go            # CLI : arguments, appels, codes de sortie
├── utils/
│   ├── process.go     # Process() : orchestre le pipeline
│   ├── modifiers.go   # (hex) (bin) (up) (low) (cap) (x, n)
│   └── format.go      # ponctuation, article a/an
├── tests/
│   ├── sample.txt     # entrées de test
│   └── expected.txt   # sorties attendues
└── README.md
```

Règle d'or : **`main` ne contient aucune logique de transformation.** Il lit, délègue, écrit, et gère les erreurs.

## 3. Séance 1 : Architecture & I/O

**Contenu :** `go mod init`, packages, visibilité (majuscule = exporté), séparation des responsabilités, `os.Args`, `os.ReadFile` / `os.WriteFile`, `os.Exit`, `os.Stderr`.

**Exemple de code : `main.go`**

```go
package main

import (
	"fmt"
	"os"

	"go-reloaded/utils"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: go run . <entree.txt> <sortie.txt>")
		os.Exit(1)
	}
	if err := run(os.Args[1], os.Args[2]); err != nil {
		fmt.Fprintln(os.Stderr, "erreur:", err)
		os.Exit(1)
	}
}

func run(in, out string) error {
	data, err := os.ReadFile(in)
	if err != nil {
		return fmt.Errorf("lecture de %q: %w", in, err)
	}
	result, err := utils.Process(string(data))
	if err != nil {
		return fmt.Errorf("traitement: %w", err)
	}
	if err := os.WriteFile(out, []byte(result), 0o644); err != nil {
		return fmt.Errorf("écriture de %q: %w", out, err)
	}
	return nil
}
```

### TP1 : Squelette de la CLI (≈ 90 min)
1. Initialiser le module `go-reloaded` et créer l'arborescence du §2.
2. Écrire `main.go` : vérifier qu'il y a exactement 2 arguments, sinon afficher l'usage et sortir avec le code 1.
3. Dans `utils/process.go`, créer `Process(s string) (string, error)` qui, pour l'instant, **renvoie le texte inchangé**.
4. Faire en sorte que l'entrée et la sortie soient différentes (refuser `in == out`).
5. Tester : fichier inexistant, un seul argument, fichier vide.

## 4. Séance 2 : Erreurs « à la Go »

**Contenu :** `error` comme valeur de retour, `if err != nil`, `fmt.Errorf` avec `%w`, `errors.Is`, messages utiles (quoi, où, pourquoi), pas de `panic` pour une erreur utilisateur, `strconv.ParseInt`.

**Exemple de code : `utils/modifiers.go` (extrait)**

```go
func hexToDec(w string) (string, error) {
	n, err := strconv.ParseInt(w, 16, 64)
	if err != nil {
		return "", fmt.Errorf("(hex) : %q n'est pas un nombre hexadécimal valide: %w", w, err)
	}
	return strconv.FormatInt(n, 10), nil
}
```

### TP2 : Conversions `(hex)` et `(bin)` (≈ 90 min)
Implémenter `hexToDec` et `binToDec`, puis les brancher sur le texte : `"1E (hex) files"` → `"30 files"` et `"10 (bin) years"` → `"2 years"`.
Cas d'erreur à traiter : mot invalide (`"zz (hex)"`), modificateur en début de texte.

## 5. Séance 3 : Parsing de texte

**Contenu :** pipeline en étapes (tokenisation → modificateurs → articles → ponctuation), `regexp`, `strings`, `unicode`, manipulation de `[]rune`, tranches (`slices`) et fenêtre « les n mots précédents ».

**Exemple de code : le pipeline**

```go
// utils/process.go
package utils

import "strings"

func Process(s string) (string, error) {
	tokens := tokenize(s)
	tokens, err := ApplyModifiers(tokens)
	if err != nil {
		return "", err
	}
	tokens = FixArticles(tokens)
	return FixPunctuation(strings.Join(tokens, " ")), nil
}
```

```go
// utils/modifiers.go
package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	tokenRe = regexp.MustCompile(`\((?:hex|bin|up|low|cap)(?:,\s*\d+)?\)|\S+`)
	modRe   = regexp.MustCompile(`^\((hex|bin|up|low|cap)(?:,\s*(\d+))?\)$`)
)

func tokenize(s string) []string { return tokenRe.FindAllString(s, -1) }

func Capitalize(w string) string {
	r := []rune(strings.ToLower(w))
	if len(r) == 0 {
		return w
	}
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func apply(kind, w string) (string, error) {
	switch kind {
	case "hex":
		return hexToDec(w)
	case "bin":
		return binToDec(w)
	case "up":
		return strings.ToUpper(w), nil
	case "low":
		return strings.ToLower(w), nil
	case "cap":
		return Capitalize(w), nil
	}
	return w, nil
}

func ApplyModifiers(tokens []string) ([]string, error) {
	out := make([]string, 0, len(tokens))
	for _, t := range tokens {
		m := modRe.FindStringSubmatch(t)
		if m == nil {
			out = append(out, t)
			continue
		}
		n := 1
		if m[2] != "" {
			n, _ = strconv.Atoi(m[2])
		}
		if len(out) == 0 {
			return nil, fmt.Errorf("modificateur %s sans mot précédent", t)
		}
		if n > len(out) {
			n = len(out)
		}
		for i := len(out) - n; i < len(out); i++ {
			v, err := apply(m[1], out[i])
			if err != nil {
				return nil, err
			}
			out[i] = v
		}
	}
	return out, nil
}
```

### TP3 : Modificateurs, ponctuation, articles (≈ 90 min)
1. Implémenter `(up)`, `(low)`, `(cap)` et leur variante `(up, 3)`.
2. Écrire `FixPunctuation` : pas d'espace avant `. , ! ? : ;`, un espace après ; les groupes `...` et `!?` restent collés.
3. Écrire `FixArticles` : `a` devient `an` devant une voyelle ou un `h`.

## 6. Séance 4 : Qualité, tests et README

**Contenu :** tests manuels avec fichiers `sample.txt` / `expected.txt` et `diff`, cas limites (vide, un seul mot, modificateur en début ou en fin), `gofmt`, `go vet`, nommage, commentaires utiles, structure d'un README.

### TP4 : Jeu de tests, README, revue croisée (≈ 90 min)
1. Créer `tests/sample.txt` (≥ 10 lignes couvrant toutes les règles) et `tests/expected.txt`.
2. Écrire un script `run_tests.sh` qui compare la sortie à l'attendu.
3. Rédiger le `README.md` (description, installation, utilisation, exemples, règles gérées, limites).
4. Revue croisée par binômes : chaque binôme applique la grille de relecture ci-dessous.
