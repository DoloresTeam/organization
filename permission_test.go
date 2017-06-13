package organization

import (
	"testing"
)

func TestFetchAllPermission(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	ps, err := org.Permissions(false, 0, nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, p := range ps.Data {
		t.Log(p)
	}
}
