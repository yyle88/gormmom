// Package gormidxname: Database index naming strategy engine for GORM optimization
// Provides intelligent index name generation from Unicode field names with pattern-based validation
// Supports multiple naming strategies including lowercase and uppercase patterns with length constraints
// Ensures cross-database compatibility by implementing intersection of various database index naming rules
//
// gormidxname: GORM 优化的数据库索引命名策略引擎
// 提供基于模式验证的 Unicode 字段名智能索引名生成
// 支持多种命名策略，包括带长度约束的小写和大写模式
// 通过实现各种数据库索引命名规则的交集来确保跨数据库兼容性
package gormidxname

import (
	"github.com/yyle88/must"
	"gorm.io/gorm/schema"
)

// IndexPatternTagEnum represents the index pattern tag type identifier
// Defines enum values for different index types with specific prefixes
// Used to distinguish between standard indexes and unique indexes
//
// IndexPatternTagEnum 代表索引模式标签类型标识符
// 为不同索引类型定义带特定前缀的枚举值
// 用于区分标准索引和唯一索引
type IndexPatternTagEnum string

const (
	IdxPatternTagName IndexPatternTagEnum = "idx" // Standard index prefix // 标准索引前缀
	UdxPatternTagName IndexPatternTagEnum = "udx" // Unique index prefix // 唯一索引前缀
)

// PatternEnum represents the index name validation pattern type
// Custom enum type defining validation strategies for index name processing
// Different databases have different index naming rules, so patterns implement intersections
// Ensures compatibility across multiple database systems with unified naming standards
//
// PatternEnum 代表索引名验证模式类型
// 自定义枚举类型，定义索引名处理的验证策略
// 由于不同数据库有不同的索引命名规则，模式实现了规则交集
// 通过统一命名标准确保多个数据库系统的兼容性
type PatternEnum string

// Pattern defines the interface for index name generation and validation
// Provides pattern identification, name validation, and index name construction
// Ensures generated index names meet specific database requirements and constraints
//
// Pattern 定义索引名生成和验证的接口
// 提供模式识别、名称验证和索引名构建
// 确保生成的索引名满足特定的数据库要求和约束
type Pattern interface {
	GetPatternEnum() PatternEnum                                                       // Get pattern type identifier // 获取模式类型标识符
	CheckIndexName(indexName string) bool                                              // Validate index name format // 验证索引名格式
	BuildIndexName(schemaIndex *schema.Index, param *BuildIndexParam) *IndexNameResult // Generate index name from schema // 从模式生成索引名
}

// BuildIndexParam contains parameters for index name construction
// Provides table, field, and column information for intelligent naming
// Used as input for pattern-based index name generation algorithms
//
// BuildIndexParam 包含索引名构建的参数
// 提供表、字段和列信息以进行智能命名
// 用作基于模式的索引名生成算法的输入
type BuildIndexParam struct {
	TableName  string // Database table name // 数据库表名
	FieldName  string // Struct field name // 结构体字段名
	ColumnName string // Database column name // 数据库列名
}

// IndexNameResult contains the result of index name generation
// Provides generated tag field name, index name, and prefix type
// Used as output from pattern-based index name construction algorithms
//
// IndexNameResult 包含索引名生成的结果
// 提供生成的标签字段名、索引名和前缀类型
// 用作基于模式的索引名构建算法的输出
type IndexNameResult struct {
	TagFieldName string              // Generated tag field name // 生成的标签字段名
	NewIndexName string              // Generated index name // 生成的索引名
	IdxUdxPrefix IndexPatternTagEnum // Index type prefix // 索引类型前缀
}

// Strategies manages multiple index naming patterns with smart defaults
// Contains default pattern selection and strategy mapping for flexible naming
// Supports pattern registration and selection based on database requirements
//
// Strategies 管理多种索引命名模式，带有智能默认值
// 包含默认模式选择和策略映射，支持灵活命名
// 支持基于数据库需求的模式注册和选择
type Strategies struct {
	defaultPattern Pattern                 // Default naming pattern // 默认命名模式
	nameStrategies map[PatternEnum]Pattern // Pattern strategy mapping // 模式策略映射
}

// NewStrategies creates a new Strategies instance with default index naming patterns
// Initializes with lowercase and uppercase patterns optimized for database compatibility
// Returns configured strategies prepared for index name generation with Unicode field support
//
// NewStrategies 创建新的策略实例，使用默认的索引命名模式
// 使用针对数据库兼容性优化的小写和大写模式进行初始化
// 返回配置好的策略，准备进行支持 Unicode 字段的索引名生成
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
