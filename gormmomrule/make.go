package gormmomrule

import (
	"strings"

	"github.com/yyle88/erero"
)

// 映射通过go字段名生成表列名，同时允许添加自定义的前缀
var makeColumnNameFunctions = map[RULE]func(fieldName string) string{
	S30:  makeS30,
	S30U: makeS30U,
	S63:  makeS63,
	S63U: makeS63U,
}

func (rule RULE) MakeName(fieldName string) string {
	if makeFunc, exist := makeColumnNameFunctions[rule]; exist {
		return makeFunc(fieldName)
	}
	panic(erero.Errorf("no validation function. fieldName=%s rule_name=%s", fieldName, string(rule)))
}

func MakeName(rule RULE, columnName string, customMakeNames map[RULE]func(string) string) string {
	if len(customMakeNames) > 0 { //优先使用自定义的函数
		if makeFunc, exist := customMakeNames[rule]; exist {
			return makeFunc(columnName)
		}
	}
	return rule.MakeName(columnName)
}

func makeS30(fieldName string) string {
	return checkLen(makeName(fieldName), 30)
}

func makeS30U(fieldName string) string {
	return strings.ToUpper(makeS30(fieldName))
}

func makeS63(fieldName string) string {
	return checkLen(makeName(fieldName), 63)
}

func makeS63U(fieldName string) string {
	return strings.ToUpper(makeS63(fieldName))
}

func checkLen(name string, size int) string {
	if len(name) > size {
		panic(erero.Errorf("column_name=%v is too long. len=%v > size=%v", name, len(name), size))
	}
	return name
}
