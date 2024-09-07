package gormmom

import (
	"unicode/utf8"

	"github.com/yyle88/gormmom/gormidxname"
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

func (cfg *Config) rewriteIndexNames(param *Param, srcChanges []*changeType) {
	var changesMap = make(map[string]*changeType, len(srcChanges))
	for _, rep := range srcChanges {
		changesMap[rep.vFieldName] = rep
	}

	schIndexes := param.sch.ParseIndexes()
	zaplog.LOG.Debug("check_indexes", zap.String("object_class", param.sch.Name), zap.String("table_name", param.sch.Table), zap.Int("index_count", len(schIndexes)))
	for _, node := range schIndexes {
		zaplog.LOG.Debug("foreach_index")
		zaplog.LOG.Debug("check_a_index", zap.String("index_name", node.Name), zap.Int("field_size", len(node.Fields)))
		if len(node.Fields) == 1 { //只检查单列索引，因为复合索引就得手写名称，因此没有问题
			rep, ok := changesMap[node.Fields[0].Name]
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
	zaplog.LOG.Debug("rewrite_single_column_index", zap.String("table_name", param.sch.Table), zap.String("field_name", change.vFieldName), zap.String("index_name", schemaIndex.Name), zap.String("index_class", schemaIndex.Class))

	newColumnName := cfg.extractSomeField(change.newTagCode, "gorm", "column")
	utils.AssertOK(newColumnName)
	zaplog.LOG.Debug("new_column_name", zap.String("name", change.vFieldName), zap.String("new_column_name", newColumnName))

	//这个是规则的枚举名称
	var whichEnumCodeName string
	switch schemaIndex.Class {
	case "":
		whichEnumCodeName = "idx"
	case "UNIQUE":
		whichEnumCodeName = "udx"
	default:
		whichEnumCodeName = ""
	}

	var idxNameEnum gormidxname.IdxNAME
	if whichEnumCodeName != "" {
		exist := false
		idxNameEnum, exist = cfg.extractIdxNameEnum(change, whichEnumCodeName)
		if !exist {
			//就是不存在时 使用默认值 的情况
			idxNameEnum = gormidxname.DEFAULT
			idxNameImp, ok := cfg.idxNameMap[idxNameEnum]
			utils.AssertOK(ok)

			match := idxNameImp.CheckIdxName(schemaIndex.Name)
			zaplog.LOG.Debug("check_idx_match", zap.Bool("match", match))
			if !match {
				//当没有配置规则，而默认规则检查不正确时，就需要把规则名设置到标签里
				change.newTagCode = cfg.newFixTagField(change.newTagCode, cfg.ruleTagName, whichEnumCodeName, string(idxNameEnum), END)
			} else {
				//当没有配置规则，而且能够满足检查时，就不做任何事情（不要破坏用户自己配置的正确索引名）
				return
			}
		}
	} else {
		//就是不确定时 使用默认值 的情况
		idxNameEnum = gormidxname.DEFAULT
	}

	idxNameImp, ok := cfg.idxNameMap[idxNameEnum]
	utils.AssertOK(ok)
	newInm := idxNameImp.GenIndexName(schemaIndex, param.sch.Table, change.vFieldName, newColumnName)
	utils.AssertOK(newInm)
	if newInm.NewIndexName == "" {
		return
	}
	if newInm.TagFieldName == "" {
		return
	}
	zaplog.LOG.Debug("compare", zap.String("which_enum_code_name", whichEnumCodeName), zap.String("enum_code_name", newInm.EnumCodeName))
	utils.AssertEquals(whichEnumCodeName, newInm.EnumCodeName)

	zaplog.LOG.Debug("new_index_name", zap.String("new_index_name", newInm.NewIndexName))
	if newInm.NewIndexName == schemaIndex.Name {
		return
	}

	zaplog.LOG.Debug("tag_field_name", zap.String("tag_field_name", newInm.TagFieldName))

	contentInGormQuotesValue, stx, etx := syntaxgo_tag.ExtractTagValueIndex(change.newTagCode, "gorm")
	utils.AssertOK(stx >= 0)
	utils.AssertOK(etx >= 0)
	utils.AssertOK(contentInGormQuotesValue) //就是排除 gorm: 以后得到的双引号里面的内容
	zaplog.LOG.Debug("gorm_tag_content", zap.String("gorm_tag_content", contentInGormQuotesValue))

	var changed = utils.FALSE()
	if !changed {
		//假如连 UTF-8 编码 都不满足，就说明这个索引名是完全错误的
		if utf8.ValidString(schemaIndex.Name) {
			//因为这个正则不能匹配非 UTF-8 编码，在前面先判断编码是否正确，编码正确以后再匹配索引名
			sfx, efx := syntaxgo_tag.ExtractFieldEqualsValueIndex(contentInGormQuotesValue, newInm.TagFieldName, schemaIndex.Name)
			if sfx > 0 && efx > 0 {
				spx := stx + sfx //把起点坐标补上前面的
				epx := stx + efx
				change.newTagCode = change.newTagCode[:spx] + newInm.NewIndexName + change.newTagCode[epx:]
				changed = true
			}
		}
	}
	if !changed {
		sfx, efx := syntaxgo_tag.ExtractNoValueFieldNameIndex(contentInGormQuotesValue, newInm.TagFieldName)
		if sfx > 0 && efx > 0 {
			spx := stx + sfx //把起点坐标补上前面的
			epx := stx + efx
			change.newTagCode = change.newTagCode[:spx] + newInm.TagFieldName + ":" + newInm.NewIndexName + change.newTagCode[epx:]
			changed = true
		}
	}
	if !changed {
		zaplog.LOG.Debug("not_change_tag", zap.String("not_change_tag", change.newTagCode))
	}

	zaplog.LOG.Debug("new_tag_string", zap.String("new_tag_string", change.newTagCode))
}

func (cfg *Config) extractIdxNameEnum(change *changeType, ruleFieldName string) (gormidxname.IdxNAME, bool) {
	var name = cfg.extractSomeField(change.newTagCode, cfg.ruleTagName, ruleFieldName)
	if name == "" {
		return gormidxname.DEFAULT, false
	}
	utils.AssertOK(name)
	zaplog.LOG.Debug("index_rule_name", zap.String("index_rule_name", name))
	return gormidxname.IdxNAME(name), true
}
