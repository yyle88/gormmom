package gormmom

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/done"
	"github.com/yyle88/gormmom/internal/utils"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/rese"
	"github.com/yyle88/rese/resb"
	"github.com/yyle88/runpath"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

/*
在 Gorm 里 unique 和 uniqueIndex 在 GORM 中确实有区别：

unique:
这个标签表示为字段添加一个唯一性约束，但它不一定会为这个约束生成一个命名的索引。数据库会确保字段值唯一，但索引的命名和管理通常由数据库引擎自动处理，因此不一定会出现一个独立命名的唯一索引。

uniqueIndex:
这个标签不仅表示字段的唯一性，还明确要求 GORM 生成一个带有名称的唯一索引。通常会按照 idx_<table_name>_<column_name> 的格式生成索引名称（具体格式可能因数据库而异）。
*/
type Example3 struct {
	V身份证号 string `gorm:"column:person_num;primaryKey"`
	V学校编号 string `gorm:"column:school_num;uniqueIndex:udx_student_unique"`
	V班级编号 string `gorm:"column:class_num;uniqueIndex:udx_student_unique"`
	V班内排名 string `gorm:"column:student_num;uniqueIndex:udx_student_unique"`
	V姓名   string `gorm:"column:name;index" mom:"mcp:S63"`      //普通索引，默认名称不正确（现在的默认名称带中文）
	V年龄   int    `gorm:"column:age;unique" mom:"mcp:S63"`      //唯一约束，而非唯一索引，默认名称正确，使用的还是 uni_param_example3_age 这个约束名
	V性别   bool   `gorm:"column:sex;uniqueIndex" mom:"mcp:S63"` //唯一索引，带名称，默认名称不正确（现在的默认名称带中文）
}

func TestDryRunMigrate(t *testing.T) {
	db := done.VPE(gorm.Open(sqlite.Open("file::memory:?cache=private"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})).Full()
	defer rese.F0(rese.P1(db.DB()).Close)

	require.NoError(t, db.Session(&gorm.Session{
		DryRun: true,
	}).AutoMigrate(&Example3{}))
}

func TestConfig_GenCode_GenIndexes(t *testing.T) {
	cfg := NewConfig(NewGormStructFromStruct[Example3](runpath.CurrentPath()), NewOptions())
	t.Log(cfg)

	newCode := cfg.GetNewCode()
	results := utils.ParseTagsTrimBackticks(newCode, &Example3{})
	t.Log(neatjsons.S(results))
	require.Equal(t, `gorm:"column:V_D359_0D54;index:idx_example3_v_d359_0d54" mom:"mcp:S63;idx:cnm;"`, resb.C1(results.Get("V姓名")))
	require.Equal(t, `gorm:"column:V_745E_849F;unique" mom:"mcp:S63;"`, resb.C1(results.Get("V年龄")))
	require.Equal(t, `gorm:"column:V_2760_2B52;uniqueIndex:udx_example3_v_2760_2b52" mom:"mcp:S63;udx:cnm;"`, resb.C1(results.Get("V性别")))
}

type Example4 struct {
	V证号 string `gorm:"primaryKey"`
	V姓名 string `gorm:"index"`
	V年龄 int    `gorm:"unique"`
	V性别 bool   `gorm:"column:sex;uniqueIndex" mom:"mcp:S63"`
}

func (*Example4) TableName() string {
	return "example4"
}

func TestConfig_GenCode_GenIndexes_Example4(t *testing.T) {
	cfg := NewConfig(NewGormStructFromStruct[Example4](runpath.CurrentPath()), NewOptions())
	t.Log(cfg)

	newCode := cfg.GetNewCode()
	results := utils.ParseTagsTrimBackticks(newCode, &Example4{})
	t.Log(neatjsons.S(results))
	require.Equal(t, `gorm:"column:v_c18b_f753;primaryKey" mom:"mcp:s63;"`, resb.C1(results.Get("V证号")))
	require.Equal(t, `gorm:"column:v_d359_0d54;index:idx_example4_v_d359_0d54" mom:"mcp:s63;idx:cnm;"`, resb.C1(results.Get("V姓名")))
	require.Equal(t, `gorm:"column:v_745e_849f;unique" mom:"mcp:s63;"`, resb.C1(results.Get("V年龄")))
	require.Equal(t, `gorm:"column:V_2760_2B52;uniqueIndex:udx_example4_v_2760_2b52" mom:"mcp:S63;udx:cnm;"`, resb.C1(results.Get("V性别")))
}

type Example6 struct {
	V账号   string `gorm:"primaryKey"`
	V身份证号 string `gorm:"column:person_num;uniqueIndex" mom:"udx:cnm;"`
	V学校编号 string `gorm:"column:school_num;index" mom:"idx:cnm;"`
}

func TestConfig_GenCode_GenIndexes_Example6(t *testing.T) {
	cfg := NewConfig(NewGormStructFromStruct[Example6](runpath.CurrentPath()), NewOptions())
	t.Log(cfg)

	newCode := cfg.GetNewCode()
	results := utils.ParseTagsTrimBackticks(newCode, &Example6{})
	t.Log(neatjsons.S(results))
	require.Equal(t, `gorm:"column:v_268d_f753;primaryKey" mom:"mcp:s63;"`, resb.C1(results.Get("V账号")))
	require.Equal(t, `gorm:"column:person_num;uniqueIndex:udx_example6_person_num" mom:"udx:cnm;"`, resb.C1(results.Get("V身份证号")))
	require.Equal(t, `gorm:"column:school_num;index:idx_example6_school_num" mom:"idx:cnm;"`, resb.C1(results.Get("V学校编号")))
}
