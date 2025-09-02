package models

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/gormcngen"
	"github.com/yyle88/gormmom"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/runpath"
	"github.com/yyle88/runpath/runtestpath"
)

func TestGen(t *testing.T) {
	objects := []interface{}{
		&Example{},
	}

	require.True(t, t.Run("GenGormMom", func(t *testing.T) {
		params := gormmom.NewGormStructs(runpath.PARENT.Path(), objects)

		cfg := gormmom.NewConfigs(params, gormmom.NewOptions())
		t.Log(cfg)

		result := cfg.GenReplaces()
		require.False(t, result.HasChange()) // 因为已经替换过，而且写到了新代码里，因此这里就只能是没有变化
	}))

	// 使用 require.True(t, t.Run(---)) 限制只有前一步成功才执行后一步的

	require.True(t, t.Run("GenGormCnm", func(t *testing.T) {
		absPath := osmustexist.FILE(runtestpath.SrcPath(t))
		t.Log(absPath)

		options := gormcngen.NewOptions().
			WithColumnClassExportable(true). //中间类型名称的样式为导出的 ExampleColumns
			WithColumnsMethodRecvName("T").
			WithColumnsCheckFieldType(true)

		cfg := gormcngen.NewConfigs(objects, options, absPath)
		cfg.Gen()
	}))
}
