package example2

import (
	"testing"

	"github.com/yyle88/gormmom"
	"github.com/yyle88/gormmom/gormmomname"
	"github.com/yyle88/runpath/runtestpath"
)

func TestGen(t *testing.T) {
	param := gormmom.NewStructSchemaInfoV2[Example](runtestpath.SrcPath(t))
	param.Validate()

	cfg := gormmom.NewConfig(param, gormmom.NewOptions().SetDefaultColumnNamePattern(gormmomname.DefaultPattern))
	t.Log(cfg)

	cfg.GenReplace()
}
