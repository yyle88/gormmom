package gormmomrule

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestName_Check(t *testing.T) {
	require.True(t, presetNameImpMap[S63].CheckName("abc"))
}

func TestName_Check_2(t *testing.T) {
	require.True(t, presetNameImpMap[S63U].CheckName("ABC"))
}

func TestName_Make(t *testing.T) {
	t.Log(presetNameImpMap[S30].GenNewCnm("v杨亦乐"))
	t.Log(presetNameImpMap[S30].GenNewCnm("v刘亦菲"))
	t.Log(presetNameImpMap[S30].GenNewCnm("v古天乐"))
}

func TestName_Make_2(t *testing.T) {
	t.Log(presetNameImpMap[S63U].GenNewCnm("v杨亦乐"))
	t.Log(presetNameImpMap[S63U].GenNewCnm("v刘亦菲的亦"))
	t.Log(presetNameImpMap[S63U].GenNewCnm("v古天乐的乐"))
}
