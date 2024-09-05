package gormidxname

import (
	"gorm.io/gorm/schema"
)

// nolint:no-doc
// 自定义枚举类型
type IdxNAME string

const (
	CNM IdxNAME = "CNM" //表示使用 column name 作为拼接索引名的规则
	CNU IdxNAME = "CNU" //表示使用 column name 作为拼接索引名的规则，但后缀是大写字母的

	DEFAULT IdxNAME = CNM
)

type IdxNameIFace interface {
	CheckIdxName(indexName string) bool
	GenIndexName(schemaIndex schema.Index, tableName string, fieldName string, columnName string) *IdxGenResType
}

var presetNameImpMap = map[IdxNAME]IdxNameIFace{
	CNM: &nameGenUseCnmImp{},
	CNU: &nameGenUseCnuImp{}, //只是凑数用的
}

func GetPresetNameImpMap() map[IdxNAME]IdxNameIFace {
	var mp = make(map[IdxNAME]IdxNameIFace, len(presetNameImpMap))
	for k, v := range presetNameImpMap {
		mp[k] = v
	}
	return mp
}

type IdxGenResType struct {
	TagFieldName string
	NewIndexName string
	EnumCodeName string
}
