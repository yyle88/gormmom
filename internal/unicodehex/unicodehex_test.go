package unicodehex

import (
	"testing"
)

func TestS2LeHex4sUppers(t *testing.T) {
	t.Log(StringToHex4Uppercase("我叫杨亦乐"))
	t.Log(StringToHex4Uppercase("刘亦菲的亦"))
	t.Log(StringToHex4Uppercase("古天乐的乐"))
}
