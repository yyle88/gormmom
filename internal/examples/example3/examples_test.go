package example3

import (
	"testing"

	"github.com/yyle88/gormmom"
	"github.com/yyle88/runpath"
)

func TestGen(t *testing.T) {
	cfg := gormmom.NewConfig()
	t.Log(cfg)

	params := gormmom.CreateParams(runpath.PARENT.Path(), []interface{}{
		&Example1{},
		&Example2{},
	})

	cfg.GenReplaces(params)
}
