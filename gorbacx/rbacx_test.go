package gorbacx

import "testing"

func TestMatch(t *testing.T) {

	//define Permission
	rbacx := rbacx()

	ut, pt := rbacx.MatchedTypes([]string{`1`})

	if len(ut) != 3 {
		t.Fatal(ut)
	}
	t.Log(ut)

	if len(pt) != 3 {
		t.Fatal(pt)
	}
	t.Log(pt)

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
	p.Add([]string{`s4`, `s5`}, false)

	if p.personTypes.Cardinality() != 5 {
		t.Fatal()
	}
}

func rbacx() *RBACX {

	p1 := NewPermission(`1`, []string{`u1`, `u2`}, []string{`p1`, `p2`})
	p1.Add([]string{`p3`}, false)

	p2 := NewPermission(`2`, []string{`u1`, `u2`}, []string{`p1`, `p2`})
	p2.Add([]string{`u3`}, true)

	r1 := NewRole(`1`, []*Permission{p1, p2})
	r2 := NewRole(`2`, []*Permission{p2})

	rbacx := New()
	rbacx.Add([]*Role{r1, r2})

	return rbacx
}
