package gormmom

import (
	"github.com/yyle88/gormmom/gormidxname"
	"github.com/yyle88/gormmom/gormmomname"
)

// Options represents configuration settings for GORM tag generation
// Contains naming strategies and behavior controls with native language field processing
// Provides customizable tag names, column naming rules, and index generation settings
// Supports smart skipping of basic fields and flexible configuration options
//
// Options 代表 GORM 标签生成的配置设置
// 包含原生语言字段处理的命名策略和行为控制
// 提供可自定义的标签名称、列命名规则和索引生成设置
// 支持智能跳过基础字段和灵活的配置选项
type Options struct {
	systemTagName          string                  // System tag name // 系统标签名称
	subTagName             string                  // Sub-tag name for column naming patterns // 列命名模式的子标签名称
	columnNamingStrategies *gormmomname.Strategies // Column naming strategies with native languages // 原生语言的列命名策略
	skipBasicColumnName    bool                    // Skip simple fields without custom configurations // 跳过没有自定义配置的简单字段
	autoIndexName          bool                    // Enable index name generation // 启用索引名称生成
	indexNamingStrategies  *gormidxname.Strategies // Index naming strategies with database optimization // 数据库优化的索引命名策略
}

// NewOptions creates a new Options instance with default settings optimized with native language processing
// Initializes with smart defaults including tag naming, column strategies, and index generation
// Returns configured options prepared for GORM tag generation with Unicode field support
//
// NewOptions 创建新的选项实例，使用针对原生语言处理优化的默认设置
// 使用智能默认值初始化，包括标签命名、列策略和索引生成
// 返回配置好的选项，准备进行支持 Unicode 字段的 GORM 标签生成
func NewOptions() *Options {
	return &Options{
		systemTagName:          "mom", // mother_tongue native_language. // 因为我也不太能熟练拼写更高级的单词，还不如返璞归真直接用口语表示
		subTagName:             "mcp", // m(mother_tongue) c(column_name) p(pattern) // 表示列名的样式标签
		columnNamingStrategies: gormmomname.NewStrategies(),
		skipBasicColumnName:    true,
		autoIndexName:          true,
		indexNamingStrategies:  gormidxname.NewStrategies(),
	}
}

// WithTagName sets the system tag name used in struct tags
// Configures the tag key name that gormmom uses to store pattern information
// Returns the Options instance to enable method chaining
//
// WithTagName 设置结构体标签中使用的系统标签名
// 配置 gormmom 用于存储模式信息的标签键名
// 返回 Options 实例以启用方法链
func (opt *Options) WithTagName(systemTagName string) *Options {
	opt.systemTagName = systemTagName
	return opt
}

// WithSubTagName sets the sub-tag name for column naming patterns
// Configures the nested tag field name within the system tag
// Returns the Options instance to enable method chaining
//
// WithSubTagName 设置列命名模式的子标签名
// 配置系统标签内的嵌套标签字段名
// 返回 Options 实例以启用方法链
func (opt *Options) WithSubTagName(subTagName string) *Options {
	opt.subTagName = subTagName
	return opt
}

// WithColumnPattern registers a custom column naming pattern
// Adds the pattern to the strategies collection for column name generation
// Returns the Options instance to enable method chaining
//
// WithColumnPattern 注册自定义的列命名模式
// 将模式添加到列名生成的策略集合中
// 返回 Options 实例以启用方法链
func (opt *Options) WithColumnPattern(pattern gormmomname.Pattern) *Options {
	opt.columnNamingStrategies.SetPattern(pattern)
	return opt
}

// WithDefaultColumnPattern sets the default column naming pattern
// Used when no specific pattern is configured in struct tags
// Returns the Options instance to enable method chaining
//
// WithDefaultColumnPattern 设置默认的列命名模式
// 当结构体标签中没有配置特定模式时使用
// 返回 Options 实例以启用方法链
func (opt *Options) WithDefaultColumnPattern(pattern gormmomname.Pattern) *Options {
	opt.columnNamingStrategies.SetDefault(pattern)
	return opt
}

// WithSkipBasicColumnName enables or disables skipping of basic column names
// When true, fields with standard ASCII column names are skipped
// Returns the Options instance to enable method chaining
//
// WithSkipBasicColumnName 启用或禁用跳过基本列名
// 当为 true 时，具有标准 ASCII 列名的字段被跳过
// 返回 Options 实例以启用方法链
func (opt *Options) WithSkipBasicColumnName(skipBasicColumnName bool) *Options {
	opt.skipBasicColumnName = skipBasicColumnName
	return opt
}

// WithAutoIndexName enables or disables index name regeneration
// When true, index names are regenerated based on configured patterns
// Returns the Options instance to enable method chaining
//
// WithAutoIndexName 启用或禁用索引名重新生成
// 当为 true 时，索引名基于配置的模式重新生成
// 返回 Options 实例以启用方法链
func (opt *Options) WithAutoIndexName(autoIndexName bool) *Options {
	opt.autoIndexName = autoIndexName
	return opt
}

// WithIndexPattern registers a custom index naming pattern
// Adds the pattern to the strategies collection for index name generation
// Returns the Options instance to enable method chaining
//
// WithIndexPattern 注册自定义的索引命名模式
// 将模式添加到索引名生成的策略集合中
// 返回 Options 实例以启用方法链
func (opt *Options) WithIndexPattern(pattern gormidxname.Pattern) *Options {
	opt.indexNamingStrategies.SetPattern(pattern)
	return opt
}

// WithDefaultIndexPattern sets the default index naming pattern
// Used when no specific pattern is configured for index generation
// Returns the Options instance to enable method chaining
//
// WithDefaultIndexPattern 设置默认的索引命名模式
// 当索引生成没有配置特定模式时使用
// 返回 Options 实例以启用方法链
func (opt *Options) WithDefaultIndexPattern(pattern gormidxname.Pattern) *Options {
	opt.indexNamingStrategies.SetDefault(pattern)
	return opt
}
