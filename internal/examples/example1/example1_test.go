package example1

import (
	"testing"

	"github.com/yyle88/gormmom"
	"github.com/yyle88/runpath/runtestpath"
)

func TestGen(t *testing.T) {

	srcPath := runtestpath.SrcPath(t)
	param := gormmom.NewStructSchemaInfoV2[Example](srcPath)
	param.Validate()

	cfg := gormmom.NewConfig(param, gormmom.NewOptions())
	t.Log(cfg)

	cfg.GenReplace()
}
