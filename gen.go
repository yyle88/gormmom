// Package gormmom: Native language programming engine that breaks down language barriers in database development
// As the smart tag generation engine of the GORM ecosystem, it empowers teams worldwide to write database models
// in native languages while auto generating database-compatible GORM tags and column names
// Supports Unicode-compatible field names in Chinese, Japanese, Korean, and additional languages
//
// gormmom: 原生语言编程引擎，打破数据库开发中的语言壁垒
// 作为 GORM 生态系统的智能标签生成引擎，它赋能全球团队使用原生语言编写数据库模型
// 同时自动生成数据库兼容的 GORM 标签和列名
// 支持 Unicode 兼容的中文、日语、韩语和其他语言字段名
package gormmom

import (
	"bytes"
	"fmt"
	"go/ast"
	"os"
	"slices"

	"github.com/yyle88/done"
	"github.com/yyle88/formatgo"
	"github.com/yyle88/gormmom/gormidxname"
	"github.com/yyle88/gormmom/gormmomname"
	"github.com/yyle88/must"
	"github.com/yyle88/rese"
	"github.com/yyle88/syntaxgo/syntaxgo_ast"
	"github.com/yyle88/syntaxgo/syntaxgo_astnode"
	"github.com/yyle88/syntaxgo/syntaxgo_search"
	"github.com/yyle88/syntaxgo/syntaxgo_tag"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
	"gorm.io/gorm/schema"
)

// Config represents the configuration for GORM tag generation operations
// Contains the target struct information and generation options with customizing output
// Provides smart mapping and deterministic generation of database-compatible tags
//
// Config 代表 GORM 标签生成操作的配置
// 包含目标结构体信息和生成选项，用于自定义输出
// 提供智能映射和确定性的数据库兼容标签生成
type Config struct {
	gormStruct *GormStruct // Target GORM struct to process // 要处理的目标 GORM 结构体
	options    *Options    // Generation options and settings // 生成选项和设置
}

// NewConfig creates a new configuration instance with GORM tag generation
// Takes the target struct and generation options, returns configured instance
// Used to initialize the generation workflow with custom settings
//
// NewConfig 创建新的 GORM 标签生成配置实例
// 接收目标结构体和生成选项，返回配置好的实例
// 用于使用自定义设置初始化生成工作流程
func NewConfig(gormStruct *GormStruct, options *Options) *Config {
	return &Config{
		gormStruct: gormStruct,
		options:    options,
	}
}

// GenReplace generates new GORM tags and replaces the original source file
// Processes the struct definition to add database-compatible tags and column names
// Formats the generated code and writes back to the source file when changes detected
// Returns the result containing the new code and change status
//
// GenReplace 生成新的 GORM 标签并替换原始源文件
// 处理结构体定义以添加数据库兼容的标签和列名
// 检测到变化时自动格式化生成的代码并写回源文件
// 返回包含新代码和变化状态的结果
func (cfg *Config) GenReplace() *NewCodeResult {
	newCode := cfg.GetNewCode()
	srcPath := must.SameNice(cfg.gormStruct.sourcePath, newCode.SrcPath)
	if newCode.HasChange() { // 只有当有变化时才写文件
		srcCode := must.Have(rese.A1(formatgo.FormatBytes(newCode.NewCode)))
		must.Done(os.WriteFile(srcPath, srcCode, 0644))
		newCode.NewCode = srcCode
	}
	return newCode
}

// GetNewCode generates new code with GORM tags without modifying the original file
// Reads the source file and processes it to add native language field mappings
// Returns the new code result with updated tags and column definitions
//
// GetNewCode 生成带有 GORM 标签的新代码而不修改原文件
// 读取源文件并处理以添加原生语言字段映射
// 返回包含更新标签和列定义的新代码结果
func (cfg *Config) GetNewCode() *NewCodeResult {
	return cfg.makeNewCode(rese.A1(os.ReadFile(cfg.gormStruct.sourcePath)))
}

func (cfg *Config) makeNewCode(sourceCode []byte) *NewCodeResult {
	astBundle := rese.C1(syntaxgo_ast.NewAstBundleV1(sourceCode))
	astFile, fileSet := astBundle.GetBundle()
	// 使用语法分析树 ast 找到结构体的代码，就像这样
	// struct {
	//    Name string
	// }
	// 结果只包含结构体的内容
	structType, ok := syntaxgo_search.FindStructTypeByName(astFile, cfg.gormStruct.structName)
	if !ok {
		const reason = "CAN NOT FIND STRUCT TYPE"
		zaplog.LOG.Panic(reason, zap.String("struct_name", cfg.gormStruct.structName))
		panic(reason)
	}
	done.Done(ast.Print(fileSet, structType))

	// 这里拿到的是要修改的操作，具体指改哪个文件哪个结构体的哪个字段，哪个位置的代码，以及新代码的内容
	modifications := cfg.collectTagModifications(structType)

	if cfg.options.renewIndexName {
		//这里增加个新逻辑，就是单列索引的索引名称不正确，需要也校正索引名，因此这个函数会补充标签内容
		cfg.correctIndexNames(modifications)

		zaplog.LOG.Debug("change_index_names")
		for _, step := range modifications {
			zaplog.LOG.Debug("check_index:", zap.String("struct_filed_name", step.structFieldName), zap.String("new_tag_code", step.newTagCode))
		}
	}

	//需要翻转下从后往前替换，因为替换以后源码会变，假如从前往后替换坐标就对不上啦，而从后往前替换则不存在这个问题
	slices.Reverse(modifications)

	//接下来替换代码，把需要 新增 或者 替换 的标签都设置到代码里
	newCode := sourceCode
	changedLineCount := 0
	for _, step := range modifications {
		oldTagCode := syntaxgo_astnode.GetCode(newCode, step.tagPosNode)
		newTagCode := []byte(step.newTagCode)
		if !bytes.Equal(oldTagCode, newTagCode) {
			zaplog.LOG.Debug("change", zap.ByteString("old_tag_code", oldTagCode), zap.ByteString("new_tag_code", newTagCode))
		} else {
			continue
		}
		newCode = syntaxgo_astnode.ChangeNodeCode(newCode, step.tagPosNode, newTagCode)
		changedLineCount++
	}
	return &NewCodeResult{
		NewCode:          newCode,
		SrcPath:          cfg.gormStruct.sourcePath,
		ChangedLineCount: changedLineCount,
	}
}

type defineTagModification struct {
	structFieldName string   //结构体的字段名
	columnName      string   //数据表中的列名
	tagPosNode      ast.Node //标签的起止位置-就是在src源码中的位置，便于后面的替换代码
	newTagCode      string   //新标签的新内容-就是标签的完整全部内容
}

// collectTagModifications analyzes struct fields and collects necessary tag modifications
// Processes each field to determine if native language column naming is needed
// Returns a collection of modifications required for proper GORM tag generation
//
// collectTagModifications 分析结构体字段并收集必要的标签修改
// 处理每个字段以确定是否需要原生语言列名
// 返回正确 GORM 标签生成所需的修改集合
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

			schemaColumn, exist := cfg.gormStruct.gormFields.Get(fieldName.Name)
			if !exist { //比如字段是 "V哈哈" 就没事 而假如是 "v哈哈" 或者 "哈哈" 就不行，因为非以大写字母开始的字段，就没有gorm的列名
				zaplog.LOG.Debug("NO SCHEMA_FIELD - MAYBE NAME IS UNEXPORTED", zap.String("struct_name", cfg.gormStruct.structName), zap.String("struct_field_name", fieldName.Name))
				continue
			}

			if fieldItem.Tag == nil {
				zaplog.LOG.Debug("NO TAG", zap.String("struct_name", cfg.gormStruct.structName), zap.String("struct_field_name", fieldName.Name))

				// 假如配置跳过简单字段，而这个字段恰好是简单字段时，就跳过（因为没有标签，也就是没有配置规则）
				if cfg.options.skipBasicColumnName && defaultPattern.CheckColumnName(schemaColumn.DBName) {
					zaplog.LOG.Debug("SKIP SIMPLE FIELD", zap.String("struct_name", cfg.gormStruct.structName), zap.String("struct_field_name", fieldName.Name))
					continue
				}

				// 假如没有标签，也就是没配置规则，但假如字段不符合默认规则，就得重新创建字段
				if !defaultPattern.CheckColumnName(schemaColumn.DBName) {
					if len(fieldItem.Names) >= 2 { //比如 a,b int 这种两个字段在一起，但其中一个字段的列名不正确时，就没法自动解决啦（其实有办法但不想实现，因为代价较大而没有收益）
						const reason = "CAN NOT HANDLE THIS SITUATION"
						zaplog.LOG.Panic(reason, zap.String("struct_name", cfg.gormStruct.structName), zap.String("struct_field_name", fieldName.Name))
						panic(reason) //这种情况下当有错时，就不处理这种情况，就需要程序员先把两个字段定义到两行里
					}
					// 需要修改标签内容
					changeTag := cfg.modifyFieldTagCorrection(schemaColumn, "``", defaultPattern.GetPatternEnum()) //这里应该走创建标签的逻辑，但和修改标签的逻辑是相同的
					// 收集标签修改操作
					results = append(results, &defineTagModification{
						structFieldName: fieldName.Name,
						columnName:      changeTag.columnName,
						tagPosNode:      syntaxgo_astnode.NewNode(fieldItem.End(), fieldItem.End()), //在尾部插入新的标签，要紧贴字段而且在换行符前面
						newTagCode:      changeTag.newTagCode,
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
					columnName:      changeTag.columnName,
					tagPosNode:      fieldItem.Tag, //完整替换原来的标签
					newTagCode:      changeTag.newTagCode,
				})
			} else {
				if cfg.options.skipBasicColumnName && defaultPattern.CheckColumnName(schemaColumn.DBName) {
					zaplog.LOG.Debug("SKIP SIMPLE FIELD", zap.String("struct_name", cfg.gormStruct.structName), zap.String("struct_field_name", fieldName.Name))
					if cfg.options.renewIndexName && cfg.hasAnyIdxTagUdxTagValue(fieldItem) {
						results = append(results, cfg.newFirstNotModification(fieldItem, schemaColumn.DBName))
					}
					continue
				}

				if !defaultPattern.CheckColumnName(schemaColumn.DBName) { //按照比较宽泛的规则也校验不过的时候就需要修正字段名
					// 需要修改标签内容
					changeTag := cfg.modifyFieldTagCorrection(schemaColumn, fieldItem.Tag.Value, defaultPattern.GetPatternEnum())
					// 收集标签修改操作
					results = append(results, &defineTagModification{
						structFieldName: fieldName.Name,
						columnName:      changeTag.columnName,
						tagPosNode:      fieldItem.Tag, //完整替换原来的标签
						newTagCode:      changeTag.newTagCode,
					})
				} else {
					zaplog.LOG.Debug("match-pattern-so-skip", zap.String("name", fieldName.Name), zap.String("tag", fieldItem.Tag.Value))
					if cfg.options.renewIndexName && cfg.hasAnyIdxTagUdxTagValue(fieldItem) {
						results = append(results, cfg.newFirstNotModification(fieldItem, schemaColumn.DBName))
					}
					continue
				}
			}
		}
	}

	zaplog.LOG.Debug("change_column_names")
	for _, rep := range results {
		zaplog.LOG.Debug("check_column:", zap.String("field_name", rep.structFieldName), zap.String("new_column_name", rep.columnName), zap.String("new_tag_code", rep.newTagCode))

		must.Nice(rep.structFieldName)
		must.Nice(rep.columnName)
		must.Nice(rep.newTagCode)
		must.Nice(rep.tagPosNode.Pos())
		must.Nice(rep.tagPosNode.End())
	}
	return results
}

// hasAnyIdxTagUdxTagValue checks if field has any index pattern tag values
// Examines field tags for idx or udx pattern configurations
// Returns true if any index pattern tag is found in the field
//
// hasAnyIdxTagUdxTagValue 检查字段是否有任何索引模式标签值
// 检查字段标签中的 idx 或 udx 模式配置
// 如果在字段中找到任何索引模式标签则返回 true
func (cfg *Config) hasAnyIdxTagUdxTagValue(fieldItem *ast.Field) bool {
	if len(fieldItem.Names) != 1 {
		return false
	}
	if fieldItem.Tag == nil {
		return false
	}
	for _, patternTagEnum := range []gormidxname.IndexPatternTagEnum{
		gormidxname.IdxPatternTagName,
		gormidxname.UdxPatternTagName,
	} {
		var name = cfg.extractTagFieldGetValue(fieldItem.Tag.Value, cfg.options.systemTagName, string(patternTagEnum))
		if name != "" {
			zaplog.LOG.Debug("match-pattern-so-has-any-tag", zap.String("pattern", string(patternTagEnum)), zap.String("name", name))
			return true
		}
	}
	return false
}

// hasOneIdxTagUdxTagValue checks if tag code contains specific index pattern
// Validates the presence of a particular index pattern in the tag string
// Returns true if the specified pattern tag is found in the code
//
// hasOneIdxTagUdxTagValue 检查标签代码是否包含特定的索引模式
// 验证标签字符串中是否存在特定的索引模式
// 如果在代码中找到指定的模式标签则返回 true
func (cfg *Config) hasOneIdxTagUdxTagValue(tagCode string, patternTagEnum gormidxname.IndexPatternTagEnum) bool {
	var name = cfg.extractTagFieldGetValue(tagCode, cfg.options.systemTagName, string(patternTagEnum))
	if name != "" {
		zaplog.LOG.Debug("match-pattern-so-has-one-tag", zap.String("pattern", string(patternTagEnum)), zap.String("name", name))
		return true
	}
	return false
}

// newFirstNotModification creates modification for fields without existing tags
// Generates initial tag modification structure for untagged fields
// Used when field requires native language column naming but has no current tags
//
// newFirstNotModification 为没有现有标签的字段创建修改
// 为未标记的字段生成初始标签修改结构
// 在字段需要原生语言列名但没有当前标签时使用
func (cfg *Config) newFirstNotModification(fieldItem *ast.Field, columnName string) *defineTagModification {
	must.Full(fieldItem)
	must.Length(fieldItem.Names, 1)
	must.Nice(fieldItem.Names[0].Name)

	return &defineTagModification{
		structFieldName: fieldItem.Names[0].Name,
		columnName:      columnName,
		tagPosNode:      fieldItem.Tag,
		newTagCode:      fieldItem.Tag.Value,
	}
}

// extractTagGetCnmPattern extracts column naming pattern from tag code
// Retrieves the column naming pattern configuration from system tag
// Returns pattern string for column name generation validation
//
// extractTagGetCnmPattern 从标签代码中提取列名模式
// 从系统标签中检索列名模式配置
// 返回用于列名生成验证的模式字符串
func (cfg *Config) extractTagGetCnmPattern(tagCode string) string {
	return cfg.extractTagFieldGetValue(tagCode, cfg.options.systemTagName, cfg.options.columnNamingSubTagName)
}

// extractTagFieldGetValue extracts nested tag field value using dual key lookup
// Performs two-level extraction from tag code using primary and secondary keys
// Returns the extracted field value or empty string if not found
//
// extractTagFieldGetValue 使用双键查找提取嵌套标签字段值
// 使用主键和辅助键从标签代码进行两级提取
// 返回提取的字段值，如果找不到则返回空字符串
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

// correctionNewTag contains the result of tag correction operations
// Holds the corrected column name and updated tag code for field processing
// Used as intermediate result during tag modification workflow
//
// correctionNewTag 包含标签纠正操作的结果
// 保存纠正的列名和更新的标签代码，用于字段处理
// 在标签修改工作流中作为中间结果使用
type correctionNewTag struct {
	columnName string // Corrected database column name // 纠正的数据库列名
	newTagCode string // Updated tag code with corrections // 带有纠正的更新标签代码
}

// modifyFieldTagCorrection applies corrections to field tag based on schema and pattern
// Processes GORM field with specified pattern to generate corrected column name and tag
// Returns corrected tag structure with updated column name and tag code
//
// modifyFieldTagCorrection 基于模式和模式对字段标签应用纠正
// 使用指定模式处理 GORM 字段，生成纠正的列名和标签
// 返回带有更新列名和标签代码的纠正标签结构
func (cfg *Config) modifyFieldTagCorrection(schemaField *schema.Field, tag string, patternType gormmomname.PatternEnum) *correctionNewTag {
	zaplog.LOG.Debug("new_fix_tag_code", zap.String("name", schemaField.Name), zap.String("tag", tag))
	//在 gorm 里修改 column 内容
	newTag := cfg.modifyGormTagWithColumn(schemaField, tag, patternType)
	//在 规则 里修改 column-name-pattern 内容
	newTag.newTagCode = cfg.modifyPatternTagWithName(newTag.newTagCode, patternType)
	//这是替换后的结果，即替换整个标签内容，获得新的完整标签内容
	zaplog.LOG.Debug("new_fix_tag_code", zap.String("name", schemaField.Name), zap.String("column_name", newTag.columnName), zap.String("new_tag_code", newTag.newTagCode))
	return newTag
}

// modifyGormTagWithColumn modifies GORM tag to include appropriate column name
// Generates database-compatible column name using specified pattern and updates GORM tag
// Returns corrected tag with correct column specification for database mapping
//
// modifyGormTagWithColumn 修改 GORM 标签以包含适当的列名
// 使用指定模式生成数据库兼容的列名并更新 GORM 标签
// 返回带有正确列规范的纠正标签，用于数据库映射
func (cfg *Config) modifyGormTagWithColumn(schemaField *schema.Field, tag string, patternType gormmomname.PatternEnum) *correctionNewTag {
	pattern := cfg.options.columnNamingStrategies.GetPattern(patternType)
	columnName := pattern.BuildColumnName(schemaField.Name)
	zaplog.LOG.Debug("new_fix_gorm_tag", zap.String("name", schemaField.Name), zap.String("column_name", columnName))
	must.Nice(columnName)

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
		return &correctionNewTag{
			columnName: columnName,
			newTagCode: tag[:p] + part + tag[p:],
		}
	}
	//设置这个标签的这个字段的值
	return &correctionNewTag{
		columnName: columnName,
		newTagCode: syntaxgo_tag.SetTagFieldValue(tag, "gorm", "column", columnName, syntaxgo_tag.INSERT_LOCATION_TOP),
	}
}

// modifyPatternTagWithName adds or updates pattern tag with specified naming pattern
// Inserts or modifies system tag to include column naming pattern specification
// Returns updated tag string with pattern information for consistent processing
//
// modifyPatternTagWithName 使用指定的命名模式添加或更新模式标签
// 插入或修改系统标签以包含列命名模式规范
// 返回带有模式信息的更新标签字符串，用于一致性处理
func (cfg *Config) modifyPatternTagWithName(tag string, patternType gormmomname.PatternEnum) string {
	zaplog.LOG.Debug("modify-pattern-tag-with-name", zap.String("column_name_pattern", string(patternType)))

	tagValue, sdx, edx := syntaxgo_tag.ExtractTagValueIndex(tag, cfg.options.systemTagName)
	if sdx < 0 || edx < 0 { //表示没找到 gorm 相关的内容
		if tagValue != "" {
			zaplog.LOG.Panic("IMPOSSIBLE")
		}
		part := fmt.Sprintf(`%s:"%s:%s;"`, cfg.options.systemTagName, cfg.options.columnNamingSubTagName, string(patternType))
		if rch := tag[len(tag)-2]; rch != ' ' && rch != '`' {
			part = " " + part //说明前面还有别的标签
		}
		p := len(tag) - 1 //插在最后一个"`"的前面，即最后一位的位置
		return tag[:p] + part + tag[p:]
	}
	//设置这个标签的这个字段的值
	return syntaxgo_tag.SetTagFieldValue(tag, cfg.options.systemTagName, cfg.options.columnNamingSubTagName, string(patternType), syntaxgo_tag.INSERT_LOCATION_TOP)
}
