package gormmomzhcn

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

func TestNewT表结构V2(t *testing.T) {
	param := NewT表结构V2[Example1](runpath.CurrentPath())
	param.Validate()
}

func TestNewT表结构V3(t *testing.T) {
	param := NewT表结构V3(runpath.CurrentPath(), &Example1{})
	param.Validate()
}

type Example2 struct {
	V姓名 string `gorm:"column:name"`
	V年龄 int    `gorm:"column:age"`
	V性别 bool   `gorm:"column:sex"`
}

func TestNewT表结构s(t *testing.T) {
	params := NewT表结构s(runpath.PARENT.Path(), []any{&Example1{}, &Example2{}})
	require.Len(t, params, 2)
}
