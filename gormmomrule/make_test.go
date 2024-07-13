package gormmomrule

import "testing"

func TestRule_MakeName(t *testing.T) {
	t.Log(S30.MakeName("v杨亦乐"))
	t.Log(S30.MakeName("v刘亦菲"))
	t.Log(S30.MakeName("v古天乐"))
}

func TestMakeName(t *testing.T) {
	t.Log(MakeName(S63, "啦啦啦", map[Rule]func(string) string{}))
}
