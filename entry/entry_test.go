package entry

import "testing"

func TestSetSubfix(t *testing.T) {

	SetSubffix(`s`)

	if subffix != `s` {
		t.FailNow()
	}
}
