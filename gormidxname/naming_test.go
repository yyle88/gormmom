package gormidxname

import "testing"

func Test_withTableCnmPattern_combineIndexName(t *testing.T) {
	t.Log(new(withTableCnmPattern).combineIndexName("idx", "student", "age"))
}
