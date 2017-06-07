package organization

import "testing"

func TestFilterRole(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
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
	rss, err := org.RoleByIDs([]string{`b47kco86h3053ecopjd0`})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rss)
}

func TestRBAC(t *testing.T) {

	org.rbacx.PrettyPrint()
}
