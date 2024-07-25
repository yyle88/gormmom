package gormmomrule

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRule_Validate(t *testing.T) {
	require.True(t, S63.Validate("abc"))
}

func TestValidate(t *testing.T) {
	require.True(t, Validate(S63, "abc", map[RULE]func(string) bool{}))
}
