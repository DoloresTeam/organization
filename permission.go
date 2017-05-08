package organization

import (
	"fmt"

	ldap "gopkg.in/ldap.v2"
)

// AddPermission to ldap server
func (org *Organization) AddPermission(name, description string, types []string, isUnit bool) error {

	dn := org.dn(generatorOID(), permissionCategory(isUnit))
	aq := ldap.NewAddRequest(dn)

	aq.Attribute(`objectClass`, []string{`permission`, `top`})

	aq.Attribute(`cn`, []string{name})
	aq.Attribute(`description`, []string{description})
	aq.Attribute(`rbacType`, types)

	return org.l.Add(aq)
}

// ModifyPermission in ldap
func (org *Organization) ModifyPermission(oid, name, description string, types []string, isUnit bool) error {

	dn := org.dn(oid, permissionCategory(isUnit))
	mq := ldap.NewModifyRequest(dn)

	if len(name) != 0 {
		mq.Replace(`cn`, []string{name})
	}
	if len(description) != 0 {
		mq.Replace(`description`, []string{description})
	}
	if len(types) > 0 {
		mq.Replace(`rbacType`, types)
	}

	err := org.l.Modify(mq)
	if err != nil {
		return err
	}

	// 更新 rbacx 内部数据
	p, _ := org.rbacx.PermissionByID(oid, isUnit)
	if p != nil {
		p.Replace(types)
	}

	return nil
}

// DelPermission in ldap
func (org *Organization) DelPermission(oid string, isUnit bool) error {

	rids, err := org.RoleByPermission(oid, isUnit)
	if err != nil {
		return err
	}

	if len(rids) > 0 {
		return fmt.Errorf(`has role reference this permission %s`, rids)
	}

	dn := org.dn(oid, role)
	dq := ldap.NewDelRequest(dn, nil)

	return org.l.Del(dq)
}

// PermissionByType all permission which contain this dolorestype
func (org *Organization) PermissionByType(dtype string, isUnit bool) ([]string, error) {
	sq := ldap.NewSearchRequest(org.parentDN(permissionCategory(isUnit)),
		ldap.ScopeSingleLevel,
		ldap.DerefAlways, 0, 0, false,
		fmt.Sprintf(`(rbacType=%s)`, dtype),
		[]string{`oid`}, nil)

	sr, err := org.l.Search(sq)
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, entry := range sr.Entries {
		ids = append(ids, entry.GetAttributeValue(`oid`))
	}

	return ids, nil
}

// AllPermissions in ldap
func (org *Organization) AllPermissions(isUnit bool) ([]map[string]interface{}, error) {
	return org.search(org.permissionSC(``, isUnit))
}

// PermissionByIDs in ldap
func (org *Organization) PermissionByIDs(ids []string, isUnit bool) ([]map[string]interface{}, error) {
	filter, err := scConvertIDsToFilter(ids)
	if err != nil {
		return nil, err
	}
	return org.search(org.permissionSC(filter, isUnit))
}
