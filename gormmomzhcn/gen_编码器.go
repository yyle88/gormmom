package gormmomzhcn

import (
	"github.com/yyle88/gormmom"
)

type T批量编码器 struct {
	configs []*T编码器
}

func NewT批量编码器(v表结构s []*T表结构, v配置项 *T配置项) *T批量编码器 {
	var configs = make([]*T编码器, 0, len(v表结构s))
	for _, v表结构 := range v表结构s {
		configs = append(configs, NewT编码器(v表结构, v配置项))
	}
	return &T批量编码器{
		configs: configs,
	}
}

func (c *T批量编码器) Gen替换源码() {
	c.GetConfigs().GenReplaces()
}

func (c *T批量编码器) GetConfigs() gormmom.Configs {
	var configs = make([]*gormmom.Config, 0, len(c.configs))
	for _, v编码器 := range c.configs {
		configs = append(configs, v编码器.GetConfig())
	}
	return configs
}

type T编码器 struct {
	config *gormmom.Config
}

func NewT编码器(v表结构 *T表结构, v配置项 *T配置项) *T编码器 {
	return &T编码器{
		config: gormmom.NewConfig(v表结构.schemaX, v配置项.options),
	}
}

func (cfg *T编码器) Gen替换源码() {
	cfg.config.GenReplace()
}

func (cfg *T编码器) Get新建代码() []byte {
	return cfg.config.GetNewCode()
}

func (cfg *T编码器) GetConfig() *gormmom.Config {
	return cfg.config
}
