package organization

import (
	"fmt"

	"github.com/DoloresTeam/organization/godn"
	"github.com/rs/xid"
)

const (
	member           = 0
	role             = 1
	unitPermission   = 2
	memberPermission = 3
	unitType         = 4
	memberType       = 5
	unit             = 6
)

func (org *Organization) dn(id string, category int) string {

	return fmt.Sprintf(`id=%s,%s`, id, org.parentDN(category))
}

func (org *Organization) parentDN(category int) string {

	baseDN := org.subffix

	switch category {
	case member:
		baseDN = godn.Member(org.subffix)
	case role:
		baseDN = godn.Role(org.subffix)
	case unitPermission:
		baseDN = godn.Permission(org.subffix, true)
	case memberPermission:
		baseDN = godn.Permission(org.subffix, false)
	case unitType:
		baseDN = godn.DoloresType(org.subffix, true)
	case memberType:
		baseDN = godn.DoloresType(org.subffix, false)
	case unit:
		baseDN = godn.Unit(org.subffix)
	}

	return baseDN
}

func typeCategory(isUnit bool) int {
	if isUnit {
		return unitType
	}
	return memberType
}

func permissionCategory(isUnit bool) int {
	if isUnit {
		return unitPermission
	}
	return memberPermission
}

func generatorID() string {
	return xid.New().String()
}
