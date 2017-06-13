package organization

import (
	"errors"
	"fmt"
	"strings"

	ldap "gopkg.in/ldap.v2"
)

// AddType desgined to add a new dolresType
func (org *Organization) AddType(name, description string, isUnit bool) (string, error) {

	id := generateNewID()
	dn := org.dn(id, typeCategory(isUnit))
	aq := ldap.NewAddRequest(dn)

	aq.Attribute(`objectClass`, []string{`doloresType`, `top`})
	aq.Attribute(`cn`, []string{name})
	aq.Attribute(`description`, []string{description})

	return id, org.l.Add(aq)
}

// ModifyType update name or description of doloresType
func (org *Organization) ModifyType(id string, name, description string, isUnit bool) error {

	dn := org.dn(id, typeCategory(isUnit))
	mq := ldap.NewModifyRequest(dn)

	if len(name) != 0 {
		mq.Replace(`cn`, []string{name})
	}
	if len(description) != 0 {
		mq.Replace(`description`, []string{description})
	}

	return org.l.Modify(mq)
}

// DelType by id
func (org *Organization) DelType(id string, isUnit bool) error {

	pIDs, e := org.PermissionByType(id, isUnit)
	if e != nil {
		return e
	}
	if len(pIDs) != 0 {
		return fmt.Errorf(`有权限规则包含此类型，请先修改权限规则 count: %d`, len(pIDs))
	}

	if isUnit {
		us, e := org.UnitIDsByTypeIDs([]string{id})
		if e != nil {
			return e
		}
		if len(us) != 0 {
			return fmt.Errorf(`有部门属于此类型，请先修改部门所属类型 count: %d`, len(pIDs))
		}
	} else {
		// 通过Type 找人
		mIDs, e := org.MemberIDsByTypeIDs([]string{id})
		if e != nil {
			return e
		}
		if len(mIDs) != 0 {
			return fmt.Errorf(`有员工属于此类型，请先修改员工所属类型 count: %d`, len(pIDs))
		}
	}

	dn := org.dn(id, typeCategory(isUnit))
	dq := ldap.NewDelRequest(dn, nil)

	return org.l.Del(dq)
}

// Types in ldap server
func (org *Organization) Types(isUnit bool, pageSize uint32, cookie []byte) (*SearchResult, error) {
	sq := &searchRequest{
		org.parentDN(typeCategory(isUnit)),
		`(objectClass=doloresType)`,
		[]string{`id`, `cn`, `description`, `modifyTimestamp`}, nil,
		pageSize,
		cookie}
	return org.search(sq)
}

// TypeByIDs ...
func (org *Organization) TypeByIDs(ids []string) ([]map[string]interface{}, error) {
	filter, err := sqConvertIDsToFilter(ids)
	if err != nil {
		return nil, err
	}
	dn := fmt.Sprintf(`ou=type,%s`, org.subffix)

	sq := ldap.NewSearchRequest(dn,
		ldap.ScopeWholeSubtree,
		ldap.DerefAlways,
		0, 0, false, filter, []string{`id`, `cn`, `description`}, nil)

	sr, err := org.l.Search(sq)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}

	for _, e := range sr.Entries {
		isUnit := false
		if strings.Contains(e.DN, `ou=unit`) {
			isUnit = true
		}
		result = append(result, map[string]interface{}{
			`id`:          e.GetAttributeValue(`id`),
			`cn`:          e.GetAttributeValue(`cn`),
			`description`: e.GetAttributeValue(`description`),
			`isUnit`:      isUnit,
		})
	}

	if result == nil {
		return nil, errors.New(`not found`)
	}
	return result, nil
}

// TypeByID ...
func (org *Organization) TypeByID(id string) (map[string]interface{}, error) {

	types, err := org.TypeByIDs([]string{id})
	if err != nil {
		return nil, err
	}
	if len(types) != 1 {
		return nil, errors.New(`found many types`)
	}

	return types[0], nil
}

// TypeByPermissionID ...
func (org *Organization) TypeByPermissionID(id string) ([]map[string]interface{}, error) {

	p, e := org.PermissionByID(id)
	if e != nil {
		return nil, e
	}
	return org.TypeByIDs(p[`rbacType`].([]string))
}
