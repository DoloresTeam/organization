package organization

import (
	"errors"

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

	pids, err := org.PermissionByType(oid, isUnit)
	if err != nil {
		return err
	}
	if len(pids) > 0 {
		return errors.New(`has Permission refrence this type,`)
	}

	dn := org.dn(generatorOID(), typeCategory(isUnit))
	dq := ldap.NewDelRequest(dn, nil)

	return org.l.Del(dq)
}

func (org *Organization) AllType(isUnit bool) ([]map[string]interface{}, error) {
	return org.search(org.typeSC(``, isUnit))
}

func (org *Organization) TypeByIDs(ids []string, isUnit bool) ([]map[string]interface{}, error) {
	filter, err := scConvertIDsToFilter(ids)
	if err != nil {
		return nil, err
	}
	return org.search(org.typeSC(filter, isUnit))
}
