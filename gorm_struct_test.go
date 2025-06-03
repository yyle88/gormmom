package gormmom

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/runpath"
)

type Example1 struct {
	Name string `gorm:"column:name"`
	Age  int    `gorm:"column:age"`
	Sex  bool   `gorm:"column:sex"`
}

func TestNewGormStructFromStruct(t *testing.T) {
	param := NewGormStructFromStruct[Example1](runpath.CurrentPath())
	t.Log(param.sourcePath)
	t.Log(param.structName)
	require.Equal(t, "Example1", param.structName)
}

func TestNewGormStructFromObject(t *testing.T) {
	param := NewGormStructFromObject(runpath.CurrentPath(), &Example1{})
	t.Log(param.sourcePath)
	t.Log(param.structName)
	require.Equal(t, "Example1", param.structName)
}

type Example2 struct {
	V姓名 string `gorm:"column:name"`
	V年龄 int    `gorm:"column:age"`
	V性别 bool   `gorm:"column:sex"`
}

func TestNewGormStructs(t *testing.T) {
	params := NewGormStructs(runpath.PARENT.Path(), []any{&Example1{}, &Example2{}})
	require.Len(t, params, 2)
	require.Equal(t, "Example1", params[0].structName)
	require.Equal(t, "Example2", params[1].structName)
}
