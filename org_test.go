package organization

import "testing"

func TestNewOrganizationWithSimpleBind(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}
	org, err := NewOrganizationWithSimpleBind(`dc=dolores,dc=org`, `localhost`, `cn=admin,dc=dolores,dc=org`, `secret`, 389)
	if err != nil {
		t.Fatal(err)
	}

	org.rbacx.PrettyPrint()
}
