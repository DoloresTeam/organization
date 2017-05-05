package entry

import (
	"strings"
	"testing"
)

func TestAddType(t *testing.T) {

	SetSubffix(`dc=test,dc=go`)

	aq, err := AddType(`TestType`, `This is TestType from type_test`, UnitType)

	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasSuffix(aq.DN, subffix) {
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

func TestModifyType(t *testing.T) {

	SetSubffix(`dc=test,dc=go`)

	_, err := ModifyType(``, ``, ``, ``)
	if err == nil {
		t.Fatal(`please check oid juddge logic`)
	}

	_, err = ModifyType(`1`, ``, ``, UnitType)
	if err == nil {
		t.FailNow()
	}

	mq, _ := ModifyType(`1`, `name`, `description string`, PersonType)
	if mq == nil {
		t.FailNow()
	}
}

func TestDelType(t *testing.T) {

	dq := DelType(`2`, UnitType)

	t.Log(dq)
}
