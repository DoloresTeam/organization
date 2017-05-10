package organization

import "testing"

func TestUnitByID(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	org, _ := neworg()

	r, e := org.UnitByIDs([]string{`2`})
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

	e := org.AddUnit(`b4996cdhfpcuikevr1lg`, `Test`, `Test add unit on parent`, `b49961dhfpcuhne4dvkg`, nil)
	if e != nil {
		t.Fatal(e)
	}
}
