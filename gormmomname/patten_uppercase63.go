package gormmomname

import (
	"regexp"
	"strings"
)

type uppercase63pattern struct{}

func (G *uppercase63pattern) IsValidColumnName(columnName string) bool {
	return regexp.MustCompile(`^[A-Z0-9_]{1,63}$`).MatchString(columnName)
}

func (G *uppercase63pattern) GenerateColumnName(fieldName string) string {
	return strings.ToUpper(ensureLength(simpleName(fieldName), 63))
}
