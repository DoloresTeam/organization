package organization

import "testing"

var org, _ = NewOrganizationWithSimpleBind(`dc=dolores,dc=store`, `dolores.store`, `cn=admin,dc=dolores,dc=store`, `dolores`, 389)

func TestNewOrganizationWithSimpleBind(t *testing.T) {
	if org == nil {
		t.Fatal(`org initial failed`)
	}
}

func BenchmarkOriganizationView(b *testing.B) {
	b.RunParallel(func(arg1 *testing.PB) {
		for arg1.Next() {
			org.OrganizationView(`b4r6e05hfpckh33hnsq0`)
		}
	})
}
