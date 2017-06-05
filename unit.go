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

func (org *Organization) ModifyUnit(id string, info map[string][]string) error {
	unit, err := org.UnitByID(id)
	if err != nil {
		return err
	}
	mq := ldap.NewModifyRequest(unit[`dn`].(string))
	for k, v := range info {
		mq.Replace(k, v)
	}
	return org.l.Modify(mq)
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
	dn := unit[`dn`].(string)

	sq := ldap.NewSearchRequest(dn, ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0, false, `(objectClass=organizationalUnit)`, []string{`id`}, nil)
	sr, err := org.l.Search(sq)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0)
	for _, e := range sr.Entries {
		ids = append(ids, e.GetAttributeValue(`id`))
	}

	return ids, nil
}

func (org *Organization) UnitByTypeIDs(ids []string) ([]map[string]interface{}, error) {
	filter, err := sqConvertArraysToFilter(`rbacType`, ids)
	if err != nil {
		return nil, err
	}
	return org.searchUnit(filter, true)
}

// DelUnitByID ...
func (org *Organization) DelUnit(id string) error {

	ids, err := org.UnitSubIDs(id)
	if err != nil {
		return err
	}

	// 通过部门ID 找员工
	mids, err := org.MemberIDsByDepartmentIDs(ids)
	if err != nil {
		return err
	}
	if len(mids) > 0 {
		return fmt.Errorf(`此部门下包含员工，请先修改员工信息 count: %d`, len(mids))
	}

	unit, err := org.UnitByID(id)
	if err != nil {
		return err
	}

	dq := ldap.NewDelRequest(unit[`dn`].(string), nil)

	return org.l.Del(dq)
}

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
