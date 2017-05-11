package organization

import "testing"

func TestUnitByID(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	org, _ := neworg()

	r, e := org.AllUnit()
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

	err := org.AddUnit(`b49kdrg6h302hrpggg8g`, map[string][]string{
		`ou`:          []string{`Test`},
		`description`: []string{`This is test unit's description`},
		`rbacType`:    []string{`b49jtn06h301mgko5jng`},
	})

	if err != nil {
		t.Fatal(err)
	}
}
