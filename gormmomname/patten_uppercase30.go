package gormmomname

import (
	"regexp"
	"strings"
)

type uppercase30pattern struct{}

func (G *uppercase30pattern) IsValidColumnName(columnName string) bool {
	return regexp.MustCompile(`^[A-Z0-9_]{1,30}$`).MatchString(columnName)
}

func (G *uppercase30pattern) GenerateColumnName(fieldName string) string {
	return strings.ToUpper(ensureLength(simpleName(fieldName), 30))
}
