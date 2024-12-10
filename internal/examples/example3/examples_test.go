package example3

import (
	"testing"

	"github.com/yyle88/gormmom"
	"github.com/yyle88/runpath"
)

func TestGen(t *testing.T) {
	params := gormmom.NewStructSchemaInfos(runpath.PARENT.Path(), []interface{}{
		&Example1{},
		&Example2{},
	})

	cfg := gormmom.NewCodeGenerationConfig(params, gormmom.NewOptions())
	t.Log(cfg)

	cfg.GenReplaces()
}
