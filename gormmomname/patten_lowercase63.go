package gormmomname

import "regexp"

type lowercase63pattern struct{}

func (G *lowercase63pattern) IsValidColumnName(columnName string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_]{1,63}$`).MatchString(columnName)
}

func (G *lowercase63pattern) GenerateColumnName(fieldName string) string {
	return ensureLength(simpleName(fieldName), 63)
}
