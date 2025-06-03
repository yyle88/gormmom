package example1

import (
	"testing"

	"github.com/yyle88/gormmom"
	"github.com/yyle88/runpath/runtestpath"
)

func TestGen(t *testing.T) {
	param := gormmom.NewGormStructFromStruct[Example](runtestpath.SrcPath(t))

	cfg := gormmom.NewConfig(param, gormmom.NewOptions())
	t.Log(cfg)

	cfg.GenReplace()
}
