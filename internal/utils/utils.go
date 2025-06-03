package utils

import (
	"go/ast"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/emirpasic/gods/v2/maps/linkedhashmap"
	"github.com/yyle88/done"
	"github.com/yyle88/gormcngen"
	"github.com/yyle88/must"
	"github.com/yyle88/rese"
	"github.com/yyle88/rese/resb"
	"github.com/yyle88/syntaxgo/syntaxgo_ast"
	"github.com/yyle88/syntaxgo/syntaxgo_reflect"
	"github.com/yyle88/syntaxgo/syntaxgo_search"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
	"gorm.io/gorm/schema"
)

// NewSchemaFieldsMap 把字段列表由 slice 转换为 map，以结构体中go的字段名为主键
func NewSchemaFieldsMap(gormSchema *schema.Schema) *linkedhashmap.Map[string, *schema.Field] {
	gormcngen.ShowSchemaEnglish(gormSchema)
	gormcngen.ShowSchemaChinese(gormSchema)

	res := linkedhashmap.New[string, *schema.Field]()
	for _, field := range gormSchema.Fields {
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

func ParseTags[T any](sourceCode []byte, structObject *T) *linkedhashmap.Map[string, string] {
	astBundle := rese.P1(syntaxgo_ast.NewAstBundleV1(sourceCode))
	astFile, fileSet := astBundle.GetBundle()
	zaplog.LOG.Debug("ast-get-package-name", zap.String("package_name", astFile.Name.Name))

	structName := syntaxgo_reflect.GetTypeNameV4(structObject)
	zaplog.LOG.Debug("reflect-get-struct-name", zap.String("name", structName))

	structType := resb.P1(syntaxgo_search.FindStructTypeByName(astFile, structName))
	must.Done(ast.Print(fileSet, structType))

	var results = linkedhashmap.New[string, string]()
	if structType.Fields != nil {
		for _, field := range structType.Fields.List {
			for _, name := range field.Names {
				if field.Tag != nil {
					results.Put(name.Name, field.Tag.Value)
				} else {
					results.Put(name.Name, "")
				}
			}
		}
	}
	return results
}

func TrimQuotes(s string) string {
	return strings.Trim(s, `"`)
}

func TrimBackticks(s string) string {
	return strings.Trim(s, "`")
}

func ParseTagsTrimBackticks[T any](sourceCode []byte, structObject *T) *linkedhashmap.Map[string, string] {
	results := linkedhashmap.New[string, string]()
	ParseTags(sourceCode, structObject).Each(func(key string, value string) {
		results.Put(key, TrimBackticks(value))
	})
	return results
}
