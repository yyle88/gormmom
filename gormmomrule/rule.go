package gormmomrule

import (
	"regexp"

	"github.com/yyle88/erero"
)

// nolint:no-doc
// 自定义枚举类型，表示使用何种字段验证方式来验证，由于不同的DB的列名规则是不同的，因此通常建议是取各种DB的交集
type Rule string

const (
	S30  Rule = "S30"
	S30U Rule = "S30U"
	S63  Rule = "S63"
	S63U Rule = "S63U"
)

// 映射验证函数
var validationFunctions = map[Rule]func(string) bool{
	S30:  regexp.MustCompile(`^[a-zA-Z0-9_]{1,30}$`).MatchString, //试过名字前带个前导空格 `gorm:"column: name;"` 也是可以的，但这个规则里不允许这种情况
	S30U: regexp.MustCompile(`^[A-Z0-9_]{1,30}$`).MatchString,
	S63:  regexp.MustCompile(`^[a-zA-Z0-9_]{1,63}$`).MatchString,
	S63U: regexp.MustCompile(`^[A-Z0-9_]{1,63}$`).MatchString,
}

func (rule Rule) Validate(columnName string) bool {
	if check, exist := validationFunctions[rule]; exist {
		return check(columnName)
	}
	panic(erero.Errorf("no validation function. column_name=%s rule_name=%s", columnName, string(rule)))
}

func Validate(rule Rule, columnName string, customValidations map[Rule]func(string) bool) bool {
	if rule == "" {
		return true
	}
	if len(customValidations) > 0 { //优先使用自定义的函数
		if check, exist := customValidations[rule]; exist {
			return check(columnName)
		}
	}
	return rule.Validate(columnName)
}
