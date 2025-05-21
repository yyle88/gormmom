package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/runpath"
)

func TestListGoFiles(t *testing.T) {
	paths := ListGoFiles(runpath.PARENT.Path())
	t.Log(neatjsons.S(paths))
	require.NotEmpty(t, paths)
}
