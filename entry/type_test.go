package entry

import (
	"strings"
	"testing"
)

func TestAddType(t *testing.T) {

	SetSubfix(`dc=test,dc=go`)

	aq, err := AddType(`TestType`, `This is TestType from type_test`, UnitType)

	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasSuffix(aq.DN, subfix) {
		t.Fatal(`Request dn is invalid`)
	}

	if len(aq.Attributes) != 3 {
		t.Failed()
	}

	t.Log(aq.Attributes)

	aq, _ = AddType(`TestType`, ``, PersonType)
	if len(aq.Attributes) != 2 {
		t.Failed()
	}

	t.Log(aq.Attributes)
}
