package enctohex

import (
	"testing"
)

func TestS2LeHex4sUppers(t *testing.T) {
	t.Log(S2LeHex4sUppers("我叫杨亦乐"))
	t.Log(S2LeHex4sUppers("刘亦菲的亦"))
	t.Log(S2LeHex4sUppers("古天乐的乐"))
}
