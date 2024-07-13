package gormmomrule

import (
	"strings"
	"unicode"

	"github.com/yyle88/gormmom/internal/encnm"
)

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
			res.WriteString(encnm.Uint32ToHex4Los(c))
			preSimple = false
		}
	}
	return res.String()
}
