package gormmomname

import "github.com/yyle88/gormmom/internal/utils"

// nolint:no-doc
// 自定义枚举类型，表示使用何种字段验证方式来验证，由于不同的DB的列名规则是不同的，因此通常建议是取各种DB的交集
type ColumnNamePattern string

const (
	Lowercase30 ColumnNamePattern = "s30"
	Uppercase30 ColumnNamePattern = "S30"
	Lowercase63 ColumnNamePattern = "s63"
	Uppercase63 ColumnNamePattern = "S63"

	DefaultPattern ColumnNamePattern = Lowercase63
)

type Naming interface {
	IsValidColumnName(columnName string) bool
	GenerateColumnName(fieldName string) string
}

var namingStrategies = map[ColumnNamePattern]Naming{
	Lowercase30: &lowercase30pattern{},
	Uppercase30: &uppercase30pattern{},
	Lowercase63: &lowercase63pattern{},
	Uppercase63: &uppercase63pattern{},
}

func GetNamingStrategies() map[ColumnNamePattern]Naming {
	return utils.CloneMap(namingStrategies)
}
