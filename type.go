package organization

import (
	"errors"
	"fmt"

	ldap "gopkg.in/ldap.v2"
)

// AddType desgined to add a new dolresType
func (org *Organization) AddType(name, description string, isUnit bool) (string, error) {

	id := generatorID()
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

	pids, err := org.PermissionByType(id, isUnit)
	if err != nil {
		return err
	}
	if len(pids) > 0 {
		return errors.New(`has Permission refrence this type,`)
	}

	dn := org.dn(id, typeCategory(isUnit))
	dq := ldap.NewDelRequest(dn, nil)

	return org.l.Del(dq)
}

// Types in ldap server
func (org *Organization) Types(isUnit bool, pageSize uint32, cookie []byte) (*SearchResult, error) {

	return org.searchType(``, org.parentDN(typeCategory(isUnit)), pageSize, cookie)
}

// TypeByIDs ...
func (org *Organization) TypeByIDs(ids []string) ([]map[string]interface{}, error) {
	filter, err := sqConvertIDsToFilter(ids)
	if err != nil {
		return nil, err
	}
	dn := fmt.Sprintf(`ou=type,%s`, org.subffix)
	r, e := org.searchType(filter, dn, 0, nil)
	if e != nil {
		return nil, e
	}
	return r.Data, nil
}

// TypeByID ...
func (org *Organization) TypeByID(id string) (map[string]interface{}, error) {

	types, err := org.TypeByIDs([]string{id})
	if err != nil {
		return nil, err
	}
	if len(types) == 0 {
		types, err = org.TypeByIDs([]string{id})
		if err != nil {
			return nil, err
		}
		if len(types) == 0 {
			return nil, errors.New(`404 Not Found`)
		}
		t := types[0]
		t[`isUnit`] = false
		return t, nil
	}

	t := types[0]
	t[`isUnit`] = true
	return t, nil
}

// TypeByPermissionID ...
func (org *Organization) TypeByPermissionID(id string) ([]map[string]interface{}, error) {

	p, e := org.PermissionByID(id)
	if e != nil {
		return nil, e
	}
	return org.TypeByIDs(p[`rbacType`].([]string))
}
