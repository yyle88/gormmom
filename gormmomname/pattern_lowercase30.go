package gormmomname

import (
	"strings"

	"github.com/yyle88/gormmom/internal/simplename"
	"github.com/yyle88/gormmom/internal/utils"
)

type Lowercase30pattern struct{}

func NewLowercase30pattern() *Lowercase30pattern {
	return &Lowercase30pattern{}
}

func (G *Lowercase30pattern) GetPatternEnum() PatternEnum {
	return "s30" //表示30个小写字符(含数字和下划线)组成的列名
}

func (G *Lowercase30pattern) CheckColumnName(columnName string) bool {
	return utils.NewCommonRegexp(30).MatchString(columnName)
}

func (G *Lowercase30pattern) BuildColumnName(fieldName string) string {
	columnName := strings.ToLower(simplename.BuildColumnName(fieldName))
	utils.MustMatchRegexp(utils.NewCommonRegexp(30), columnName)
	return columnName
}
