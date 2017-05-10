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
// info can include custom field
func (org *Organization) AddMember(commonName, realName string, ptypeID string, roleIDs, unitIDs []string, info map[string][]string) error {

	if len(commonName) == 0 {
		return errorWithPropertyName(`commonName`)
	}
	if len(realName) == 0 {
		return errorWithPropertyName(`realName`)
	}
	if len(ptypeID) == 0 {
		return errorWithPropertyName(`ptypeID`)
	}
	if len(roleIDs) == 0 {
		return errorWithPropertyName(`roleIDs`)
	}
	if len(unitIDs) == 0 {
		return errorWithPropertyName(`unitIDs`)
	}

	//
	ps, _ := org.TypeByIDs([]string{ptypeID}, false)
	if len(ps) != 1 {
		return errors.New(`invalidate member type`)
	}

	rs, _ := org.RoleByIDs(roleIDs)
	if len(rs) != len(roleIDs) {
		return errors.New(`invalid roles`)
	}

	us, _ := org.UnitByIDs(unitIDs)
	if len(us) != len(unitIDs) {
		return errors.New(`invalid units`)
	}

	aq := ldap.NewAddRequest(org.dn(generatorID(), member))

	aq.Attribute(`objectClass`, []string{`member`, `memberExtended`, `top`})

	aq.Attribute(`cn`, []string{commonName})
	aq.Attribute(`sn`, []string{realName})

	aq.Attribute(`rbacRole`, roleIDs)
	aq.Attribute(`rbacType`, []string{ptypeID})
	aq.Attribute(`unitID`, unitIDs)

	return org.l.Add(aq)
}

// DelMember by id
func (org *Organization) DelMember(id string) error {
	if len(id) == 0 {
		return errors.New(`member id is empty`)
	}

	dq := ldap.NewDelRequest(fmt.Sprintf(`id=%s,%s`, id, org.parentDN(member)), nil)

	return org.l.Del(dq)
}
