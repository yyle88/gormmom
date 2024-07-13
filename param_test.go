package gormmom

import (
	"testing"

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
	param := NewParamV3(runpath.Current(), paramExample{})
	param.Validate()
}
