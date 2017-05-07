package organization

import "testing"

func TestAddRole(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	org, err := neworg()
	if err != nil {
		t.Fatal(err)
	}

	err = org.AddRole(`Role1`, `Test add Role 1`, []string{`b46rr2thfpcik527pgm0`}, []string{`b46rr2thfpcik527pgmg`})
	if err != nil {
		t.Fatal(err)
	}
}

func TestFilterRole(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	org, err := neworg()
	if err != nil {
		t.Fatal(err)
	}

	rs, err := org.AllRoles()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(`all roles`)
	for _, r := range rs {
		t.Log(r)
	}

	t.Log(`filter role ` + `b47kco86h3053ecopjd0`)
	rs, err = org.RoleByIDs([]string{`b47kco86h3053ecopjd0`})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rs)
}
