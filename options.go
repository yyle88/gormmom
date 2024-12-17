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

func (opt *Options) WithNamingTagName(namingTagName string) *Options {
	opt.namingTagName = namingTagName
	return opt
}

func (opt *Options) WithColumnNamePatternFieldName(columnNamePatternFieldName string) *Options {
	opt.columnNamePatternFieldName = columnNamePatternFieldName
	return opt
}

func (opt *Options) WithDefaultColumnNamePattern(pattern gormmomname.ColumnNamePattern) *Options {
	opt.defaultColumnNamePattern = pattern
	return opt
}

func (opt *Options) WithColumnNamingStrategies(pattern gormmomname.ColumnNamePattern, naming gormmomname.Naming) *Options {
	opt.columnNamingStrategies[pattern] = naming
	return opt
}

func (opt *Options) WithSkipBasicNaming(skipBasicNaming bool) *Options {
	opt.skipBasicNaming = skipBasicNaming
	return opt
}

func (opt *Options) WithRenewIndexName(renewIndexName bool) *Options {
	opt.renewIndexName = renewIndexName
	return opt
}

func (opt *Options) WithIndexNamingStrategies(pattern gormidxname.IndexNamePattern, naming gormidxname.Naming) *Options {
	opt.indexNamingStrategies[pattern] = naming
	return opt
}
