# Go Reloaded

Go Reloaded est un outil en ligne de commande écrit en Go. Il lit un fichier
texte, applique des transformations et écrit le résultat dans un fichier de
sortie.

## Prérequis

- Go installé (version indiquée dans [`go.mod`](./go.mod)).
- Aucune dépendance externe : le projet utilise uniquement la bibliothèque
  standard de Go.

## Utilisation

À la racine du projet, exécutez :

```sh
go run . <fichier-entree> <fichier-sortie>
```

Par exemple :

```sh
go run . sample.txt result.txt
```

Le fichier d’entrée doit exister. Le fichier de sortie est créé ou remplacé.

## Transformations prises en charge

Les marqueurs s’écrivent après le mot ou les mots à transformer :

| Marqueur | Résultat |
| --- | --- |
| `(hex)` | Convertit le nombre hexadécimal précédent en décimal. |
| `(bin)` | Convertit le nombre binaire précédent en décimal. |
| `(up)` | Met le mot précédent en majuscules. |
| `(low)` | Met le mot précédent en minuscules. |
| `(cap)` | Met la première lettre du mot précédent en majuscule. |
| `(up, N)` | Met les `N` mots précédents en majuscules. |
| `(low, N)` | Met les `N` mots précédents en minuscules. |
| `(cap, N)` | Met une majuscule initiale aux `N` mots précédents. |

Exemple :

```text
Simply add 42 (hex) and 10 (bin). This is exciting (up, 2)!
```

Résultat :

```text
Simply add 66 and 2. This IS EXCITING!
```

Le programme normalise également les signes `. , ! ? : ;` en les collant au
texte précédent et en ajoutant un espace après le groupe de ponctuation si un
mot suit. Les groupes comme `...`, `!?` et `!!` restent regroupés. Les espaces
autour des apostrophes de citation sont supprimés, et l’article `a` devient
`an` devant un mot commençant par une voyelle ou un `h`.

## Tests

Pour exécuter tous les tests depuis la racine du projet :

```sh
go test ./...
```

Pour afficher le détail des tests :

```sh
go test ./... -v
```

## Organisation du projet

- `main.go` : lit les arguments, charge le texte et écrit le résultat.
- `parser/` : découpe le texte en jetons.
- `transform/` : applique les marqueurs et normalise le texte.
- `utils/` : reconstruit le texte à partir des jetons.
- `test/` : contient les tests du traitement du texte.