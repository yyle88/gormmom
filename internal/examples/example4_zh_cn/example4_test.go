package example4_zh_cn

import (
	"testing"

	"github.com/yyle88/gormmom/gormmomname"
	"github.com/yyle88/gormmom/gormmomzhcn"
	"github.com/yyle88/runpath/runtestpath"
)

func TestGen(t *testing.T) {
	param := gormmomzhcn.NewT表结构V2[Example](runtestpath.SrcPath(t))
	param.Validate()

	cfg := gormmomzhcn.NewT编码器(param, gormmomzhcn.NewT配置项().With默认列名样式(gormmomname.NewUppercase63pattern()))
	t.Log(cfg)

	cfg.Gen替换源码()
}
