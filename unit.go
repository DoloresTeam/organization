package organization

import (
	"errors"
	"fmt"

	ldap "gopkg.in/ldap.v2"
)

// UnitAttributes ...
var UnitAttributes = [...]string{`id`, `ou`, `description`, `priority`}

// AddUnit to ldap
func (org *Organization) AddUnit(parentID string, info map[string][]string) (string, error) {

	id := generateNewID()

	var dn string
	if len(parentID) == 0 {
		dn = org.dn(id, unit)
	} else {
		filter := fmt.Sprintf(`(id=%s)`, parentID)

		sq := ldap.NewSearchRequest(org.parentDN(unit),
			ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0, true, filter, []string{`id`}, nil)
		sr, err := org.Search(sq)
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

	err := org.Add(aq)
	if err != nil {
		return ``, err
	}

	unit, err := org.UnitByID(id)
	if err != nil {
		return ``, err
	}
	go org.logAddUnit(unit)
	return id, nil
}

// ModifyUnit ...
func (org *Organization) ModifyUnit(id string, info map[string][]string) error {
	oUnit, err := org.UnitByID(id)
	if err != nil {
		return err
	}
	mq := ldap.NewModifyRequest(oUnit[`dn`].(string))
	for k, v := range info {
		mq.Replace(k, v)
	}
	err = org.Modify(mq)
	if err != nil {
		return err
	}
	nUnit, err := org.UnitByID(oUnit[`id`].(string))
	if err != nil {
		return err
	}
	return org.logModifyUnit(oUnit, nUnit)
}

// DelUnit ...
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

	err = org.Del(dq)
	if err != nil {
		return err
	}
	go org.logDelUnit(unit[`id`].(string), unit[`rbacType`].(string))
	return nil
}

// UnitByID ...
func (org *Organization) UnitByID(id string) (map[string]interface{}, error) {

	us, e := org.UnitByIDs([]string{id})
	if e != nil {
		return nil, e
	}
	if len(us) > 1 {
		return nil, fmt.Errorf(`%s found many units`, id)
	}
	if len(us) == 0 {
		return nil, fmt.Errorf(`%s unit does't exist`, id)
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

// UnitSubIDs ....
func (org *Organization) UnitSubIDs(id string) ([]string, error) {

	unit, err := org.UnitByID(id)
	if err != nil {
		return nil, err
	}
	dn := unit[`dn`].(string)

	sq := ldap.NewSearchRequest(dn, ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0, false, `(objectClass=organizationalUnit)`, []string{`id`}, nil)
	sr, err := org.Search(sq)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0)
	for _, e := range sr.Entries {
		ids = append(ids, e.GetAttributeValue(`id`))
	}

	return ids, nil
}

// UnitIDsByTypeIDs ...
func (org *Organization) UnitIDsByTypeIDs(ids []string) ([]string, error) {
	filter, err := sqConvertArraysToFilter(`rbacType`, ids)
	if err != nil {
		return nil, err
	}
	sq := ldap.NewSearchRequest(org.parentDN(unit), ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0, false, filter, []string{`id`}, nil)
	sr, err := org.Search(sq)
	if err != nil {
		return nil, err
	}
	unitIDs := make([]string, 0)
	for _, e := range sr.Entries {
		unitIDs = append(unitIDs, e.GetAttributeValue(`id`))
	}

	return unitIDs, nil
}

// AllUnit ...
func (org *Organization) AllUnit() ([]map[string]interface{}, error) {
	return org.searchUnit(``, true)
}
