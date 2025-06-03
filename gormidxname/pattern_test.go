package gormidxname

import (
	"testing"

	"github.com/yyle88/gormmom/internal/simpleindexname"
	"github.com/yyle88/must"
)

func TestPattern(t *testing.T) {
	must.Same(simpleindexname.IdxPatternTagName, IdxPatternTagName)
	must.Same(simpleindexname.UdxPatternTagName, UdxPatternTagName)
}
