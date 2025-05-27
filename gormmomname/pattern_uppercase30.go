package gormmomname

import (
	"regexp"
	"strings"

	"github.com/yyle88/gormmom/internal/simplename"
)

type Uppercase30pattern struct{}

func NewUppercase30pattern() *Uppercase30pattern {
	return &Uppercase30pattern{}
}

func (G *Uppercase30pattern) GetPatternEnum() PatternEnum {
	return "S30" //表示30个大写字符(含数字和下划线)组成的列名
}

func (G *Uppercase30pattern) CheckColumnName(columnName string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_]{1,30}$`).MatchString(columnName)
}

func (G *Uppercase30pattern) BuildColumnName(fieldName string) string {
	columnName := simplename.BuildColumnName(fieldName)
	simplename.CheckLength(columnName, 30)
	return strings.ToUpper(columnName)
}
