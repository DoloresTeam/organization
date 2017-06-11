package organization

import (
	"testing"
)

func TestAddType(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	id, err := org.AddType(`TestAndType`, `This is TestDescription`, true)
	if err != nil {
		t.Fatal(err)
	}

	err = org.ModifyType(id, `modify type name`, `modifytype name`, true)
	if err != nil {
		t.Fatal(err)
	}

	err = org.DelType(id, true)
	if err != nil {
		t.Fatal(err)
	}
}
