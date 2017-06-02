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

	id := generatorID()
	dn := org.dn(id, permissionCategory(isUnit))
	aq := ldap.NewAddRequest(dn)

	aq.Attribute(`objectClass`, []string{`permission`, `top`})

	aq.Attribute(`cn`, []string{name})
	aq.Attribute(`description`, []string{description})
	aq.Attribute(`rbacType`, types)

	return id, org.l.Add(aq)
}

// ModifyPermission in ldap
func (org *Organization) ModifyPermission(id, name, description string, types []string, isUnit bool) error {

	dn := org.dn(id, permissionCategory(isUnit))
	mq := ldap.NewModifyRequest(dn)

	if len(name) != 0 {
		mq.Replace(`cn`, []string{name})
	}
	if len(description) != 0 {
		mq.Replace(`description`, []string{description})
	}
	if types != nil {
		mq.Replace(`rbacType`, types)
	}

	err := org.l.Modify(mq)
	if err != nil {
		return err
	}

	return nil
}

// DelPermission in ldap
func (org *Organization) DelPermission(id string, isUnit bool) error {

	rids, err := org.RoleIDsByPermissionID(id, isUnit)
	if err != nil {
		return err
	}

	if len(rids) > 0 {
		return fmt.Errorf(`有角色/岗位引用当前权限 count %d`, len(rids))
	}

	dn := org.dn(id, permissionCategory(isUnit))
	fmt.Print(dn)
	dq := ldap.NewDelRequest(dn, nil)

	return org.l.Del(dq)
}

// PermissionByType all permission which contain this dolorestype
func (org *Organization) PermissionByType(dtype string, isUnit bool) ([]string, error) {
	sq := ldap.NewSearchRequest(org.parentDN(permissionCategory(isUnit)),
		ldap.ScopeSingleLevel,
		ldap.DerefAlways, 0, 0, false,
		fmt.Sprintf(`(rbacType=%s)`, dtype),
		[]string{`id`}, nil)

	sr, err := org.l.Search(sq)
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
	if len(rs.Data) != 1 {
		return nil, errors.New(`found many results`)
	}

	p := rs.Data[0]
	dn := p[`dn`].(string)
	p[`isUnit`] = strings.Contains(dn, `ou=unit`)
	delete(p, `dn`)

	return p, nil
}
