package organization

import "testing"

var c = make(chan []string, 0)
var org, _ = NewOrganizationWithSimpleBind(`dc=dolores,dc=store`, `dolores.store`, `cn=admin,dc=dolores,dc=store`, `dolores`, 389, c)

func TestNewOrganizationWithSimpleBind(t *testing.T) {
	if org == nil {
		t.Fatal(`org initial failed`)
	}
}

func BenchmarkOriganizationView(b *testing.B) {
	b.RunParallel(func(arg1 *testing.PB) {
		for arg1.Next() {
			_, _, _, err := org.OrganizationView(`b4r6e05hfpckh33hnsq0`)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
