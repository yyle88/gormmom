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

// GormStruct represents a GORM struct with its location and field information
// Contains the source file path, struct name, and comprehensive field mappings
// Provides structured access to GORM schema and field definitions for processing
// Maintains ordered field mapping using linked hash map for deterministic generation
//
// GormStruct 代表一个 GORM 结构体及其位置和字段信息
// 包含源文件路径、结构体名称和全面的字段映射
// 为处理提供对 GORM 模式和字段定义的结构化访问
// 使用链式哈希映射维护有序字段映射，确保确定性生成
type GormStruct struct {
	sourcePath string                                    // Source file path where struct is defined // 定义结构体的源文件路径
	structName string                                    // Name of the target struct // 目标结构体的名称
	gormSchema *schema.Schema                            // GORM schema information // GORM 模式信息
	gormFields *linkedhashmap.Map[string, *schema.Field] // Ordered field mapping for deterministic processing // 确定性处理的有序字段映射
}

// NewGormStruct creates a new GormStruct instance with field information
// Reads and processes struct field information from the source file and GORM schema
// Builds ordered field mapping for deterministic tag generation and processing
// Returns configured GormStruct prepared for native language tag processing
//
// NewGormStruct 创建新的 GormStruct 实例并读取字段信息
// 从源文件和 GORM 模式中读取和处理结构体字段信息
// 构建有序字段映射以进行确定性标签生成和处理
// 返回配置好的 GormStruct，准备进行原生语言标签处理
func NewGormStruct(sourcePath string, structName string, gormSchema *schema.Schema) *GormStruct {
	zaplog.LOG.Debug("new-struct-schema-info", zap.String("struct_name", structName), zap.String("source_path", sourcePath))

	// Validate table name for ASCII compatibility before processing
	// Ensures database compatibility and prevents downstream issues
	//
	// 在处理前验证表名的 ASCII 兼容性
	// 确保数据库兼容性并避免下游问题
	utils.ValidateTableName(gormSchema.Table, structName)

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
