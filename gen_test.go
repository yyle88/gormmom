package gormmom

import (
	"testing"
	"time"

	"github.com/yyle88/gormmom/gormmomname"
	"github.com/yyle88/runpath"
)

func TestMain(m *testing.M) {
	m.Run()
}

type Example struct {
	ID   int32  `gorm:"column:id; primaryKey;" json:"id"`
	V名称  string `gorm:"type:text" mom:"mcp:S63"`
	V字段  string `gorm:"column: some_field" mom:"mcp:S63;"`
	V性别  string
	V特殊  string `gorm:"column:特殊;type:int32" mom:"mcp:S63;"`
	V年龄  int    `json:"age"` //理论上不要直接给model添加json标签，因为那是view层的逻辑，但实际上假如非这样做也能处理
	Rank int32  ``           //看看这种情况是啥效果
	V身高  int32  ``           //看看这种情况是啥效果
	V体重  int32  ``           //看看这种情况是啥效果
	// v啥呀       string    //这个是小写字母开头的所以不是导出字段
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func TestGetNewCode(t *testing.T) {
	cfg := NewConfig(NewSchemaX2[Example](runpath.CurrentPath()), NewOptions())
	t.Log(cfg)

	newCode := cfg.GetNewCode()
	t.Log(string(newCode))
}

func TestGetNewCode_S63(t *testing.T) {
	cfg := NewConfig(NewSchemaX2[Example](runpath.CurrentPath()), NewOptions().WithDefaultColumnPattern(gormmomname.NewUppercase63pattern()))
	t.Log(cfg)

	newCode := cfg.GetNewCode()
	t.Log(string(newCode))
}
