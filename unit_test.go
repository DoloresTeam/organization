package organization

import "testing"

func TestUnitByID(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	org, _ := neworg()

	r, e := org.UnitByIDs([]string{`b4ds0t5hfpcr4h3thtd0`})
	if e != nil {
		t.Fatal(e)
	}

	t.Log(r)
}

func TestOrganizationUnitByMemberID(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	org, err := neworg()
	if err != nil {
		t.Fatal(err)
	}

	r, e := org.OrganizationUnitByMemberID(`b49kehg6h302jg98oi70`)
	if e != nil {
		t.Fatal(e)
	}

	for _, v := range r {
		t.Log(v[`id`])
	}
}

func TestAddUnit(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	org, _ := neworg()

	err := org.AddUnit(``, map[string][]string{
		`ou`:          []string{`iOS-Developer`},
		`description`: []string{`iOS-DEveloper is a test ou`},
		`rbacType`:    []string{`b4drradhfpcqnna2pvh0`},
	})

	if err != nil {
		t.Fatal(err)
	}
}
