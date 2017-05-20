package gorbacx

import "testing"

func TestPrettyPrint(t *testing.T) {

	rbacx := rbacx()
	rbacx.PrettyPrint()
}

func TestMatch(t *testing.T) {

	//define Permission
	rbacx := rbacx()

	ut := rbacx.MatchedTypes([]string{`1`}, true)

	if len(ut) != 6 {
		t.Fatal(ut)
	}
	t.Log(ut)
}

func TestFind(t *testing.T) {

	rbacx := rbacx()

	role, err := rbacx.RoleByID(`1`)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(role)

	p, err := rbacx.PermissionByID(`1`)
	if err != nil {
		t.Fatal(err)
	}
	p.Add([]string{`s1`, `s2`, `s3`})

	if p.types.Cardinality() != 6 {
		t.Fatal()
	}
	t.Log(p.types)
}

func rbacx() *RBACX {

	p1 := NewPermission(`1`, []string{`p1`, `p2`})
	p1.Add([]string{`p3`})

	p2 := NewPermission(`2`, []string{`u1`, `u2`})
	p2.Add([]string{`u3`})

	r1 := NewRole(`1`, []*Permission{p1, p2}, []*Permission{p1, p2})
	r2 := NewRole(`2`, []*Permission{p2}, []*Permission{p1})

	rbacx := New()
	rbacx.Add([]*Role{r1, r2})

	return rbacx
}
