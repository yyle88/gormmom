package gormmom

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/runpath"
)

type paramExample struct {
	Name string `gorm:"name"`
	Age  int    `gorm:"age"`
	Sex  bool   `gorm:"sex"`
}

func TestNewParamV2(t *testing.T) {
	param := NewParamV2[paramExample](runpath.Current())
	param.Validate()
}

func TestNewParamV3(t *testing.T) {
	param := NewParamV3(runpath.Current(), &paramExample{})
	param.Validate()
}

type paramExample2 struct {
	V姓名 string `gorm:"column:name"`
	V年龄 int    `gorm:"column:age"`
	V性别 bool   `gorm:"column:sex"`
}

func TestNewParams(t *testing.T) {
	params := NewParams(runpath.PARENT.Path(), []any{&paramExample{}, &paramExample2{}})
	require.Len(t, params, 2)
}
