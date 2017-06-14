package organization

import (
	"errors"
	"fmt"
	"strings"

	ldap "gopkg.in/ldap.v2"
)

// AddPermission to ldap server
func (org *Organization) AddPermission(name, description string, types []string, isUnit bool) (string, error) {

	ts, _ := org.TypeByIDs(types)
	if len(ts) != len(types) {
		return ``, errors.New(`invalid types`)
	}

	id := generateNewID()
	dn := org.dn(id, permissionCategory(isUnit))
	aq := ldap.NewAddRequest(dn)

	aq.Attribute(`objectClass`, []string{`permission`, `top`})

	aq.Attribute(`cn`, []string{name})
	aq.Attribute(`description`, []string{description})
	aq.Attribute(`rbacType`, types)

	return id, org.Add(aq)
}

// ModifyPermission in ldap
func (org *Organization) ModifyPermission(id, name, description string, types []string) error {

	op, err := org.PermissionByID(id)
	if err != nil {
		return err
	}

	mq := ldap.NewModifyRequest(op[`dn`].(string))

	if len(name) != 0 {
		mq.Replace(`cn`, []string{name})
	}
	if len(description) != 0 {
		mq.Replace(`description`, []string{description})
	}
	if types != nil {
		mq.Replace(`rbacType`, types)
	}

	err = org.Modify(mq)
	if err != nil {
		return err
	}

	rids, _ := org.RoleIDsByPermissionID(id)
	if len(rids) != 0 { // 有角色应用改权限
		org.refreshRBACIfNeeded(op[`rbacType`].([]string), types)
	}

	return nil
}

// DelPermission in ldap
func (org *Organization) DelPermission(id string) error {

	rids, err := org.RoleIDsByPermissionID(id)
	if err != nil {
		return err
	}

	if len(rids) > 0 {
		return fmt.Errorf(`有角色/岗位引用当前权限 count %d`, len(rids))
	}

	op, err := org.PermissionByID(id)
	if err != nil {
		return err
	}

	dq := ldap.NewDelRequest(op[`dn`].(string), nil)

	return org.Del(dq)
}

// PermissionByType all permission which contain this dolorestype
func (org *Organization) PermissionByType(dtype string, isUnit bool) ([]string, error) {
	sq := ldap.NewSearchRequest(org.parentDN(permissionCategory(isUnit)),
		ldap.ScopeSingleLevel,
		ldap.DerefAlways, 0, 0, false,
		fmt.Sprintf(`(rbacType=%s)`, dtype),
		[]string{`id`}, nil)

	sr, err := org.Search(sq)
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, entry := range sr.Entries {
		ids = append(ids, entry.GetAttributeValue(`id`))
	}

	return ids, nil
}

// Permissions in ldap
func (org *Organization) Permissions(isUnit bool, pageSize uint32, cookie []byte) (*SearchResult, error) {
	return org.searchPermission(``, org.parentDN(permissionCategory(isUnit)), pageSize, cookie)
}

// PermissionByIDs in ldap
func (org *Organization) PermissionByIDs(ids []string) (*SearchResult, error) {
	filter, err := sqConvertIDsToFilter(ids)
	if err != nil {
		return nil, err
	}
	dn := fmt.Sprintf(`ou=permission, %s`, org.subffix)

	return org.searchPermission(filter, dn, 0, nil)
}

// PermissionByID in ldap
func (org *Organization) PermissionByID(id string) (map[string]interface{}, error) {

	rs, e := org.PermissionByIDs([]string{id})
	if e != nil {
		return nil, e
	}
	if len(rs.Data) == 0 {
		return nil, fmt.Errorf(`[%s] permission doesn't exist`, id)
	}
	if len(rs.Data) > 1 {
		return nil, fmt.Errorf(`[%s] found many permission`, id)
	}

	p := rs.Data[0]
	dn := p[`dn`].(string)
	p[`isUnit`] = strings.Contains(dn, `ou=unit`)

	return p, nil
}
