package organization

import (
	"fmt"

	"github.com/DoloresTeam/organization/gorbacx"

	ldap "gopkg.in/ldap.v2"
)

func (org *Organization) AddPermission(name, description string, types []string, isUnit bool) error {

	dn := org.dn(generatorOID(), permissionCategory(isUnit))
	aq := ldap.NewAddRequest(dn)

	aq.Attribute(`objectClass`, []string{`permission`, `top`})

	aq.Attribute(`cn`, []string{name})
	aq.Attribute(`description`, []string{description})
	aq.Attribute(`rbacType`, types)

	return org.l.Add(aq)
}

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

func (org *Organization) DelPermission(oid string, isUnit bool) error {

	rids, err := org.RoleByPermission(oid, isUnit)
	if err != nil {
		return err
	}

	if len(rids) > 0 {
		return fmt.Errorf(`has role reference this permission %s`, rids)
	}

	dn := org.dn(oid, ROLE)
	dq := ldap.NewDelRequest(dn, nil)

	return org.l.Del(dq)
}

func (org *Organization) PermissionByType(dtype string, isUnit bool) ([]string, error) {
	sq := ldap.NewSearchRequest(org.parentDN(permissionCategory(isUnit)),
		ldap.ScopeSingleLevel,
		ldap.NeverDerefAliases, 0, 0, false,
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

func (org *Organization) PermissionByID(oid string, isUnit bool) (*gorbacx.Permission, error) {
	sq := ldap.NewSearchRequest(org.parentDN(permissionCategory(isUnit)),
		ldap.ScopeSingleLevel,
		ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(`(oid=%s)`, oid),
		[]string{`oid`, `rbacType`}, nil)

	sr, err := org.l.Search(sq)
	if err != nil || len(sr.Entries) != 1 {
		return nil, err
	}

	entry := sr.Entries[0]
	p := gorbacx.NewPermission(oid, entry.GetAttributeValues(`rbacType`))

	return p, nil
}
