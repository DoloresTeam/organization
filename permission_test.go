package organization

import "testing"

func TestSearchPermission(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	org, err := neworg()
	if err != nil {
		t.Fatal(err)
	}

	_, err = org.PermissionByType(`b45v085hfpcidvk1m8fg`, true)
}

func TestAddPermission(t *testing.T) {

	org, err := neworg()
	if err != nil {
		t.Fatal(err)
	}

	err = org.AddPermission(`Test`, `This is Test Permission`, []string{`1`, `2`}, true)
	if err != nil {
		t.Fatal(err)
	}

	err = org.AddPermission(`Test`, `This is Test Permission`, []string{`1`, `2`}, false)
	if err != nil {
		t.Fatal(err)
	}
}
