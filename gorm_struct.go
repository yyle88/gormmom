package gormmom

import (
	"github.com/emirpasic/gods/v2/maps/linkedhashmap"
	"github.com/yyle88/gormmom/internal/utils"
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

// GormStruct 结构体的位置和结构体里的字段信息
type GormStruct struct {
	sourcePath string
	structName string
	gormSchema *schema.Schema
	gormFields *linkedhashmap.Map[string, *schema.Field]
}

// NewGormStruct 读取结构体字段信息
func NewGormStruct(sourcePath string, structName string, gormSchema *schema.Schema) *GormStruct {
	zaplog.LOG.Debug("new-struct-schema-info", zap.String("struct_name", structName), zap.String("source_path", sourcePath))

	return &GormStruct{
		sourcePath: must.Nice(sourcePath),
		structName: must.Nice(structName),
		gormSchema: must.Full(gormSchema),
		gormFields: must.Full(utils.NewSchemaFieldsMap(gormSchema)), //这里提前把列做成map方便使用
	}
}

// NewGormStructFromStruct 使用泛型创建参数信息。T 只能传类型名称而非带指针的类型名
func NewGormStructFromStruct[StructType any](sourcePath string) *GormStruct {
	return NewGormStructFromObject(sourcePath, new(StructType))
}

// NewGormStructFromObject 使用对象创建参数信息 object 传对象或者对象指针都行
func NewGormStructFromObject(sourcePath string, object interface{}) *GormStruct {
	return NewGormStruct(sourcePath, syntaxgo_reflect.GetTypeNameV3(object), utils.ParseSchema(object))
}

func NewGormStructs(root string, objects []interface{}) []*GormStruct {
	var objectMap = linkedhashmap.New[string, any]() // 使用有序map来存储对象，避免乱序执行导致每次执行结果不同
	for idx, object := range objects {
		structName := syntaxgo_reflect.GetTypeNameV3(object) // 获取结构体名称
		if structName == "" {                                // 这里不允许获取不到名称的
			zaplog.LOG.Panic("object doesn't have struct name", zap.Int("idx", idx))
		}
		objectMap.Put(structName, object)
	}

	var results = make([]*GormStruct, 0, len(objects))
	for _, sourcePath := range utils.ListGoFiles(root) {
		astBundle := rese.P1(syntaxgo_ast.NewAstBundleV4(sourcePath))
		astFile, _ := astBundle.GetBundle()
		for _, structName := range objectMap.Keys() {
			if _, ok := syntaxgo_search.FindStructDeclarationByName(astFile, structName); !ok {
				continue //说明这个结构体的定义不在这个文件里
			}
			// 得到相应的结构体对象
			oneObject := resb.V1(objectMap.Get(structName))
			// 得到结构体类型的定义，代码文件路径，结构体名称，内部字段列表
			results = append(results, NewGormStructFromObject(sourcePath, oneObject))
			// 移除，以确保只处理一次，这样也能避免重复搜索代码块
			objectMap.Remove(structName)
		}
	}
	return results
}
