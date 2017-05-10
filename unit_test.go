package organization

import "testing"

func TestUnitByID(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	org, _ := neworg()

	r, e := org.UnitByID(`2`)
	if e != nil {
		t.Fatal(e)
	}

	t.Log(r)
}

func TestAddUnit(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	org, _ := neworg()

	e := org.AddUnit(`2`, `Test`, `Test add unit on parent`, `b49845lhfpcu1phd6ql0`)
	if e != nil {
		t.Fatal(e)
	}
}
