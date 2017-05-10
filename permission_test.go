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

	_, _ = org.PermissionByType(`b45v085hfpcidvk1m8fg`, true)
}

func TestAddPermission(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	org, err := neworg()
	if err != nil {
		t.Fatal(err)
	}

	err = org.AddPermission(`Test`, `This is Test Permission`, []string{`b49961dhfpcuhne4dvkg`, `b49bk6dhfpcvsehgvm1g`}, true)
	if err != nil {
		t.Fatal(err)
	}

	err = org.AddPermission(`Test`, `This is Test Permission`, []string{`b49961dhfpcuhne4dvl0`, `b49bk6dhfpcvsehgvm20`}, false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFetchAllPermission(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	org, err := neworg()
	if err != nil {
		t.Fatal(err)
	}

	ps, err := org.AllPermissions(true)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(`all permissions`)
	for _, p := range ps {
		t.Log(p)
	}
}
