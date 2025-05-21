package simpleindex

import (
	"regexp"
	"strings"

	"github.com/yyle88/must"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
	"gorm.io/gorm/schema"
)

type BuildIndexParam struct {
	TableName  string
	FieldName  string
	ColumnName string
}

type IndexNameResult struct {
	TagFieldName string
	NewIndexName string
	IdxUdxPrefix string
}

func BuildIndexName(schemaIndex *schema.Index, param *BuildIndexParam) *IndexNameResult {
	zaplog.LOG.Debug(
		"new_index_name",
		zap.String("table_name", param.TableName),
		zap.String("field_name", param.FieldName),
		zap.String("column_name", param.ColumnName),
	)

	var enumCodeName string
	var tagFieldName string
	var newIndexName string
	switch schemaIndex.Class {
	case "":
		enumCodeName = "idx"
		tagFieldName = "index"
		newIndexName = mergeIndexName("idx", param.TableName, param.ColumnName)
	case "UNIQUE":
		enumCodeName = "udx"
		tagFieldName = "uniqueIndex"
		newIndexName = mergeIndexName("udx", param.TableName, param.ColumnName)
	default:
		newIndexName = mergeIndexName("idx", param.TableName, param.ColumnName)

		if newIndexName != schemaIndex.Name { //这种情况暂时没有遇到，依然是暂不处理
			zaplog.LOG.Warn("new_index_name", zap.String("new_index_name", newIndexName))
		}

		if !regexp.MustCompile(`^[a-zA-Z0-9_]{1,63}$`).MatchString(schemaIndex.Name) {
			zaplog.LOG.Warn("idx_not_match", zap.String("old_index_name", schemaIndex.Name))
		}

		return &IndexNameResult{
			TagFieldName: "", //这种情况就不处理啦，打出告警日志让开发者手动解决
			NewIndexName: newIndexName,
			IdxUdxPrefix: "",
		}
	}
	must.Nice(newIndexName)

	return &IndexNameResult{
		TagFieldName: tagFieldName,
		NewIndexName: newIndexName,
		IdxUdxPrefix: enumCodeName,
	}
}

func mergeIndexName(prefix string, tableName string, suffix string) string {
	return strings.ReplaceAll(strings.Join([]string{prefix, tableName, suffix}, "_"), ".", "_")
}

func CheckLength(name string, maxLength int) {
	if len(name) > maxLength {
		zaplog.LOG.Panic("INDEX-NAME-IS-TOO-LONG", zap.String("index_name", name), zap.Int("index_name_length", len(name)), zap.Int("max_length", maxLength))
	}
}
