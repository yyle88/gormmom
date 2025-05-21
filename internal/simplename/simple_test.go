package simplename

import "testing"

func TestBuildSimpleName(t *testing.T) {
	t.Log(BuildSimpleName("abc"))
	t.Log(BuildSimpleName("ABC"))
	t.Log(BuildSimpleName("AaBbCc"))
	t.Log(BuildSimpleName("FieldNameExample世界"))
	t.Log(BuildSimpleName("世界ABC"))
	t.Log(BuildSimpleName("姓名"))
	t.Log(BuildSimpleName("V性别"))
	t.Log(BuildSimpleName("(特殊)（情况）"))
	t.Log(BuildSimpleName("abc*xyz"))
}
