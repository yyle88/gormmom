package example2

import (
	"testing"

	"github.com/yyle88/gormmom"
	"github.com/yyle88/gormmom/gormmomrule"
	"github.com/yyle88/runpath/runtestpath"
)

func TestGen(t *testing.T) {
	cfg := gormmom.NewConfig().SetDefaultRule(gormmomrule.DEFAULT)
	t.Log(cfg)

	srcPath := runtestpath.SrcPath(t)
	param := gormmom.NewParamV2[Example](srcPath)
	param.CheckParam()

	cfg.GenReplace(param)
}
