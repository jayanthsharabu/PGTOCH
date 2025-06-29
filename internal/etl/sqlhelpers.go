package etl

import (
	"regexp"
	"strings"
)

func QuoteIdentifier(identifier string) string {
	escaped := strings.ReplaceAll(identifier, `"`, `""`)
	return `"` + escaped + `"`
}

func IsValidIdentifier(identifier string) bool {
	pattern := regexp.MustCompile(`^[a-zA-Z0-9_\.]+$`)
	return pattern.MatchString(identifier)
}
