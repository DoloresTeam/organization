package organization

import (
	"errors"
	"fmt"

	ldap "gopkg.in/ldap.v2"
)

// AddUnit to ldap
func (org *Organization) AddUnit(parentID string, info map[string][]string) (string, error) {

	id := generatorID()

	var dn string
	if len(parentID) == 0 {
		dn = org.dn(id, unit)
	} else {
		filter := fmt.Sprintf(`(id=%s)`, parentID)

		sq := ldap.NewSearchRequest(org.parentDN(unit),
			ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0, true, filter, []string{`id`}, nil)
		sr, err := org.l.Search(sq)
		if err != nil {
			return ``, err
		}
		if len(sr.Entries) != 1 {
			return ``, errors.New(`parent id invalid`)
		}

		dn = fmt.Sprintf(`id=%s,%s`, id, sr.Entries[0].DN)
	}

	aq := ldap.NewAddRequest(dn)

	aq.Attribute(`objectClass`, []string{`organizationalUnit`, `unit`, `top`})

	for k, v := range info {
		aq.Attribute(k, v)
	}

	return id, org.l.Add(aq)
}

// UnitByID ...
func (org *Organization) UnitByID(id string) (map[string]interface{}, error) {

	us, e := org.UnitByIDs([]string{id})
	if e != nil {
		return nil, e
	}
	if len(us) != 1 {
		return nil, errors.New(`found many units`)
	}
	return us[0], nil
}

// UnitByIDs ...
func (org *Organization) UnitByIDs(ids []string) ([]map[string]interface{}, error) {

	filter, err := sqConvertIDsToFilter(ids)
	if err != nil {
		return nil, err
	}

	return org.searchUnit(filter, true)
}

func (org *Organization) UnitSubIDs(id string) ([]string, error) {

	unit, err := org.UnitByID(id)
	if err != nil {
		return nil, err
	}

	return org.searchSubUnitIDs(unit[`dn`].(string))
}

// DelUnitByID ...
// func (org *Organization) DelUnitByID(id string) error {
//
// 	// 所有子部门都会被删除
// 	// 如果部门下有人，那么不允许删除
// }

// AllUnit ...
func (org *Organization) AllUnit() ([]map[string]interface{}, error) {
	return org.searchUnit(``, true)
}

// OrganizationUnitByMemberID ...
func (org *Organization) OrganizationUnitByMemberID(id string) ([]map[string]interface{}, error) {

	filter, err := org.filterConditionByMemberID(id, true)
	if err != nil {
		return nil, err
	}
	return org.searchUnit(filter, false)
}

func (org *Organization) filterConditionByMemberID(id string, isUnit bool) (string, error) {

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
