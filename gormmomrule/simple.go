package gormmomrule

import (
	"strings"
	"unicode"

	"github.com/yyle88/erero"
	"github.com/yyle88/gormmom/internal/enctohex"
)

// 这是个简单的替换逻辑，能把特殊符号转换为相应的字母（但似乎也没有这个必要，因为字段名也不会包含这些字符）
var replacementsMap = map[string]string{
	"(": "x",
	")": "x",
	"（": "x",
	"）": "x",
	"%": "p",
	".": "t",
	"|": "n",
	"、": "c",
	"/": "k",
	":": "c",
	"：": "c",
}

func makeName(fieldName string) string {
	var res strings.Builder
	var preSimple = true
	for i, c := range fieldName {
		chs := string(c)
		if replacement, ok := replacementsMap[chs]; ok {
			if i > 0 { //非首个时需要添加下划线以转换成蛇形的
				res.WriteRune('_')
			}
			res.WriteString(replacement)
			preSimple = false
		} else if c <= unicode.MaxASCII {
			nowSimple := true
			if unicode.IsUpper(c) {
				if i > 0 { //非首个时需要添加下划线以转换成蛇形的
					res.WriteRune('_')
				}
				res.WriteRune(unicode.ToLower(c))
			} else if unicode.IsLower(c) {
				if !preSimple {
					res.WriteRune('_')
				}
				res.WriteRune(c)
			} else if unicode.IsDigit(c) {
				if !preSimple {
					res.WriteRune('_')
				}
				res.WriteRune(c)
			} else if c == '_' {
				res.WriteRune(c)
			} else {
				if i > 0 { //非首个时需要添加下划线以转换成蛇形的
					res.WriteRune('_')
				}
				res.WriteString("x")
				nowSimple = false //这种情况不能当作普通字符，而是要当成特殊字符
			}
			preSimple = nowSimple //这表示前一次是个普通的简单字符
		} else {
			if i > 0 { //非首个时需要添加下划线以转换成蛇形的
				res.WriteRune('_')
			} else {
				res.WriteString("column")
				res.WriteRune('_')
			}
			res.WriteString(enctohex.Uint32ToHex4Los(c))
			preSimple = false
		}
	}
	return res.String()
}

func checkLen(name string, size int) string {
	if len(name) > size {
		panic(erero.Errorf("column_name=%v is too long. len=%v > size=%v", name, len(name), size))
	}
	return name
}
