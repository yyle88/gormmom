package gormmom

import (
	"github.com/yyle88/gormmom/gormidxname"
	"github.com/yyle88/gormmom/gormmomname"
)

type Options struct {
	namingTagName              string
	columnNamePatternFieldName string
	defaultColumnNamePattern   gormmomname.ColumnNamePattern //默认检查规则
	columnNamingStrategies     map[gormmomname.ColumnNamePattern]gormmomname.Naming
	skipBasicNaming            bool //是否跳过简单字段，有的字段虽然没有配置名称或者规则，但是它满足简单字段，就也不做任何处理
	renewIndexName             bool
	indexNamingStrategies      map[gormidxname.IndexNamePattern]gormidxname.Naming
}

func NewOptions() *Options {
	return &Options{
		namingTagName:              "mom",
		columnNamePatternFieldName: "naming",
		defaultColumnNamePattern:   gormmomname.DefaultPattern, //默认检查规则，就是查看是不是63个合法字符（即字母数组下划线等）
		columnNamingStrategies:     gormmomname.GetNamingStrategies(),
		skipBasicNaming:            true,
		renewIndexName:             true,
		indexNamingStrategies:      gormidxname.GetNamingStrategies(),
	}
}

func (cfg *Options) SetNamingTagName(namingTagName string) *Options {
	cfg.namingTagName = namingTagName
	return cfg
}

func (cfg *Options) SetColumnNamePatternFieldName(columnNamePatternFieldName string) *Options {
	cfg.columnNamePatternFieldName = columnNamePatternFieldName
	return cfg
}

func (cfg *Options) SetDefaultColumnNamePattern(pattern gormmomname.ColumnNamePattern) *Options {
	cfg.defaultColumnNamePattern = pattern
	return cfg
}

func (cfg *Options) SetColumnNamingStrategies(pattern gormmomname.ColumnNamePattern, naming gormmomname.Naming) *Options {
	cfg.columnNamingStrategies[pattern] = naming
	return cfg
}

func (cfg *Options) SetSkipBasicNaming(skipBasicNaming bool) *Options {
	cfg.skipBasicNaming = skipBasicNaming
	return cfg
}

func (cfg *Options) SetRenewIndexName(renewIndexName bool) *Options {
	cfg.renewIndexName = renewIndexName
	return cfg
}

func (cfg *Options) SetIndexNamingStrategies(pattern gormidxname.IndexNamePattern, naming gormidxname.Naming) *Options {
	cfg.indexNamingStrategies[pattern] = naming
	return cfg
}
