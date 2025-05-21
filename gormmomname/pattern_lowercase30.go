package gormmomname

import (
	"regexp"
	"strings"

	"github.com/yyle88/gormmom/internal/simplename"
)

type Lowercase30pattern struct{}

func NewLowercase30pattern() *Lowercase30pattern {
	return &Lowercase30pattern{}
}

func (G *Lowercase30pattern) GetPatternEnum() PatternEnum {
	return "s30" //表示30个小写字符(含数字和下划线)组成的列名
}

func (G *Lowercase30pattern) CheckColumnName(columnName string) bool {
	//当列名前带个前导空格 `gorm:"column: name;"` 时，在gorm中也是可以用的，但有可能存在问题，因此该规则里不允许这种情况
	return regexp.MustCompile(`^[a-zA-Z0-9_]{1,30}$`).MatchString(columnName)
}

func (G *Lowercase30pattern) BuildColumnName(fieldName string) string {
	columnName := simplename.BuildSimpleName(fieldName)
	simplename.CheckLength(columnName, 30)
	return strings.ToLower(columnName)
}
