package gormmom

import (
	"github.com/yyle88/gormmom/gormidxname"
	"github.com/yyle88/gormmom/gormmomname"
)

type Options struct {
	tagName                string
	columnNamingSubTagName string
	columnNamingStrategies *gormmomname.Strategies
	skipBasicColumnName    bool //是否跳过简单字段，有的字段虽然没有配置名称或者规则，但是它满足简单字段，就也不做任何处理
	renewIndexName         bool
	indexNamingStrategies  *gormidxname.Strategies
}

func NewOptions() *Options {
	return &Options{
		tagName:                "mom", // mother_tongue native_language. // 因为我也不太能熟练拼写更高级的单词，还不如返璞归真直接用口语表示
		columnNamingSubTagName: "mcp", // m(mother_tongue) c(column_name) p(pattern) // 表示列名的样式标签
		columnNamingStrategies: gormmomname.NewStrategies(),
		skipBasicColumnName:    true,
		renewIndexName:         true,
		indexNamingStrategies:  gormidxname.NewStrategies(),
	}
}

func (opt *Options) WithTagName(tagName string) *Options {
	opt.tagName = tagName
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
