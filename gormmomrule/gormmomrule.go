package gormmomrule

// nolint:no-doc
// 自定义枚举类型，表示使用何种字段验证方式来验证，由于不同的DB的列名规则是不同的，因此通常建议是取各种DB的交集
type MomRULE string

const (
	S30  MomRULE = "S30"
	S30U MomRULE = "S30U"
	S63  MomRULE = "S63"
	S63U MomRULE = "S63U"

	DEFAULT MomRULE = S63
)

type CnmMakeIFace interface {
	CheckName(columnName string) bool
	GenNewCnm(fieldName string) string
}

var presetNameImpMap = map[MomRULE]CnmMakeIFace{
	S30:  &nameS30Imp{},
	S30U: &nameS30UImp{},
	S63:  &nameS63Imp{},
	S63U: &nameS63UImp{},
}

func GetPresetCnmMakeMap() map[MomRULE]CnmMakeIFace {
	var mp = make(map[MomRULE]CnmMakeIFace, len(presetNameImpMap))
	for k, v := range presetNameImpMap {
		mp[k] = v
	}
	return mp
}
