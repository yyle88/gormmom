package gormmomname

import (
	"strings"

	"github.com/yyle88/gormmom/internal/simplename"
	"github.com/yyle88/gormmom/internal/utils"
)

type Uppercase30pattern struct{}

func NewUppercase30pattern() *Uppercase30pattern {
	return &Uppercase30pattern{}
}

func (G *Uppercase30pattern) GetPatternEnum() PatternEnum {
	return "S30" //表示30个大写字符(含数字和下划线)组成的列名
}

func (G *Uppercase30pattern) CheckColumnName(columnName string) bool {
	return utils.NewCommonRegexp(30).MatchString(columnName)
}

func (G *Uppercase30pattern) BuildColumnName(fieldName string) string {
	columnName := strings.ToUpper(simplename.BuildColumnName(fieldName))
	utils.MustMatchRegexp(utils.NewCommonRegexp(30), columnName)
	return columnName
}
