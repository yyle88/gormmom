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
	"github.com/yyle88/rese"
	"github.com/yyle88/syntaxgo/syntaxgo_ast"
	"github.com/yyle88/syntaxgo/syntaxgo_astnode"
	"github.com/yyle88/syntaxgo/syntaxgo_search"
	"github.com/yyle88/syntaxgo/syntaxgo_tag"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
	"gorm.io/gorm/schema"
)

type Configs []*Config

func NewConfigs(schemaXs []*SchemaX, options *Options) Configs {
	var configs = make([]*Config, 0, len(schemaXs))
	for _, schemaX := range schemaXs {
		configs = append(configs, NewConfig(schemaX, options))
	}
	return configs
}

func (configs Configs) GenReplaces() {
	for _, config := range configs {
		config.GenReplace()
	}
}

type Config struct {
	schemaX *SchemaX
	options *Options
}

func NewConfig(schemaX *SchemaX, options *Options) *Config {
	return &Config{
		schemaX: schemaX,
		options: options,
	}
}

func (cfg *Config) GenReplace() {
	utils.WriteFile(cfg.schemaX.sourcePath, rese.A1(formatgo.FormatBytes(cfg.GetNewCode())))
}

func (cfg *Config) GetNewCode() []byte {
	cfg.schemaX.Validate()

	sourceCode := rese.A1(os.ReadFile(cfg.schemaX.sourcePath))
	astBundle := rese.C1(syntaxgo_ast.NewAstBundleV1(sourceCode))
	astFile, fileSet := astBundle.GetBundle()
	// 使用语法分析树 ast 找到结构体的代码，就像这样
	// struct {
	//    Name string
	// }
	// 结果只包含结构体的内容
	astStructType, ok := syntaxgo_search.FindStructTypeByName(astFile, cfg.schemaX.structName)
	if !ok {
		const reason = "CAN NOT FIND STRUCT TYPE"
		zaplog.LOG.Panic(reason, zap.String("struct_name", cfg.schemaX.structName))
		panic(reason)
	}
	done.Done(ast.Print(fileSet, astStructType))

	// 这里拿到的是要修改的操作，具体指改哪个文件哪个结构体的哪个字段，哪个位置的代码，以及新代码的内容
	sourceModifications := cfg.collectTagModifications(astStructType)

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

func (cfg *Config) collectTagModifications(structType *ast.StructType) []*defineTagModification {
	var results []*defineTagModification

	// 默认的样式配置
	defaultPattern := cfg.options.columnNamingStrategies.GetDefault()

	// 遍历结构体字段
	for _, fieldItem := range structType.Fields.List {
		// 打印字段名称和类型
		for _, fieldName := range fieldItem.Names {
			zaplog.LOG.Debug("--")
			zaplog.LOG.Debug("process", zap.String("struct_field_name:", fieldName.Name))
			zaplog.LOG.Debug("--")

			schemaColumn, exist := cfg.schemaX.schColumns.Get(fieldName.Name)
			if !exist { //比如字段是 "V哈哈" 就没事 而假如是 "v哈哈" 或者 "哈哈" 就不行，因为非以大写字母开始的字段，就没有gorm的列名
				zaplog.LOG.Debug("NO SCHEMA_FIELD - MAYBE NAME IS UNEXPORTED", zap.String("struct_name", cfg.schemaX.structName), zap.String("struct_field_name", fieldName.Name))
				continue
			}

			if fieldItem.Tag == nil {
				zaplog.LOG.Debug("NO TAG", zap.String("struct_name", cfg.schemaX.structName), zap.String("struct_field_name", fieldName.Name))

				// 假如配置跳过简单字段，而这个字段恰好是简单字段时，就跳过（因为没有标签，也就是没有配置规则）
				if cfg.options.skipBasicColumnName && defaultPattern.CheckColumnName(schemaColumn.DBName) {
					zaplog.LOG.Debug("SKIP SIMPLE FIELD", zap.String("struct_name", cfg.schemaX.structName), zap.String("struct_field_name", fieldName.Name))
					continue
				}

				// 假如没有标签，也就是没配置规则，但假如字段不符合默认规则，就得重新创建字段
				if !defaultPattern.CheckColumnName(schemaColumn.DBName) {
					if len(fieldItem.Names) >= 2 { //比如 a,b int 这种两个字段在一起，但其中一个字段的列名不正确时，就没法自动解决啦（其实有办法但不想实现，因为代价较大而没有收益）
						const reason = "CAN NOT HANDLE THIS SITUATION"
						zaplog.LOG.Panic(reason, zap.String("struct_name", cfg.schemaX.structName), zap.String("struct_field_name", fieldName.Name))
						panic(reason) //这种情况下当有错时，就不处理这种情况，就需要程序员先把两个字段定义到两行里
					}
					// 需要修改标签内容
					changeTag := cfg.modifyFieldTagCorrection(schemaColumn, "``", defaultPattern.GetPatternEnum()) //这里应该走创建标签的逻辑，但和修改标签的逻辑是相同的
					// 收集标签修改操作
					results = append(results, &defineTagModification{
						structFieldName: fieldName.Name,
						previousTagNode: syntaxgo_astnode.NewNode(fieldItem.End(), fieldItem.End()), //在尾部插入新的标签，要紧贴字段而且在换行符前面
						modifiedTagCode: changeTag,
					})
				}
			} else if patternName := cfg.extractTagGetCnmPattern(fieldItem.Tag.Value); patternName != "" {
				patternEnum := gormmomname.PatternEnum(patternName)
				zaplog.LOG.Debug("process", zap.String("column_name_pattern", string(patternEnum)))
				// 需要修改标签内容
				changeTag := cfg.modifyFieldTagCorrection(schemaColumn, fieldItem.Tag.Value, patternEnum)
				// 收集标签修改操作
				results = append(results, &defineTagModification{
					structFieldName: fieldName.Name,
					previousTagNode: fieldItem.Tag, //完整替换原来的标签
					modifiedTagCode: changeTag,
				})
			} else {
				if cfg.options.skipBasicColumnName && defaultPattern.CheckColumnName(schemaColumn.DBName) {
					zaplog.LOG.Debug("SKIP SIMPLE FIELD", zap.String("struct_name", cfg.schemaX.structName), zap.String("struct_field_name", fieldName.Name))
					continue
				}

				if !defaultPattern.CheckColumnName(schemaColumn.DBName) { //按照比较宽泛的规则也校验不过的时候就需要修正字段名
					// 需要修改标签内容
					changeTag := cfg.modifyFieldTagCorrection(schemaColumn, fieldItem.Tag.Value, defaultPattern.GetPatternEnum())
					// 收集标签修改操作
					results = append(results, &defineTagModification{
						structFieldName: fieldName.Name,
						previousTagNode: fieldItem.Tag, //完整替换原来的标签
						modifiedTagCode: changeTag,
					})
				} else {
					zaplog.LOG.Debug("match-pattern-so-skip", zap.String("name", fieldName.Name), zap.String("tag", fieldItem.Tag.Value))
				}
			}
		}
	}

	zaplog.LOG.Debug("change_column_names")
	for _, rep := range results {
		zaplog.LOG.Debug("check_column:", zap.String("name", rep.structFieldName), zap.String("code", rep.modifiedTagCode))
	}
	return results
}

func (cfg *Config) extractTagGetCnmPattern(tagCode string) string {
	return cfg.extractTagFieldGetValue(tagCode, cfg.options.tagName, cfg.options.columnNamingSubTagName)
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

func (cfg *Config) modifyFieldTagCorrection(schemaField *schema.Field, tag string, patternType gormmomname.PatternEnum) string {
	zaplog.LOG.Debug("new_fix_tag_code", zap.String("name", schemaField.Name), zap.String("tag", tag))
	//在 gorm 里修改 column 内容
	newTag := cfg.modifyGormTagWithColumn(schemaField, tag, patternType)
	//在 规则 里修改 column-name-pattern 内容
	newTag = cfg.modifyPatternTagWithName(newTag, patternType)
	//这是替换后的结果，即替换整个标签内容，获得新的完整标签内容
	zaplog.LOG.Debug("new_fix_tag_code", zap.String("name", schemaField.Name), zap.String("new_tag", newTag))
	return newTag
}

func (cfg *Config) modifyGormTagWithColumn(schemaField *schema.Field, tag string, patternType gormmomname.PatternEnum) string {
	pattern := cfg.options.columnNamingStrategies.GetPattern(patternType)
	columnName := pattern.BuildColumnName(schemaField.Name)
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

func (cfg *Config) modifyPatternTagWithName(tag string, patternType gormmomname.PatternEnum) string {
	zaplog.LOG.Debug("modify-pattern-tag-with-name", zap.String("column_name_pattern", string(patternType)))

	tagValue, sdx, edx := syntaxgo_tag.ExtractTagValueIndex(tag, cfg.options.tagName)
	if sdx < 0 || edx < 0 { //表示没找到 gorm 相关的内容
		if tagValue != "" {
			zaplog.LOG.Panic("IMPOSSIBLE")
		}
		part := fmt.Sprintf(`%s:"%s:%s;"`, cfg.options.tagName, cfg.options.columnNamingSubTagName, string(patternType))
		if rch := tag[len(tag)-2]; rch != ' ' && rch != '`' {
			part = " " + part //说明前面还有别的标签
		}
		p := len(tag) - 1 //插在最后一个"`"的前面，即最后一位的位置
		return tag[:p] + part + tag[p:]
	}
	//设置这个标签的这个字段的值
	return syntaxgo_tag.SetTagFieldValue(tag, cfg.options.tagName, cfg.options.columnNamingSubTagName, string(patternType), syntaxgo_tag.INSERT_LOCATION_TOP)
}
