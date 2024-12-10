package gormmom

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/runpath"
)

type Example1 struct {
	Name string `gorm:"name"`
	Age  int    `gorm:"age"`
	Sex  bool   `gorm:"sex"`
}

func TestNewStructSchemaInfoV2(t *testing.T) {
	param := NewStructSchemaInfoV2[Example1](runpath.CurrentPath())
	param.Validate()
}

func TestNewStructSchemaInfoV3(t *testing.T) {
	param := NewStructSchemaInfoV3(runpath.CurrentPath(), &Example1{})
	param.Validate()
}

type Example2 struct {
	V姓名 string `gorm:"column:name"`
	V年龄 int    `gorm:"column:age"`
	V性别 bool   `gorm:"column:sex"`
}

func TestNewStructSchemaInfos(t *testing.T) {
	params := NewStructSchemaInfos(runpath.PARENT.Path(), []any{&Example1{}, &Example2{}})
	require.Len(t, params, 2)
}
