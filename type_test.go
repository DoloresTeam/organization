package organization

import (
	"testing"
)

// func TestAddType(t *testing.T) {
//
// 	if testing.Short() {
// 		t.SkipNow()
// 	}
//
// 	org, err := neworg()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	id, err := org.AddType(`TestAndType`, `This is TestDescription`, true)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Log(id)
//
// 	id, err = org.AddType(`TestAndType`, `This is TestDescription`, false)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Log(id)
// }

func TestAllTypes(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	types, err := org.Types(false, 0, nil)
	if err != nil {
		t.Fatal(err)
	}

	// types, err := org.TypeByIDs([]string{`b46otklhfpcs0pe51am0`, `b46otklhfpcs0pe51am0`}, true)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(`fileter id:b46otklhfpcs0pe51am0 b46otklhfpcs0pe51am0`)
	for _, ty := range types.Data {
		t.Log(ty)
	}
}

func TestModifyType(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	err := org.ModifyType(`b4f8h9dhfpci6tfa4jj0`, `Modify Type Name`, `Modify Type Description`, true)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDelType(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	err := org.DelType(`b45v01dhfpcidf9rgtag`, true)
	if err != nil {
		t.Fatal(err)
	}
}
