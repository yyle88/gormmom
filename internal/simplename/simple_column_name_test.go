package simplename

import "testing"

func TestBuildColumnName(t *testing.T) {
	t.Log(BuildColumnName("abc"))
	t.Log(BuildColumnName("ABC"))
	t.Log(BuildColumnName("AaBbCc"))
	t.Log(BuildColumnName("FieldNameExample世界"))
	t.Log(BuildColumnName("世界ABC"))
	t.Log(BuildColumnName("姓名"))
	t.Log(BuildColumnName("V性别"))
	t.Log(BuildColumnName("(特殊)（情况）"))
	t.Log(BuildColumnName("abc*xyz"))
}
