package gormidxname

import (
	"github.com/yyle88/gormmom/internal/utils"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
	"gorm.io/gorm/schema"
	"regexp"
	"strings"
)

type nameGenFromColumnName struct{}

func (x *nameGenFromColumnName) CheckIdxName(indexName string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_]{1,63}$`).MatchString(indexName)
}

func (x *nameGenFromColumnName) GenIndexName(schemaIndex schema.Index, tableName string, fieldName string, columnName string) *IdxGenResType {
	zaplog.LOG.Warn(
		"new_index_name",
		zap.String("table_name", tableName),
		zap.String("field_name", fieldName),
		zap.String("column_name", columnName),
	)

	var enumCodeName string
	var tagFieldName string
	var newIndexName string
	switch schemaIndex.Class {
	case "":
		enumCodeName = "idx"
		tagFieldName = "index"
		newIndexName = x.makeIndexName("idx", tableName, columnName)
	case "UNIQUE":
		enumCodeName = "udx"
		tagFieldName = "uniqueIndex"
		newIndexName = x.makeIndexName("udx", tableName, columnName)
	default:
		newIndexName = x.makeIndexName("idx", tableName, columnName)

		if newIndexName != schemaIndex.Name { //这种情况暂时没有遇到，依然是暂不处理
			zaplog.LOG.Warn("new_index_name", zap.String("new_index_name", newIndexName))
		}

		if !x.CheckIdxName(schemaIndex.Name) {
			zaplog.LOG.Warn("idx_not_match", zap.String("old_index_name", schemaIndex.Name))
		}

		return &IdxGenResType{
			TagFieldName: "", //这种情况就不处理啦，打出告警日志让开发者手动解决
			NewIndexName: newIndexName,
			EnumCodeName: "",
		}
	}
	utils.AssertOK(newIndexName)

	return &IdxGenResType{
		TagFieldName: tagFieldName,
		NewIndexName: newIndexName,
		EnumCodeName: enumCodeName,
	}
}

func (x *nameGenFromColumnName) makeIndexName(prefix string, tableName string, nameSuffix string) string {
	return strings.ReplaceAll(strings.Join([]string{prefix, tableName, nameSuffix}, "_"), ".", "_")
}
