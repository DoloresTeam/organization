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

func TestAddPermission(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	id, err := org.AddPermission(`Test`, `This is Test Permission`, []string{`b4oejsdhfpcjdr8fq6p0`}, true)
	if err != nil {
		t.Fatal(err)
	}

	// 测试 修改
	err = org.ModifyPermission(id, ``, ``, []string{`b4oejsdhfpcjdr8fq6p0`})
	if err != nil {
		t.Fatal(err)
	}

	// 删除测试添加的权限
	err = org.DelPermission(id)
	if err != nil {
		t.Fatal(err)
	}
}

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
