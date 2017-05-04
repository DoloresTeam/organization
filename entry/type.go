package entry

import (
	"errors"
	"fmt"

	ldap "gopkg.in/ldap.v2"
)

const (
	UnitType   = `unit`
	PersonType = `person`
)

func AddType(name, description string, which string) (*ldap.AddRequest, error) {

	aq := ldap.NewAddRequest(dn(generatorID(), which))

	aq.Attribute(`objectClass`, []string{`doloresType`, `top`})
	if len(name) == 0 {
		return nil, errors.New(`name must be not nil`)
	}

	aq.Attribute(`cn`, []string{name})
	aq.Attribute(`description`, []string{description})

	return aq, nil
}

func dn(oid, which string) string {

	return fmt.Sprintf(`oid=%s,ou=%s,%s`, oid, which, subfix)
}
