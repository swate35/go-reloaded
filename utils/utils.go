package utils

import "strings"

// Reconstruct reconstitue la chaîne finale à partir des jetons transformés.
func Reconstruct(tokens []string) string {
	return strings.Join(tokens, "")
}
