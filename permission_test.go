package organization

import (
	"errors"
	"testing"
)

func TestSearchPermission(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	ps, _ := org.PermissionByType(`b4drqelhfpcqn7f7du50`, true)
	if len(ps) == 0 {
		t.Fatal(errors.New(`no permission`))
	}
}

// func TestAddPermission(t *testing.T) {
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
// 	_, err = org.AddPermission(`Test`, `This is Test Permission`, []string{`b4drqelhfpcqn7f7du50`}, true)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	_, err = org.AddPermission(`Test`, `This is Test Permission`, []string{`b4drqelhfpcqn7f7du5g`}, false)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

func TestFetchAllPermission(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	ps, err := org.Permissions(false, 0, nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(`all permissions`)
	for _, p := range ps.Data {
		t.Log(p)
	}
}
