package gormidxname

import (
	"regexp"
	"strings"

	"github.com/yyle88/gormmom/internal/simpleindexname"
	"gorm.io/gorm/schema"
)

type Uppercase63pattern struct{}

func NewUppercase63pattern() *Uppercase63pattern {
	return &Uppercase63pattern{}
}

func (G *Uppercase63pattern) GetPatternEnum() PatternEnum {
	return "CNM" //表示使用 column name 作为拼接索引名的规则，但后缀是大写字母的
}

func (G *Uppercase63pattern) CheckIndexName(indexName string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_]{1,63}$`).MatchString(indexName)
}

func (G *Uppercase63pattern) BuildIndexName(schemaIndex *schema.Index, param *BuildIndexParam) *IndexNameResult {
	result := simpleindexname.BuildIndexName(schemaIndex, &simpleindexname.BuildIndexParam{
		TableName:  param.TableName,
		FieldName:  param.FieldName,
		ColumnName: strings.ToUpper(param.ColumnName),
	})
	simpleindexname.CheckLength(result.NewIndexName, 63)
	return &IndexNameResult{
		TagFieldName: result.TagFieldName,
		NewIndexName: result.NewIndexName,
		IdxUdxPrefix: IndexPatternTagEnum(result.IdxUdxPrefix),
	}
}
