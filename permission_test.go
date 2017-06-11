package organization

import (
	"testing"
)

func TestPermission(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	id, err := org.AddPermission(`Test`, `This is Test Permission`, []string{`b4oejsdhfpcjdr8fq6p0`}, true)
	if err != nil {
		t.Fatal(err)
	}

	// 测试 修改
	err = org.ModifyPermission(id, ``, ``, []string{`b4oejsdhfpcjdr8fq6p0`, `b4rts55hfpclmh1obi2g`})
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

	for _, p := range ps.Data {
		t.Log(p)
	}
}
