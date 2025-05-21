package gormmom

import (
	"github.com/emirpasic/gods/v2/maps/linkedhashmap"
	"github.com/yyle88/gormmom/internal/utils"
	"github.com/yyle88/must"
	"github.com/yyle88/rese"
	"github.com/yyle88/syntaxgo/syntaxgo_ast"
	"github.com/yyle88/syntaxgo/syntaxgo_reflect"
	"github.com/yyle88/syntaxgo/syntaxgo_search"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
	"gorm.io/gorm/schema"
)

// SchemaX 结构体的位置和结构体里的字段信息
type SchemaX struct {
	sourcePath string
	structName string
	sch        *schema.Schema
	schColumns *linkedhashmap.Map[string, *schema.Field]
}

// NewSchemaX 读取结构体字段信息
func NewSchemaX(sourcePath string, structName string, sch *schema.Schema) *SchemaX {
	zaplog.LOG.Debug("new-struct-schema-info", zap.String("struct_name", structName), zap.String("source_path", sourcePath))

	return &SchemaX{
		sourcePath: sourcePath,
		structName: structName,
		sch:        sch,
		schColumns: utils.NewSchemaFieldsMap(sch), //这里提前把列做成map方便使用
	}
}

// NewSchemaX2 使用泛型创建参数信息。T 只能传类型名称而非带指针的类型名
func NewSchemaX2[T any](sourcePath string) *SchemaX {
	return NewSchemaX3(sourcePath, new(T))
}

// NewSchemaX3 使用对象创建参数信息 object 传对象或者对象指针都行
func NewSchemaX3(sourcePath string, object interface{}) *SchemaX {
	return NewSchemaX(sourcePath, syntaxgo_reflect.GetTypeNameV3(object), utils.ParseSchema(object))
}

func (a *SchemaX) Validate() {
	must.Nice(a.sourcePath)
	must.Nice(a.structName)
	must.Full(a.sch)
	must.Full(a.schColumns)
}

func NewSchemaXs(root string, objects []interface{}) []*SchemaX {
	var schemaXs = make([]*SchemaX, 0, len(objects))
	var paths = utils.ListGoFiles(root)
	var exists = make(map[string]bool, len(objects)) //记住已经处理的数据
	for _, sourcePath := range paths {
		astBundle := rese.P1(syntaxgo_ast.NewAstBundleV4(sourcePath))

		astFile, _ := astBundle.GetBundle()

		for idx, object := range objects {
			structName := syntaxgo_reflect.GetTypeNameV3(object)
			if structName == "" {
				zaplog.LOG.Warn("object doesn't have struct name", zap.Int("idx", idx))
				continue
			}

			if exists[structName] {
				continue //跳过已经处理的，按结构体名称，确保只处理一次
			}
			if _, ok := syntaxgo_search.FindStructDeclarationByName(astFile, structName); !ok {
				continue //说明这个结构体的定义不在这个文件里
			}
			exists[structName] = true

			schemaXs = append(schemaXs, NewSchemaX3(sourcePath, object))
		}
	}
	return schemaXs
}
