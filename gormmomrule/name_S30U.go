package gormmomrule

import (
	"regexp"
	"strings"
)

type nameS30UImp struct {
}

func (G *nameS30UImp) GenNewCnm(fieldName string) string {
	return strings.ToUpper(checkLen(makeName(fieldName), 30))
}

func (G *nameS30UImp) CheckName(columnName string) bool {
	return regexp.MustCompile(`^[A-Z0-9_]{1,30}$`).MatchString(columnName)
}
