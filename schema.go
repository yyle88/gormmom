package gormmom

import (
	"github.com/yyle88/erero"
	"github.com/yyle88/gormmom/internal/utils"
	"github.com/yyle88/rese"
	"github.com/yyle88/syntaxgo/syntaxgo_ast"
	"github.com/yyle88/syntaxgo/syntaxgo_reflect"
	"github.com/yyle88/syntaxgo/syntaxgo_search"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
	"gorm.io/gorm/schema"
)

type SchemaCache struct {
	sourcePath string
	structName string
	sch        *schema.Schema
	schColumns map[string]*schema.Field
}

// NewSchemaCache 创建参数信息
func NewSchemaCache(sourcePath string, structName string, sch *schema.Schema) *SchemaCache {
	zaplog.LOG.Debug("new-struct-schema-info", zap.String("struct_name", structName), zap.String("source_path", sourcePath))

	return &SchemaCache{
		sourcePath: sourcePath,
		structName: structName,
		sch:        sch,
		schColumns: utils.NewSchemaFieldsMap(sch), //这里提前把列做成map方便使用
	}
}

// NewSchemaCacheV2 使用泛型创建参数信息。T 只能传类型名称而非带指针的类型名
func NewSchemaCacheV2[T any](sourcePath string) *SchemaCache {
	return NewSchemaCacheV3(sourcePath, utils.Newp[T]())
}

// NewSchemaCacheV3 使用对象创建参数信息 object 传对象或者对象指针都行
func NewSchemaCacheV3(sourcePath string, object interface{}) *SchemaCache {
	return NewSchemaCache(sourcePath, syntaxgo_reflect.GetTypeNameV3(object), utils.ParseSchema(object))
}

func (a *SchemaCache) Validate() {
	if a.sourcePath == "" {
		panic(erero.New("a.source_path is none"))
	}
	if a.structName == "" {
		panic(erero.New("a.struct_name is none"))
	}
	if a.sch == nil {
		panic(erero.New("a.sch is none"))
	}
	if a.schColumns == nil {
		panic(erero.New("a.sch_columns is none"))
	}
}

func NewSchemaCaches(root string, objects []interface{}) []*SchemaCache {
	var schemaCaches = make([]*SchemaCache, 0, len(objects))
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

			schemaCaches = append(schemaCaches, NewSchemaCacheV3(sourcePath, object))
		}
	}
	return schemaCaches
}
