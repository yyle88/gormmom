package gormmomrule

import "regexp"

type nameS30Imp struct {
}

func (G *nameS30Imp) GenNewCnm(fieldName string) string {
	return checkLen(makeName(fieldName), 30)
}

func (G *nameS30Imp) CheckName(columnName string) bool {
	//当列名前带个前导空格 `gorm:"column: name;"` 时，在gorm中也是可以用的，但有可能存在问题，因此该规则里不允许这种情况
	return regexp.MustCompile(`^[a-zA-Z0-9_]{1,30}$`).MatchString(columnName)
}
