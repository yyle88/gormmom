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

func (opt *T配置项) With命名总标签名(tagName string) *T配置项 {
	opt.options.WithTagName(tagName)
	return opt
}

func (opt *T配置项) With列名模式标签(columnNamingSubTagName string) *T配置项 {
	opt.options.WithColumnNamingSubTagName(columnNamingSubTagName)
	return opt
}

func (opt *T配置项) With列名命名样式(pattern gormmomname.Pattern) *T配置项 {
	opt.options.WithColumnPattern(pattern)
	return opt
}

func (opt *T配置项) With默认列名样式(pattern gormmomname.Pattern) *T配置项 {
	opt.options.WithDefaultColumnPattern(pattern)
	return opt
}

func (opt *T配置项) With跳过基本列名(skipBasicColumnName bool) *T配置项 {
	opt.options.WithSkipBasicColumnName(skipBasicColumnName)
	return opt
}

func (opt *T配置项) With重新命名索引(renewIndexName bool) *T配置项 {
	opt.options.WithRenewIndexName(renewIndexName)
	return opt
}

func (opt *T配置项) With索引命名样式(pattern gormidxname.Pattern) *T配置项 {
	opt.options.WithIndexPattern(pattern)
	return opt
}

func (opt *T配置项) With默认索引样式(pattern gormidxname.Pattern) *T配置项 {
	opt.options.WithDefaultIndexPattern(pattern)
	return opt
}

func (opt *T配置项) GetOptions() *gormmom.Options {
	return opt.options
}
