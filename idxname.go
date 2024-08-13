package gormmom

import (
	"strings"

	"github.com/yyle88/gormmom/gormmomrule"
	"github.com/yyle88/gormmom/internal/utils"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
	"gorm.io/gorm/schema"
)

func (cfg *Config) CheckIndexes(params []*Param) {
	for _, param := range params {
		schIndexes := param.sch.ParseIndexes()
		zaplog.LOG.Debug("check_indexes", zap.String("object_class", param.sch.Name), zap.String("table_name", param.sch.Table), zap.Int("index_count", len(schIndexes)))
		for _, node := range schIndexes {
			zaplog.LOG.Debug("foreach_index")
			zaplog.LOG.Debug("check_a_index", zap.String("index_name", node.Name), zap.Int("field_size", len(node.Fields)))
			if len(node.Fields) == 1 { //只检查单列索引，因为复合索引就得手写名称，因此没有问题
				cfg.checkSingleColumnIndex(node.Name, node.Fields[0].Name)
			} else {
				cfg.checkCompositeIndex(node.Name, node.Fields)
			}
		}
	}
}

func (cfg *Config) checkSingleColumnIndex(indexName string, fieldName string) {
	zaplog.LOG.Debug("check_single_column_index", zap.String("index_name", indexName), zap.String("field_name", fieldName))
}

func (cfg *Config) checkCompositeIndex(indexName string, fields []schema.IndexOption) {
	zaplog.LOG.Debug("check_composite_index", zap.String("index_name", indexName), zap.Int("field_size", len(fields)))
}

func (cfg *Config) rewriteIndexNames(param *Param, replacers []*changeType) {
	var replacersMap = make(map[string]*changeType, len(replacers))
	for _, rep := range replacers {
		replacersMap[rep.name] = rep
	}

	schIndexes := param.sch.ParseIndexes()
	zaplog.LOG.Debug("check_indexes", zap.String("object_class", param.sch.Name), zap.String("table_name", param.sch.Table), zap.Int("index_count", len(schIndexes)))
	for _, node := range schIndexes {
		zaplog.LOG.Debug("foreach_index")
		zaplog.LOG.Debug("check_a_index", zap.String("index_name", node.Name), zap.Int("field_size", len(node.Fields)))
		if len(node.Fields) == 1 { //只检查单列索引，因为复合索引就得手写名称，因此没有问题
			rep, ok := replacersMap[node.Fields[0].Name]
			if ok {
				cfg.rewriteSingleColumnIndex(param, node, rep)
			} else {
				cfg.checkSingleColumnIndex(node.Name, node.Fields[0].Name)
			}
		} else {
			cfg.checkCompositeIndex(node.Name, node.Fields)
		}
	}
}

func (cfg *Config) rewriteSingleColumnIndex(param *Param, schemaIndex schema.Index, change *changeType) {
	zaplog.LOG.Debug("rewrite_single_column_index", zap.String("table_name", param.sch.Table), zap.String("field_name", change.name), zap.String("index_name", schemaIndex.Name), zap.String("index_class", schemaIndex.Class))

	newColumnName := cfg.extractSomeField(change.code, "gorm", "column")
	utils.AssertOK(newColumnName)
	zaplog.LOG.Debug("new_column_name", zap.String("new_column_name", newColumnName))

	indexRuleName := cfg.extractSomeField(change.code, cfg.ruleTagName, "idx_rule")
	if indexRuleName != "" {
		zaplog.LOG.Panic("NOT IMPLEMENTED") //现在没空整太多用不到的，就暂时只支持相同的规则，而不支持独立的规则
	} else {
		indexRuleName = cfg.extractRuleField(change.code)
		utils.AssertOK(indexRuleName)
	}
	zaplog.LOG.Debug("index_rule_name", zap.String("index_rule_name", indexRuleName))

	rule := gormmomrule.RULE(indexRuleName)
	subIndexName := gormmomrule.MakeName(rule, change.name, cfg.nameFuncMap)
	utils.AssertOK(subIndexName)
	zaplog.LOG.Debug("sub_index_name", zap.String("sub_index_name", subIndexName))

	var tagFieldName string
	var newIndexName string
	switch schemaIndex.Class {
	case "":
		tagFieldName = "index"
		newIndexName = makeIndexName("idx", param.sch.Table, subIndexName)
	case "UNIQUE":
		tagFieldName = "uniqueIndex"
		newIndexName = makeIndexName("udx", param.sch.Table, subIndexName)
	default:
		newIndexName = makeIndexName("idx", param.sch.Table, subIndexName)

		if newIndexName != schemaIndex.Name { //这种情况暂时没有遇到，依然是暂不处理
			zaplog.LOG.Warn("new_index_name", zap.String("new_index_name", newIndexName))
		}

		return //这种情况就不处理啦，打出告警日志让开发者手动解决
	}
	utils.AssertOK(newIndexName)
	zaplog.LOG.Debug("new_index_name", zap.String("new_index_name", newIndexName))
	if newIndexName == schemaIndex.Name {
		return
	}

	zaplog.LOG.Debug("tag_field_name", zap.String("tag_field_name", tagFieldName))

	//这里还是有些问题
	//change.code = cfg.newFixTagField(change.code, "gorm", tagFieldName, newIndexName, END)

	zaplog.LOG.Debug("new_gorm_tag", zap.String("new_gorm_tag", change.code))
}

func makeIndexName(prefix string, tableName string, subIndexName string) string {
	return strings.ReplaceAll(strings.Join([]string{prefix, tableName, subIndexName}, "_"), ".", "_")
}
