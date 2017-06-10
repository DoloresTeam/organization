package organization

import "testing"

func TestRole(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	id, err := org.AddRole(`Test`, `Test and Role`,
		[]string{`b4ofic5hfpcjdr8fq6qg`, `b4rrr85hfpcj0qpeire0`},
		[]string{`b4ohonthfpckql08mas0`, `b4rtt2lhfpclmh1obi30`})
	if err != nil {
		t.Fatal(err)
	}

	err = org.ModifyRole(id, ``, ``, []string{`b4ofic5hfpcjdr8fq6qg`}, []string{`b4ohonthfpckql08mas0`})
	if err != nil {
		t.Fatal(err)
	}

	err = org.DelRole(id)
	if err != nil {
		t.Fatal(err)
	}
}
