package gormidxname

import (
	"github.com/yyle88/gormmom/internal/utils"
	"gorm.io/gorm/schema"
)

// nolint:no-doc
// 自定义枚举类型
type IndexNamePattern string

const (
	WithTableCnm IndexNamePattern = "cnm" //表示使用 column name 作为拼接索引名的规则
	UppercaseCnm IndexNamePattern = "CNM" //表示使用 column name 作为拼接索引名的规则，但后缀是大写字母的

	DefaultPattern IndexNamePattern = WithTableCnm
)

type Naming interface {
	IsValidIndexName(indexName string) bool
	GenerateIndexName(schemaIndex *schema.Index, tableName string, fieldName string, columnName string) *GenerateIndexNameResult
}

type GenerateIndexNameResult struct {
	IndexTagFieldName string
	NewIndexName      string
	IdxUdxPrefix      string
}

var namingStrategies = map[IndexNamePattern]Naming{
	WithTableCnm: &withTableCnmPattern{},
	UppercaseCnm: &uppercaseCnmPattern{},
}

func GetNamingStrategies() map[IndexNamePattern]Naming {
	return utils.CloneMap(namingStrategies)
}
