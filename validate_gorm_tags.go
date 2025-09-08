package gormmom

import (
	"fmt"

	"github.com/yyle88/erero"
	"github.com/yyle88/gormmom/gormidxname"
	"github.com/yyle88/gormmom/gormmomname"
	"github.com/yyle88/gormmom/internal/utils"
	"github.com/yyle88/tern/zerotern"
	"gorm.io/gorm/schema"
)

// ValidateGormTags validates GORM tags to ensure naming conventions
// Checks to find non-ASCII characters in index names and tag issues
// Returns each detected issue as an error
//
// ValidateGormTags 验证 GORM 标签的命名规范
// 检查索引名中的非 ASCII 字符和其他标签问题
// 如果发现问题标签则返回错误
func (configs Configs) ValidateGormTags() error {
	for _, config := range configs {
		if err := config.validateGormTags(); err != nil {
			return erero.Wro(err)
		}
	}
	return nil
}

func (cfg *Config) validateGormTags() error {
	// Use GORM schema fields for validation along with naming strategies
	// 使用 GORM 模式字段结合命名策略进行验证
	structName := cfg.gormStruct.structName
	gormSchema := cfg.gormStruct.gormSchema

	// Check schema table name
	// 检查模式表名
	utils.ValidateTableName(gormSchema.Table, structName)

	// Validate column names in each field
	// 验证每个字段的列名
	for _, field := range gormSchema.Fields {
		if err := cfg.validateColumnNaming(structName, field); err != nil {
			return erero.Wro(err)
		}
	}

	// Validate index names (parse indexes once)
	// 验证索引名（只解析一次索引）
	if err := cfg.validateAllIndexNames(structName, gormSchema); err != nil {
		return erero.Wro(err)
	}

	return nil
}

// validateColumnNaming checks column names using columnNamingStrategies
// validateColumnNaming 使用 columnNamingStrategies 检查列名
func (cfg *Config) validateColumnNaming(structName string, field *schema.Field) error {
	// Extract pattern from mom tag, with default fallback
	// 从 mom 标签中提取模式，有默认值作为后备
	patternName := zerotern.VF(cfg.extractTagGetCnmPattern(fmt.Sprintf("`%s`", field.Tag)), func() string {
		return string(cfg.options.columnNamingStrategies.GetDefault().GetPatternEnum())
	})

	// Get the pattern validation instance
	// 获取模式验证实例
	pattern := cfg.options.columnNamingStrategies.GetPattern(gormmomname.PatternEnum(patternName))

	// Use pattern's CheckColumnName method for validation
	// 使用模式的 CheckColumnName 方法进行验证
	if !pattern.CheckColumnName(field.DBName) {
		return erero.Errorf("field %s.%s column '%s' not match pattern '%s'", structName, field.Name, field.DBName, patternName)
	}
	return nil
}

// validateAllIndexNames checks all index names using indexNamingStrategies
// Iterates through each index and validates naming against configured patterns
// Prevents corrupted index names like "idx_products_p名���" from being generated
//
// validateAllIndexNames 使用 indexNamingStrategies 检查所有索引名
// 遍历每个索引并根据配置的模式验证命名
// 防止生成损坏的索引名如 "idx_products_p名���"
func (cfg *Config) validateAllIndexNames(structName string, gormSchema *schema.Schema) error {
	// Parse indexes from the schema once
	// 一次性从模式中解析索引
	indexes := gormSchema.ParseIndexes()

	for _, index := range indexes {
		for _, indexOption := range index.Fields {
			patternName := zerotern.VF(cfg.extractFieldIndexPattern(indexOption.Field, index), func() string {
				return string(cfg.options.indexNamingStrategies.GetDefault().GetPatternEnum())
			})

			// Get the pattern validation instance
			// 获取模式验证实例
			pattern := cfg.options.indexNamingStrategies.GetPattern(gormidxname.PatternEnum(patternName))

			// Use pattern's CheckIndexName method for validation
			// 使用模式的 CheckIndexName 方法进行验证
			if !pattern.CheckIndexName(index.Name) {
				return erero.Errorf("struct %s index '%s' not match pattern '%s'", structName, index.Name, patternName)
			}
		}
	}

	return nil
}

// extractFieldIndexPattern extracts index pattern from field's mom tag based on index type
// extractFieldIndexPattern 根据索引类型从字段的 mom 标签中提取索引模式
func (cfg *Config) extractFieldIndexPattern(field *schema.Field, index *schema.Index) string {
	// Get the struct field tag to extract pattern information
	// 获取结构体字段标签以提取模式信息
	tagCode := fmt.Sprintf("`%s`", field.Tag)

	// Determine index pattern tag based on index type
	// 根据索引类型确定索引模式标签
	var patternTagEnum gormidxname.IndexPatternTagEnum
	if index.Class == "UNIQUE" {
		patternTagEnum = gormidxname.UdxPatternTagName
	} else {
		patternTagEnum = gormidxname.IdxPatternTagName
	}

	return cfg.extractTagFieldGetValue(tagCode, cfg.options.systemTagName, string(patternTagEnum))
}
