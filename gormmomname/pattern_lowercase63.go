package gormmomname

import (
	"strings"

	"github.com/yyle88/gormmom/internal/simplename"
	"github.com/yyle88/gormmom/internal/utils"
)

type Lowercase63pattern struct{}

func NewLowercase63pattern() *Lowercase63pattern {
	return &Lowercase63pattern{}
}

func (G *Lowercase63pattern) GetPatternEnum() PatternEnum {
	return "s63" //表示63个小写字符(含数字和下划线)组成的列名
}

func (G *Lowercase63pattern) CheckColumnName(columnName string) bool {
	return utils.NewCommonRegexp(63).MatchString(columnName)
}

func (G *Lowercase63pattern) BuildColumnName(fieldName string) string {
	columnName := strings.ToLower(simplename.BuildColumnName(fieldName))
	utils.MustMatchRegexp(utils.NewCommonRegexp(63), columnName)
	return columnName
}
