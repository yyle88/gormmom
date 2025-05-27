package gormidxname

import (
	"regexp"
	"strings"

	"github.com/yyle88/gormmom/internal/simpleindexname"
	"gorm.io/gorm/schema"
)

type Lowercase63pattern struct{}

func NewLowercase63pattern() *Lowercase63pattern {
	return &Lowercase63pattern{}
}

func (G *Lowercase63pattern) GetPatternEnum() PatternEnum {
	return "cnm" //表示使用 column name 作为拼接索引名的规则，但后缀是小写字母的
}

func (G *Lowercase63pattern) CheckIndexName(indexName string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_]{1,63}$`).MatchString(indexName)
}

func (G *Lowercase63pattern) BuildIndexName(schemaIndex *schema.Index, param *BuildIndexParam) *IndexNameResult {
	result := simpleindexname.BuildIndexName(schemaIndex, &simpleindexname.BuildIndexParam{
		TableName:  param.TableName,
		FieldName:  param.FieldName,
		ColumnName: strings.ToLower(param.ColumnName),
	})
	simpleindexname.CheckLength(result.NewIndexName, 63)
	return &IndexNameResult{
		TagFieldName: result.TagFieldName,
		NewIndexName: result.NewIndexName,
		IdxUdxPrefix: result.IdxUdxPrefix,
	}
}
