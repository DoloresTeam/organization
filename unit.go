package organization

import (
	"errors"
	"fmt"

	ldap "gopkg.in/ldap.v2"
)

// AddUnit to ldap
func (org *Organization) AddUnit(parentID string, info map[string][]string) error {

	id := generatorID()

	var dn string
	if len(parentID) == 0 {
		dn = org.dn(id, unit)
	} else {
		_, err := org.UnitByIDs([]string{parentID})
		if err != nil {
			return errors.New(`parent must be not nil`)
		}
		dn = fmt.Sprintf(`id=%s,id=%s,%s`, id, parentID, org.parentDN(unit))
	}

	aq := ldap.NewAddRequest(dn)

	aq.Attribute(`objectClass`, []string{`organizationalUnit`, `unit`, `top`})

	for k, v := range info {
		aq.Attribute(k, v)
	}

	return org.l.Add(aq)
}

// UnitByIDs ...
func (org *Organization) UnitByIDs(ids []string) ([]map[string]interface{}, error) {

	filter, err := scConvertIDsToFilter(ids)
	if err != nil {
		return nil, err
	}

	result, err := org.search(org.unitSC(filter, true))

	if err != nil {
		return nil, err
	}

	return result, nil
}

// AllUnit ...
func (org *Organization) AllUnit() ([]map[string]interface{}, error) {

	return org.search(org.unitSC(``, false))
}
