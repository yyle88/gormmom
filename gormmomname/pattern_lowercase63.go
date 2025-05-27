package gormmomname

import (
	"regexp"
	"strings"

	"github.com/yyle88/gormmom/internal/simplename"
)

type Lowercase63pattern struct{}

func NewLowercase63pattern() *Lowercase63pattern {
	return &Lowercase63pattern{}
}

func (G *Lowercase63pattern) GetPatternEnum() PatternEnum {
	return "s63" //表示63个小写字符(含数字和下划线)组成的列名
}

func (G *Lowercase63pattern) CheckColumnName(columnName string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_]{1,63}$`).MatchString(columnName)
}

func (G *Lowercase63pattern) BuildColumnName(fieldName string) string {
	columnName := simplename.BuildColumnName(fieldName)
	simplename.CheckLength(columnName, 63)
	return strings.ToLower(columnName)
}
