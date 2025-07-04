package gormmomname

import (
	"github.com/yyle88/must"
)

//type CaseMode string
//
//const (
//	Lowercase CaseMode = "LOWERCASE"
//	Uppercase CaseMode = "UPPERCASE"
//)

// nolint:no-doc
// 自定义枚举类型，表示使用何种字段验证方式来验证，由于不同的DB的列名规则是不同的，因此通常建议是取各种DB的交集
type PatternEnum string

type Pattern interface {
	GetPatternEnum() PatternEnum
	CheckColumnName(columnName string) bool
	BuildColumnName(fieldName string) string
}

type Strategies struct {
	defaultPattern Pattern
	nameStrategies map[PatternEnum]Pattern
}

func NewStrategies() *Strategies {
	nameStrategies := make(map[PatternEnum]Pattern)
	for _, pattern := range []Pattern{
		NewLowercase30pattern(),
		NewUppercase30pattern(),
		NewLowercase63pattern(),
		NewUppercase63pattern(),
	} {
		nameStrategies[pattern.GetPatternEnum()] = pattern
	}
	return &Strategies{
		defaultPattern: NewLowercase63pattern(),
		nameStrategies: nameStrategies,
	}
}

func (s *Strategies) GetDefault() Pattern {
	return s.defaultPattern
}

func (s *Strategies) SetDefault(pattern Pattern) {
	s.defaultPattern = pattern
	s.SetPattern(pattern) // 同时也要把它设置到 map 里面
}

func (s *Strategies) GetPattern(patternEnum PatternEnum) Pattern {
	namePattern, ok := s.nameStrategies[patternEnum]
	must.True(ok)
	return namePattern
}

func (s *Strategies) SetPattern(pattern Pattern) {
	s.nameStrategies[pattern.GetPatternEnum()] = pattern
}
