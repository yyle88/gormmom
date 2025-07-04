package gormmomname

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLowercase63pattern_CheckColumnName(t *testing.T) {
	pattern := NewLowercase63pattern()

	require.True(t, pattern.CheckColumnName("abc"))
}

func TestUppercase63pattern_CheckColumnName(t *testing.T) {
	pattern := NewUppercase63pattern()

	require.True(t, pattern.CheckColumnName("ABC"))
}

func TestLowercase30pattern_BuildColumnName(t *testing.T) {
	pattern := NewLowercase30pattern()

	t.Log(pattern.BuildColumnName("v杨亦乐"))
	t.Log(pattern.BuildColumnName("v刘亦菲"))
	t.Log(pattern.BuildColumnName("v古天乐"))
}

func TestUppercase63pattern_BuildColumnName(t *testing.T) {
	pattern := NewUppercase63pattern()

	t.Log(pattern.BuildColumnName("v杨亦乐"))
	t.Log(pattern.BuildColumnName("v刘亦菲的亦"))
	t.Log(pattern.BuildColumnName("v古天乐的乐"))
}
