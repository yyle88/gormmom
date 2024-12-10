package gormmomname

import "testing"

func Test_simpleCreateColumnName(t *testing.T) {
	t.Log(simpleCreateColumnName("abc"))
	t.Log(simpleCreateColumnName("ABC"))
	t.Log(simpleCreateColumnName("AaBbCc"))
	t.Log(simpleCreateColumnName("FieldNameExample世界"))
	t.Log(simpleCreateColumnName("世界ABC"))
	t.Log(simpleCreateColumnName("姓名"))
	t.Log(simpleCreateColumnName("V性别"))
	t.Log(simpleCreateColumnName("(特殊)（情况）"))
	t.Log(simpleCreateColumnName("abc*xyz"))
}
