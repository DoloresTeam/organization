package organization

import (
	"fmt"

	ldap "gopkg.in/ldap.v2"
)

func (org *Organization) AddType(name, description string, isUnit bool) error {

	dn := org.dn(generatorOID(), typeCategory(isUnit))
	aq := ldap.NewAddRequest(dn)

	aq.Attribute(`objectClass`, []string{`doloresType`, `top`})
	aq.Attribute(`cn`, []string{name})
	aq.Attribute(`description`, []string{description})

	return org.l.Add(aq)
}

func (org *Organization) ModifyType(oid string, name, description string, isUnit bool) error {

	dn := org.dn(generatorOID(), typeCategory(isUnit))
	mq := ldap.NewModifyRequest(dn)

	if len(name) != 0 {
		mq.Add(`cn`, []string{name})
	}
	if len(description) != 0 {
		mq.Add(`description`, []string{description})
	}

	return org.l.Modify(mq)
}

func (org *Organization) DelType(oid string, isUnit bool) error {

	// 删除类型之前需要确认没有任何一个权限有引用这个类型
	pid := org.rbacx.HasRefrenceType(oid, isUnit)
	if len(pid) != 0 {
		return fmt.Errorf(`has other permission: [%s] refrence this type`, pid)
	}

	dn := org.dn(generatorOID(), typeCategory(isUnit))
	dq := ldap.NewDelRequest(dn, nil)

	return org.l.Del(dq)
}

func (org *Organization) AllTypes(isUnit bool) (*ldap.SearchResult, error) {

	sq := ldap.NewSearchRequest(org.parentDN(typeCategory(isUnit)),
		ldap.ScopeSingleLevel,
		ldap.NeverDerefAliases, 0, 0, false,
		`(objectClass=doloresType)`,
		[]string{`oid`, `cn`, `description`}, nil)

	return org.l.Search(sq)
}
