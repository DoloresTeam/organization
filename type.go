package organization

import (
	"errors"

	ldap "gopkg.in/ldap.v2"
)

// AddType desgined to add a new dolresType
func (org *Organization) AddType(name, description string, isUnit bool) error {

	dn := org.dn(generatorID(), typeCategory(isUnit))
	aq := ldap.NewAddRequest(dn)

	aq.Attribute(`objectClass`, []string{`doloresType`, `top`})
	aq.Attribute(`cn`, []string{name})
	aq.Attribute(`description`, []string{description})

	return org.l.Add(aq)
}

// ModifyType update name or description of doloresType
func (org *Organization) ModifyType(id string, name, description string, isUnit bool) error {

	dn := org.dn(generatorID(), typeCategory(isUnit))
	mq := ldap.NewModifyRequest(dn)

	if len(name) != 0 {
		mq.Add(`cn`, []string{name})
	}
	if len(description) != 0 {
		mq.Add(`description`, []string{description})
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

	dn := org.dn(generatorID(), typeCategory(isUnit))
	dq := ldap.NewDelRequest(dn, nil)

	return org.l.Del(dq)
}

// AllType in ldap server
func (org *Organization) AllType(isUnit bool) ([]map[string]interface{}, error) {
	return org.search(org.typeSC(``, isUnit))
}

// TypeByIDs ...
func (org *Organization) TypeByIDs(ids []string, isUnit bool) ([]map[string]interface{}, error) {
	filter, err := scConvertIDsToFilter(ids)
	if err != nil {
		return nil, err
	}
	return org.search(org.typeSC(filter, isUnit))
}
