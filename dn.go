package organization

import (
	"fmt"

	"github.com/DoloresTeam/organization/godn"
	"github.com/rs/xid"
)

const (
	person           = 0
	role             = 1
	unitPermission   = 2
	personPermission = 3
	unitType         = 4
	personType       = 5
	unit             = 6
)

func (org *Organization) dn(oid string, category int) string {

	return fmt.Sprintf(`oid=%s,%s`, oid, org.parentDN(category))
}

func (org *Organization) parentDN(category int) string {

	baseDN := org.subffix

	switch category {
	case person:
		baseDN = godn.Person(org.subffix)
	case role:
		baseDN = godn.Role(org.subffix)
	case unitPermission:
		baseDN = godn.Permission(org.subffix, true)
	case personPermission:
		baseDN = godn.Permission(org.subffix, false)
	case unitType:
		baseDN = godn.DoloresType(org.subffix, true)
	case personType:
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
	return personType
}

func permissionCategory(isUnit bool) int {
	if isUnit {
		return unitPermission
	}
	return personPermission
}

func generatorOID() string {
	return xid.New().String()
}
