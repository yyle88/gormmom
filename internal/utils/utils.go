package utils

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/emirpasic/gods/v2/maps/linkedhashmap"
	"github.com/yyle88/done"
	"github.com/yyle88/gormcngen"
	"github.com/yyle88/must"
	"gorm.io/gorm/schema"
)

// NewSchemaFieldsMap 把字段列表由 slice 转换为 map，以结构体中go的字段名为主键
func NewSchemaFieldsMap(sch *schema.Schema) *linkedhashmap.Map[string, *schema.Field] {
	gormcngen.ShowSchemaEnglish(sch)
	gormcngen.ShowSchemaChinese(sch)

	res := linkedhashmap.New[string, *schema.Field]()
	for _, field := range sch.Fields {
		res.Put(field.Name, field) //键是Go结构体成员名称
	}
	return res
}

func WriteFile(path string, data []byte) {
	must.Done(os.WriteFile(path, data, 0644))
}

// ListGoFiles 获取指定目录下的所有 .go 文件路径（不递归子目录）
func ListGoFiles(root string) []string {
	var paths []string
	for _, one := range done.VAE(os.ReadDir(root)).Nice() {
		// 检查是否是文件和扩展名为 .go
		if !one.IsDir() && filepath.Ext(one.Name()) == ".go" {
			paths = append(paths, filepath.Join(root, one.Name()))
		}
	}
	return paths
}

func ParseSchema(object interface{}) *schema.Schema {
	return done.VCE(schema.Parse(object, &sync.Map{}, &schema.NamingStrategy{
		SingularTable: false, //和默认值相同
		NoLowerCase:   false, //和默认值相同
	})).Nice()
}
