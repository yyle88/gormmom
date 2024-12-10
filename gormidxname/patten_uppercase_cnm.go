package gormidxname

import (
	"strings"

	"gorm.io/gorm/schema"
)

type uppercaseCnmPattern struct{}

func (G *uppercaseCnmPattern) IsValidIndexName(indexName string) bool {
	return new(withTableCnmPattern).IsValidIndexName(indexName)
}

func (G *uppercaseCnmPattern) GenerateIndexName(schemaIndex schema.Index, tableName string, fieldName string, columnName string) *GenerateIndexNameResult {
	return new(withTableCnmPattern).GenerateIndexName(schemaIndex, tableName, fieldName, strings.ToUpper(columnName))
}
