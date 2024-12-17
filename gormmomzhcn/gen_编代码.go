package gormmomzhcn

import (
	"github.com/yyle88/gormmom"
)

type T批量编码器 struct {
	configBatch *gormmom.ConfigBatch
}

func NewT批量编码器(v表结构s []*T表结构, v配置项 *T配置项) *T批量编码器 {
	var schemaCaches = make([]*gormmom.SchemaCache, 0, len(v表结构s))
	for _, v表结构 := range v表结构s {
		schemaCaches = append(schemaCaches, v表结构.schemaCache)
	}
	return &T批量编码器{
		configBatch: gormmom.NewConfigBatch(schemaCaches, v配置项.options),
	}
}

func (c *T批量编码器) Gen替换源码() {
	c.configBatch.GenReplaces()
}

func (c *T批量编码器) ConfigBatch() *gormmom.ConfigBatch {
	return c.configBatch
}

type T编码器 struct {
	config *gormmom.Config
}

func NewT编码器(v表结构 *T表结构, v配置项 *T配置项) *T编码器 {
	return &T编码器{
		config: gormmom.NewConfig(v表结构.schemaCache, v配置项.options),
	}
}

func (cfg *T编码器) Gen替换源码() {
	cfg.config.GenReplace()
}

func (cfg *T编码器) Get新建代码() []byte {
	return cfg.config.GetNewCode()
}

func (cfg *T编码器) Config() *gormmom.Config {
	return cfg.config
}
