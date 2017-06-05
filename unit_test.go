package organization

import "testing"

func TestUnitByID(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	_, e := org.UnitByIDs([]string{`b4ds0t5hfpcr4h3thtd0`})
	if e != nil {
		t.Fatal(e)
	}

	ids, err := org.UnitSubIDs(`b4g14m0m20mgfdkk5bk0`)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(ids)
}

func TestOrganizationUnitByMemberID(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	r, e := org.OrganizationUnitByMemberID(`b4ju9dthfpcjdopdqcl0`)
	if e != nil {
		t.Fatal(e)
	}

	for _, v := range r {
		t.Log(v[`id`])
	}
}

// func TestAddUnit(t *testing.T) {
//
// 	if testing.Short() {
// 		t.SkipNow()
// 	}
//
// 	_, err := org.AddUnit(`b4jv4llhfpcjtv6o07og`, map[string][]string{
// 		`ou`:          []string{`iOS-Tester-sub`},
// 		`description`: []string{`iOS-Tester-sub is a test ou`},
// 		`rbacType`:    []string{`b4drradhfpcqnna2pvh0`},
// 	})
//
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }
