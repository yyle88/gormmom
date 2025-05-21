package gormmomname

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestName_CheckColumnName(t *testing.T) {
	pattern := NewLowercase63pattern()

	require.True(t, pattern.CheckColumnName("abc"))
}

func TestName_CheckColumnName_2(t *testing.T) {
	pattern := NewUppercase63pattern()

	require.True(t, pattern.CheckColumnName("ABC"))
}

func TestName_BuildColumnName(t *testing.T) {
	pattern := NewLowercase30pattern()

	t.Log(pattern.BuildColumnName("v杨亦乐"))
	t.Log(pattern.BuildColumnName("v刘亦菲"))
	t.Log(pattern.BuildColumnName("v古天乐"))
}

func TestName_BuildColumnName_2(t *testing.T) {
	pattern := NewUppercase63pattern()

	t.Log(pattern.BuildColumnName("v杨亦乐"))
	t.Log(pattern.BuildColumnName("v刘亦菲的亦"))
	t.Log(pattern.BuildColumnName("v古天乐的乐"))
}
