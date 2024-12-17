package gormmom

import (
	"fmt"
	"go/ast"
	"os"
	"slices"

	"github.com/yyle88/done"
	"github.com/yyle88/formatgo"
	"github.com/yyle88/gormmom/gormmomname"
	"github.com/yyle88/gormmom/internal/utils"
	"github.com/yyle88/syntaxgo/syntaxgo_ast"
	"github.com/yyle88/syntaxgo/syntaxgo_astnode"
	"github.com/yyle88/syntaxgo/syntaxgo_search"
	"github.com/yyle88/syntaxgo/syntaxgo_tag"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
	"gorm.io/gorm/schema"
)

type ConfigBatch struct {
	configs []*Config
}

func NewConfigBatch(schemaCaches []*SchemaCache, options *Options) *ConfigBatch {
	var configs []*Config
	for _, param := range schemaCaches {
		configs = append(configs, NewConfig(param, options))
	}
	return &ConfigBatch{configs: configs}
}

func (c *ConfigBatch) GenReplaces() {
	for _, cfg := range c.configs {
		cfg.GenReplace()
	}
}

type Config struct {
	schemaCache *SchemaCache
	options     *Options
}

func NewConfig(schemaCache *SchemaCache, options *Options) *Config {
	return &Config{
		schemaCache: schemaCache,
		options:     options,
	}
}

func (cfg *Config) GenReplace() {
	utils.MustWriteFile(cfg.schemaCache.sourcePath, done.VAE(formatgo.FormatBytes(cfg.GetNewCode())).Nice())
}

func (cfg *Config) GetNewCode() []byte {
	cfg.schemaCache.Validate()

	sourceCode := done.VAE(os.ReadFile(cfg.schemaCache.sourcePath)).Nice()

	astBundle := done.VCE(syntaxgo_ast.NewAstBundleV1(sourceCode)).Nice()

	astFile, fileSet := astBundle.GetBundle()

	structContent, ok := syntaxgo_search.FindStructTypeByName(astFile, cfg.schemaCache.structName)
	if !ok {
		const reason = "CAN NOT FIND STRUCT TYPE"
		zaplog.LOG.Panic(reason, zap.String("struct_name", cfg.schemaCache.structName))
		panic(reason)
	}
	done.Done(ast.Print(fileSet, structContent))

	sourceModifications := cfg.collectTagModifications(structContent)

	if cfg.options.renewIndexName {
		//这里增加个新逻辑，就是单列索引的索引名称不正确，需要也校正索引名，因此这个函数会补充标签内容
		cfg.correctIndexNames(sourceModifications)

		zaplog.LOG.Debug("change_index_names")
		for _, rep := range sourceModifications {
			zaplog.LOG.Debug("check_index:", zap.String("name", rep.structFieldName), zap.String("code", rep.modifiedTagCode))
		}
	}

	//需要翻转下从后往前替换，因为替换以后源码会变，假如从前往后替换坐标就对不上啦，而从后往前替换则不存在这个问题
	slices.Reverse(sourceModifications)

	//接下来替换代码，把需要 新增 或者 替换 的标签都设置到代码里
	newCode := sourceCode
	for _, step := range sourceModifications {
		newCode = syntaxgo_astnode.ChangeNodeCode(newCode, step.previousTagNode, []byte(step.modifiedTagCode))
	}
	return newCode
}

type defineTagModification struct {
	structFieldName string   //结构体的字段名
	previousTagNode ast.Node //标签的起止位置-就是在src源码中的位置，便于后面的替换代码
	modifiedTagCode string   //新标签的新内容-就是标签的完整全部内容
}

func (cfg *Config) collectTagModifications(structContent *ast.StructType) []*defineTagModification {
	var defineTagModifications []*defineTagModification

	// 遍历结构体的字段
	for _, field := range structContent.Fields.List {
		// 打印字段名称和类型
		for _, nameIdent := range field.Names {
			zaplog.LOG.Debug("--")
			zaplog.LOG.Debug("process", zap.String("struct_field_name:", nameIdent.Name))
			zaplog.LOG.Debug("--")

			schemaField, exist := cfg.schemaCache.schColumns[nameIdent.Name]
			if !exist { //比如字段是 "V哈哈" 就没事 而假如是 "v哈哈" 或者 "哈哈" 就不行，因为非以大写字母开始的字段，就没有gorm的列名
				zaplog.LOG.Debug("NO SCHEMA_FIELD - MAYBE NAME IS UNEXPORTED", zap.String("struct_name", cfg.schemaCache.structName), zap.String("struct_field_name", nameIdent.Name))
				continue
			}

			if field.Tag == nil {
				zaplog.LOG.Debug("NO TAG", zap.String("struct_name", cfg.schemaCache.structName), zap.String("struct_field_name", nameIdent.Name))
				if cfg.options.skipBasicNaming && cfg.options.columnNamingStrategies[gormmomname.DefaultPattern].IsValidColumnName(schemaField.DBName) {
					zaplog.LOG.Debug("SKIP SIMPLE FIELD", zap.String("struct_name", cfg.schemaCache.structName), zap.String("struct_field_name", nameIdent.Name))
					continue
				}
				columnNamePattern := cfg.options.defaultColumnNamePattern
				if !cfg.options.columnNamingStrategies[columnNamePattern].IsValidColumnName(schemaField.DBName) {
					if len(field.Names) >= 2 { //比如 a,b int 这种两个字段在一起，但其中一个字段的列名不正确时，就没法自动解决啦（其实有办法但不想实现，因为代价较大而没有收益）
						const reason = "CAN NOT HANDLE THIS SITUATION"
						zaplog.LOG.Panic(reason, zap.String("struct_name", cfg.schemaCache.structName), zap.String("struct_field_name", nameIdent.Name))
						panic(reason) //这种情况下当有错时，就不处理这种情况，就需要程序员先把两个字段定义到两行里
					}
					changeTag := cfg.modifyFieldTagCorrection(schemaField, "``", columnNamePattern) //这里应该走创建标签的逻辑，但和修改标签的逻辑是相同的
					defineTagModifications = append(defineTagModifications, &defineTagModification{
						structFieldName: nameIdent.Name,
						previousTagNode: syntaxgo_astnode.NewNode(field.End(), field.End()), //在尾部插入新的标签，要紧贴字段而且在换行符前面
						modifiedTagCode: changeTag,
					})
				}
			} else {
				if cnmPattern := cfg.extractTagGetCnmPattern(field.Tag.Value); cnmPattern != "" {
					columnNamePattern := gormmomname.ColumnNamePattern(cnmPattern)
					zaplog.LOG.Debug("process", zap.String("column_name_pattern", string(columnNamePattern)))
					changeTag := cfg.modifyFieldTagCorrection(schemaField, field.Tag.Value, columnNamePattern)
					defineTagModifications = append(defineTagModifications, &defineTagModification{
						structFieldName: nameIdent.Name,
						previousTagNode: field.Tag, //完整替换原来的标签
						modifiedTagCode: changeTag,
					})
				} else {
					if cfg.options.skipBasicNaming && cfg.options.columnNamingStrategies[gormmomname.DefaultPattern].IsValidColumnName(schemaField.DBName) {
						zaplog.LOG.Debug("SKIP SIMPLE FIELD", zap.String("struct_name", cfg.schemaCache.structName), zap.String("struct_field_name", nameIdent.Name))
						continue
					}
					columnNamePattern := cfg.options.defaultColumnNamePattern
					if !cfg.options.columnNamingStrategies[columnNamePattern].IsValidColumnName(schemaField.DBName) { //按照比较宽泛的规则也校验不过的时候就需要修正字段名
						changeTag := cfg.modifyFieldTagCorrection(schemaField, field.Tag.Value, columnNamePattern)
						defineTagModifications = append(defineTagModifications, &defineTagModification{
							structFieldName: nameIdent.Name,
							previousTagNode: field.Tag, //完整替换原来的标签
							modifiedTagCode: changeTag,
						})
					} else {
						zaplog.LOG.Debug("match-pattern-so-skip", zap.String("name", nameIdent.Name), zap.String("tag", field.Tag.Value))
					}
				}
			}
		}
	}

	zaplog.LOG.Debug("change_column_names")
	for _, rep := range defineTagModifications {
		zaplog.LOG.Debug("check_column:", zap.String("name", rep.structFieldName), zap.String("code", rep.modifiedTagCode))
	}
	return defineTagModifications
}

func (cfg *Config) extractTagGetCnmPattern(tagCode string) string {
	return cfg.extractTagFieldGetValue(tagCode, cfg.options.namingTagName, cfg.options.columnNamePatternFieldName)
}

func (cfg *Config) extractTagFieldGetValue(tagCode string, key1 string, key2 string) string {
	tagValue := syntaxgo_tag.ExtractTagValue(tagCode, key1)
	if tagValue == "" {
		return ""
	}
	tagField := syntaxgo_tag.ExtractTagField(tagValue, key2, syntaxgo_tag.EXCLUDE_WHITESPACE_PREFIX)
	if tagField == "" {
		return ""
	}
	return tagField
}

func (cfg *Config) modifyFieldTagCorrection(schemaField *schema.Field, tag string, columnNamePattern gormmomname.ColumnNamePattern) string {
	zaplog.LOG.Debug("new_fix_tag_code", zap.String("name", schemaField.Name), zap.String("tag", tag))
	//在 gorm 里修改 column 内容
	newTag := cfg.modifyGormTagWithColumn(schemaField, tag, columnNamePattern)
	//在 规则 里修改 column-name-pattern 内容
	newTag = cfg.modifyPatternTagWithName(newTag, columnNamePattern)
	//这是替换后的结果，即替换整个标签内容，获得新的完整标签内容
	zaplog.LOG.Debug("new_fix_tag_code", zap.String("name", schemaField.Name), zap.String("new_tag", newTag))
	return newTag
}

func (cfg *Config) modifyGormTagWithColumn(schemaField *schema.Field, tag string, columnNamePattern gormmomname.ColumnNamePattern) string {
	var columnName = cfg.options.columnNamingStrategies[columnNamePattern].GenerateColumnName(schemaField.Name)
	zaplog.LOG.Debug("new_fix_gorm_tag", zap.String("name", schemaField.Name), zap.String("column_name", columnName))

	tagValue, sdx, edx := syntaxgo_tag.ExtractTagValueIndex(tag, "gorm")
	if sdx < 0 || edx < 0 { //表示没找到 gorm 相关的内容
		if tagValue != "" {
			zaplog.LOG.Panic("IMPOSSIBLE")
		}
		part := fmt.Sprintf(`gorm:"column:%s;"`, columnName)
		if tag[1] != ' ' && tag[1] != '`' {
			part += " " //说明后面还有别的标签
		}
		p := 1 //插在第一个"`"的后面，即第一位的位置
		return tag[:p] + part + tag[p:]
	}
	//设置这个标签的这个字段的值
	return syntaxgo_tag.SetTagFieldValue(tag, "gorm", "column", columnName, syntaxgo_tag.INSERT_LOCATION_TOP)
}

func (cfg *Config) modifyPatternTagWithName(tag string, columnNamePattern gormmomname.ColumnNamePattern) string {
	zaplog.LOG.Debug("modify-pattern-tag-with-name", zap.String("column_name_pattern", string(columnNamePattern)))

	tagValue, sdx, edx := syntaxgo_tag.ExtractTagValueIndex(tag, cfg.options.namingTagName)
	if sdx < 0 || edx < 0 { //表示没找到 gorm 相关的内容
		if tagValue != "" {
			zaplog.LOG.Panic("IMPOSSIBLE")
		}
		part := fmt.Sprintf(`%s:"%s:%s;"`, cfg.options.namingTagName, cfg.options.columnNamePatternFieldName, string(columnNamePattern))
		if rch := tag[len(tag)-2]; rch != ' ' && rch != '`' {
			part = " " + part //说明前面还有别的标签
		}
		p := len(tag) - 1 //插在最后一个"`"的前面，即最后一位的位置
		return tag[:p] + part + tag[p:]
	}
	//设置这个标签的这个字段的值
	return syntaxgo_tag.SetTagFieldValue(tag, cfg.options.namingTagName, cfg.options.columnNamePatternFieldName, string(columnNamePattern), syntaxgo_tag.INSERT_LOCATION_TOP)
}
