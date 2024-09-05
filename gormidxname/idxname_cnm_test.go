package gormidxname

import "testing"

func Test_nameGenUseCnmImp_makeIndexName(t *testing.T) {
	t.Log(new(nameGenUseCnmImp).makeIndexName("idx", "tab", "age"))
}
