package example3

import (
	"testing"

	"github.com/yyle88/gormmom"
	"github.com/yyle88/runpath"
)

func TestGen(t *testing.T) {
	params := gormmom.NewGormStructs(runpath.PARENT.Path(), []interface{}{
		&Example{},
		&Example2{},
	})

	cfg := gormmom.NewConfigs(params, gormmom.NewOptions())
	t.Log(cfg)

	cfg.GenReplaces()
}
