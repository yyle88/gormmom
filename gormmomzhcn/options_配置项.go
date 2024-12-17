package gormmomzhcn

import (
	"github.com/yyle88/gormmom"
	"github.com/yyle88/gormmom/gormidxname"
	"github.com/yyle88/gormmom/gormmomname"
)

type T配置项 struct {
	options *gormmom.Options
}

func NewT配置项() *T配置项 {
	return &T配置项{
		options: gormmom.NewOptions(),
	}
}

func (opt *T配置项) With命名总标签名(namingTagName string) *T配置项 {
	opt.options.WithNamingTagName(namingTagName)
	return opt
}

func (opt *T配置项) With列名模式字段(columnNamePatternFieldName string) *T配置项 {
	opt.options.WithColumnNamePatternFieldName(columnNamePatternFieldName)
	return opt
}

func (opt *T配置项) With默认列名模式(pattern gormmomname.ColumnNamePattern) *T配置项 {
	opt.options.WithDefaultColumnNamePattern(pattern)
	return opt
}

func (opt *T配置项) With列名命名策略(pattern gormmomname.ColumnNamePattern, naming gormmomname.Naming) *T配置项 {
	opt.options.WithColumnNamingStrategies(pattern, naming)
	return opt
}

func (opt *T配置项) With跳过基本命名(skipBasicNaming bool) *T配置项 {
	opt.options.WithSkipBasicNaming(skipBasicNaming)
	return opt
}

func (opt *T配置项) With重命名索引名(renewIndexName bool) *T配置项 {
	opt.options.WithRenewIndexName(renewIndexName)
	return opt
}

func (opt *T配置项) With索引命名策略(pattern gormidxname.IndexNamePattern, naming gormidxname.Naming) *T配置项 {
	opt.options.WithIndexNamingStrategies(pattern, naming)
	return opt
}

func (opt *T配置项) Options() *gormmom.Options {
	return opt.options
}
