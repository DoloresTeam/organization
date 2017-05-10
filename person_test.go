package organization

import "testing"

func TestAddPerson(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	org, err := neworg()
	if err != nil {
		t.Fatal(err)
	}

	err = org.AddPerson(`Kevin.Gong`, `巩祥`, `b49961dhfpcuhne4dvl0`,
		[]string{`b49bu5lhfpcvvurbmns0`},
		[]string{`b499755hfpcul32vtt9g`},
		nil)
	if err != nil {
		t.Fatal(err)
	}
}
