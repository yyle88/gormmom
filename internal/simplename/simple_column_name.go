package simplename

import (
	"unicode"

	"github.com/yyle88/gormmom/internal/unicodehex"
	"github.com/yyle88/printgo"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

// 这是个简单的替换逻辑，能把特殊符号转换为相应的字母（这允许在自定义标签里包含特殊符号）
var punctuationMap = map[string]string{
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

func BuildColumnName(fieldName string) string {
	var pts = printgo.NewPTS()
	var preSimple = true
	for i, c := range fieldName {
		if replacement, ok := punctuationMap[string(c)]; ok {
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

func CheckLength(name string, maxLength int) {
	if len(name) > maxLength {
		zaplog.LOG.Panic("COLUMN-NAME-IS-TOO-LONG", zap.String("column_name", name), zap.Int("column_name_length", len(name)), zap.Int("max_length", maxLength))
	}
}
