package gormidxname

import (
	"regexp"
	"strings"

	"github.com/yyle88/must"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
	"gorm.io/gorm/schema"
)

type withTableCnmPattern struct{}

func (G *withTableCnmPattern) IsValidIndexName(indexName string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_]{1,63}$`).MatchString(indexName)
}

func (G *withTableCnmPattern) GenerateIndexName(schemaIndex *schema.Index, tableName string, fieldName string, columnName string) *GenerateIndexNameResult {
	zaplog.LOG.Debug(
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
		newIndexName = G.combineIndexName("idx", tableName, columnName)
	case "UNIQUE":
		enumCodeName = "udx"
		tagFieldName = "uniqueIndex"
		newIndexName = G.combineIndexName("udx", tableName, columnName)
	default:
		newIndexName = G.combineIndexName("idx", tableName, columnName)

		if newIndexName != schemaIndex.Name { //这种情况暂时没有遇到，依然是暂不处理
			zaplog.LOG.Warn("new_index_name", zap.String("new_index_name", newIndexName))
		}

		if !G.IsValidIndexName(schemaIndex.Name) {
			zaplog.LOG.Warn("idx_not_match", zap.String("old_index_name", schemaIndex.Name))
		}

		return &GenerateIndexNameResult{
			IndexTagFieldName: "", //这种情况就不处理啦，打出告警日志让开发者手动解决
			NewIndexName:      newIndexName,
			IdxUdxPrefix:      "",
		}
	}
	must.Nice(newIndexName)

	return &GenerateIndexNameResult{
		IndexTagFieldName: tagFieldName,
		NewIndexName:      newIndexName,
		IdxUdxPrefix:      enumCodeName,
	}
}

func (G *withTableCnmPattern) combineIndexName(prefix string, tableName string, nameSuffix string) string {
	return strings.ReplaceAll(strings.Join([]string{prefix, tableName, nameSuffix}, "_"), ".", "_")
}
