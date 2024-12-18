package example3enus

import (
	"testing"

	"github.com/yyle88/gormcngen"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/runpath/runtestpath"
)

func TestGenerate(t *testing.T) {
	absPath := osmustexist.FILE(runtestpath.SrcPath(t))
	t.Log(absPath)

	options := gormcngen.NewOptions().
		WithColumnClassExportable(true). //中间类型名称的样式为非导出的 exampleColumns
		WithColumnsMethodRecvName("T").
		WithColumnsCheckFieldType(true)

	cfg := gormcngen.NewConfigs([]interface{}{&Example1{}, &Example2{}}, options, absPath)
	cfg.Gen()
}
