package gormmomrule

import (
	"regexp"
	"strings"
)

type nameS63UImp struct {
}

func (G *nameS63UImp) GenNewCnm(fieldName string) string {
	return strings.ToUpper(checkLen(makeName(fieldName), 63))
}

func (G *nameS63UImp) CheckName(columnName string) bool {
	return regexp.MustCompile(`^[A-Z0-9_]{1,63}$`).MatchString(columnName)
}
