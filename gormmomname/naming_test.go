package gormmomname

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestName_IsValidNamePattern(t *testing.T) {
	require.True(t, namingStrategies[Lowercase63].IsValidColumnName("abc"))
}

func TestName_Check_IsValidNamePattern(t *testing.T) {
	require.True(t, namingStrategies[Uppercase63].IsValidColumnName("ABC"))
}

func TestName_GenerateColumnName(t *testing.T) {
	t.Log(namingStrategies[Lowercase30].GenerateColumnName("v杨亦乐"))
	t.Log(namingStrategies[Lowercase30].GenerateColumnName("v刘亦菲"))
	t.Log(namingStrategies[Lowercase30].GenerateColumnName("v古天乐"))
}

func TestName_Make_GenerateColumnName(t *testing.T) {
	t.Log(namingStrategies[Uppercase63].GenerateColumnName("v杨亦乐"))
	t.Log(namingStrategies[Uppercase63].GenerateColumnName("v刘亦菲的亦"))
	t.Log(namingStrategies[Uppercase63].GenerateColumnName("v古天乐的乐"))
}
