package gormmomname

import (
	"regexp"
	"strings"

	"github.com/yyle88/gormmom/internal/simplename"
)

type Uppercase63pattern struct{}

func NewUppercase63pattern() *Uppercase63pattern {
	return &Uppercase63pattern{}
}

func (G *Uppercase63pattern) GetPatternEnum() PatternEnum {
	return "S63" //表示63个大写字符(含数字和下划线)组成的列名
}

func (G *Uppercase63pattern) CheckColumnName(columnName string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_]{1,63}$`).MatchString(columnName)
}

func (G *Uppercase63pattern) BuildColumnName(fieldName string) string {
	columnName := simplename.BuildSimpleName(fieldName)
	simplename.CheckLength(columnName, 63)
	return strings.ToUpper(columnName)
}
