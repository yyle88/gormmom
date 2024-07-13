package gormmomrule

import "testing"

func Test_makeName(t *testing.T) {
	t.Log(makeName("abc"))
	t.Log(makeName("ABC"))
	t.Log(makeName("AaBbCc"))
	t.Log(makeName("FieldNameExample世界"))
	t.Log(makeName("世界ABC"))
	t.Log(makeName("姓名"))
	t.Log(makeName("V性别"))
	t.Log(makeName("(特殊)（情况）"))
	t.Log(makeName("abc*xyz"))
}
