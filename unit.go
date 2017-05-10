package organization

import (
	"errors"
	"fmt"

	ldap "gopkg.in/ldap.v2"
)

// AddUnit to ldap
func (org *Organization) AddUnit(parentID, name, description, utypeID string, info map[string][]string) error {

	// 验证参数有效性
	if len(name) == 0 {
		return errors.New(`unit name must be not empty`)
	}

	// 部门类型是不是正确
	if len(utypeID) > 0 {
		ts, _ := org.TypeByIDs([]string{utypeID}, true)
		if len(ts) != 1 {
			return errors.New(`invalid utype`)
		}
	} else {
		return errors.New(`utypeID must be not empty`)
	}

	oid := generatorOID()

	var dn string
	if len(parentID) == 0 {
		dn = org.dn(oid, unit)
	} else {
		_, err := org.UnitByIDs([]string{parentID})
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

// UnitByIDs ...
func (org *Organization) UnitByIDs(ids []string) ([]map[string]interface{}, error) {

	filter, err := scConvertIDsToFilter(ids)
	if err != nil {
		return nil, err
	}

	result, err := org.search(org.unitSC(filter, mapper{
		`ou`: `name`,
	}))

	if err != nil {
		return nil, err
	}

	return result, nil
}
