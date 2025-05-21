package example1_en_us

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
		WithColumnClassExportable(true). //中间类型名称的样式为导出的 ExampleColumns
		WithColumnsMethodRecvName("T").
		WithColumnsCheckFieldType(true)

	cfg := gormcngen.NewConfigs([]interface{}{&Example{}}, options, absPath)
	cfg.Gen()
}
