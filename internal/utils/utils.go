package utils

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/yyle88/done"
	"github.com/yyle88/gormcngen"
	"github.com/yyle88/zaplog"
	"gorm.io/gorm/schema"
)

// NewSchemaFieldsMap 把字段列表由 slice 转换为 map，以结构体中go的字段名为主键
func NewSchemaFieldsMap(sch *schema.Schema) map[string]*schema.Field {
	gormcngen.ShowSchemaMessage(sch)

	var mp = make(map[string]*schema.Field, len(sch.Fields))
	for _, field := range sch.Fields {
		mp[field.Name] = field ////go结构体成员名称
	}
	return mp
}

func WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

// LsGoFiles 获取指定目录下的所有 .go 文件路径（不递归子目录）
func LsGoFiles(root string) []string {
	var paths []string
	var entries = done.VAE(os.ReadDir(root)).Nice()
	for _, one := range entries {
		// 检查是否是文件和扩展名为 .go
		if !one.IsDir() && filepath.Ext(one.Name()) == ".go" {
			paths = append(paths, filepath.Join(root, one.Name()))
		}
	}
	return paths
}

func AssertOK[T comparable](v T) {
	var zero T
	if v == zero {
		zaplog.LOG.Panic("assert not ok arg is zero then panic")
	}
}

func AssertEquals[T comparable](a, b T) {
	if a != b {
		panic(errors.New("not equals"))
	}
}

func FALSE() bool {
	return false
}
