package organization

import "testing"

func TestUnitByID(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	org, _ := neworg()

	r, e := org.UnitByIDs([]string{`b4ju9dthfpcjdopdqcl0`})
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

	r, e := org.OrganizationUnitByMemberID(`b4ju9dthfpcjdopdqcl0`)
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

	_, err := org.AddUnit(`b4jv4llhfpcjtv6o07og`, map[string][]string{
		`ou`:          []string{`iOS-Tester-sub`},
		`description`: []string{`iOS-Tester-sub is a test ou`},
		`rbacType`:    []string{`b4drradhfpcqnna2pvh0`},
	})

	if err != nil {
		t.Fatal(err)
	}
}
