package models

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/gormcngen"
	"github.com/yyle88/gormmom"
	"github.com/yyle88/runpath"
	"github.com/yyle88/runpath/runtestpath"
)

func TestGenGormMomAndCnm(t *testing.T) {
	objects := []interface{}{
		&T사용자{},
	}

	require.True(t, t.Run("GenGormMom", func(t *testing.T) {
		params := gormmom.ParseObjects(runpath.PARENT.Path(), objects)
		cfg := gormmom.NewConfigs(params, gormmom.NewOptions())
		t.Log("第一步：使用 gormmom 生成韩语字段的 mom 标签")

		result := cfg.Generate()
		require.False(t, result.HasChange()) // 因为已经替换过，而且写到了新代码里，因此这里就只能是没有变化
		require.NoError(t, cfg.ValidateGormTags())
	}))

	require.True(t, t.Run("GenGormCnm", func(t *testing.T) {
		absPath := runtestpath.SrcPath(t)
		t.Log("第二步：使用 gormcngen 生成类型安全的列方法")

		options := gormcngen.NewOptions().
			WithColumnClassExportable(true).
			WithColumnsMethodRecvName("T").
			WithColumnsCheckFieldType(true)

		cfg := gormcngen.NewConfigs(objects, options, absPath)
		cfg.Gen()
		t.Log("gormcngen 生成完成")
	}))
}
