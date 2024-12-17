package gormmomname

import (
	"unicode"

	"github.com/yyle88/erero"
	"github.com/yyle88/gormmom/internal/unicodehex"
	"github.com/yyle88/printgo"
)

// 这是个简单的替换逻辑，能把特殊符号转换为相应的字母（但似乎也没有这个必要，因为字段名也不会包含这些字符）
var punctuationReplacementMap = map[string]string{
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

func simpleName(fieldName string) string {
	var pts = printgo.NewPTS()
	var preSimple = true
	for i, c := range fieldName {
		chs := string(c)
		if replacement, ok := punctuationReplacementMap[chs]; ok {
			if i > 0 { //非首个时需要添加下划线以转换成蛇形的
				pts.WriteRune('_')
			}
			pts.WriteString(replacement)
			preSimple = false
		} else if c <= unicode.MaxASCII {
			nowSimple := true
			if unicode.IsUpper(c) {
				if i > 0 { //非首个时需要添加下划线以转换成蛇形的
					pts.WriteRune('_')
				}
				pts.WriteRune(unicode.ToLower(c))
			} else if unicode.IsLower(c) {
				if !preSimple {
					pts.WriteRune('_')
				}
				pts.WriteRune(c)
			} else if unicode.IsDigit(c) {
				if !preSimple {
					pts.WriteRune('_')
				}
				pts.WriteRune(c)
			} else if c == '_' {
				pts.WriteRune(c)
			} else {
				if i > 0 { //非首个时需要添加下划线以转换成蛇形的
					pts.WriteRune('_')
				}
				pts.WriteString("x")
				nowSimple = false //这种情况不能当作普通字符，而是要当成特殊字符
			}
			preSimple = nowSimple //这表示前一次是个普通的简单字符
		} else {
			if i > 0 { //非首个时需要添加下划线以转换成蛇形的
				pts.WriteRune('_')
			} else {
				pts.WriteString("column")
				pts.WriteRune('_')
			}
			pts.WriteString(unicodehex.Uint32ToHex4Lowercase(c))
			preSimple = false
		}
	}
	return pts.String()
}

func ensureLength(name string, size int) string {
	if len(name) > size {
		panic(erero.Errorf("column_name=%v is too long. len=%v > size=%v", name, len(name), size))
	}
	return name
}
