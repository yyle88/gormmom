package gormmom

import (
	"strings"
	"unicode/utf8"

	"github.com/yyle88/gormmom/gormmomrule"
	"github.com/yyle88/gormmom/internal/utils"
	"github.com/yyle88/syntaxgo/syntaxgo_tag"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
	"gorm.io/gorm/schema"
)

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
	zaplog.LOG.Debug("new_column_name", zap.String("name", change.name), zap.String("new_column_name", newColumnName))

	var whichRuleFieldName string
	switch schemaIndex.Class {
	case "":
		whichRuleFieldName = "idx"
	case "UNIQUE":
		whichRuleFieldName = "udx"
	default:
		whichRuleFieldName = ""
	}

	var ruleEnum gormmomrule.RULE
	if whichRuleFieldName != "" {
		rule, exist := cfg.extractIndexRuleEnum(change, whichRuleFieldName)
		if !exist { //就是不存在 使用默认值 的情况
			match := rule.Validate(schemaIndex.Name)
			zaplog.LOG.Debug("check_idx_match", zap.Bool("match", match))
			if !match {
				//当没有配置规则，而默认规则检查不正确时，就需要把规则名设置到标签里
				change.code = cfg.newFixTagField(change.code, cfg.ruleTagName, whichRuleFieldName, string(rule), END)
			} else {
				//当没有配置规则，而且能够满足检查时，就不做任何事情（不要破坏用户自己配置的正确索引名）
				return
			}
		}
		ruleEnum = rule
	} else {
		ruleEnum = gormmomrule.DEFAULT
	}

	subIndexName := gormmomrule.MakeName(ruleEnum, change.name, cfg.nameFuncMap)
	utils.AssertOK(subIndexName)
	zaplog.LOG.Debug("sub_index_name", zap.String("sub_index_name", subIndexName))

	var tagFieldName string
	var newIndexName string
	switch schemaIndex.Class {
	case "":
		utils.AssertOK(whichRuleFieldName)

		tagFieldName = "index"
		newIndexName = makeIndexName("idx", param.sch.Table, subIndexName)
	case "UNIQUE":
		utils.AssertOK(whichRuleFieldName)

		tagFieldName = "uniqueIndex"
		newIndexName = makeIndexName("udx", param.sch.Table, subIndexName)
	default:
		newIndexName = makeIndexName("idx", param.sch.Table, subIndexName)

		if newIndexName != schemaIndex.Name { //这种情况暂时没有遇到，依然是暂不处理
			zaplog.LOG.Warn("new_index_name", zap.String("new_index_name", newIndexName))
		}

		if !ruleEnum.Validate(schemaIndex.Name) {
			zaplog.LOG.Warn("idx_not_match", zap.String("old_index_name", schemaIndex.Name))
		}

		return //这种情况就不处理啦，打出告警日志让开发者手动解决
	}
	utils.AssertOK(newIndexName)
	zaplog.LOG.Debug("new_index_name", zap.String("new_index_name", newIndexName))
	if newIndexName == schemaIndex.Name {
		return
	}

	zaplog.LOG.Debug("tag_field_name", zap.String("tag_field_name", tagFieldName))

	contentInGormQuotesValue, stx, etx := syntaxgo_tag.ExtractTagValueIndex(change.code, "gorm")
	utils.AssertOK(stx >= 0)
	utils.AssertOK(etx >= 0)
	utils.AssertOK(contentInGormQuotesValue) //就是排除 gorm: 以后得到的双引号里面的内容
	zaplog.LOG.Debug("gorm_tag_content", zap.String("gorm_tag_content", contentInGormQuotesValue))

	var changed = utils.FALSE()
	if !changed {
		//假如连 UTF-8 编码 都不满足，就说明这个索引名是完全错误的
		if utf8.ValidString(schemaIndex.Name) {
			//因为这个正则不能匹配非 UTF-8 编码，在前面先判断编码是否正确，编码正确以后再匹配索引名
			sfx, efx := syntaxgo_tag.ExtractFieldEqualsValueIndex(contentInGormQuotesValue, tagFieldName, schemaIndex.Name)
			if sfx > 0 && efx > 0 {
				spx := stx + sfx //把起点坐标补上前面的
				epx := stx + efx
				change.code = change.code[:spx] + newIndexName + change.code[epx:]
				changed = true
			}
		}
	}
	if !changed {
		sfx, efx := syntaxgo_tag.ExtractNoValueFieldNameIndex(contentInGormQuotesValue, tagFieldName)
		if sfx > 0 && efx > 0 {
			spx := stx + sfx //把起点坐标补上前面的
			epx := stx + efx
			change.code = change.code[:spx] + tagFieldName + ":" + newIndexName + change.code[epx:]
			changed = true
		}
	}
	if !changed {
		zaplog.LOG.Debug("not_change_tag", zap.String("not_change_tag", change.code))
	}

	zaplog.LOG.Debug("new_tag_string", zap.String("new_tag_string", change.code))
}

func (cfg *Config) extractIndexRuleEnum(change *changeType, ruleFieldName string) (gormmomrule.RULE, bool) {
	var name = cfg.extractSomeField(change.code, cfg.ruleTagName, ruleFieldName)
	if name == "" {
		return gormmomrule.DEFAULT, false
	}
	utils.AssertOK(name)
	zaplog.LOG.Debug("index_rule_name", zap.String("index_rule_name", name))
	return gormmomrule.RULE(name), true
}

func makeIndexName(prefix string, tableName string, subIndexName string) string {
	return strings.ReplaceAll(strings.Join([]string{prefix, tableName, subIndexName}, "_"), ".", "_")
}
