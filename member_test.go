package organization

import "testing"

func TestAddMember(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	org, err := neworg()
	if err != nil {
		t.Fatal(err)
	}

	err = org.AddMember(map[string][]string{
		`cn`:       []string{`巩祥`},
		`sn`:       []string{`Kevin.Gong`},
		`email`:    []string{`aoxianglele@icloud.com`},
		`title`:    []string{`Developer`},
		`rbacRole`: []string{`b49jug06h301nm494sd0`},
		`rbacType`: []string{`b49jtn06h301mgko5jo0`},
		`unitID`:   []string{`b49kdrg6h302hrpggg8g`},
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestFetchMemberRoles(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	org, err := neworg()
	if err != nil {
		t.Fatal(err)
	}

	ids, err := org.RoleIDsByMemberID(`b49kehg6h302jg98oi70`)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ids)
}

func TestFetchOrgMemberByID(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	org, err := neworg()
	if err != nil {
		t.Fatal(err)
	}

	ids, err := org.OrganizationMemberByMemberID(`b49kehg6h302jg98oi70`)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ids)
}
