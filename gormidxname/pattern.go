package gormidxname

import (
	"github.com/yyle88/must"
	"gorm.io/gorm/schema"
)

// nolint:no-doc
// 自定义枚举类型
type PatternEnum string

type Pattern interface {
	GetPatternEnum() PatternEnum
	CheckIndexName(indexName string) bool
	BuildIndexName(schemaIndex *schema.Index, param *BuildIndexParam) *IndexNameResult
}

type BuildIndexParam struct {
	TableName  string
	FieldName  string
	ColumnName string
}

type IndexNameResult struct {
	TagFieldName string
	NewIndexName string
	IdxUdxPrefix string
}

type Strategies struct {
	defaultPattern Pattern
	nameStrategies map[PatternEnum]Pattern
}

func NewStrategies() *Strategies {
	nameStrategies := make(map[PatternEnum]Pattern)
	for _, pattern := range []Pattern{
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
