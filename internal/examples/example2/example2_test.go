package example2

import (
	"testing"

	"github.com/yyle88/gormmom"
	"github.com/yyle88/gormmom/gormmomrule"
	"github.com/yyle88/runpath/runtestpath"
)

func TestGen(t *testing.T) {
	cfg := gormmom.NewConfig().SetDftRule(gormmomrule.S63U)
	t.Log(cfg)

	srcPath := runtestpath.SrcPath(t)
	param := gormmom.NewParamV2[Example](srcPath)
	param.Validate()

	cfg.GenReplace(param)
}
