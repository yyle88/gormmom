package example2enus

import (
	"testing"

	"github.com/yyle88/gormmom"
	"github.com/yyle88/gormmom/gormmomname"
	"github.com/yyle88/runpath/runtestpath"
)

func TestGen(t *testing.T) {
	param := gormmom.NewSchemaCacheV2[Example](runtestpath.SrcPath(t))
	param.Validate()

	cfg := gormmom.NewConfig(param, gormmom.NewOptions().WithDefaultColumnNamePattern(gormmomname.Uppercase63))
	t.Log(cfg)

	cfg.GenReplace()
}
