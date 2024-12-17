package gormmomname

import "testing"

func Test_simpleName(t *testing.T) {
	t.Log(simpleName("abc"))
	t.Log(simpleName("ABC"))
	t.Log(simpleName("AaBbCc"))
	t.Log(simpleName("FieldNameExample世界"))
	t.Log(simpleName("世界ABC"))
	t.Log(simpleName("姓名"))
	t.Log(simpleName("V性别"))
	t.Log(simpleName("(特殊)（情况）"))
	t.Log(simpleName("abc*xyz"))
}
