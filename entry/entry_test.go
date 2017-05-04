package entry

import "testing"

func TestSetSubfix(t *testing.T) {

	SetSubfix(`s`)

	if subfix != `s` {
		t.FailNow()
	}
}
