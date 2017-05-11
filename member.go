package organization

import (
	"errors"
	"fmt"

	ldap "gopkg.in/ldap.v2"
)

func errorWithPropertyName(p string) error {
	return fmt.Errorf(`%s must be not empty`, p)
}

// AddMember to ldap server
func (org *Organization) AddMember(info map[string][]string) error {

	aq := ldap.NewAddRequest(org.dn(generatorID(), member))

	aq.Attribute(`objectClass`, []string{`member`, `top`})

	for k, v := range info {
		aq.Attribute(k, v)
	}

	return org.l.Add(aq)
}

// ModifyMember ...
func (org *Organization) ModifyMember(id string, info map[string][]string) error {

	dn := org.dn(id, member)

	mq := ldap.NewModifyRequest(dn)

	for k, v := range info {
		mq.Replace(k, v)
	}

	return org.l.Modify(mq)
}

// DelMember by id
func (org *Organization) DelMember(id string) error {
	if len(id) == 0 {
		return errors.New(`member id is empty`)
	}

	dq := ldap.NewDelRequest(fmt.Sprintf(`id=%s,%s`, id, org.parentDN(member)), nil)

	return org.l.Del(dq)
}
