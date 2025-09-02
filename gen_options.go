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
	columnNamingSubTagName string                  // Sub-tag name for column naming patterns // 列命名模式的子标签名称
	columnNamingStrategies *gormmomname.Strategies // Column naming strategies with native languages // 原生语言的列命名策略
	skipBasicColumnName    bool                    // Skip simple fields without custom configurations // 跳过没有自定义配置的简单字段
	renewIndexName         bool                    // Enable index name generation // 启用索引名称生成
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
		columnNamingSubTagName: "mcp", // m(mother_tongue) c(column_name) p(pattern) // 表示列名的样式标签
		columnNamingStrategies: gormmomname.NewStrategies(),
		skipBasicColumnName:    true,
		renewIndexName:         true,
		indexNamingStrategies:  gormidxname.NewStrategies(),
	}
}

func (opt *Options) WithTagName(systemTagName string) *Options {
	opt.systemTagName = systemTagName
	return opt
}

func (opt *Options) WithColumnNamingSubTagName(columnNamingSubTagName string) *Options {
	opt.columnNamingSubTagName = columnNamingSubTagName
	return opt
}

func (opt *Options) WithColumnPattern(pattern gormmomname.Pattern) *Options {
	opt.columnNamingStrategies.SetPattern(pattern)
	return opt
}

func (opt *Options) WithDefaultColumnPattern(pattern gormmomname.Pattern) *Options {
	opt.columnNamingStrategies.SetDefault(pattern)
	return opt
}

func (opt *Options) WithSkipBasicColumnName(skipBasicColumnName bool) *Options {
	opt.skipBasicColumnName = skipBasicColumnName
	return opt
}

func (opt *Options) WithRenewIndexName(renewIndexName bool) *Options {
	opt.renewIndexName = renewIndexName
	return opt
}

func (opt *Options) WithIndexPattern(pattern gormidxname.Pattern) *Options {
	opt.indexNamingStrategies.SetPattern(pattern)
	return opt
}

func (opt *Options) WithDefaultIndexPattern(pattern gormidxname.Pattern) *Options {
	opt.indexNamingStrategies.SetDefault(pattern)
	return opt
}
