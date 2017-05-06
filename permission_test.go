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
