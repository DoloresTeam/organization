package organization

import (
	"errors"
	"fmt"

	ldap "gopkg.in/ldap.v2"
)

// AddUnit to ldap
func (org *Organization) AddUnit(parentID string, info map[string][]string) error {

	id := generatorID()

	var dn string
	if len(parentID) == 0 {
		dn = org.dn(id, unit)
	} else {
		_, err := org.UnitByIDs([]string{parentID})
		if err != nil {
			return errors.New(`parent must be not nil`)
		}
		dn = fmt.Sprintf(`id=%s,id=%s,%s`, id, parentID, org.parentDN(unit))
	}

	aq := ldap.NewAddRequest(dn)

	aq.Attribute(`objectClass`, []string{`organizationalUnit`, `unit`, `top`})

	for k, v := range info {
		aq.Attribute(k, v)
	}

	return org.l.Add(aq)
}

// UnitByIDs ...
func (org *Organization) UnitByIDs(ids []string) ([]map[string]interface{}, error) {

	dn := org.parentDN(unit)
	filter, err := sqConvertIDsToFilter(ids)
	if err != nil {
		return nil, err
	}

	sq := &searchRequest{dn, filter, []string{`id`}, nil, 0, nil}
	r, err := org.search(sq)
	if err != nil {
		return nil, err
	}

	return r.Data, nil
}

// OrganizationUnitByMemberID ...
func (org *Organization) OrganizationUnitByMemberID(id string) ([]map[string]interface{}, error) {

	filter, err := org.filterByMemberID(id, true)
	if err != nil {
		return nil, err
	}
	r, e := org.searchUnit(filter, false, 0, nil)
	if e != nil {
		return nil, e
	}
	return r.Data, nil
}

func (org *Organization) filterByMemberID(id string, isUnit bool) (string, error) {

	// 通过id 拿到所有的 角色
	roleIDs, err := org.RoleIDsByMemberID(id)
	if err != nil {
		return ``, err
	}

	types := org.rbacx.MatchedTypes(roleIDs, isUnit)

	filter, err := sqConvertArraysToFilter(`rbacType`, types)
	if err != nil {
		return ``, err
	}

	return filter, nil
}
