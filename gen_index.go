package gormmom

import (
	"unicode/utf8"

	"github.com/yyle88/gormmom/gormidxname"
	"github.com/yyle88/must"
	"github.com/yyle88/syntaxgo/syntaxgo_tag"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
	"gorm.io/gorm/schema"
)

func (cfg *Config) validateSingleColumnIndex(indexName string, fieldName string) {
	zaplog.LOG.Debug("validate-single-column-index", zap.String("index_name", indexName), zap.String("field_name", fieldName))
}

func (cfg *Config) validateCompositeIndex(indexName string, fields []schema.IndexOption) {
	zaplog.LOG.Debug("validate-composite-index", zap.String("index_name", indexName), zap.Int("field_size", len(fields)))
}

func (cfg *Config) correctIndexNames(srcChanges []*defineTagModification) {
	var mapTagModifications = make(map[string]*defineTagModification, len(srcChanges))
	for _, rep := range srcChanges {
		mapTagModifications[rep.structFieldName] = rep
	}

	schemaIndexes := cfg.schemaX.sch.ParseIndexes()
	zaplog.LOG.Debug("check_indexes", zap.String("object_class", cfg.schemaX.sch.Name), zap.String("table_name", cfg.schemaX.sch.Table), zap.Int("index_count", len(schemaIndexes)))
	for _, node := range schemaIndexes {
		zaplog.LOG.Debug("foreach_index")
		zaplog.LOG.Debug("check_a_index", zap.String("index_name", node.Name), zap.Int("field_size", len(node.Fields)))
		if len(node.Fields) == 1 { //只检查单列索引，因为复合索引就得手写名称，因此没有问题
			rep, ok := mapTagModifications[node.Fields[0].Name]
			if ok {
				cfg.rewriteSingleColumnIndex(node, rep)
			} else {
				cfg.validateSingleColumnIndex(node.Name, node.Fields[0].Name)
			}
		} else {
			cfg.validateCompositeIndex(node.Name, node.Fields)
		}
	}
}

func (cfg *Config) rewriteSingleColumnIndex(schemaIndex *schema.Index, modification *defineTagModification) {
	zaplog.LOG.Debug("rewrite_single_column_index", zap.String("table_name", cfg.schemaX.sch.Table), zap.String("field_name", modification.structFieldName), zap.String("index_name", schemaIndex.Name), zap.String("index_class", schemaIndex.Class))

	columnName := cfg.extractTagFieldGetValue(modification.modifiedTagCode, "gorm", "column")
	must.Nice(columnName)
	zaplog.LOG.Debug("new_column_name", zap.String("name", modification.structFieldName), zap.String("new_column_name", columnName))

	//这个是规则的枚举名称
	var patternTagName string
	switch schemaIndex.Class {
	case "":
		patternTagName = "idx"
	case "UNIQUE":
		patternTagName = "udx"
	default:
		patternTagName = ""
	}

	defaultPattern := cfg.options.indexNamingStrategies.GetDefault()

	var patternEnum gormidxname.PatternEnum
	if patternTagName != "" {
		exist := false
		patternEnum, exist = cfg.resolveIndexPattern(modification, patternTagName)
		if !exist {
			patternEnum = defaultPattern.GetPatternEnum()

			match := defaultPattern.CheckIndexName(schemaIndex.Name)
			zaplog.LOG.Debug("check_idx_match", zap.Bool("match", match))
			if !match {
				//当没有配置规则，而默认规则检查不正确时，就需要把规则名设置到标签里
				modification.modifiedTagCode = syntaxgo_tag.SetTagFieldValue(modification.modifiedTagCode, cfg.options.tagName, patternTagName, string(patternEnum), syntaxgo_tag.INSERT_LOCATION_END)
			} else {
				//当没有配置规则，而且能够满足检查时，就不做任何事情（不要破坏用户自己配置的正确索引名）
				return
			}
		}
	} else {
		//就是不确定时 使用默认值 的情况
		patternEnum = defaultPattern.GetPatternEnum()
	}

	pattern := cfg.options.indexNamingStrategies.GetPattern(patternEnum)

	indexNameResult := must.Nice(pattern.BuildIndexName(schemaIndex, &gormidxname.BuildIndexParam{
		TableName:  cfg.schemaX.sch.Table,
		FieldName:  modification.structFieldName,
		ColumnName: columnName,
	}))
	if indexNameResult.NewIndexName == "" {
		return
	}
	if indexNameResult.TagFieldName == "" {
		return
	}
	zaplog.LOG.Debug("compare", zap.String("which_enum_code_name", patternTagName), zap.String("enum_code_name", indexNameResult.IdxUdxPrefix))
	must.Equals(patternTagName, indexNameResult.IdxUdxPrefix)

	zaplog.LOG.Debug("new_index_name", zap.String("new_index_name", indexNameResult.NewIndexName))
	if indexNameResult.NewIndexName == schemaIndex.Name {
		return
	}

	zaplog.LOG.Debug("tag_field_name", zap.String("tag_field_name", indexNameResult.TagFieldName))

	gormTagContent, stx, etx := syntaxgo_tag.ExtractTagValueIndex(modification.modifiedTagCode, "gorm")
	must.TRUE(stx >= 0)
	must.TRUE(etx >= 0)
	must.Nice(gormTagContent) //就是排除 gorm: 以后得到的双引号里面的内容
	zaplog.LOG.Debug("gorm_tag_content", zap.String("gorm_tag_content", gormTagContent))

	var changed = false
	//假如连 UTF-8 编码 都不满足，就说明这个索引名是完全错误的
	if utf8.ValidString(schemaIndex.Name) {
		//因为这个正则不能匹配非 UTF-8 编码，在前面先判断编码是否正确，编码正确以后再匹配索引名
		sfx, efx := syntaxgo_tag.ExtractFieldEqualsValueIndex(gormTagContent, indexNameResult.TagFieldName, schemaIndex.Name)
		if sfx > 0 && efx > 0 {
			spx := stx + sfx //把起点坐标补上前面的
			epx := stx + efx
			modification.modifiedTagCode = modification.modifiedTagCode[:spx] + indexNameResult.NewIndexName + modification.modifiedTagCode[epx:]
			changed = true
		}
	}
	if !changed {
		sfx, efx := syntaxgo_tag.ExtractNoValueFieldNameIndex(gormTagContent, indexNameResult.TagFieldName)
		if sfx > 0 && efx > 0 {
			spx := stx + sfx //把起点坐标补上前面的
			epx := stx + efx
			modification.modifiedTagCode = modification.modifiedTagCode[:spx] + indexNameResult.TagFieldName + ":" + indexNameResult.NewIndexName + modification.modifiedTagCode[epx:]
			changed = true
		}
	}
	if !changed {
		zaplog.LOG.Debug("not_change_tag", zap.String("not_change_tag", modification.modifiedTagCode))
	}

	zaplog.LOG.Debug("new_tag_string", zap.String("new_tag_string", modification.modifiedTagCode))
}

func (cfg *Config) resolveIndexPattern(modification *defineTagModification, patternTagName string) (gormidxname.PatternEnum, bool) {
	var name = cfg.extractTagFieldGetValue(modification.modifiedTagCode, cfg.options.tagName, patternTagName)
	if name == "" {
		defaultPattern := cfg.options.indexNamingStrategies.GetDefault()
		return defaultPattern.GetPatternEnum(), false
	}
	zaplog.LOG.Debug("resolve-index-pattern", zap.String("index_name_pattern", name))
	return gormidxname.PatternEnum(name), true
}
