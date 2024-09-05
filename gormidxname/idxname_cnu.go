package gormidxname

import (
	"strings"

	"gorm.io/gorm/schema"
)

type nameGenUseCnuImp struct{}

func (G *nameGenUseCnuImp) CheckIdxName(indexName string) bool {
	return new(nameGenUseCnmImp).CheckIdxName(indexName)
}

func (G *nameGenUseCnuImp) GenIndexName(schemaIndex schema.Index, tableName string, fieldName string, columnName string) *IdxGenResType {
	return new(nameGenUseCnmImp).GenIndexName(schemaIndex, tableName, fieldName, strings.ToUpper(columnName))
}
