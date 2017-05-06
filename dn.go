package organization

import (
	"fmt"

	"github.com/DoloresTeam/organization/godn"
	"github.com/rs/xid"
)

const (
	PERSON           = 0
	ROLE             = 1
	UNITPERSIMISSON  = 2
	PERSONPERMISSION = 3
	UNITTYPE         = 4
	PERSONTYPE       = 5
	UNIT             = 6
)

func (org *Organization) dn(oid string, category int) string {

	return fmt.Sprintf(`oid=%s,%s`, oid, org.parentDN(category))
}

func (org *Organization) parentDN(category int) string {

	baseDN := org.subffix

	switch category {
	case PERSON:
		baseDN = godn.Person(org.subffix)
	case ROLE:
		baseDN = godn.Role(org.subffix)
	case UNITPERSIMISSON:
		baseDN = godn.Permission(org.subffix, true)
	case PERSONPERMISSION:
		baseDN = godn.Permission(org.subffix, false)
	case UNITTYPE:
		baseDN = godn.DoloresType(org.subffix, true)
	case PERSONTYPE:
		baseDN = godn.DoloresType(org.subffix, false)
	case UNIT:
		baseDN = godn.Unit(org.subffix)
	}

	return baseDN
}

func typeCategory(isUnit bool) int {
	if isUnit {
		return UNITTYPE
	}
	return PERSONTYPE
}

func permissionCategory(isUnit bool) int {
	if isUnit {
		return UNITPERSIMISSON
	}
	return PERSONPERMISSION
}

func generatorOID() string {
	return xid.New().String()
}
