package organization

import "testing"

func TestUnit(t *testing.T) {
	// if testing.Short() {
	// 	t.SkipNow()
	// }

	id, err := org.AddUnit(``, map[string][]string{
		`ou`:          []string{`Test Add Unit`},
		`description`: []string{`This is description`},
		`rbacType`:    []string{`b4oejsdhfpcjdr8fq6p0`},
	})
	if err != nil {
		t.Fatal(err)
	}

	err = org.ModifyUnit(id, map[string][]string{
		`rbacType`: []string{`b4rts55hfpclmh1obi2g`},
	})
	if err != nil {
		t.Fatal(err)
	}

	err = org.DelUnit(id)
	if err != nil {
		t.Fatal(err)
	}
}
