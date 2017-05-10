package organization

import (
	"errors"
	"fmt"

	ldap "gopkg.in/ldap.v2"
)

// AddUnit to ldap
func (org *Organization) AddUnit(parentID, name, description, utypeID string) error {

	// 验证参数有效性
	if len(name) == 0 {
		return errors.New(`unit name must be not nil`)
	}

	// 部门类型是不是正确
	if len(utypeID) > 0 {
		ts, _ := org.TypeByIDs([]string{utypeID}, true)
		if len(ts) != 1 {
			return errors.New(`invalid utype`)
		}
	}

	oid := generatorOID()

	var dn string
	if len(parentID) == 0 {
		dn = org.dn(oid, unit)
	} else {
		_, err := org.UnitByID(parentID)
		if err != nil {
			return errors.New(`parent must be not nil`)
		}
		dn = fmt.Sprintf(`oid=%s,oid=%s,%s`, oid, parentID, org.parentDN(unit))
	}

	aq := ldap.NewAddRequest(dn)

	aq.Attribute(`objectClass`, []string{`organizationalUnit`, `unitExtended`, `top`})
	aq.Attribute(`ou`, []string{name})
	if len(description) > 0 {
		aq.Attribute(`description`, []string{description})
	}
	if len(utypeID) > 0 {
		aq.Attribute(`rbacType`, []string{utypeID})
	}

	return org.l.Add(aq)
}

// UnitByID ...
func (org *Organization) UnitByID(id string) (map[string]interface{}, error) {

	if len(id) == 0 {
		return nil, errors.New(`id must be not empty`)
	}

	result, err := org.search(org.unitSC(fmt.Sprintf(`(oid=%s)`, id), mapper{
		`ou`: `name`,
	}))

	if err != nil {
		return nil, err
	}
	if len(result) == 1 {
		return result[0], nil
	}

	return nil, errors.New(`404 Not Found`)
}
