package organization

import "testing"

var org, _ = NewOrganizationWithSimpleBind(`dc=dolores,dc=store`, `dolores.store`, `cn=admin,dc=dolores,dc=store`, `dolores`, 389)

func TestNewOrganizationWithSimpleBind(t *testing.T) {
}
