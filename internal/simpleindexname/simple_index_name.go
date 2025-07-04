package simpleindexname

import (
	"strings"

	"github.com/yyle88/gormmom/internal/utils"
	"github.com/yyle88/must"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
	"gorm.io/gorm/schema"
)

const (
	IdxPatternTagName = "idx"
	UdxPatternTagName = "udx"
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
		enumCodeName = IdxPatternTagName
		tagFieldName = "index"
		newIndexName = mergeIndexName("idx", param.TableName, param.ColumnName)
	case "UNIQUE":
		enumCodeName = UdxPatternTagName
		tagFieldName = "uniqueIndex"
		newIndexName = mergeIndexName("udx", param.TableName, param.ColumnName)
	default:
		newIndexName = mergeIndexName("idx", param.TableName, param.ColumnName)

		if newIndexName != schemaIndex.Name { //这种情况暂时没有遇到，依然是暂不处理
			zaplog.LOG.Warn("new_index_name", zap.String("new_index_name", newIndexName))
		}

		if !utils.NewCommonRegexp(63).MatchString(schemaIndex.Name) {
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
