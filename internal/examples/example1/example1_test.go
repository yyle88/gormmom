package example1

import (
	"testing"

	"github.com/yyle88/gormmom"
	"github.com/yyle88/runpath/runtestpath"
)

func TestGen(t *testing.T) {
	param := gormmom.NewSchemaCacheV2[Example](runtestpath.SrcPath(t))
	param.Validate()

	cfg := gormmom.NewConfig(param, gormmom.NewOptions())
	t.Log(cfg)

	cfg.GenReplace()
}
