// Package gormmomname: Native language column name generation strategies for database compatibility
// Provides pattern-based column name generation from Unicode field names to database-safe identifiers
// Supports multiple naming strategies including lowercase and uppercase patterns with length constraints
// Ensures cross-database compatibility by implementing intersection of various database naming rules
//
// gormmomname: 数据库兼容的原生语言列名生成策略
// 提供基于模式的列名生成，将 Unicode 字段名转换为数据库安全标识符
// 支持多种命名策略，包括带长度约束的小写和大写模式
// 通过实现各种数据库命名规则的交集来确保跨数据库兼容性
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

// PatternEnum represents the column name validation pattern type
// Custom enum type defining validation strategies for field name processing
// Different databases have different column naming rules, so patterns implement intersections
// Ensures compatibility across multiple database systems with unified naming standards
//
// PatternEnum 代表列名验证模式类型
// 自定义枚举类型，定义字段名处理的验证策略
// 由于不同数据库有不同的列命名规则，模式实现了规则交集
// 通过统一命名标准确保多个数据库系统的兼容性
type PatternEnum string

// Pattern defines the interface for column name generation and validation
// Provides pattern identification, name validation, and field-to-column conversion
// Ensures generated column names meet specific database requirements and constraints
//
// Pattern 定义列名生成和验证的接口
// 提供模式识别、名称验证和字段到列的转换
// 确保生成的列名满足特定的数据库要求和约束
type Pattern interface {
	GetPatternEnum() PatternEnum             // Get pattern type identifier // 获取模式类型标识符
	CheckColumnName(columnName string) bool  // Validate column name format // 验证列名格式
	BuildColumnName(fieldName string) string // Generate column name from field name // 从字段名生成列名
}

// Strategies manages multiple column naming patterns with smart defaults
// Contains default pattern selection and strategy mapping for flexible naming
// Supports pattern registration and selection based on requirements
//
// Strategies 管理多种列命名模式，带有智能默认值
// 包含默认模式选择和策略映射，支持灵活命名
// 支持基于需求的模式注册和选择
type Strategies struct {
	defaultPattern Pattern                 // Default naming pattern // 默认命名模式
	nameStrategies map[PatternEnum]Pattern // Pattern strategy mapping // 模式策略映射
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
