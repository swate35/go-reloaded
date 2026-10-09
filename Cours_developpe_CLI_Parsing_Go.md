# Support de cours : Application CLI & Parsing en Go

---

# SECTION 1 : Architecture d'une CLI et notion de package (3h)

## 1.1 Objectifs de la séance
À la fin de la séance, vous saurez :
- expliquer la différence entre **module**, **package** et **fichier** ;
- structurer un projet en séparant `main` et `utils` ;
- écrire un programme qui lit et écrit des fichiers à partir d'arguments en ligne de commande ;
- expliquer ce que fait `go run .`.

## 1.2 Qu'est-ce qu'une CLI ?
Une **CLI** (Command Line Interface) est un programme qu'on lance depuis un terminal, sans interface graphique. Il reçoit des **arguments**, fait un traitement, puis rend la main avec un **code de sortie**.

```
go run . entree.txt sortie.txt
          └───────┬───────┘
         arguments du programme
```

Un bon programme CLI respecte des conventions :

| Convention | Signification |
|---|---|
| Résultat sur la **sortie standard** (`stdout`) | Ce que le programme produit |
| Erreurs sur la **sortie d'erreur** (`stderr`) | Messages à destination de l'utilisateur, séparés du résultat |
| Code de sortie `0` | Tout s'est bien passé |
| Code de sortie différent de `0` (souvent `1`) | Une erreur est survenue |
| Message d'usage si les arguments sont faux | Aide l'utilisateur à corriger sa commande |

## 1.3 Module, package, fichier : trois niveaux d'organisation

```
MODULE  go-reloaded            ← défini par go.mod, = tout le projet
 ├── PACKAGE main              ← la RACINE du projet (pas de dossier "main/") : le programme exécutable
 │     └── main.go             ← FICHIER (contient "package main" et func main())
 └── PACKAGE utils             ← dossier utils/ : une bibliothèque interne
       ├── process.go          ← FICHIER
       ├── modifiers.go        ← FICHIER
       └── format.go           ← FICHIER
```

> **Attention :** il n'y a pas de dossier nommé `main/`. `main.go` est placé directement à la racine du projet. « Package main » désigne la ligne `package main` écrite dans ce fichier, pas un dossier. Pour de gros projets, on peut placer le point d'entrée dans `cmd/<nom-du-programme>/main.go` (lancé avec `go run ./cmd/<nom-du-programme>`), mais ce n'est pas nécessaire ici.

### Le module
- C'est l'ensemble du projet, décrit par `go.mod` (créé par `go mod init go-reloaded`).
- Son nom (`go-reloaded`) sert de **préfixe** à tous les chemins d'import internes.

### Le package
Un package est **un dossier** contenant des fichiers `.go` qui partagent le même nom de package. C'est l'unité de réutilisation et d'encapsulation en Go.

Règles à connaître :

1. **Un dossier = un package.** Tous les fichiers d'un même dossier commencent par la même ligne `package xxx`.
2. **Plusieurs fichiers, un seul espace de noms.** Une fonction définie dans `modifiers.go` est utilisable dans `format.go` sans import, tant qu'ils sont dans le même package.
3. **Le chemin d'import = nom du module + chemin du dossier.** Le dossier `utils/` du module `go-reloaded` s'importe avec `"go-reloaded/utils"`.
4. **Le nom du package est celui utilisé dans le code** : après `import "go-reloaded/utils"`, on écrit `utils.Process(...)`.
5. **Par convention**, le nom du package est identique au nom du dossier, en minuscules, court, sans tiret ni majuscule.
6. **Imports circulaires interdits** : si `a` importe `b`, alors `b` ne peut pas importer `a`. Le compilateur refuse. C'est un bon signal que la découpe est mauvaise.
7. **Import inutilisé = erreur de compilation**. De même pour une variable déclarée mais jamais utilisée.

### Le package spécial `main`
- Un package nommé `main` qui contient une fonction `func main()` produit un **programme exécutable**.
- Tout autre package (`utils`, etc.) est une **bibliothèque** : on ne peut pas le lancer seul, seulement l'importer.

### La visibilité : exporté ou privé
Go n'a pas de mots-clés `public` / `private`. Tout repose sur la **première lettre** du nom :

```go
package utils

func Process(s string) (string, error) { ... } // Majuscule : EXPORTÉ, utilisable depuis main
func tokenize(s string) []string        { ... } // minuscule : PRIVÉ, visible seulement dans utils
```

| Nom | Visible depuis un autre package ? |
|---|---|
| `Process`, `Capitalize`, `ErrEmpty` | Oui |
| `tokenize`, `modRe`, `apply` | Non |

> **Bonne pratique :** exporter le minimum. Une API réduite est plus facile à comprendre et à modifier sans casser le reste du programme.

## 1.4 La séparation des responsabilités

Principe : **chaque partie du programme a une seule raison de changer.**

| Partie | Responsabilité | Ne doit PAS faire |
|---|---|---|
| `main.go` | Lire les arguments, appeler le traitement, écrire le fichier, afficher les erreurs, choisir le code de sortie | Aucune règle de transformation |
| `utils/process.go` | Enchaîner les étapes du traitement | Lire ou écrire des fichiers |
| `utils/modifiers.go` | Appliquer `(hex)`, `(up)`, etc. | Connaître d'où vient le texte |
| `utils/format.go` | Corriger ponctuation et articles | Gérer les erreurs de fichiers |

**Test mental très utile :** `Process` reçoit une `string` et renvoie une `string`. Elle ne sait pas si le texte vient d'un fichier, d'un test ou du clavier. On peut donc la tester sans créer aucun fichier.

**Contre-exemple (à éviter) :**

```go
// Mauvais : tout est mélangé dans main
func main() {
    data, _ := os.ReadFile(os.Args[1])
    words := strings.Fields(string(data))
    for i, w := range words {
        if w == "(up)" { /* ... 80 lignes de logique ... */ }
    }
    os.WriteFile(os.Args[2], []byte(strings.Join(words, " ")), 0644)
}
```
Problèmes : impossible à tester, erreurs ignorées, aucune réutilisation, difficile à lire.

## 1.5 Lire et écrire des fichiers

### Arguments : `os.Args`
`os.Args` est un `[]string`. L'élément `[0]` est le nom du programme, les suivants sont les arguments.

```go
// go run . entree.txt sortie.txt
// os.Args[0] -> chemin du binaire temporaire
// os.Args[1] -> "entree.txt"
// os.Args[2] -> "sortie.txt"
```
D'où le test `len(os.Args) != 3`.

### Lecture et écriture simples

```go
data, err := os.ReadFile("entree.txt")           // data est un []byte
err = os.WriteFile("sortie.txt", data, 0o644)    // permissions
```

- `[]byte` et `string` se convertissent : `string(data)` et `[]byte(texte)`.
- `0o644` est un nombre **octal** de permissions Unix : lecture/écriture pour le propriétaire, lecture seule pour les autres.
- `os.WriteFile` **écrase** le fichier s'il existe déjà. C'est pour cela qu'on vérifie que entrée et sortie sont différentes.

### Pourquoi `run()` et pas tout dans `main()` ?
`os.Exit` termine le programme **immédiatement**, sans exécuter les `defer`. Si toute la logique est dans `main`, on mélange logique et sortie. Le schéma recommandé est donc :

```go
func main() {
    if err := run(); err != nil {
        fmt.Fprintln(os.Stderr, "erreur:", err)
        os.Exit(1)        // un seul endroit où on quitte
    }
}

func run() error {        // toute la logique, retourne une erreur
    ...
}
```

## 1.6 Que fait `go run .` ?
1. `.` désigne le **package du dossier courant** (ici `main`).
2. Go compile ce package et ses dépendances (retrouvées via `go.mod`).
3. Il lance l'exécutable temporaire en lui passant les arguments qui suivent.
4. Il supprime l'exécutable.

Pour obtenir un binaire permanent : `go build -o go-reloaded .`.

## 1.7 Pièges fréquents (S1)

| Piège | Explication |
|---|---|
| `package utils` dans un dossier mais `package util` dans un autre fichier du même dossier | Erreur : un dossier = un seul package |
| Appeler `utils.process()` | Fonction privée (minuscule) : inaccessible depuis `main` |
| Oublier `go mod init` | L'import `"go-reloaded/utils"` échoue |
| Écrire le résultat avec `fmt.Println` pour un message d'erreur | Les erreurs vont sur `os.Stderr` |
| Lire et écrire le même fichier | Perte des données d'origine |

## 1.8 Questions de compréhension (corrigées)

1. *Quelle est la différence entre un module et un package ?*
   → Le module est le projet entier (un `go.mod`). Un package est un dossier de code dans ce module.
2. *Pourquoi `Process` commence-t-elle par une majuscule ?*
   → Pour être exportée et appelable depuis `main`.
3. *Que se passe-t-il si `utils` importe `main` ?*
   → On ne peut pas importer `main`, et de toute façon un import circulaire est refusé.
4. *Pourquoi `Process` ne lit-elle pas elle-même le fichier ?*
   → Pour rester testable et réutilisable, et respecter la séparation des responsabilités.
5. *`go run . a b` : que vaut `len(os.Args)` ?* → 3.

## 1.9 Exercices supplémentaires (avec corrigés)

**Exercice A.** Créer un package `greet` avec une fonction exportée `Hello(name string) string` qui retourne `"Bonjour, <name> !"` et une fonction privée `clean(name string) string` qui supprime les espaces autour du nom (chercher une fonction trim). Utiliser `Hello` depuis `main`.

```go
// greet/greet.go
package greet

import "strings"

func clean(name string) string { return strings.TrimSpace(name) }

func Hello(name string) string { return "Bonjour, " + clean(name) + " !" }
```
```go
// main.go
package main

import (
	"fmt"
	"go-reloaded/greet"
)

func main() { fmt.Println(greet.Hello("  Ada ")) } // Bonjour, Ada !
```

**Exercice B.** Écrire un programme qui affiche le nombre de lignes d'un fichier passé en argument.

```go
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run . <fichier>")
		os.Exit(1)
	}
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "erreur:", err)
		os.Exit(1)
	}
	fmt.Println(strings.Count(string(data), "\n"))
}
```

---

# SECTION 2 : La gestion des erreurs en Go (3h)

## 2.1 Objectifs
Apprendre à :
- expliquer pourquoi Go n'a pas d'exceptions et ce que cela implique ;
- écrire des fonctions qui retournent une erreur et traiter ces retours ;
- enrichir une erreur avec du contexte (`%w`) et la tester (`errors.Is`, `errors.As`) ;
- rédiger des messages d'erreur utiles.

## 2.2 Le principe : une erreur est une valeur
En Go, une erreur n'est pas « lancée » : elle est **retournée** comme n'importe quelle autre valeur. Le type `error` est une simple interface :

```go
type error interface {
    Error() string
}
```
Tout type qui possède une méthode `Error() string` est une erreur. `nil` signifie **« pas d'erreur »**.

Comparaison de philosophie :

| Exceptions (Java, Python) | Go |
|---|---|
| `try / catch`, l'erreur « remonte » invisiblement | L'erreur est retournée et **visible dans la signature** |
| On peut oublier de la traiter | On doit décider quoi en faire à chaque appel |
| Flux de contrôle caché | Flux de contrôle lisible de haut en bas |

## 2.3 Le schéma de base

```go
data, err := os.ReadFile(path)
if err != nil {
    return fmt.Errorf("lecture de %q: %w", path, err)
}
// ici, on est sûr que data est valide
```

Règles d'or :
1. **Tester `err` immédiatement** après l'appel.
2. **Ne jamais ignorer une erreur avec `_`** sauf cas justifié et commenté.
3. **Sortir tôt** (`return` dans le `if err != nil`) pour éviter l'imbrication ; le chemin « normal » reste aligné à gauche.
4. En cas d'erreur, ne pas utiliser les autres valeurs de retour : elles ne sont pas fiables.

## 2.4 Créer des erreurs

| Besoin | Outil |
|---|---|
| Message fixe | `errors.New("fichier vide")` |
| Message avec valeurs | `fmt.Errorf("valeur %q invalide", w)` |
| Ajouter du contexte à une erreur existante | `fmt.Errorf("contexte: %w", err)` |

### Le verbe `%w` (wrapping)
`%w` **enveloppe** l'erreur d'origine : le message est enrichi, mais l'erreur interne reste récupérable.

```go
n, err := strconv.ParseInt(w, 16, 64)
if err != nil {
    return fmt.Errorf("(hex) : %q invalide: %w", w, err)
}
```
Avec `%v` à la place de `%w`, le texte serait identique mais l'erreur d'origine serait **perdue** pour `errors.Is` et `errors.As`.

### Schéma de la chaîne d'erreurs
```
main        : "erreur: traitement: (hex) : "zz" invalide: strconv.ParseInt: parsing "zz": invalid syntax"
  └ run     : "traitement: ..."                       (contexte ajouté)
    └ Process / apply : "(hex) : "zz" invalide: ..."  (contexte ajouté)
      └ strconv       : "strconv.ParseInt: parsing "zz": invalid syntax"  (erreur d'origine)
```
Chaque couche ajoute **où** et **quoi** ; le message final se lit comme une phrase.

## 2.5 Tester une erreur : `errors.Is` et `errors.As`

### Erreurs sentinelles et `errors.Is`
Une **erreur sentinelle** est une valeur d'erreur déclarée une fois, comparée plus tard.

```go
// utils/errors.go
var ErrNoPreviousWord = errors.New("aucun mot précédent")
```
```go
// dans le code appelant
_, err := utils.Process(texte)
if errors.Is(err, utils.ErrNoPreviousWord) {
    fmt.Fprintln(os.Stderr, "le texte ne doit pas commencer par un modificateur")
}
```
`errors.Is` traverse toute la chaîne d'enveloppes `%w`. Autre exemple fréquent :

```go
_, err := os.ReadFile(path)
if errors.Is(err, os.ErrNotExist) {
    fmt.Fprintf(os.Stderr, "le fichier %q n'existe pas\n", path)
}
```
> N'utilisez jamais `err == ErrXxx` ni `err.Error() == "..."` sur une erreur potentiellement enveloppée.

### Types d'erreurs personnalisés et `errors.As`
Quand on veut transporter plusieurs informations (mot, modificateur, position), on crée un type :

```go
type ModifierError struct {
    Modifier string
    Word     string
    Err      error
}

func (e *ModifierError) Error() string {
    return fmt.Sprintf("(%s) sur %q: %v", e.Modifier, e.Word, e.Err)
}

func (e *ModifierError) Unwrap() error { return e.Err }
```
```go
var me *utils.ModifierError
if errors.As(err, &me) {
    fmt.Fprintln(os.Stderr, "mot fautif :", me.Word)
}
```
`errors.As` cherche dans la chaîne une erreur de ce **type** et la place dans `me`.

| Outil | Question posée |
|---|---|
| `errors.Is(err, X)` | « Est-ce l'erreur X (ou enveloppant X) ? » |
| `errors.As(err, &t)` | « Contient-elle une erreur de ce type ? Donne-la-moi. » |

## 2.6 Où traiter l'erreur ? Une seule fois, en haut
Deux attitudes possibles dans une fonction intermédiaire : **retourner** l'erreur (avec contexte) ou **la traiter**. Jamais les deux.

```go
// Mauvais : on journalise ET on retourne -> message en double
if err != nil {
    fmt.Fprintln(os.Stderr, err)
    return err
}

// Bon : on enrichit et on remonte, main affichera une seule fois
if err != nil {
    return fmt.Errorf("traitement: %w", err)
}
```
Règle pratique : les fonctions de `utils` **retournent**, seul `main` **affiche et quitte**.

## 2.7 `panic` : à réserver aux bugs
`panic` arrête brutalement le programme avec une trace. Elle est destinée aux erreurs de **programmation** (index hors limites, `nil` déréférencé), pas aux erreurs **attendues** (fichier absent, saisie invalide).

| Situation | Réponse |
|---|---|
| Fichier introuvable | `error` |
| Texte mal formé | `error` |
| Argument manquant | message d'usage + `os.Exit(1)` |
| Cas « impossible » qui révèle un bug | `panic` (rare) |

## 2.8 Écrire de bons messages d'erreur
Un bon message répond à trois questions : **quoi** (action), **où** (donnée concernée), **pourquoi** (cause).

| Mauvais | Bon |
|---|---|
| `erreur` | `lecture de "in.txt": open in.txt: no such file or directory` |
| `Invalid input!!` | `(hex) : "zz" n'est pas un nombre hexadécimal valide` |
| `Erreur lors de la lecture.` | `écriture de "out.txt": permission denied` |

Conventions Go : message en **minuscule**, sans point final, sans retour à la ligne (les messages sont chaînés avec `: `).

## 2.9 Pièges fréquents (S2)

| Piège | Correction |
|---|---|
| `result, _ := fonction()` | Récupérer et tester l'erreur |
| `fmt.Errorf("... %v", err)` quand on veut garder la cause | Utiliser `%w` |
| Retourner `err` sans contexte sur 5 niveaux | Ajouter `fmt.Errorf("contexte: %w", err)` |
| Utiliser la valeur retournée alors que `err != nil` | Faire un `return` dans le `if` |
| `os.Exit` dans une fonction de `utils` | Retourner une erreur, laisser `main` décider |
| Variable `err` « masquée » (`:=` dans un bloc interne) | Attention au *shadowing* : relire les portées |

## 2.10 Exercices (avec corrigés)

**Exercice C.** Écrire `parseCount(s string) (int, error)` qui convertit `s` en entier strictement positif. Retourner une erreur claire si `s` n'est pas un nombre ou si le nombre est ≤ 0.

```go
func parseCount(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("compteur %q invalide: %w", s, err)
	}
	if n <= 0 {
		return 0, fmt.Errorf("compteur %d invalide: doit être supérieur à 0", n)
	}
	return n, nil
}
```

**Exercice D.** Définir une erreur sentinelle `ErrEmptyInput` et faire échouer `Process` quand le texte ne contient aucun mot. Dans `main`, afficher un message spécifique pour ce cas.

```go
// utils/errors.go
var ErrEmptyInput = errors.New("texte vide")

// utils/process.go
func Process(s string) (string, error) {
	tokens := tokenize(s)
	if len(tokens) == 0 {
		return "", ErrEmptyInput
	}
	// ...
}

// main.go : dans run()
result, err := utils.Process(string(data))
if err != nil {
	if errors.Is(err, utils.ErrEmptyInput) {
		return fmt.Errorf("%q ne contient aucun texte à traiter", in)
	}
	return fmt.Errorf("traitement: %w", err)
}
```
*Remarque pédagogique :* c'est un choix de conception de refuser un fichier vide. Dans le TP1, on l'acceptait. Faire discuter la classe sur ce choix.

**Exercice E (questions).**
1. *Quelle est la différence entre `%v` et `%w` ?* → `%w` conserve l'erreur d'origine pour `errors.Is/As`, `%v` ne conserve que le texte.
2. *Pourquoi `main` est-il le seul à appeler `os.Exit` ?* → Pour centraliser la sortie et garder les autres fonctions testables.
3. *`os.ReadFile` échoue : peut-on utiliser `data` ?* → Non, la valeur n'est pas fiable.

---

# SECTION 3 : Parser du texte (3h)

## 3.1 Objectifs
L'étudiant sait :
- distinguer octets, caractères (runes) et mots en Go ;
- découper un texte en jetons (*tokens*) ;
- appliquer des règles de transformation en pipeline ;
- manipuler des tranches (*slices*) en toute sécurité ;
- utiliser `strings`, `strconv`, `unicode` et `regexp`.

## 3.2 Qu'est-ce que « parser » ?
**Parser** = transformer un texte brut en une structure exploitable, puis agir dessus selon des règles. Notre approche se fait en **étapes**, chacune simple et testable :

```
texte brut
   │  1. Tokenisation       (découper en mots et modificateurs)
   ▼
[]string de jetons
   │  2. Modificateurs      (hex, bin, up, low, cap)
   ▼
[]string transformés
   │  3. Articles           (a → an)
   ▼
[]string
   │  4. Assemblage + ponctuation
   ▼
texte final
```
Avantage : en cas de bug, on identifie **quelle étape** est fautive en affichant l'état entre deux étapes.

## 3.3 Chaînes, octets et runes : le point qui piège tout le monde

- Une `string` Go est une suite **immuable d'octets** encodés en UTF-8.
- `len(s)` renvoie le nombre d'**octets**, pas de caractères.
- Un caractère (*rune*) peut occuper 1 à 4 octets.

```go
s := "été"
fmt.Println(len(s))          // 5 (é = 2 octets)
fmt.Println(len([]rune(s)))  // 3
fmt.Println(s[0])            // 195 : un octet, pas une lettre !

for i, r := range s {        // range sur une string parcourt des RUNES
    fmt.Println(i, string(r)) // positions 0, 2, 4
}
```

| Opération | Quand l'utiliser |
|---|---|
| `s[i]`, `s[a:b]` | Texte uniquement ASCII, ou découpage sur des positions déjà sûres |
| `[]rune(s)` | Accéder ou modifier un caractère à une position |
| `for _, r := range s` | Parcourir caractère par caractère |
| `utf8.RuneCountInString(s)` | Compter les caractères |

Conséquence directe sur `Capitalize` :

```go
func Capitalize(w string) string {
	r := []rune(strings.ToLower(w))   // runes pour gérer "école" -> "École"
	if len(r) == 0 {
		return w                      // garde-fou : mot vide
	}
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
```

## 3.4 Le package `strings` : boîte à outils

| Fonction | Rôle | Exemple |
|---|---|---|
| `strings.Fields(s)` | Découpe sur les espaces (tous types) | `"a  b\n c"` → `[a b c]` |
| `strings.Split(s, sep)` | Découpe sur un séparateur précis | `"a,b"` → `[a b]` |
| `strings.Join(t, sep)` | Assemble | `[a b]`, `" "` → `"a b"` |
| `strings.ToUpper/ToLower` | Casse | |
| `strings.TrimSpace` | Retire espaces aux extrémités | |
| `strings.Contains/HasPrefix/HasSuffix` | Tests | |
| `strings.Builder` | Construction efficace d'une grosse chaîne | |

## 3.5 Conversions numériques : `strconv`

```go
n, err := strconv.ParseInt("1E", 16, 64)  // base 16 -> 30
n, err  = strconv.ParseInt("1010", 2, 64) // base 2  -> 10
s := strconv.FormatInt(30, 10)            // 30 -> "30"
n2, err := strconv.Atoi("42")             // string -> int
```
Toujours traiter `err` : `"zz"` n'est pas un nombre hexadécimal.

## 3.6 Expressions régulières avec `regexp`
Une regexp décrit un **motif** de texte. Go utilise la syntaxe RE2 (pas de retours arrière ni de *lookahead*).

### Éléments de syntaxe utiles

| Motif | Signification |
|---|---|
| `\S+` | Un ou plusieurs caractères non-espace (un « mot ») |
| `\s+` | Un ou plusieurs espaces |
| `\(` `\)` | Parenthèses littérales (échappées) |
| `(hex\|bin\|up)` | Groupe avec alternative |
| `(?:...)` | Groupe non capturant |
| `\d+` | Un ou plusieurs chiffres |
| `?` `*` `+` | 0 ou 1, 0 ou plus, 1 ou plus |
| `^` `$` | Début, fin de chaîne |

### Fonctions principales
```go
re := regexp.MustCompile(`\((?:hex|bin|up|low|cap)(?:,\s*\d+)?\)|\S+`)
tokens := re.FindAllString("Ready (up, 2) go", -1) // [Ready (up, 2) go]

m := modRe.FindStringSubmatch("(up, 3)") // m[0]="(up, 3)", m[1]="up", m[2]="3"
s = spaceBefore.ReplaceAllString(s, "$1")  // $1 = premier groupe capturé
```

> **Pourquoi pas `strings.Fields` pour tokeniser ?** Parce que `(up, 2)` contient un espace : il serait coupé en `(up,` et `2)`. La regexp reconnaît d'abord un modificateur complet, puis seulement un mot ordinaire.

> **Piège :** `"$1a"` est lu comme le groupe nommé `1a`. Écrire `"${1}a"` pour lever l'ambiguïté.

> **Bonne pratique :** compiler les regexp **une seule fois** (variable de package avec `regexp.MustCompile`), jamais dans une boucle.

## 3.7 Tranches (*slices*) : manipuler une fenêtre de mots

Une `slice` est une vue sur un tableau : `t[a:b]` va de `a` (inclus) à `b` (exclu).

```go
out := []string{"this", "is", "so", "exciting"}
n := 2
fenetre := out[len(out)-n:]   // ["so", "exciting"] : les 2 derniers
```

Dans `ApplyModifiers`, on construit `out` au fur et à mesure et on modifie ses **derniers éléments** quand on rencontre un modificateur :

```go
for i := len(out) - n; i < len(out); i++ {
    out[i] = transformer(out[i])
}
```

### Les gardes-fous indispensables
| Cas limite | Risque | Protection |
|---|---|---|
| `(up)` en tout début de texte | `out` est vide, rien à modifier | Retourner une erreur |
| `(up, 5)` avec seulement 3 mots avant | `len(out)-n` négatif → `panic` | `if n > len(out) { n = len(out) }` |
| `(up, 0)` | Aucune modification | Décider : accepter ou rejeter |
| Mot vide | Index `[0]` impossible | Tester `len(r) == 0` |

## 3.8 Le pipeline complet, étape par étape

Texte d'entrée : `it (cap) was a amazing day (up, 2) .`

| Étape | État |
|---|---|
| Tokenisation | `[it (cap) was a amazing day (up, 2) .]` |
| Modificateurs | `[It was a amazing DAY .]` → `(up, 2)` agit sur `amazing` et `day` : `[It was a AMAZING DAY .]` |
| Articles | `a` suivi de `AMAZING` (commence par `a`) → `[It was an AMAZING DAY .]` |
| Assemblage | `It was an AMAZING DAY .` |
| Ponctuation | `It was an AMAZING DAY.` |

*Exercice de discussion :* que se passe-t-il si l'on applique la correction `a/an` **avant** les modificateurs ? Pourquoi l'ordre des étapes est-il important ?
*Réponse :* certains modificateurs (`hex`, `bin`) transforment les mots en nombres ; l'ordre détermine ce que chaque étape « voit ».

## 3.9 Améliorer `FixArticles` pour l'UTF-8
La version vue en TP utilise `[:1]` (premier **octet**). Pour un mot commençant par `é`, cela coupe un caractère en deux. Version correcte :

```go
import "unicode/utf8"

func startsWithVowelOrH(w string) bool {
	r, _ := utf8.DecodeRuneInString(w)       // premier caractère complet
	return strings.ContainsRune("aeiouhAEIOUH", r)
}

func FixArticles(tokens []string) []string {
	for i := 0; i < len(tokens)-1; i++ {
		switch tokens[i] {
		case "a":
			if startsWithVowelOrH(tokens[i+1]) { tokens[i] = "an" }
		case "A":
			if startsWithVowelOrH(tokens[i+1]) { tokens[i] = "An" }
		}
	}
	return tokens
}
```
Le `switch` supprime aussi les appels répétés à `ToLower`.

## 3.10 Pièges fréquents (S3)

| Piège | Explication |
|---|---|
| Supprimer des éléments d'une slice **pendant** qu'on la parcourt avec `range` | Décalage d'indices : construire une nouvelle slice `out` |
| Oublier de gérer `(up, n)` avec `n` trop grand | `panic: slice bounds out of range` |
| Regexp compilée dans une boucle | Lent et inutile |
| `len(s)` pour compter les caractères | Compte des octets |
| `"3.14"` transformé en `"3. 14"` par `FixPunctuation` | Limite connue : à traiter en extension |
| `ToUpper` sur un seul octet d'un caractère accentué | Passer par `[]rune` |

## 3.11 Exercices (avec corrigés)

**Exercice F.** Écrire `CountWords(s string) int` qui compte les mots d'un texte.
```go
func CountWords(s string) int { return len(strings.Fields(s)) }
```

**Exercice G.** Écrire `Reverse(s string) string` qui inverse une chaîne **en respectant l'UTF-8**.
```go
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
// Reverse("été") -> "été" ; Reverse("abc") -> "cba"
```

**Exercice H.** Ajouter le modificateur `(rev)` qui inverse le mot précédent. Indication : ajouter `rev` dans les deux regexp et un `case "rev"` dans `apply`.
```go
tokenRe = regexp.MustCompile(`\((?:hex|bin|up|low|cap|rev)(?:,\s*\d+)?\)|\S+`)
modRe   = regexp.MustCompile(`^\((hex|bin|up|low|cap|rev)(?:,\s*(\d+))?\)$`)
// dans apply :
case "rev":
	return Reverse(w), nil
```

**Exercice I.** Gérer les apostrophes : `' awesome '` → `'awesome'`.
```go
var quotes = regexp.MustCompile(`'\s*([^']*?)\s*'`)

func FixQuotes(s string) string {
	return quotes.ReplaceAllString(s, "'$1'")
}
```
*Limite à discuter :* cette version ne distingue pas les apostrophes de mots (`don't`) des guillemets simples.

---

# SECTION 4 : Qualité, tests et documentation (3h)

## 4.1 Objectifs
L'étudiant sait :
- concevoir un jeu de tests couvrant cas nominaux et cas limites ;
- vérifier une sortie avec `diff` ;
- utiliser `gofmt` et `go vet` ;
- rédiger un README exploitable ;
- relire le code d'un pair avec une grille.

## 4.2 Pourquoi tester ?
Un test répond à la question : « le programme fait-il toujours ce que je pense ? ». Sans test, chaque modification peut casser une règle sans qu'on s'en aperçoive (régression).

## 4.3 Concevoir des cas de test
Pour chaque règle, tester trois familles de cas :

| Famille | Exemple pour `(up)` |
|---|---|
| **Nominal** | `go (up)` → `GO` |
| **Limite** | `(up)` seul ; `(up, 5)` avec 2 mots ; texte vide |
| **Erreur** | `zz (hex)` → erreur explicite |

Ajouter aussi des **cas combinés** : plusieurs modificateurs sur une même ligne, ponctuation collée aux mots, accents.

## 4.4 Tests manuels avec fichiers
Principe : une entrée, une sortie **attendue écrite à la main**, une comparaison automatique.

```bash
go run . tests/sample.txt /tmp/result.txt
diff -u tests/expected.txt /tmp/result.txt && echo OK
```
- `diff` n'affiche rien si les fichiers sont identiques.
- Le `-u` rend les différences lisibles (lignes précédées de `-` et `+`).
- Écrire l'attendu **avant** d'exécuter, sinon on risque de valider un résultat faux par habitude.

## 4.5 Aller plus loin : tests Go automatisés (bonus)
Go intègre un outil de test. Un fichier `xxx_test.go` du même dossier est exécuté par `go test ./...`. Le style **table-driven** est idiomatique :

```go
// utils/process_test.go
package utils

import "testing"

func TestProcess(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"hex", "1E (hex) files", "30 files"},
		{"bin", "10 (bin) years", "2 years"},
		{"up", "go (up) !", "GO!"},
		{"up n", "so exciting (up, 2) .", "SO EXCITING."},
		{"article", "a amazing rock", "an amazing rock"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Process(tt.in)
			if err != nil {
				t.Fatalf("erreur inattendue: %v", err)
			}
			if got != tt.want {
				t.Errorf("Process(%q) = %q, attendu %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestProcessErrors(t *testing.T) {
	if _, err := Process("zz (hex)"); err == nil {
		t.Error("une erreur était attendue pour un hexadécimal invalide")
	}
}
```
Intérêt : ajouter un cas = ajouter **une ligne** dans le tableau.

## 4.6 Les outils de qualité

| Commande | Rôle |
|---|---|
| `gofmt -l .` | Liste les fichiers mal formatés (`gofmt -w .` les corrige) |
| `go vet ./...` | Détecte des erreurs probables (formats `Printf` faux, code inaccessible…) |
| `go build ./...` | Vérifie que tout compile |
| `go test ./...` | Lance les tests |

Habitude à prendre : lancer `gofmt` et `go vet` **avant chaque commit**.

## 4.7 Lisibilité du code

| Pratique | Détail |
|---|---|
| Noms explicites | `hexToDec` plutôt que `conv1` ; noms courts acceptés pour de petites portées (`i`, `n`, `w`) |
| Fonctions courtes | Une fonction = une idée ; si elle dépasse ~40 lignes, la découper |
| Commentaires | Expliquer le **pourquoi**, pas le quoi. Commentaire de doc sur chaque fonction exportée : `// Process applique ...` |
| Pas de code mort | Supprimer le code commenté, il est dans Git |
| Constantes | Éviter les « nombres magiques » : `const maxBase = 36` |

## 4.8 Le README

Un README permet à quelqu'un qui découvre le projet de l'installer et de l'utiliser **sans vous poser de question**.

| Section | Contenu |
|---|---|
| Titre + description | En une ou deux phrases : ce que fait l'outil |
| Prérequis | Version de Go |
| Installation | Commandes de clonage |
| Utilisation | Commande exacte avec arguments |
| Exemples | Entrée → sortie, au moins 3 |
| Règles gérées | Liste des modificateurs et corrections |
| Tests | Comment les lancer |
| Structure | Arborescence commentée |
| Limites connues | Honnêteté sur ce qui n'est pas géré |
| Auteurs | Noms |

*Test du README :* le faire suivre par un camarade qui n'a jamais vu le projet. S'il bloque, le README est incomplet.

## 4.9 Git : bonnes pratiques d'équipe
- Commits **petits et fréquents**, message clair à l'impératif : `Ajoute le modificateur (cap)`.
- Un `.gitignore` pour exclure binaires et fichiers temporaires.
- Ne jamais commiter un fichier de sortie généré.
- Une branche par fonctionnalité si le travail est collectif.

## 4.10 Grille de revue de code (binômes)

| Critère | Points d'attention | Oui / Non |
|---|---|---|
| Architecture | `main` sans logique métier ; packages cohérents | |
| Erreurs | Aucune erreur ignorée ; contexte ajouté avec `%w` | |
| Messages | Clairs, en minuscule, indiquant quoi/où/pourquoi | |
| Parsing | Cas limites gérés (vide, `n` trop grand, UTF-8) | |
| Lisibilité | Noms, fonctions courtes, commentaires utiles | |
| Qualité | `gofmt` et `go vet` sans alerte | |
| Tests | `sample.txt` / `expected.txt` présents et passent | |
| README | Installable et utilisable tel quel | |

Consigne de revue : chaque relecteur rédige **3 points forts** et **3 points d'amélioration**, avec référence au fichier et à la ligne.

## 4.11 Mini-évaluation de fin de séance (corrigée)

1. *Pourquoi écrire le fichier attendu avant d'exécuter le programme ?*
   → Pour ne pas valider par erreur une sortie fausse.
2. *Que fait `gofmt -l .` ?* → Liste les fichiers dont le formatage n'est pas conforme.
3. *Citez deux cas limites pour `(low, n)`.* → `n` supérieur au nombre de mots précédents ; modificateur en début de texte.
4. *Quel est l'avantage des tests table-driven ?* → Ajouter un cas ne demande qu'une ligne, et chaque cas est nommé et isolé.
5. *Citez trois sections indispensables d'un README.* → Description, installation, utilisation.

---

# Annexe : Fiche mémo Go pour le projet

```go
// Arguments
os.Args                              // []string

// Fichiers
data, err := os.ReadFile(p)
err = os.WriteFile(p, data, 0o644)

// Erreurs
errors.New("msg")
fmt.Errorf("contexte: %w", err)
errors.Is(err, cible)
errors.As(err, &typed)

// Chaînes
strings.Fields(s)   strings.Join(t, " ")   strings.ToUpper(s)   strings.TrimSpace(s)
[]rune(s)           string(runes)           unicode.ToUpper(r)

// Nombres
strconv.ParseInt(s, base, 64)   strconv.FormatInt(n, 10)   strconv.Atoi(s)

// Regexp
re := regexp.MustCompile(`motif`)
re.FindAllString(s, -1)   re.FindStringSubmatch(s)   re.ReplaceAllString(s, "$1")

// Sortie
fmt.Fprintln(os.Stderr, "erreur:", err)
os.Exit(1)
```
