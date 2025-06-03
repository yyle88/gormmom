package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/rese"
	"github.com/yyle88/rese/resb"
	"github.com/yyle88/runpath"
)

func TestListGoFiles(t *testing.T) {
	paths := ListGoFiles(runpath.PARENT.Path())
	t.Log(neatjsons.S(paths))
	require.NotEmpty(t, paths)
}

type Example struct {
	Name    string `json:"name"`
	Account string `gorm:"column:account"`
	Score   string `yaml:"score"`
}

func TestParseTags(t *testing.T) {
	data := rese.A1(os.ReadFile(runpath.Path()))
	tags := ParseTags(data, &Example{})
	t.Log(neatjsons.S(tags))
	require.Equal(t, 3, tags.Size())
	keys := tags.Keys()
	require.Equal(t, "Name", keys[0])
	require.Equal(t, `json:"name"`, TrimBackticks(resb.C1(tags.Get(keys[0]))))
	require.Equal(t, "Account", keys[1])
	require.Equal(t, `gorm:"column:account"`, TrimBackticks(resb.C1(tags.Get(keys[1]))))
	require.Equal(t, "Score", keys[2])
	require.Equal(t, `yaml:"score"`, TrimBackticks(resb.C1(tags.Get(keys[2]))))
}

func TestParseTagsTrimBackticks(t *testing.T) {
	data := rese.A1(os.ReadFile(runpath.Path()))
	tags := ParseTagsTrimBackticks(data, &Example{})
	t.Log(neatjsons.S(tags))
	require.Equal(t, `gorm:"column:account"`, resb.C1(tags.Get("Account")))
}
