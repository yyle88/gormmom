package gormmomzhcn

import (
	"github.com/yyle88/gormmom"
	"gorm.io/gorm/schema"
)

type T表结构 struct {
	schemaCache *gormmom.SchemaCache
}

// NewT表结构 创建参数信息
func NewT表结构(sourcePath string, structName string, sch *schema.Schema) *T表结构 {
	return &T表结构{
		schemaCache: gormmom.NewSchemaCache(sourcePath, structName, sch),
	}
}

// NewT表结构V2 使用泛型创建参数信息。T 只能传类型名称而非带指针的类型名
func NewT表结构V2[T any](sourcePath string) *T表结构 {
	return &T表结构{
		schemaCache: gormmom.NewSchemaCacheV2[T](sourcePath),
	}
}

// NewT表结构V3 使用对象创建参数信息 object 传对象或者对象指针都行
func NewT表结构V3(sourcePath string, object interface{}) *T表结构 {
	return &T表结构{
		schemaCache: gormmom.NewSchemaCacheV3(sourcePath, object),
	}
}

func (a *T表结构) Validate() {
	a.schemaCache.Validate()
}

func NewT表结构s(root string, objects []interface{}) []*T表结构 {
	structSchemas := gormmom.NewSchemaCaches(root, objects)

	results := make([]*T表结构, 0, len(structSchemas))
	for _, x := range structSchemas {
		results = append(results, &T表结构{
			schemaCache: x,
		})
	}
	return results
}

func (a *T表结构) SchemaCache() *gormmom.SchemaCache {
	return a.schemaCache
}
