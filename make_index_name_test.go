package gormmom

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/done"
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
type paramExample3 struct {
	V身份证号 string `gorm:"column:person_num;primaryKey"`
	V学校编号 string `gorm:"column:school_num;uniqueIndex:udx_student_unique"`
	V班级编号 string `gorm:"column:class_num;uniqueIndex:udx_student_unique"`
	V班内排名 string `gorm:"column:student_num;uniqueIndex:udx_student_unique"`
	V姓名   string `gorm:"column:name;index" mom:"rule:S63"`      //普通索引，默认名称不正确（现在的默认名称带中文）
	V年龄   int    `gorm:"column:age;unique" mom:"rule:S63"`      //唯一约束，而非唯一索引，默认名称正确，使用的还是 uni_param_example3_age 这个约束名
	V性别   bool   `gorm:"column:sex;uniqueIndex" mom:"rule:S63"` //唯一索引，带名称，默认名称不正确（现在的默认名称带中文）
}

func TestDryRunMigrate(t *testing.T) {
	db := done.VCE(gorm.Open(sqlite.Open("file::memory:?cache=private"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})).Nice()
	defer func() {
		done.Done(done.VCE(db.DB()).Nice().Close())
	}()

	require.NoError(t, db.Session(&gorm.Session{
		DryRun: true,
	}).AutoMigrate(&paramExample3{}))
}

func TestConfig_GenCode_GenIndexes(t *testing.T) {

	srcPath := runpath.CurrentPath()
	param := NewStructSchemaInfoV2[paramExample3](srcPath)

	cfg := NewConfig(param, NewOptions())
	t.Log(cfg)

	newCode := cfg.GenerateCode()
	t.Log(string(newCode))
}

type paramExample4 struct {
	V证号 string `gorm:"primaryKey"`
	V姓名 string `gorm:"index"`
	V年龄 int    `gorm:"unique"`
	V性别 bool   `gorm:"column:sex;uniqueIndex" mom:"rule:S63"`
}

func (*paramExample4) TableName() string {
	return "tbn88"
}

func TestConfig_GenCode_GenIndexes_2(t *testing.T) {
	srcPath := runpath.CurrentPath()
	param := NewStructSchemaInfoV2[paramExample4](srcPath)

	cfg := NewConfig(param, NewOptions())
	t.Log(cfg)

	newCode := cfg.GenerateCode()
	t.Log(string(newCode))
}
