package gormidxname

import (
	"strings"

	"github.com/yyle88/gormmom/internal/simpleindexname"
	"github.com/yyle88/gormmom/internal/utils"
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
	return utils.NewCommonRegexp(63).MatchString(indexName)
}

func (G *Lowercase63pattern) BuildIndexName(schemaIndex *schema.Index, param *BuildIndexParam) *IndexNameResult {
	result := simpleindexname.BuildIndexName(schemaIndex, &simpleindexname.BuildIndexParam{
		TableName:  param.TableName,
		FieldName:  param.FieldName,
		ColumnName: strings.ToLower(param.ColumnName),
	})
	utils.MustMatchRegexp(utils.NewCommonRegexp(63), result.NewIndexName)
	return &IndexNameResult{
		TagFieldName: result.TagFieldName,
		NewIndexName: result.NewIndexName,
		IdxUdxPrefix: IndexPatternTagEnum(result.IdxUdxPrefix),
	}
}
