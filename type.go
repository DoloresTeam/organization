package organization

import (
	ldap "gopkg.in/ldap.v2"

	"github.com/DoloresTeam/organization/entry"
)

func (org *Organization) AddType(name, description string, isUnitType bool) error {

	which := entry.PersonType
	if isUnitType {
		which = entry.UnitType
	}

	aq, err := entry.AddType(name, description, which)
	if err != nil {
		return err
	}

	return org.l.Add(aq)
}

func (org *Organization) ModifyType(oid string, name, description string, isUnitType bool) error {

	which := entry.PersonType
	if isUnitType {
		which = entry.UnitType
	}

	mq, err := entry.ModifyType(oid, name, description, which)
	if err != nil {
		return err
	}

	return org.l.Modify(mq)
}

func (org *Organization) DelType(oid string, isUnitType bool) error {

	// 删除类型之前需要确认没有任何一个权限有引用这个类型
	return UnimplementError
}

func (org *Organization) allTypes(isUnitType bool) (*ldap.SearchResult, error) {

	which := entry.PersonType
	if isUnitType {
		which = entry.UnitType
	}

	sq := entry.FetchAllTypes(which)

	return org.l.Search(sq)
}
