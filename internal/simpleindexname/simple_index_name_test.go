package simpleindexname

import "testing"

func TestName_mergeIndexName(t *testing.T) {
	t.Log(mergeIndexName("idx", "student", "age"))
}
