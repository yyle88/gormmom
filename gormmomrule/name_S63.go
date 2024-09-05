package gormmomrule

import "regexp"

type nameS63Imp struct {
}

func (G *nameS63Imp) GenNewCnm(fieldName string) string {
	return checkLen(makeName(fieldName), 63)
}

func (G *nameS63Imp) CheckName(columnName string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_]{1,63}$`).MatchString(columnName)
}
