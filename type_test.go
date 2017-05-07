package organization

import (
	"testing"
)

func TestAddType(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	org, err := neworg()
	if err != nil {
		t.Fatal(err)
	}

	err = org.AddType(`TestAndType`, `This is TestDescription`, true)
	if err != nil {
		t.Fatal(err)
	}

	err = org.AddType(`TestAndType`, `This is TestDescription`, false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAllTypes(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	org, err := neworg()
	if err != nil {
		t.Fatal(err)
	}

	_, err = org.AllTypes(false)
	if err != nil {
		t.Fatal(err)
	}

	types, err := org.AllTypes(true)
	if err != nil {
		t.Fatal(err)
	}
	for _, ty := range types {
		t.Log(ty)
	}
}

func TestDelType(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	org, err := neworg()
	if err != nil {
		t.Fatal(err)
	}

	err = org.DelType(`b45v01dhfpcidf9rgtag`, true)
	if err != nil {
		t.Fatal(err)
	}
}

func neworg() (*Organization, error) {

	return NewOrganizationWithSimpleBind(`dc=dolores,dc=org`, `localhost`, `cn=admin,dc=dolores,dc=org`, `secret`, 389)
}
