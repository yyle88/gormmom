package gormmom

import (
	"unicode/utf8"

	"github.com/yyle88/gormmom/gormidxname"
	"github.com/yyle88/gormmom/internal/utils"
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

func (cfg *Config) correctIndexNames(srcChanges []*TagModification) {
	var mapTagModifications = make(map[string]*TagModification, len(srcChanges))
	for _, rep := range srcChanges {
		mapTagModifications[rep.structFieldName] = rep
	}

	schemaIndexes := cfg.param.sch.ParseIndexes()
	zaplog.LOG.Debug("check_indexes", zap.String("object_class", cfg.param.sch.Name), zap.String("table_name", cfg.param.sch.Table), zap.Int("index_count", len(schemaIndexes)))
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

func (cfg *Config) rewriteSingleColumnIndex(schemaIndex schema.Index, change *TagModification) {
	zaplog.LOG.Debug("rewrite_single_column_index", zap.String("table_name", cfg.param.sch.Table), zap.String("field_name", change.structFieldName), zap.String("index_name", schemaIndex.Name), zap.String("index_class", schemaIndex.Class))

	columnName := cfg.extractTagFieldGetValue(change.modifiedTagCode, "gorm", "column")
	must.Nice(columnName)
	zaplog.LOG.Debug("new_column_name", zap.String("name", change.structFieldName), zap.String("new_column_name", columnName))

	//这个是规则的枚举名称
	var indexPrefixCate string
	switch schemaIndex.Class {
	case "":
		indexPrefixCate = "idx"
	case "UNIQUE":
		indexPrefixCate = "udx"
	default:
		indexPrefixCate = ""
	}

	var indexNamePattern gormidxname.IndexNamePattern
	if indexPrefixCate != "" {
		exist := false
		indexNamePattern, exist = cfg.resolveIndexPattern(change, indexPrefixCate, cfg.options)
		if !exist {
			//就是不存在时 使用默认值 的情况
			indexNamePattern = gormidxname.DefaultPattern
			indexNaming, ok := cfg.options.indexNamingStrategies[indexNamePattern]
			must.TRUE(ok)

			match := indexNaming.IsValidIndexName(schemaIndex.Name)
			zaplog.LOG.Debug("check_idx_match", zap.Bool("match", match))
			if !match {
				//当没有配置规则，而默认规则检查不正确时，就需要把规则名设置到标签里
				change.modifiedTagCode = syntaxgo_tag.SetTagFieldValue(change.modifiedTagCode, cfg.options.namingTagName, indexPrefixCate, string(indexNamePattern), syntaxgo_tag.INSERT_LOCATION_END)
			} else {
				//当没有配置规则，而且能够满足检查时，就不做任何事情（不要破坏用户自己配置的正确索引名）
				return
			}
		}
	} else {
		//就是不确定时 使用默认值 的情况
		indexNamePattern = gormidxname.DefaultPattern
	}

	namingImp, ok := cfg.options.indexNamingStrategies[indexNamePattern]
	must.TRUE(ok)
	idxRes := must.Nice(namingImp.GenerateIndexName(schemaIndex, cfg.param.sch.Table, change.structFieldName, columnName))
	if idxRes.NewIndexName == "" {
		return
	}
	if idxRes.IndexTagFieldName == "" {
		return
	}
	zaplog.LOG.Debug("compare", zap.String("which_enum_code_name", indexPrefixCate), zap.String("enum_code_name", idxRes.IdxUdxPrefix))
	must.Equals(indexPrefixCate, idxRes.IdxUdxPrefix)

	zaplog.LOG.Debug("new_index_name", zap.String("new_index_name", idxRes.NewIndexName))
	if idxRes.NewIndexName == schemaIndex.Name {
		return
	}

	zaplog.LOG.Debug("tag_field_name", zap.String("tag_field_name", idxRes.IndexTagFieldName))

	gormTagContent, stx, etx := syntaxgo_tag.ExtractTagValueIndex(change.modifiedTagCode, "gorm")
	must.TRUE(stx >= 0)
	must.TRUE(etx >= 0)
	must.Nice(gormTagContent) //就是排除 gorm: 以后得到的双引号里面的内容
	zaplog.LOG.Debug("gorm_tag_content", zap.String("gorm_tag_content", gormTagContent))

	var changed = utils.NewBoolean(false)
	if !changed {
		//假如连 UTF-8 编码 都不满足，就说明这个索引名是完全错误的
		if utf8.ValidString(schemaIndex.Name) {
			//因为这个正则不能匹配非 UTF-8 编码，在前面先判断编码是否正确，编码正确以后再匹配索引名
			sfx, efx := syntaxgo_tag.ExtractFieldEqualsValueIndex(gormTagContent, idxRes.IndexTagFieldName, schemaIndex.Name)
			if sfx > 0 && efx > 0 {
				spx := stx + sfx //把起点坐标补上前面的
				epx := stx + efx
				change.modifiedTagCode = change.modifiedTagCode[:spx] + idxRes.NewIndexName + change.modifiedTagCode[epx:]
				changed = true
			}
		}
	}
	if !changed {
		sfx, efx := syntaxgo_tag.ExtractNoValueFieldNameIndex(gormTagContent, idxRes.IndexTagFieldName)
		if sfx > 0 && efx > 0 {
			spx := stx + sfx //把起点坐标补上前面的
			epx := stx + efx
			change.modifiedTagCode = change.modifiedTagCode[:spx] + idxRes.IndexTagFieldName + ":" + idxRes.NewIndexName + change.modifiedTagCode[epx:]
			changed = true
		}
	}
	if !changed {
		zaplog.LOG.Debug("not_change_tag", zap.String("not_change_tag", change.modifiedTagCode))
	}

	zaplog.LOG.Debug("new_tag_string", zap.String("new_tag_string", change.modifiedTagCode))
}

func (cfg *Config) resolveIndexPattern(change *TagModification, ruleFieldName string, options *Options) (gormidxname.IndexNamePattern, bool) {
	var name = cfg.extractTagFieldGetValue(change.modifiedTagCode, options.namingTagName, ruleFieldName)
	if name == "" {
		return gormidxname.DefaultPattern, false
	}
	zaplog.LOG.Debug("index_rule_name", zap.String("index_rule_name", name))
	return gormidxname.IndexNamePattern(name), true
}
