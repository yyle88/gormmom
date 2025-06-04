package gormmom

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/gormmom/gormmomname"
	"github.com/yyle88/gormmom/internal/utils"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/rese/resb"
	"github.com/yyle88/runpath"
)

func TestMain(m *testing.M) {
	m.Run()
}

type Example struct {
	ID   int32  `gorm:"column:id; primaryKey;" json:"id"`
	V名称  string `gorm:"type:text" mom:"mcp:s63"`
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
	cfg := NewConfig(NewGormStructFromStruct[Example](runpath.CurrentPath()), NewOptions())
	t.Log(cfg)

	newCode := cfg.GetNewCode()
	t.Log(newCode.SrcPath)
	t.Log(newCode.ChangedLineCount)

	require.Equal(t, 7, newCode.ChangedLineCount)

	results := utils.ParseTagsTrimBackticks(newCode.NewCode, &Example{})
	t.Log(neatjsons.S(results))
	require.Equal(t, `gorm:"column:v_0d54_f079;type:text" mom:"mcp:s63;"`, resb.C1(results.Get("V名称")))
	require.Equal(t, `gorm:"column:V_575B_B56B;" mom:"mcp:S63;"`, resb.C1(results.Get("V字段")))
	require.Equal(t, `gorm:"column:v_2760_2b52;" mom:"mcp:s63;"`, resb.C1(results.Get("V性别")))
	require.Equal(t, `gorm:"column:V_7972_8A6B;type:int32" mom:"mcp:S63;"`, resb.C1(results.Get("V特殊")))
}

func TestGetNewCode_S63(t *testing.T) {
	cfg := NewConfig(NewGormStructFromStruct[Example](runpath.CurrentPath()), NewOptions().WithDefaultColumnPattern(gormmomname.NewUppercase63pattern()))
	t.Log(cfg)

	newCode := cfg.GetNewCode()
	t.Log(newCode.SrcPath)
	t.Log(newCode.ChangedLineCount)

	require.Equal(t, 7, newCode.ChangedLineCount)

	results := utils.ParseTagsTrimBackticks(newCode.NewCode, &Example{})
	t.Log(neatjsons.S(results))
	require.Equal(t, `gorm:"column:V_575B_B56B;" mom:"mcp:S63;"`, resb.C1(results.Get("V字段")))
}

type Example5 struct {
	V嘿哈 string `gorm:"column:;type:text"`
}

func TestGetNewCode_Example5(t *testing.T) {
	cfg := NewConfig(NewGormStructFromStruct[Example5](runpath.CurrentPath()), NewOptions().WithDefaultColumnPattern(gormmomname.NewLowercase30pattern()))
	t.Log(cfg)

	newCode := cfg.GetNewCode()
	t.Log(newCode.SrcPath)
	t.Log(newCode.ChangedLineCount)

	require.Equal(t, 1, newCode.ChangedLineCount)

	results := utils.ParseTagsTrimBackticks(newCode.NewCode, &Example5{})
	t.Log(neatjsons.S(results))
	require.Equal(t, `gorm:"column:v_3f56_c854;type:text" mom:"mcp:s30;"`, resb.C1(results.Get("V嘿哈")))
}
